package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/kexec"
	"golang.org/x/sys/unix"
)

var (
	cmdline     = "/proc/cmdline"
	sysfirmware = "/sys/firmware/efi"
)

// Bootinfo is written by the installer in the target os to tell us
// which kernel, initrd and cmdline must be used for kexec
type Bootinfo struct {
	Initrd  string `yaml:"initrd"`
	Cmdline string `yaml:"cmdline"`
	Kernel  string `yaml:"kernel"`
}

// ParseCmdline will put each key=value pair from /proc/cmdline into a map.
func ParseCmdline() (map[string]string, error) {
	cmdline, err := ioutil.ReadFile(cmdline)
	if err != nil {
		return nil, fmt.Errorf("unable to read %s: %v", cmdline, err)
	}

	cmdLineValues := strings.Split(string(cmdline), " ")
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
	kernel, err := os.OpenFile(info.Kernel, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("could not open kernel: %s error: %v", info.Kernel, err)
	}
	defer kernel.Close()

	ramfs, err := os.OpenFile(info.Initrd, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("could not open initrd: %s error: %v", info.Initrd, err)
	}
	defer ramfs.Close()

	if err := kexec.FileLoad(kernel, ramfs, info.Cmdline); err != nil {
		return fmt.Errorf("could not execute kexec load: %v error: %s", info, err)
	}

	err = kexec.Reboot()
	if err != nil {
		return fmt.Errorf("could not fire kexec reboot info: %v error: %v", info, err)
	}
	return nil
}

// Reboot reboots the the server
func Reboot() error {
	if err := unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART); err != nil {
		return fmt.Errorf("unable to reboot error %v", err.Error())
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
