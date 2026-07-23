package firmware

import (
	"bufio"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/metal-stack/metal-hammer/pkg/os/command"
)

const (
	// iceDriver is the kernel driver of the E810 network card family.
	iceDriver = "ice"

	nvmUpdateLog = "/tmp/nvmupdate.log"

	// nvmUpdateResult is the machine readable XML report the update utility writes with -o.
	// Its top level PowerCycleRequired and RebootRequired counters state whether a card was
	// flashed and how the new firmware is activated, see parseActivation.
	nvmUpdateResult = "/tmp/nvmupdate.xml"

	// nvmUpdateCfg is the inventory configuration of the update utility, it is part
	// of the nvm update package next to the utility itself.
	nvmUpdateCfg = "nvmupdate.cfg"

	// nvmVersionFile contains the version of the packaged firmware, written by the docker build.
	nvmVersionFile = "version"

	// nvmMinorDigits is the width of the minor of an intel nvm version, see normalizeNVMVersion.
	nvmMinorDigits = 2

	// nvmUpdateTimeout bounds the flash of all cards of this machine, a hanging update
	// utility must not block the provisioning of this machine forever. It is generous,
	// a utility which is terminated while it writes an nvm can leave a card behind which
	// is not usable anymore, so this must only ever hit a utility which is really stuck.
	nvmUpdateTimeout = 20 * time.Minute

	// ethtoolTimeout bounds the firmware version lookup of a single interface.
	ethtoolTimeout = 30 * time.Second
)

var _ updater = &intel{}

type intel struct {
	name string
	log  *slog.Logger

	// sysClassNet is where the network interfaces of this machine are listed.
	sysClassNet string

	// nvmDir is where the intel nvm update package is placed in the initrd, see Dockerfile.
	nvmDir string

	// detected reports whether detect() already ran, cards, desiredVersion and detectErr
	// are its result and must only be read after it did.
	detected       bool
	detectErr      error
	cards          []card
	desiredVersion string
}

func newIntel(log *slog.Logger) *intel {
	return &intel{
		name:        "intel nics",
		log:         log,
		sysClassNet: "/sys/class/net",
		nvmDir:      "/intel",
	}
}

// card is a network card detected in this machine.
type card struct {
	// interfaceName is the name of one of the network interfaces of this card, e.g. eth0
	interfaceName string
	// address is the pci address of this card without the pci function, e.g. 0000:51:00
	address string
	// nvmVersion is the normalized version of the firmware flashed onto this card, e.g. 4.60
	nvmVersion string
}

func (i *intel) String() string {
	return i.name
}

func (i *intel) nvmPath(name string) string {
	return filepath.Join(i.nvmDir, name)
}

// applicable returns true if at least one e810 based network card is present and the
// intel nvm update package is part of this initrd.
func (i *intel) applicable(ctx context.Context) bool {
	if err := i.detect(ctx); err != nil {
		i.log.Info("intel", "message", "the firmware of intel nics is not updated on this machine", "reason", err)
		return false
	}

	return len(i.cards) > 0
}

// detect inspects the nvm update package and the cards of this machine once and caches
// the result, so every method can call it and this updater does not depend on the order
// in which its methods are used.
func (i *intel) detect(ctx context.Context) error {
	if !i.detected {
		i.detected = true
		i.detectErr = i.inspect(ctx)
	}

	return i.detectErr
}

// inspect reads the version of the packaged firmware and detects the cards of this machine.
func (i *intel) inspect(ctx context.Context) error {
	// the utility, its inventory configuration and the version file are all written by the
	// docker build, a package which is missing one of them can not be used.
	for _, name := range []string{command.NVMUpdate, nvmUpdateCfg, nvmVersionFile} {
		if _, err := os.Stat(i.nvmPath(name)); err != nil {
			return fmt.Errorf("the nvm update package is not part of this initrd %w", err)
		}
	}

	// the version of the packaged firmware is the version the cards are expected to be on,
	// therefore there is no firmware version to maintain in this repository.
	version, err := os.ReadFile(i.nvmPath(nvmVersionFile))
	if err != nil {
		return fmt.Errorf("unable to read the version of the nvm update package %w", err)
	}

	i.desiredVersion, err = normalizeNVMVersion(strings.TrimSpace(string(version)))
	if err != nil {
		return fmt.Errorf("unable to read the version of the nvm update package %w", err)
	}

	i.cards, err = i.iceCards(ctx)
	if err != nil {
		return fmt.Errorf("unable to detect ice cards %w", err)
	}

	return nil
}

// iceCards returns all network cards which are bound to the ice driver, together with the
// firmware version they are running. A card with multiple ports shows up with one interface
// per port, all ports share the same firmware, so every card is only reported once.
// A card which can not be inspected is logged and skipped, it must not keep the remaining,
// healthy cards of this machine from being updated.
func (i *intel) iceCards(ctx context.Context) ([]card, error) {
	entries, err := os.ReadDir(i.sysClassNet)
	if err != nil {
		return nil, fmt.Errorf("unable to read %s %w", i.sysClassNet, err)
	}

	var cards []card
	seen := map[string]bool{}
	for _, entry := range entries {
		devicePath := filepath.Join(i.sysClassNet, entry.Name(), "device")

		driver, err := os.Readlink(filepath.Join(devicePath, "driver"))
		if err != nil {
			// virtual interfaces like lo do not have a device/driver link at all
			continue
		}
		if filepath.Base(driver) != iceDriver {
			continue
		}

		address, err := pciAddress(devicePath)
		if err != nil {
			i.log.Warn("intel", "message", "skipping interface", "interface", entry.Name(), "error", err)
			continue
		}
		if seen[address] {
			continue
		}

		version, err := i.nvmVersion(ctx, entry.Name())
		if err != nil {
			// another port of this card may still report the version, so this address
			// is only marked as seen once it could be inspected.
			i.log.Warn("intel", "message", "skipping interface", "interface", entry.Name(), "pci", address, "error", err)
			continue
		}
		seen[address] = true

		i.log.Info("intel", "interface", entry.Name(), "pci", address, "nvm version", version)
		cards = append(cards, card{interfaceName: entry.Name(), address: address, nvmVersion: version})
	}

	return cards, nil
}

// pciAddress returns the pci address of a card without the pci function, which differs
// per port, e.g. a device link pointing to 0000:51:00.1 results in 0000:51:00.
func pciAddress(devicePath string) (string, error) {
	device, err := os.Readlink(devicePath)
	if err != nil {
		return "", fmt.Errorf("unable to detect the pci address of %s %w", devicePath, err)
	}

	address, _, found := strings.Cut(filepath.Base(device), ".")
	if !found {
		return "", fmt.Errorf("unable to detect the pci address of %s, %q is no pci address", devicePath, filepath.Base(device))
	}

	return address, nil
}

// nvmVersion returns the firmware version of the card the given interface belongs to.
func (i *intel) nvmVersion(ctx context.Context, interfaceName string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, ethtoolTimeout)
	defer cancel()

	stdout, stderr, err := run(ctx, i.log, command.Ethtool, "-i", interfaceName)
	if err != nil {
		return "", fmt.Errorf("unable to get driver information of %s %w %s", interfaceName, err, strings.TrimSpace(stderr))
	}

	version, err := parseNVMVersion(interfaceName, stdout)
	if err != nil {
		return "", err
	}

	return normalizeNVMVersion(version)
}

// parseNVMVersion parses the output of "ethtool -i <interface>" which looks like:
//
//	driver: ice
//	version: 6.8.0-generic
//	firmware-version: 4.60 0x8001e8e4 1.3653.0
//	expansion-rom-version:
//	bus-info: 0000:51:00.0
//
// the first field of the firmware-version is the nvm version.
func parseNVMVersion(interfaceName, output string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		key, value, found := strings.Cut(scanner.Text(), ":")
		if !found || strings.TrimSpace(key) != "firmware-version" {
			continue
		}

		version, _, _ := strings.Cut(strings.TrimSpace(value), " ")
		if version != "" {
			return version, nil
		}
	}

	return "", fmt.Errorf("unable to detect firmware version of %s", interfaceName)
}

// normalizeNVMVersion brings an intel nvm version into a form which can be compared as semver.
// The minor of an nvm version is a fixed two digit field: ethtool always reports it padded
// ("4.60"), while the version derived from the name of the update package can be unpadded
// ("4.8" for E810_NVMUpdatePackage_v4_8_Linux.tar.gz). Compared as semver the unpadded 4.8
// would wrongly be older than 4.60, therefore the minor is padded before any comparison.
func normalizeNVMVersion(version string) (string, error) {
	major, minor, found := strings.Cut(version, ".")
	if !found || !isDigits(major) || !isDigits(minor) {
		return "", fmt.Errorf("%q is no nvm version, expected a numeric major and minor like 4.80", version)
	}

	for len(minor) < nvmMinorDigits {
		minor += "0"
	}

	return major + "." + minor, nil
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// current returns the lowest firmware version of all e810 cards in this machine,
// so a machine with cards on different firmware versions is considered outdated.
func (i *intel) current(ctx context.Context) (string, error) {
	if err := i.detect(ctx); err != nil {
		return "", err
	}

	var (
		lowest        *semver.Version
		lowestVersion string
	)

	for _, c := range i.cards {
		version, err := semver.NewVersion(c.nvmVersion)
		if err != nil {
			// a card with a version which can not be interpreted is left out of the
			// comparison instead of hiding the versions of all other cards.
			i.log.Warn("intel", "message", "unable to parse firmware version", "interface", c.interfaceName, "version", c.nvmVersion, "error", err)
			continue
		}

		if lowest == nil || version.LessThan(lowest) {
			lowest = version
			lowestVersion = c.nvmVersion
		}
	}

	if lowest == nil {
		return "", fmt.Errorf("unable to detect the firmware version of any %s based network card", iceDriver)
	}

	return lowestVersion, nil
}

// desired is the NVM version of the update package which is baked into this initrd.
func (i *intel) desired(ctx context.Context) (string, error) {
	if err := i.detect(ctx); err != nil {
		return "", err
	}

	if i.desiredVersion == "" {
		return "", fmt.Errorf("the version of the nvm update package is unknown")
	}

	return i.desiredVersion, nil
}

// update flashes the firmware of all intel network cards for which the
// nvm update package contains an image.
// One invocation covers all cards, the utility inventories all adapters itself
// and skips the ones which are already up to date, therefore this must not be
// called once per card.
// It reports the restart the flashed firmware needs to become active, or ActivationNone
// when no card took an update. A card can be behind the packaged version without being
// updatable at all, e.g. an oem locked nvm, and such a card stays behind forever.
func (i *intel) update(ctx context.Context) (Activation, error) {
	if err := i.detect(ctx); err != nil {
		return ActivationNone, err
	}

	// the result of an earlier attempt must not be reported as the result of this run
	for _, f := range []string{nvmUpdateLog, nvmUpdateResult} {
		if err := os.Remove(f); err != nil && !errors.Is(err, os.ErrNotExist) {
			return ActivationNone, fmt.Errorf("unable to remove %s of a previous nvm update %w", f, err)
		}
	}

	ctx, cancel := context.WithTimeout(ctx, nvmUpdateTimeout)
	defer cancel()

	stdout, stderr, err := run(ctx, i.log, i.nvmPath(command.NVMUpdate), "-u", "-s", "-l", nvmUpdateLog, "-o", nvmUpdateResult, "-c", i.nvmPath(nvmUpdateCfg), "-a", i.nvmDir)
	i.log.Info("intel", "nvmupdate stdout", stdout, "nvmupdate stderr", stderr)

	// the log file carries the human readable per card details, print it in any case
	if logfile, logErr := os.ReadFile(nvmUpdateLog); logErr != nil {
		i.log.Warn("intel", "message", "unable to read nvmupdate log", "log", nvmUpdateLog, "error", logErr)
	} else {
		i.log.Info("intel", "nvmupdate log", string(logfile))
	}

	if err != nil {
		return ActivationNone, fmt.Errorf("unable to update intel firmware %w %s", err, strings.TrimSpace(stderr))
	}

	report, err := os.ReadFile(nvmUpdateResult)
	if err != nil {
		return ActivationNone, fmt.Errorf("unable to read the nvm update result %s %w", nvmUpdateResult, err)
	}

	activation, err := parseActivation(report)
	if err != nil {
		return ActivationNone, err
	}

	if activation == ActivationNone {
		i.log.Warn("intel", "message", "the nvm update utility did not flash any card, the firmware of this machine stays as it is")
		return ActivationNone, nil
	}

	i.log.Info("intel", "message", "firmware updated, a restart is required to activate it", "activation", activation)
	return activation, nil
}

// parseActivation reads the XML report the nvm update utility writes with -o and reports
// how the freshly flashed firmware is activated. The report carries two top level counters:
// PowerCycleRequired for firmware which sits in an inactive bank and only becomes active
// once the card lost power, and RebootRequired for firmware which is live after a warm
// reboot.
func parseActivation(report []byte) (Activation, error) {
	// the utility pads the values with spaces, e.g. "<PowerCycleRequired> 1 </...>",
	// encoding/xml trims that before parsing the int.
	var result struct {
		PowerCycleRequired int `xml:"PowerCycleRequired"`
		RebootRequired     int `xml:"RebootRequired"`
	}
	if err := xml.Unmarshal(report, &result); err != nil {
		return ActivationNone, fmt.Errorf("unable to parse nvm update result %w", err)
	}

	switch {
	case result.PowerCycleRequired > 0:
		return ActivationPowerCycle, nil
	case result.RebootRequired > 0:
		return ActivationWarmReboot, nil
	default:
		return ActivationNone, nil
	}
}
