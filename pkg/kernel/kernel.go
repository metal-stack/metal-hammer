package kernel

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/u-root/u-root/pkg/boot/kexec"
	"go.uber.org/zap"
	"golang.org/x/sys/unix"
	"gopkg.in/yaml.v2"
)

var (
	cmdline     = "/proc/cmdline"
	sysfirmware = "/sys/firmware/efi"
)

// Bootinfo is written by the installer in the target os to tell us
// which kernel, initrd and cmdline must be used for kexec
type Bootinfo struct {
	Initrd       string `yaml:"initrd"`
	Cmdline      string `yaml:"cmdline"`
	Kernel       string `yaml:"kernel"`
	BootloaderID string `yaml:"bootloader_id"`
}

// ReadBootinfo read boot-info.yaml which was written by the OS install.sh
// to get all information required to do kexec.
func ReadBootinfo(file string) (*Bootinfo, error) {
	bi, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not read boot-info.yaml %w", err)
	}

	info := &Bootinfo{}
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

	cmdLineValues := strings.Split(string(cmdLine), " ")
	envmap := make(map[string]string)
	for _, v := range cmdLineValues {
		keyValue := strings.Split(v, "=")
		if len(keyValue) == 2 {
			key := strings.TrimSpace(keyValue[0])
			value := strings.TrimSpace(keyValue[1])
			envmap[key] = value
		}
	}
	return envmap, nil
}

// RunKexec boot into the new kernel given in Bootinfo
func RunKexec(info *Bootinfo) error {
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
// from https://github.com/gokrazy/gokrazy
func Watchdog(log *zap.SugaredLogger) {
	f, err := os.OpenFile("/dev/watchdog", os.O_WRONLY, 0)
	if err != nil {
		log.Errorw("watchdog", "disabling hardware watchdog, as it could not be opened.", err)
		return
	}
	defer f.Close()
	// timeout in seconds after which a reboot will be triggered if no write to /dev/watchdog was made.
	timeout := uint32(60)
	if _, _, errno := unix.Syscall(unix.SYS_IOCTL, f.Fd(), unix.WDIOC_SETTIMEOUT, uintptr(unsafe.Pointer(&timeout))); errno != 0 {
		log.Errorw("watchdog", "set timeout failed", errno)
	}

	for {
		if _, _, errno := unix.Syscall(unix.SYS_IOCTL, f.Fd(), unix.WDIOC_KEEPALIVE, 0); errno != 0 {
			log.Errorw("watchdog", "hardware watchdog ping failed", errno)
		}
		time.Sleep(10 * time.Second)
	}
}

// AutoReboot will start a timer and reboot after given duration a random variation spread is added
func AutoReboot(log *zap.SugaredLogger, after, spread time.Duration, callback func()) {
	log.Infow("autoreboot set to", "after", after, "spread", spread)
	spreadMinutes, err := rand.Int(rand.Reader, big.NewInt(int64(spread.Minutes())))
	if err != nil {
		log.Warnw("autoreboot", "unable to calculate spread, disable spread", err)
		spread = time.Duration(0)
	}
	spread = time.Minute * time.Duration(spreadMinutes.Int64())
	after = after + spread

	log.Infow("autoreboot with spread", "after", after)
	rebootTimer := time.NewTimer(after)
	<-rebootTimer.C
	log.Infow("autoreboot", "timeout reached", "rebooting in 10sec")
	callback()
	time.Sleep(10 * time.Second)
	err = Reboot()
	if err != nil {
		log.Errorw("autoreboot", "unable to reboot, error", err)
	}
}
