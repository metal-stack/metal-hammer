package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	img "git.f-i-ts.de/cloud-native/metal/metal-hammer/cmd/image"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/cmd/storage"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/models"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg"
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
	// Hostname of the device
	Hostname string `yaml:"hostname"`
	// IPAddress is expected to be in the form without mask
	IPAddress string `yaml:"ipaddress"`
	// DeviceUUID is the unique UUID for this device, usually the board serial.
	DeviceUUID string `yaml:"deviceuuid"`
	// must be calculated from the last 4 byte of the IPAddress
	ASN string `yaml:"asn"`
	// SSHPublicKey of the user
	SSHPublicKey string `yaml:"sshpublickey"`
	// Password is the password for the metal user.
	Password string `yaml:"password"`
	// Devmode passes mode of installation.
	Devmode bool `yaml:"devmode"`
}

// Install a given image to the disk by using genuinetools/img
func (h *Hammer) Install(deviceWithToken *models.ModelsMetalDeviceWithPhoneHomeToken) (*pkg.Bootinfo, error) {
	device := deviceWithToken.Device
	phtoken := deviceWithToken.PhoneHomeToken
	image := *device.Allocation.Image.URL

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

	info, err := h.install(prefix, device, *phtoken)
	if err != nil {
		return nil, err
	}

	storage.UnMountAll(prefix)

	return info, nil
}

// install will execute /install.sh in the pulled docker image which was extracted onto disk
// to finish installation e.g. install mbr, grub, write network and filesystem config
func (h *Hammer) install(prefix string, device *models.ModelsMetalDevice, phoneHomeToken string) (*pkg.Bootinfo, error) {
	log.Info("install image", "image", device.Allocation.Image.URL)

	err := h.writeInstallerConfig(device)
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

	err = h.writeUserData(device)
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

	info, err := readBootInfo()
	if err != nil {
		return nil, errors.Wrap(err, "unable to read boot-info.yaml")
	}

	files := []string{info.Kernel, info.Initrd}
	tmp := "/tmp"
	for _, f := range files {
		if f == "" {
			// initrd can be empty.
			continue
		}
		src := path.Join(prefix, f)
		dest := path.Join(tmp, filepath.Base(f))
		_, err := copy(src, dest)
		if err != nil {
			log.Error("could not copy", "src", src, "dest", dest, "error", err)
			return nil, err
		}
	}
	info.Kernel = path.Join(tmp, filepath.Base(info.Kernel))
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

func (h *Hammer) writeUserData(device *models.ModelsMetalDevice) error {
	configdir := path.Join(prefix, "etc", "metal")
	destination := path.Join(configdir, "userdata")

	base64UserData := device.Allocation.UserData
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

func (h *Hammer) writeInstallerConfig(device *models.ModelsMetalDevice) error {
	log.Info("write installation configuration")
	configdir := path.Join(prefix, "etc", "metal")
	err := os.MkdirAll(configdir, 0755)
	if err != nil {
		return errors.Wrapf(err, "mkdir of %s target os failed", configdir)
	}
	destination := path.Join(configdir, "install.yaml")

	var ipaddress string
	var asn int64
	if *device.Allocation.Cidr == "dhcp" {
		ipaddress = *device.Allocation.Cidr
	} else {
		ip, _, err := net.ParseCIDR(*device.Allocation.Cidr)
		if err != nil {
			return errors.Wrap(err, "unable to parse ip from device.ip")
		}

		asn, err = ipToASN(*device.Allocation.Cidr)
		if err != nil {
			return errors.Wrap(err, "unable to parse ip from device.ip")
		}
		ipaddress = ip.String()
	}

	// FIXME
	sshPubkeys := strings.Join(device.Allocation.SSHPubKeys, "\n")
	y := &InstallerConfig{
		Hostname:     *device.Allocation.Hostname,
		SSHPublicKey: sshPubkeys,
		IPAddress:    ipaddress,
		DeviceUUID:   h.Spec.DeviceUUID,
		ASN:          fmt.Sprintf("%d", asn),
		Devmode:      h.Spec.DevMode,
		Password:     h.Spec.ConsolePassword,
	}
	yamlContent, err := yaml.Marshal(y)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(destination, yamlContent, 0600)
}

func readBootInfo() (*pkg.Bootinfo, error) {
	bi, err := ioutil.ReadFile(path.Join(prefix, "etc", "metal", "boot-info.yaml"))
	if err != nil {
		return nil, errors.Wrap(err, "could not read boot-info.yaml")
	}

	info := &pkg.Bootinfo{}
	err = yaml.Unmarshal(bi, info)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal boot-info.yaml")
	}
	return info, nil
}
