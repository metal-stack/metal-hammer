package cmd

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/metal-stack/metal-hammer/cmd/utils"

	log "github.com/inconshreveable/log15"
	img "github.com/metal-stack/metal-hammer/cmd/image"
	"github.com/metal-stack/metal-hammer/cmd/storage"
	"github.com/metal-stack/metal-hammer/metal-core/models"
	"github.com/metal-stack/metal-hammer/pkg/kernel"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// InstallerConfig contains configuration items which are
// consumed by the install.sh of the individual target OS.
type InstallerConfig struct {
	// Hostname of the machine
	Hostname string `yaml:"hostname"`
	// Networks all networks connected to this machine
	Networks []*models.ModelsV1MachineNetwork `yaml:"networks"`
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
	// Timestamp is the the timestamp of installer config creation.
	Timestamp string `yaml:"timestamp"`
	// Nics are the network interfaces of this machine including their neighbors.
	Nics []*models.ModelsV1MachineNicExtended `yaml:"nics"`
}

// Install a given image to the disk by using genuinetools/img
func (h *Hammer) Install(machine *models.ModelsV1MachineResponse, nics []*models.ModelsV1MachineNicExtended) (*kernel.Bootinfo, error) {
	err := h.Disk.Partition()
	if err != nil {
		return nil, err
	}

	err = h.Disk.MountPartitions(h.ChrootPrefix)
	if err != nil {
		return nil, err
	}

	image := machine.Allocation.Image.URL

	err = img.Pull(image, h.OsImageDestination)
	if err != nil {
		return nil, err
	}

	err = img.Burn(h.ChrootPrefix, image, h.OsImageDestination)
	if err != nil {
		return nil, err
	}

	err = storage.MountSpecialFilesystems(h.ChrootPrefix)
	if err != nil {
		return nil, err
	}

	info, err := h.install(h.ChrootPrefix, machine, nics)
	if err != nil {
		return nil, err
	}

	storage.UnMountAll(h.ChrootPrefix)

	return info, nil
}

// install will execute /install.sh in the pulled docker image which was extracted onto disk
// to finish installation e.g. install mbr, grub, write network and filesystem config
func (h *Hammer) install(prefix string, machine *models.ModelsV1MachineResponse, nics []*models.ModelsV1MachineNicExtended) (*kernel.Bootinfo, error) {
	log.Info("install", "image", machine.Allocation.Image.URL)

	err := h.writeInstallerConfig(machine, nics)
	if err != nil {
		return nil, errors.Wrap(err, "writing configuration install.yaml failed")
	}

	err = h.writeDiskConfig()
	if err != nil {
		return nil, errors.Wrap(err, "writing configuration disk.json failed")
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
	log.Info("finish running install.sh")

	err = os.Remove(path.Join(prefix, "install.sh"))
	if err != nil {
		log.Warn("unable to remove install.sh, ignoring...", "error")
	}

	info, err := kernel.ReadBootinfo(path.Join(prefix, "etc", "metal", "boot-info.yaml"))
	if err != nil {
		return info, errors.Wrap(err, "unable to read boot-info.yaml")
	}

	err = h.EnsureBootOrder(info.BootloaderID)
	if err != nil {
		return info, errors.Wrap(err, "unable to ensure boot order")
	}

	tmp := "/tmp"
	_, err = utils.Copy(path.Join(prefix, info.Kernel), path.Join(tmp, filepath.Base(info.Kernel)))
	if err != nil {
		log.Error("install", "could not copy kernel", "error", err)
		return info, err
	}
	info.Kernel = path.Join(tmp, filepath.Base(info.Kernel))

	if info.Initrd == "" {
		return info, nil
	}

	_, err = utils.Copy(path.Join(prefix, info.Initrd), path.Join(tmp, filepath.Base(info.Initrd)))
	if err != nil {
		log.Error("install", "could not copy initrd", "error", err)
		return info, err
	}
	info.Initrd = path.Join(tmp, filepath.Base(info.Initrd))

	return info, nil
}

func (h *Hammer) writeDiskConfig() error {
	configdir := path.Join(h.ChrootPrefix, "etc", "metal")
	destination := path.Join(configdir, "disk.json")
	j, err := json.MarshalIndent(h.Disk, "", "  ")
	if err != nil {
		return errors.Wrap(err, "unable to marshal to json")
	}
	return ioutil.WriteFile(destination, j, 0600)
}

func (h *Hammer) writeUserData(machine *models.ModelsV1MachineResponse) error {
	configdir := path.Join(h.ChrootPrefix, "etc", "metal")
	destination := path.Join(configdir, "userdata")

	base64UserData := machine.Allocation.UserData
	if base64UserData != "" {
		userdata, err := base64.StdEncoding.DecodeString(base64UserData)
		if err != nil {
			log.Info("install", "base64 decode of userdata failed, using plain text", err)
			userdata = []byte(base64UserData)
		}
		return ioutil.WriteFile(destination, userdata, 0600)
	}
	return nil
}

func (h *Hammer) writeInstallerConfig(machine *models.ModelsV1MachineResponse, nics []*models.ModelsV1MachineNicExtended) error {
	log.Info("write installation configuration")
	configdir := path.Join(h.ChrootPrefix, "etc", "metal")
	err := os.MkdirAll(configdir, 0755)
	if err != nil {
		return errors.Wrapf(err, "mkdir of %s target os failed", configdir)
	}
	destination := path.Join(configdir, "install.yaml")

	alloc := machine.Allocation

	sshPubkeys := strings.Join(alloc.SSHPubKeys, "\n")
	cmdline, err := kernel.ParseCmdline()
	if err != nil {
		return errors.Wrap(err, "unable to get kernel cmdline map")
	}

	console, ok := cmdline["console"]
	if !ok {
		console = "ttyS0"
	}

	y := &InstallerConfig{
		Hostname:     *alloc.Hostname,
		SSHPublicKey: sshPubkeys,
		Networks:     alloc.Networks,
		MachineUUID:  h.Spec.MachineUUID,
		Devmode:      h.Spec.DevMode,
		Password:     h.Spec.ConsolePassword,
		Console:      console,
		Timestamp:    time.Now().Format(time.RFC3339),
		Nics:         nicsWithNeighbors(nics),
	}
	yamlContent, err := yaml.Marshal(y)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(destination, yamlContent, 0600)
}

func nicsWithNeighbors(nics []*models.ModelsV1MachineNicExtended) []*models.ModelsV1MachineNicExtended {
	result := []*models.ModelsV1MachineNicExtended{}
	for _, nic := range nics {
		for _, neigh := range nic.Neighbors {
			if *neigh.Mac != "" {
				result = append(result, nic)
			}
		}
	}
	return result
}
