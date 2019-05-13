package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/pkg/errors"

	img "git.f-i-ts.de/cloud-native/metal/metal-hammer/cmd/image"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/cmd/storage"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/models"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/kernel"
	log "github.com/inconshreveable/log15"
	"gopkg.in/yaml.v2"
)

const (
	prefix             = "/rootfs"
	osImageDestination = "/tmp/os.tgz"
)

// InstallerConfig contains configuration items which are
// consumed by the install.sh of the individual target OS.
type InstallerConfig struct {
	// Hostname of the machine
	Hostname string `yaml:"hostname"`
	// IPAddress is expected to be in the form without mask
	IPAddress string `yaml:"ipaddress"`
	// must be calculated from the last 4 byte of the IPAddress
	ASN string `yaml:"asn"`
	// Networks all networks connected to this machine
	Networks []Network `yaml:"networks"`
	// MachineUUID is the unique UUID for this machine, usually the board serial.
	MachineUUID string `yaml:"machineuuid"`
	// SSHPublicKey of the user
	SSHPublicKey string `yaml:"sshpublickey"`
	// Password is the password for the metal user.
	Password string `yaml:"password"`
	// Devmode passes mode of installation.
	Devmode bool `yaml:"devmode"`
	// Console specifies where the kernel should connect its console to.
	Console string `yaml:"console"`
}

type Network struct {
	Ips       []string `yaml:"ips"`
	Networkid *string  `yaml:"networkid"`
	Primary   *bool    `yaml:"primary"`
	Prefixes  []string `yaml:"prefixes"`
	Vrf       *int64   `yaml:"vrf"`
	ASN       *int64   `yaml:"asn"`
	Nat       *bool    `yaml:"nat"`
}

// Install a given image to the disk by using genuinetools/img
func (h *Hammer) Install(machine *models.ModelsV1MachineWaitResponse) (*kernel.Bootinfo, error) {
	phtoken := machine.PhoneHomeToken
	image := machine.Allocation.Image.URL

	err := h.Disk.Partition()
	if err != nil {
		return nil, err
	}

	err = h.Disk.MountPartitions(prefix)
	if err != nil {
		return nil, err
	}

	err = img.Pull(image, osImageDestination)
	if err != nil {
		return nil, err
	}

	err = img.Burn(prefix, image, osImageDestination)
	if err != nil {
		return nil, err
	}

	err = storage.MountSpecialFilesystems(prefix)
	if err != nil {
		return nil, err
	}

	info, err := h.install(prefix, machine, *phtoken)
	if err != nil {
		return nil, err
	}

	storage.UnMountAll(prefix)

	return info, nil
}

// install will execute /install.sh in the pulled docker image which was extracted onto disk
// to finish installation e.g. install mbr, grub, write network and filesystem config
func (h *Hammer) install(prefix string, machine *models.ModelsV1MachineWaitResponse, phoneHomeToken string) (*kernel.Bootinfo, error) {
	log.Info("install", "image", machine.Allocation.Image.URL)

	err := h.writeInstallerConfig(machine)
	if err != nil {
		return nil, errors.Wrap(err, "writing configuration install.yaml failed")
	}

	err = h.writeDiskConfig()
	if err != nil {
		return nil, errors.Wrap(err, "writing configuration disk.json failed")
	}

	err = h.writePhoneHomeToken(phoneHomeToken)
	if err != nil {
		return nil, errors.Wrap(err, "writing phoneHome.jwt failed")
	}

	err = h.writeUserData(machine)
	if err != nil {
		return nil, errors.Wrap(err, "writing userdata failed")
	}

	log.Info("running /install.sh on", "prefix", prefix)
	err = os.Chdir(prefix)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to chdir to: %s error", prefix)
	}
	cmd := exec.Command("/install.sh")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	// these syscalls are required to execute the command in a chroot env.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid:    uint32(0),
			Gid:    uint32(0),
			Groups: []uint32{0},
		},
		Chroot: prefix,
	}
	if err := cmd.Run(); err != nil {
		return nil, errors.Wrap(err, "running install.sh in chroot failed")
	}

	err = os.Chdir("/")
	if err != nil {
		return nil, errors.Wrap(err, "unable to chdir to: / error")
	}
	log.Info("finish running /install.sh")

	err = os.Remove(path.Join(prefix, "/install.sh"))
	if err != nil {
		log.Warn("unable to remove install.sh, ignoring...", "error")
	}

	info, err := kernel.ReadBootinfo(path.Join(prefix, "etc", "metal", "boot-info.yaml"))
	if err != nil {
		return nil, errors.Wrap(err, "unable to read boot-info.yaml")
	}

	tmp := "/tmp"
	_, err = copy(path.Join(prefix, info.Kernel), path.Join(tmp, filepath.Base(info.Kernel)))
	if err != nil {
		log.Error("install", "could not copy kernel", "error", err)
		return nil, err
	}
	info.Kernel = path.Join(tmp, filepath.Base(info.Kernel))

	if info.Initrd == "" {
		return info, nil
	}

	_, err = copy(path.Join(prefix, info.Initrd), path.Join(tmp, filepath.Base(info.Initrd)))
	if err != nil {
		log.Error("install", "could not copy initrd", "error", err)
		return nil, err
	}
	info.Initrd = path.Join(tmp, filepath.Base(info.Initrd))

	return info, nil
}

func (h *Hammer) writeDiskConfig() error {
	configdir := path.Join(prefix, "etc", "metal")
	destination := path.Join(configdir, "disk.json")
	j, err := json.MarshalIndent(h.Disk, "", "  ")
	if err != nil {
		return errors.Wrap(err, "unable to marshal to json")
	}
	return ioutil.WriteFile(destination, j, 0600)
}

func (h *Hammer) writePhoneHomeToken(phoneHomeToken string) error {
	configdir := path.Join(prefix, "etc", "metal")
	destination := path.Join(configdir, "phoneHome.jwt")
	return ioutil.WriteFile(destination, []byte(phoneHomeToken), 0600)
}

func (h *Hammer) writeUserData(machine *models.ModelsV1MachineWaitResponse) error {
	configdir := path.Join(prefix, "etc", "metal")
	destination := path.Join(configdir, "userdata")

	base64UserData := machine.Allocation.UserData
	if base64UserData != "" {
		userdata, err := base64.StdEncoding.DecodeString(base64UserData)
		if err != nil {
			log.Error("install", "writing userdata failed", err)
			return nil
		}
		return ioutil.WriteFile(destination, userdata, 0600)
	}
	return nil
}

func (h *Hammer) writeInstallerConfig(machine *models.ModelsV1MachineWaitResponse) error {
	log.Info("write installation configuration")
	configdir := path.Join(prefix, "etc", "metal")
	err := os.MkdirAll(configdir, 0755)
	if err != nil {
		return errors.Wrapf(err, "mkdir of %s target os failed", configdir)
	}
	destination := path.Join(configdir, "install.yaml")

	var ipaddress string
	var asn int64
	allocation := machine.Allocation
	var networks []Network
	for _, nw := range allocation.Networks {
		if *nw.Primary && len(nw.Ips) > 0 {
			// Keep IP and ASN for backward compatibility with os install.sh
			// TODO can be removed from InstallConfig struct once install.sh
			// can create all network configuration from the Networks struct.
			ipaddress = nw.Ips[0]
			asn = *nw.Asn
		} else {
			log.Warn("install no default network with ips found")
		}
		network := Network{
			Ips:       nw.Ips,
			Networkid: nw.Networkid,
			Primary:   nw.Primary,
			Prefixes:  nw.Prefixes,
			Vrf:       nw.Vrf,
			ASN:       nw.Asn,
			Nat:       nw.Nat,
		}
		networks = append(networks, network)
	}

	sshPubkeys := strings.Join(machine.Allocation.SSHPubKeys, "\n")
	cmdline, err := kernel.ParseCmdline()
	if err != nil {
		return errors.Wrap(err, "unable to get kernel cmdline map")
	}

	console, ok := cmdline["console"]
	if !ok {
		console = "ttyS0"
	}

	y := &InstallerConfig{
		Hostname:     *machine.Allocation.Hostname,
		SSHPublicKey: sshPubkeys,
		IPAddress:    ipaddress,
		ASN:          fmt.Sprintf("%d", asn),
		Networks:     networks,
		MachineUUID:  h.Spec.MachineUUID,
		Devmode:      h.Spec.DevMode,
		Password:     h.Spec.ConsolePassword,
		Console:      console,
	}
	yamlContent, err := yaml.Marshal(y)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(destination, yamlContent, 0600)
}
