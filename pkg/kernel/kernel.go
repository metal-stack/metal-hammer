package kernel

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"strings"
	"time"

	"github.com/metal-stack/metal-hammer/pkg/api"

	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/watchdog"
	"golang.org/x/sys/unix"
	"gopkg.in/yaml.v3"
)

var (
	cmdline     = "/proc/cmdline"
	sysfirmware = "/sys/firmware/efi"
)

// ReadBootinfo read boot-info.yaml which was written by the OS install.sh
// to get all information required to do kexec.
func ReadBootinfo(file string) (*api.Bootinfo, error) {
	bi, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not read boot-info.yaml %w", err)
	}

	info := &api.Bootinfo{}
	err = yaml.Unmarshal(bi, info)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal boot-info.yaml %w", err)
	}
	return info, nil
}

// ParseCmdline will put each key=value pair from /proc/cmdline into a map.
func ParseCmdline() (map[string]string, error) {
	cmdLine, err := os.ReadFile(cmdline)
	if err != nil {
		return nil, fmt.Errorf("unable to read %s %w", cmdLine, err)
	}

	cmdLineValues := strings.Fields(string(cmdLine))
	envmap := make(map[string]string)
	for _, v := range cmdLineValues {
		key, value, found := strings.Cut(v, "=")
		if found {
			key := strings.TrimSpace(key)
			value := strings.TrimSpace(value)
			envmap[key] = value
		}
	}
	return envmap, nil
}

// RunKexec boot into the new kernel given in Bootinfo
func RunKexec(info *api.Bootinfo) error {
	if info != nil {
		kernel, err := os.OpenFile(info.Kernel, os.O_RDONLY, 0)
		if err != nil {
			return fmt.Errorf("could not open kernel: %s %w", info.Kernel, err)
		}
		defer kernel.Close()

		// Initrd can be empty, then we pass an empty pointer to kexec.FileLoad
		var ramfs *os.File
		if info.Initrd != "" {
			ramfs, err = os.OpenFile(info.Initrd, os.O_RDONLY, 0)
			if err != nil {
				return fmt.Errorf("could not open initrd: %s %w", info.Initrd, err)
			}
			defer ramfs.Close()
		}

		if err := kexec.FileLoad(kernel, ramfs, info.Cmdline); err != nil {
			return fmt.Errorf("could not execute kexec load: %v %w", info, err)
		}
	}

	err := kexec.Reboot()
	if err != nil {
		return fmt.Errorf("could not fire kexec reboot info: %v %w", info, err)
	}
	return nil
}

// Reboot reboots the the server
func Reboot() error {
	if err := unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART); err != nil {
		return fmt.Errorf("unable to reboot %w", err)
	}
	return nil
}

// Firmware returns either efi or bios, depending on the boot method.
func Firmware() string {
	_, err := os.Stat(sysfirmware)
	if os.IsNotExist(err) {
		return "bios"
	}
	return "efi"
}

// Watchdog periodically pings kernel software watchdog.
func Watchdog(log *slog.Logger) {
	wd, err := watchdog.Open(watchdog.Dev)
	if err != nil {
		log.Error("watchdog", "disabling hardware watchdog, as it could not be opened.", err)
		return
	}
	defer wd.MagicClose()

	for {
		if err := wd.KeepAlive(); err != nil {
			log.Error("watchdog", "keepalive failed", err)
		}
		time.Sleep(10 * time.Second)
	}
}

// AutoReboot will start a timer and reboot after given duration a random variation spread is added
func AutoReboot(log *slog.Logger, after, spread time.Duration, callback func()) {
	log.Info("autoreboot set to", "after", after.String(), "spread", spread.String())
	spreadMinutes := rand.N(spread) // nolint:gosec
	after = after + spreadMinutes

	log.Info("autoreboot with spread", "after", after.String())
	rebootTimer := time.NewTimer(after)
	<-rebootTimer.C
	log.Info("autoreboot", "timeout reached", "rebooting in 10sec")
	callback()
	time.Sleep(10 * time.Second)
	err := Reboot()
	if err != nil {
		log.Error("autoreboot", "unable to reboot, error", err)
	}
}
