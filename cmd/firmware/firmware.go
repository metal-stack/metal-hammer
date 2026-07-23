package firmware

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Masterminds/semver/v3"
	"github.com/metal-stack/metal-hammer/pkg/os"
)

// Activation is the kind of restart a freshly flashed firmware needs to become active.
// The values are ordered from cheapest to strongest, so the requirements of several
// updaters can be combined by taking the maximum, see Result.
type Activation int

const (
	// ActivationNone means nothing was flashed, so no restart is required.
	ActivationNone Activation = iota

	// ActivationWarmReboot means the new firmware is live after a kernel reboot.
	ActivationWarmReboot

	// ActivationPowerCycle means the new firmware sits in an inactive bank which only
	// gets activated once the component lost power, a warm reboot is not sufficient.
	ActivationPowerCycle
)

func (a Activation) String() string {
	switch a {
	case ActivationWarmReboot:
		return "reboot"
	case ActivationPowerCycle:
		return "power cycle"
	default:
		return "none"
	}
}

// updater check if a firmware update is required and updates
// the firmware if required.
type updater interface {
	fmt.Stringer

	// applicable reports whether the hardware and the required tooling
	// are present on this machine.
	applicable(ctx context.Context) bool

	// current is the firmware version this machine is running.
	current(ctx context.Context) (string, error)

	// desired is the firmware version this machine should be running.
	desired(ctx context.Context) (string, error)

	// update flashes the firmware and reports the restart the flashed firmware needs
	// to become active. It returns ActivationNone when no component took an update,
	// so there is nothing to activate.
	update(ctx context.Context) (Activation, error)
}

// Firmware take care of firmware management
type Firmware struct {
	updaters []updater
	log      *slog.Logger
}

// New create a new Firmware manager with all Updaters.
// The raidcontroller updater in raid.go is not implemented yet and therefore not registered.
func New(log *slog.Logger) *Firmware {
	return &Firmware{
		updaters: []updater{newIntel(log)},
		log:      log,
	}
}

// Result is the outcome of a run of Update.
type Result struct {
	// Activation is the strongest restart required to activate the firmware that was
	// flashed in this run, i.e. the maximum over all flashed updaters. It is
	// ActivationNone when nothing was flashed.
	Activation Activation

	// NotFlashed names the updaters whose components are behind the desired version
	// but did not take the update. Those are retried on every boot without ever
	// converging, so this is worth reporting to the outside, see UpdateFirmware.
	NotFlashed []string
}

// Update run updates for all firmwares found.
// A firmware which can not be inspected or updated is logged and skipped
func (f *Firmware) Update(ctx context.Context) Result {
	var result Result
	for _, u := range f.updaters {
		if !u.applicable(ctx) {
			f.log.Info("firmware", "name", u.String(), "message", "not applicable, skipping")
			continue
		}

		required, current, desired, err := updateRequired(ctx, u)
		if err != nil {
			f.log.Error("firmware", "name", u.String(), "message", "unable to detect if an update is required", "error", err)
			continue
		}

		f.log.Info("firmware", "name", u.String(), "current", current, "desired", desired, "update required", required)
		if !required {
			continue
		}

		activation, err := u.update(ctx)
		if err != nil {
			f.log.Error("firmware", "name", u.String(), "message", "unable to update", "error", err)
			continue
		}
		if activation == ActivationNone {
			// the components stay on their current version, there is nothing to activate.
			// this is not an error, a component can be behind the shipped version without
			// being updatable, e.g. a vendor locked firmware.
			f.log.Warn("firmware", "name", u.String(), "message", "no component took the update, keeping the current firmware")
			result.NotFlashed = append(result.NotFlashed, u.String())
			continue
		}

		// several updaters may flash in one run, the machine has to satisfy the strongest
		// activation requirement of all of them, so a single power cycle also activates a
		// firmware which would have been happy with a warm reboot.
		result.Activation = max(result.Activation, activation)
	}
	return result
}

// updateRequired reports whether the firmware of this updater is behind the desired version.
// Firmware versions are compared as semver, a version which is ahead of the desired one
// is never downgraded.
func updateRequired(ctx context.Context, u updater) (required bool, current, desired string, err error) {
	current, err = u.current(ctx)
	if err != nil {
		return false, "", "", fmt.Errorf("unable to get current version %w", err)
	}

	desired, err = u.desired(ctx)
	if err != nil {
		return false, current, "", fmt.Errorf("unable to get desired version %w", err)
	}

	currentVersion, err := semver.NewVersion(current)
	if err != nil {
		return false, current, desired, fmt.Errorf("unable to parse current version %q %w", current, err)
	}

	desiredVersion, err := semver.NewVersion(desired)
	if err != nil {
		return false, current, desired, fmt.Errorf("unable to parse desired version %q %w", desired, err)
	}

	return currentVersion.LessThan(desiredVersion), current, desired, nil
}

// run execute a command with arguments, returns its stdout and stderr and an error
func run(ctx context.Context, log *slog.Logger, command string, args ...string) (stdout, stderr string, err error) {
	stdout, stderr, err = os.ExecuteCommandWithOutput(ctx, command, args...)

	log.Debug("run", "command", command, "args", args, "stdout", stdout, "stderr", stderr, "error", err)
	return stdout, stderr, err
}
