package ipmi

// IPMI Wiki
// https://www.thomas-krenn.com/de/wiki/IPMI_Konfiguration_unter_Linux_mittels_ipmitool
//
// Oder:
// https://wiki.hetzner.de/index.php/IPMI

import (
	"bufio"
	"fmt"
	log "github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
)

// Privilege of a ipmitool user
type Privilege int

const (
	// Callback ipmi privilege
	Callback = Privilege(1)
	// User ipmi privilege
	User = Privilege(2)
	// Operator ipmi privilege
	Operator = Privilege(3)
	// Administrator ipmi privilege
	Administrator = Privilege(4)
	// OEM ipmi privilege
	OEM = Privilege(5)
	// NoAccess ipmi privilege
	NoAccess = Privilege(15)
)

// Ipmi defines methods to interact with ipmi
type Ipmi interface {
	DevicePresent() bool
	Run(arg ...string) (string, error)
	CreateUser(username, password string, uid int, privilege Privilege) error
	GetLanConfig() (LanConfig, error)
	EnableUEFI(bootdev Bootdev, persistent bool) error
	UEFIEnabled() bool
	BootOptionsPersistent() bool
	GetSession() (Session, error)
}

// Ipmitool is used to query and modify the IPMI based BMC from the host os.
type Ipmitool struct {
	Command string
}

// LanConfig contains the config of ipmi.
// tag must contain first column name of ipmitool lan print command output
// to get the second column value be parsed into the field.
type LanConfig struct {
	IP  string `ipmitool:"IP Address"`
	Mac string `ipmitool:"MAC Address"`
}

// Session information of the current ipmi session
type Session struct {
	UserID    string `ipmitool:"user id"`
	Privilege string `ipmitool:"privilege level"`
}

// Bootdev specifies from which device to boot
type Bootdev string

const (
	// PXE boot server via PXE
	PXE = Bootdev("pxe")
	// Disk boot server from hard disk
	Disk = Bootdev("disk")
)

// New create a new Ipmitool with the default command
func New() Ipmi {
	return &Ipmitool{Command: "ipmitool"}
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
func (i *Ipmitool) Run(args ...string) (string, error) {
	path, err := exec.LookPath(i.Command)
	if err != nil {
		return "", errors.Wrapf(err, "unable to locate program:%s in path", i.Command)
	}
	cmd := exec.Command(path, args...)
	output, err := cmd.Output()

	log.Debug("run ipmitool", "args", args, "output", string(output), "error", err)
	return string(output), err
}

// GetLanConfig returns the LanConfig
func (i *Ipmitool) GetLanConfig() (LanConfig, error) {
	config := LanConfig{}
	cmdOutput, err := i.Run("lan", "print")
	if err != nil {
		return config, errors.Errorf("unable to execute ipmitool info:%v", cmdOutput)
	}
	lanConfigMap := output2Map(cmdOutput)
	from(config, lanConfigMap)
	return config, nil
}

// GetSession returns the Session info
func (i *Ipmitool) GetSession() (Session, error) {
	session := Session{}
	cmdOutput, err := i.Run("session", "info", "all")
	if err != nil {
		return session, errors.Errorf("unable to execute ipmitool info:%v", cmdOutput)
	}
	sessionMap := output2Map(cmdOutput)
	from(session, sessionMap)
	return session, nil
}

// CreateUser create a ipmi user with password and privilege level
func (i *Ipmitool) CreateUser(username, password string, uid int, privilege Privilege) error {
	out, err := i.Run("user", "set", "name", string(uid), username)
	if err != nil {
		return errors.Errorf("unable to create user %s info: %v", username, out)
	}
	out, err = i.Run("user", "set", "password", string(uid), password)
	if err != nil {
		return errors.Errorf("unable to set password for user %s info: %v", username, out)
	}
	channelnumber := "1"
	out, err = i.Run("channel", "setaccess", channelnumber, string(uid), "link=on", "ipmi=on", "callin=on", fmt.Sprintf("privilege=%d", int(privilege)))
	if err != nil {
		return errors.Errorf("unable to set privilege for user %s info: %v", username, out)
	}
	out, err = i.Run("user", "enable", string(uid))
	if err != nil {
		return errors.Errorf("unable to enable user %s info: %v", username, out)
	}
	return nil
}

// EnableUEFI set the firmware to boot with uefi for given bootdev,
// bootdev can be one of pxe|disk
// if persistent is set to true this will last for every subsequent boot, not only the next.
func (i *Ipmitool) EnableUEFI(bootdev Bootdev, persistent bool) error {
	options := "options=efiboot"
	if persistent {
		options = options + ",persistent"
	}

	out, err := i.Run("chassis", "bootdev", string(bootdev), options)
	if err != nil {
		return errors.Errorf("unable to enable uefi on:%s persistent:%t info:%v", bootdev, persistent, out)
	}
	return nil
}

// UEFIEnabled returns true if the firmware is set to boot with uefi, otherwise false
func (i *Ipmitool) UEFIEnabled() bool {
	enabled := i.matchBootParam("BIOS EFI boot")
	log.Info("ipmi", "uefi enabled", enabled)
	return enabled
}

// BootOptionsPersistent returns true of the boot parameters are set persistent.
func (i *Ipmitool) BootOptionsPersistent() bool {
	persistent := i.matchBootParam("Options apply to all future boots")
	log.Info("ipmi", "boot params persistent", persistent)
	return persistent
}

func (i *Ipmitool) matchBootParam(parameter string) bool {
	// Brainfuck magic dunno where this is documented.
	// some light can be get from https://bugs.launchpad.net/ironic/+bug/1611306
	output, err := i.Run("chassis", "bootparam", "get", "5")
	if err != nil {
		return false
	}
	log.Info("ipmi", "bootparams", output)
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, parameter) {
			return true
		}
	}
	return false
}

func output2Map(cmdOutput string) map[string]string {
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

// from uses reflection to fill a struct based on the tags on it.
func from(target interface{}, input map[string]string) {
	log.Debug("from", "target", target, "input", input)
	val := reflect.ValueOf(target).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		ipmitoolKey := tag.Get("ipmitool")
		valueField.SetString(input[ipmitoolKey])
	}
}
