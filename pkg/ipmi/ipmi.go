package ipmi

// IPMI Wiki
// https://www.thomas-krenn.com/de/wiki/IPMI_Konfiguration_unter_Linux_mittels_ipmitool

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
)

// Privilege of a ipmitool user
type Privilege int

const (
	Callback      = Privilege(1)
	User          = Privilege(2)
	Operator      = Privilege(3)
	Administrator = Privilege(4)
	OEM           = Privilege(5)
	NoAccess      = Privilege(15)
)

// Ipmi defines methods to interact with ipmi
type Ipmi interface {
	DevicePresent() bool
	Run(arg ...string) (string, error)
	CreateUser(username, password string, uid int, privilege Privilege) error
	GetLanConfig() (LanConfig, error)
	EnableUEFI(bootdev string, persistent bool) error
	UEFIEnabled() bool
	BootOptionsPersitent() bool
}

// Ipmitool is used to query and modify the IPMI based BMC from the host os.
type Ipmitool struct {
	Command string
}

// New create a new Ipmitool with the default command
func New() Ipmi {
	return &Ipmitool{Command: "ipmitool"}
}

// LanConfig contains the config of ipmi.
// tag must contain first column name of ipmitool lan print command output
// to get the second column value be parsed into the field.
type LanConfig struct {
	IP  string `ipmitool:"IP Address"`
	Mac string `ipmitool:"MAC Address"`
}

// DevicePresent returns true if the character device which is required to talk to the BMC is present.
func (i *Ipmitool) DevicePresent() bool {
	const ipmiDevicePrefix = "/dev/ipmi*"
	matches, err := filepath.Glob(ipmiDevicePrefix)
	if err != nil {
		return false
	}
	return len(matches) > 0
}

// Run execute ipmitool
func (i *Ipmitool) Run(arg ...string) (string, error) {
	path, err := exec.LookPath(i.Command)
	if err != nil {
		return "", fmt.Errorf("unable to locate program:%s in path info:%v", i.Command, err)
	}
	cmd := exec.Command(path, arg...)
	output, err := cmd.Output()

	return string(output), err
}

// GetLanConfig returns the LanConfig
func (i *Ipmitool) GetLanConfig() (LanConfig, error) {
	config := LanConfig{}

	cmdOutput, err := i.Run("lan", "print")
	if err != nil {
		return config, fmt.Errorf("unable to execute ipmitool info:%v", err)
	}
	lanConfigMap := getLanConfig(cmdOutput)

	config.from(lanConfigMap)

	return config, nil
}

// CreateUser create a ipmi user with password and privilege level
func (i *Ipmitool) CreateUser(username, password string, uid int, privilege Privilege) error {
	_, err := i.Run("user", "set", "name", string(uid), username)
	if err != nil {
		return fmt.Errorf("unable to create user %s info: %v", username, err)
	}
	_, err = i.Run("user", "set", "password", string(uid), password)
	if err != nil {
		return fmt.Errorf("unable to set password for user %s info: %v", username, err)
	}
	channelnumber := "1"
	_, err = i.Run("channel", "setaccess", channelnumber, string(uid), "link=on", "ipmi=on", "callin=on", fmt.Sprintf("privilege=%d", int(privilege)))
	if err != nil {
		return fmt.Errorf("unable to set privilege for user %s info: %v", username, err)
	}
	_, err = i.Run("user", "enable", string(uid))
	if err != nil {
		return fmt.Errorf("unable to enable user %s info: %v", username, err)
	}
	return nil
}

// EnableUEFI set the firmware to boot with uefi for given bootdev,
// bootdev can be one of pxe|disk
// if persistent is set to true this will last for every subsequent boot, not only the next.
func (i *Ipmitool) EnableUEFI(bootdev string, persistent bool) error {
	options := "options=efiboot"
	if persistent {
		options = options + ",persistent"
	}

	_, err := i.Run("chassis", "bootdev", bootdev, options)
	if err != nil {
		return fmt.Errorf("unable to enable uefi on:%s persistent:%t info:%v", bootdev, persistent, err)
	}
	return nil
}

// UEFIEnabled returns true if the firmware is set to boot with uefi, otherwise false
func (i *Ipmitool) UEFIEnabled() bool {
	return i.matchBootParam("BIOS EFI boot")
}

// BootOptionsPersitent returns true of the boot parameters are set persistent.
func (i *Ipmitool) BootOptionsPersitent() bool {
	return i.matchBootParam("Options apply to all future boots")
}

func (i *Ipmitool) matchBootParam(parameter string) bool {
	// Brainfuck magic dunno where this is documented.
	// some light can be get from https://bugs.launchpad.net/ironic/+bug/1611306
	output, err := i.Run("chassis", "bootparam", "get", "5")
	if err != nil {
		return false
	}
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, parameter) {
			return true
		}
	}
	return false
}

func getLanConfig(cmdOutput string) map[string]string {
	result := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(cmdOutput))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(strings.Join(parts[1:], ""))
		result[key] = value
	}
	return result
}

// from uses reflection to fill the LanConfig struct based on the tags on it.
func (c *LanConfig) from(output map[string]string) {
	val := reflect.ValueOf(c).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		ipmitoolKey := tag.Get("ipmitool")
		valueField.SetString(output[ipmitoolKey])
	}
}
