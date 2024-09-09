package install

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/afero"
)

type operatingsystem string

const (
	osUbuntu    = operatingsystem("ubuntu")
	osDebian    = operatingsystem("debian")
	osAlmalinux = operatingsystem("almalinux")
)

func (o operatingsystem) BootloaderID() string {
	switch o {
	case osAlmalinux:
		return string(o)
	case osDebian, osUbuntu:
		return fmt.Sprintf("metal-%s", o)
	default:
		return fmt.Sprintf("metal-%s", o)
	}
}

func (o operatingsystem) SudoGroup() string {
	switch o {
	case osAlmalinux:
		return "wheel"
	case osDebian, osUbuntu:
		return "sudo"
	default:
		return "sudo"
	}
}

func (o operatingsystem) Initramdisk(kernversion string) string {
	switch o {
	case osAlmalinux:
		return fmt.Sprintf("initramfs-%s.img", kernversion)
	case osDebian, osUbuntu:
		return fmt.Sprintf("initrd.img-%s", kernversion)
	default:
		return fmt.Sprintf("initrd.img-%s", kernversion)
	}
}
func (o operatingsystem) NeedUpdateInitRamfs() bool {
	switch o {
	case osAlmalinux:
		return false
	case osDebian, osUbuntu:
		return true
	default:
		return true
	}
}

func (o operatingsystem) GrubInstallCmd() string {
	switch o {
	case osAlmalinux:
		return "" // no execution required
	case osDebian, osUbuntu:
		return "grub-install"
	default:
		return "grub-install"
	}
}

func operatingSystemFromString(s string) (operatingsystem, error) {
	unquoted, err := strconv.Unquote(s)
	if err == nil {
		s = unquoted
	}

	switch operatingsystem(strings.ToLower(s)) {
	case osUbuntu:
		return osUbuntu, nil
	case osDebian:
		return osDebian, nil
	case osAlmalinux:
		return osAlmalinux, nil
	default:
		return operatingsystem(""), fmt.Errorf("unsupported operating system: %s", s)
	}
}

func detectOS(fs afero.Fs) (operatingsystem, error) {
	content, err := afero.ReadFile(fs, "/etc/os-release")
	if err != nil {
		return operatingsystem(""), err
	}

	env := map[string]string{}
	for _, line := range strings.Split(string(content), "\n") {
		k, v, found := strings.Cut(line, "=")
		if found {
			env[k] = v
		}
	}

	if os, ok := env["ID"]; ok {
		oss, err := operatingSystemFromString(os)
		if err != nil {
			return operatingsystem(""), err
		}
		return oss, nil
	}

	return operatingsystem(""), fmt.Errorf("unable to detect OS")
}
