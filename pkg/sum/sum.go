package sum

import (
	"encoding/xml"
	"fmt"
	"github.com/metal-stack/metal-hammer/pkg/kernel"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/inconshreveable/log15"
	"github.com/metal-stack/metal-hammer/pkg/os/command"
	"github.com/pkg/errors"
	"golang.org/x/net/html/charset"
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
  <Menu name="Security">
    <Menu name="SMC Secure Boot Configuration">
      <Setting name="Secure Boot" selectedOption="Enabled" type="Option"/>
    </Menu>
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
  <Menu name="Security">
    <Menu name="SMC Secure Boot Configuration">
      <Setting name="Secure Boot" selectedOption="Enabled" type="Option"/>
    </Menu>
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

type Menu struct {
	XMLName  xml.Name `xml:"Menu"`
	Name     string   `xml:"name,attr"`
	Settings []struct {
		XMLName        xml.Name `xml:"Setting"`
		Name           string   `xml:"name,attr"`
		Order          string   `xml:"order,attr,omitempty"`
		SelectedOption string   `xml:"selectedOption,attr"`
	} `xml:"Setting"`
	Menus []Menu `xml:"Menu"`
}

type BiosCfg struct {
	XMLName xml.Name `xml:"BiosCfg"`
	Menus   []Menu   `xml:"Menu"`
}

// Sum defines methods to interact with Supermicro Update Manager (SUM)
type Sum interface {
	UpdateBIOS() (bool, error)
	EnsureBootOrder(bootloaderID string) error
}

// sum is used to modify the BIOS config from the host OS.
type sum struct {
	bootloaderID          string
	biosCfgXML            string
	biosCfg               BiosCfg
	machineType           machineType
	uefiNetworkBootOption string
	secureBootEnabled     bool
}

func New() (Sum, error) {
	return &sum{}, nil
}

// UpdateBIOS updates BIOS to UEFI boot and disables CSM-module if required.
// If returns whether machine needs to be rebooted or not.
func (s *sum) UpdateBIOS() (bool, error) {
	firmware := kernel.Firmware()
	log.Info("firmware", "is", firmware)

	err := s.prepare()
	if err != nil {
		return false, err
	}

	if firmware == "efi" && (s.machineType == s2 || s.secureBootEnabled) { // we cannot disable csm-support for S2 servers yet
		return false, nil
	}

	fragment := uefiBootXMLFragmentTemplates[s.machineType]
	fragment = strings.ReplaceAll(fragment, "UEFI_NETWORK_BOOT_OPTION", s.uefiNetworkBootOption)

	return true, s.changeBiosCfg(fragment)
}

// EnsureBootOrder ensures BIOS boot order so that boot from the given allocated OS image is attempted before PXE boot.
func (s *sum) EnsureBootOrder(bootloaderID string) error {
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

	return s.changeBiosCfg(fragment)
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
	s.determineSecureBoot()

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
	for _, menu := range s.biosCfg.Menus {
		if menu.Name != "Boot" {
			continue
		}
		for _, setting := range menu.Settings {
			if setting.Name == "UEFI Boot Option #1" { // not available in BigTwin BIOS
				s.machineType = s2
				return
			}
		}
	}

	s.machineType = bigTwin
}

func (s *sum) determineSecureBoot() {
	if s.machineType == s2 { // secure boot option is not available in S2 BIOS
		return
	}
	for _, menu := range s.biosCfg.Menus {
		if menu.Name != "Security" {
			continue
		}
		for _, m := range menu.Menus {
			if m.Name != "SMC Secure Boot Configuration" {
				continue
			}
			for _, setting := range m.Settings {
				if setting.Name == "Secure Boot" {
					s.secureBootEnabled = setting.SelectedOption == "Enabled"
					return
				}
			}
		}
	}
}

func (s *sum) unmarshalBiosCfg() error {
	s.biosCfg = BiosCfg{}
	decoder := xml.NewDecoder(strings.NewReader(s.biosCfgXML))
	decoder.CharsetReader = charset.NewReaderLabel
	return decoder.Decode(&s.biosCfg)
}

func (s *sum) findUEFINetworkBootOption() error {
	for _, menu := range s.biosCfg.Menus {
		if menu.Name != "Boot" {
			continue
		}
		for _, setting := range menu.Settings {
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
	for _, menu := range s.biosCfg.Menus {
		if menu.Name != "Boot" {
			continue
		}
		for _, setting := range menu.Settings {
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

func (s *sum) changeBiosCfg(fragment string) error {
	err := ioutil.WriteFile(biosCfgUpdateXML, []byte(fragment), 0644)
	if err != nil {
		return err
	}

	return s.execute("-c", "ChangeBiosCfg", "--file", biosCfgUpdateXML)
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
