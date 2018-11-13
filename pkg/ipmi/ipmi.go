package ipmi

// IPMI Wiki
// https://www.thomas-krenn.com/de/wiki/IPMI_Konfiguration_unter_Linux_mittels_ipmitool

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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
	Run(arg ...string) (string, error)
	CreateUser(username, password string, uid int, privilege Privilege) error
	GetLanConfig() (LanConfig, error)
}

type ipmitool struct{}

func New() Ipmi {
	return &ipmitool{}
}

// LanConfig contains the config of ipmi.
// tag must contain first column name of ipmitool lan print command output
// to get the second column value be parsed into the field.
type LanConfig struct {
	IP  string `ipmitool:"IP Address"`
	Mac string `ipmitool:"MAC Address"`
}

// GetLanConfig returns the LanConfig
func (i *ipmitool) GetLanConfig() (LanConfig, error) {
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
func (i *ipmitool) CreateUser(username, password string, uid int, privilege Privilege) error {
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

var ipmitoolCommand = "ipmitool"

func (i *ipmitool) Run(arg ...string) (string, error) {
	path, err := exec.LookPath(ipmitoolCommand)
	if err != nil {
		return "", fmt.Errorf("unable to locate program:%s in path info:%v", ipmitoolCommand, err)
	}
	cmd := exec.Command(path, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	output, err := cmd.Output()

	return string(output), err
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
