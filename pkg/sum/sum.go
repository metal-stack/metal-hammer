package sum

import (
	"encoding/xml"
	"fmt"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/os/command"
	log "github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type machineType int

const (
	bigTwin machineType = iota
	s2
)

const (
	biosCfgXML       = "biosCfg.xml"
	biosCfgUpdateXML = "biosCfgUpdate.xml"
)

var (
	// SUM does not complain or fail if more boot options are given than actually available
	uefiBootXMLFragmentTemplates = map[machineType]string{
		bigTwin: `<?xml version="1.0" encoding="ISO-8859-1" standalone="yes"?>
<BiosCfg>
  <Menu name="Boot">
    <Setting name="Boot mode select" selectedOption="UEFI" type="Option"/>
    <Setting name="LEGACY to EFI support" selectedOption="Disabled" type="Option"/>
    <Setting name="Boot Option #1" order="1" selectedOption="UEFI_NETWORK_BOOT_OPTION" type="Option"/>
    <Setting name="Boot Option #2" order="1" selectedOption="Disabled" type="Option"/>
    <Setting name="Boot Option #3" order="1" selectedOption="Disabled" type="Option"/>
    <Setting name="Boot Option #4" order="1" selectedOption="Disabled" type="Option"/>
    <Setting name="Boot Option #5" order="1" selectedOption="Disabled" type="Option"/>
    <Setting name="Boot Option #6" order="1" selectedOption="Disabled" type="Option"/>
    <Setting name="Boot Option #7" order="1" selectedOption="Disabled" type="Option"/>
    <Setting name="Boot Option #8" order="1" selectedOption="Disabled" type="Option"/>
    <Setting name="Boot Option #9" order="1" selectedOption="Disabled" type="Option"/>
  </Menu>
</BiosCfg>`,
		s2: `<?xml version="1.0" encoding="ISO-8859-1" standalone="yes"?>
<BiosCfg>
  <Menu name="Boot">
    <Setting name="Boot mode select" selectedOption="UEFI" type="Option"/>
    <Setting name="Legacy to EFI support" selectedOption="Disabled" type="Option"/>
    <Setting name="UEFI Boot Option #1" selectedOption="UEFI_NETWORK_BOOT_OPTION" type="Option"/>
    <Setting name="UEFI Boot Option #2" selectedOption="Disabled" type="Option"/>
    <Setting name="UEFI Boot Option #3" selectedOption="Disabled" type="Option"/>
    <Setting name="UEFI Boot Option #4" selectedOption="Disabled" type="Option"/>
    <Setting name="UEFI Boot Option #5" selectedOption="Disabled" type="Option"/>
    <Setting name="UEFI Boot Option #6" selectedOption="Disabled" type="Option"/>
    <Setting name="UEFI Boot Option #7" selectedOption="Disabled" type="Option"/>
    <Setting name="UEFI Boot Option #8" selectedOption="Disabled" type="Option"/>
    <Setting name="UEFI Boot Option #9" selectedOption="Disabled" type="Option"/>
  </Menu>
</BiosCfg>`,
	}

	bootOrderXMLFragmentTemplates = map[machineType]string{
		bigTwin: `<?xml version="1.0" encoding="ISO-8859-1" standalone="yes"?>
<BiosCfg>
  <Menu name="Boot">
    <Setting name="Boot Option #1" order="1" selectedOption="UEFI Hard Disk:BOOTLOADER_ID" type="Option"/>
    <Setting name="Boot Option #2" order="1" selectedOption="UEFI_NETWORK_BOOT_OPTION" type="Option"/>
    <Setting name="Boot Option #3" order="1" selectedOption="Disabled" type="Option"/>
    <Setting name="Boot Option #4" order="1" selectedOption="Disabled" type="Option"/>
    <Setting name="Boot Option #5" order="1" selectedOption="Disabled" type="Option"/>
    <Setting name="Boot Option #6" order="1" selectedOption="Disabled" type="Option"/>
    <Setting name="Boot Option #7" order="1" selectedOption="Disabled" type="Option"/>
    <Setting name="Boot Option #8" order="1" selectedOption="Disabled" type="Option"/>
    <Setting name="Boot Option #9" order="1" selectedOption="Disabled" type="Option"/>
  </Menu>
</BiosCfg>`,
		s2: `<?xml version="1.0" encoding="ISO-8859-1" standalone="yes"?>
<BiosCfg>
  <Menu name="Boot">
    <Setting name="UEFI Boot Option #1" selectedOption="UEFI Hard Disk:BOOTLOADER_ID" type="Option"/>
    <Setting name="UEFI Boot Option #2" selectedOption="UEFI_NETWORK_BOOT_OPTION" type="Option"/>
    <Setting name="UEFI Boot Option #3" selectedOption="Disabled" type="Option"/>
    <Setting name="UEFI Boot Option #4" selectedOption="Disabled" type="Option"/>
    <Setting name="UEFI Boot Option #5" selectedOption="Disabled" type="Option"/>
    <Setting name="UEFI Boot Option #6" selectedOption="Disabled" type="Option"/>
    <Setting name="UEFI Boot Option #7" selectedOption="Disabled" type="Option"/>
    <Setting name="UEFI Boot Option #8" selectedOption="Disabled" type="Option"/>
    <Setting name="UEFI Boot Option #9" selectedOption="Disabled" type="Option"/>
  </Menu>
</BiosCfg>`,
	}
)

type BiosCfg struct {
	XMLName xml.Name `xml:"BiosCfg"`
	Menu    []struct {
		XMLName xml.Name `xml:"Menu"`
		Name    string   `xml:"name,attr"`
		Setting []struct {
			XMLName        xml.Name `xml:"Setting"`
			Name           string   `xml:"name,attr"`
			Order          string   `xml:"order,attr,omitempty"`
			SelectedOption string   `xml:"selectedOption,attr"`
		}
	}
}

// Sum defines methods to interact with Supermicro Update Manager (SUM)
type Sum interface {
	EnsureUEFIBoot(reboot bool) error
	EnsureBootOrder(bootloaderID string, reboot bool) error
}

// sum is used to modify the BIOS config from the host OS.
type sum struct {
	bootloaderID          string
	biosCfgXML            string
	biosCfg               BiosCfg
	machineType           machineType
	uefiNetworkBootOption string
}

func New() (Sum, error) {
	return &sum{}, nil
}

// EnsureUEFIBoot updates BIOS to UEFI boot.
func (s *sum) EnsureUEFIBoot(reboot bool) error {
	err := s.prepare()
	if err != nil {
		return err
	}

	fragment := uefiBootXMLFragmentTemplates[s.machineType]
	fragment = strings.ReplaceAll(fragment, "UEFI_NETWORK_BOOT_OPTION", s.uefiNetworkBootOption)

	return s.changeBiosCfg(fragment, reboot)
}

// EnsureBootOrder ensures BIOS boot order so that boot from the given allocated OS image is attempted before PXE boot.
func (s *sum) EnsureBootOrder(bootloaderID string, reboot bool) error {
	s.bootloaderID = bootloaderID

	err := s.prepare()
	if err != nil {
		log.Warn("BIOS updates for this machine type are intentionally not supported, skipping EnsureBootOrder", "error", err)
		return nil
	}

	ok := s.bootOrderProperlySet()
	if ok {
		log.Info("sum", "message", "boot order is already configured")
		return nil
	}

	fragment := bootOrderXMLFragmentTemplates[s.machineType]
	fragment = strings.ReplaceAll(fragment, "BOOTLOADER_ID", s.bootloaderID)
	fragment = strings.ReplaceAll(fragment, "UEFI_NETWORK_BOOT_OPTION", s.uefiNetworkBootOption)

	return s.changeBiosCfg(fragment, reboot)
}

func (s *sum) prepare() error {
	err := s.getCurrentBiosCfg()
	if err != nil {
		return err
	}

	err = s.unmarshalBiosCfg()
	if err != nil {
		return errors.Wrapf(err, "unable to unmarshal BIOS configuration:\n%s", s.biosCfgXML)
	}

	s.determineMachineType()

	return s.findUEFINetworkBootOption()
}

func (s *sum) getCurrentBiosCfg() error {
	err := s.execute("-c", "GetCurrentBiosCfg", "--file", biosCfgXML)
	if err != nil {
		return errors.Wrapf(err, "unable to get BIOS configuration via:%s -c GetCurrentBiosCfg --file %s", command.SUM, biosCfgXML)
	}

	bb, err := ioutil.ReadFile(biosCfgXML)
	if err != nil {
		return errors.Wrapf(err, "unable to read file:%s", biosCfgXML)
	}

	s.biosCfgXML = string(bb)
	return nil
}

func (s *sum) determineMachineType() {
	for _, menu := range s.biosCfg.Menu {
		if menu.Name != "Boot" {
			continue
		}
		for _, setting := range menu.Setting {
			if setting.Name == "UEFI Boot Option #1" { // not available in BigTwin BIOS
				s.machineType = s2
				return
			}
		}
	}

	s.machineType = bigTwin
}

func (s *sum) unmarshalBiosCfg() error {
	s.biosCfg = BiosCfg{}
	decoder := xml.NewDecoder(strings.NewReader(s.biosCfgXML))
	decoder.CharsetReader = charset.NewReaderLabel
	return decoder.Decode(&s.biosCfg)
}

func (s *sum) findUEFINetworkBootOption() error {
	for _, menu := range s.biosCfg.Menu {
		if menu.Name != "Boot" {
			continue
		}
		for _, setting := range menu.Setting {
			if strings.Contains(setting.SelectedOption, "UEFI Network") {
				s.uefiNetworkBootOption = setting.SelectedOption
				return nil
			}
		}
	}

	return fmt.Errorf("cannot find PXE boot option in BIOS configuration:\n%s\n", s.biosCfgXML)
}

func (s *sum) bootOrderProperlySet() bool {
	if !s.checkBootOptionAt(1, s.bootloaderID) {
		return false
	}
	if !s.checkBootOptionAt(2, s.uefiNetworkBootOption) {
		return false
	}
	for i := 2; i <= 9; i++ {
		if !s.checkBootOptionAt(i, "Disabled") {
			return false
		}
	}
	return true
}

func (s *sum) checkBootOptionAt(index int, bootOption string) bool {
	for _, menu := range s.biosCfg.Menu {
		if menu.Name != "Boot" {
			continue
		}
		for _, setting := range menu.Setting {
			switch s.machineType {
			case bigTwin:
				if setting.Order != "1" {
					continue
				}
				if setting.Name != fmt.Sprintf("Boot Option #%d", index) {
					continue
				}
			case s2:
				if setting.Name != fmt.Sprintf("UEFI Boot Option #%d", index) {
					continue
				}
			}

			return strings.Contains(setting.SelectedOption, bootOption)
		}
	}

	return false
}

func (s *sum) changeBiosCfg(fragment string, reboot bool) error {
	err := ioutil.WriteFile(biosCfgUpdateXML, []byte(fragment), 0644)
	if err != nil {
		return err
	}

	args := []string{"-c", "ChangeBiosCfg", "--file", biosCfgUpdateXML}
	if reboot {
		args = append(args, "--reboot")
	}

	return s.execute(args...)
}

func (s *sum) execute(args ...string) error {
	cmd := exec.Command(command.SUM, args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid:    uint32(0),
			Gid:    uint32(0),
			Groups: []uint32{0},
		},
	}
	return cmd.Run()
}
