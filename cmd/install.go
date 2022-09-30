package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/metal-stack/metal-hammer/cmd/utils"
	"github.com/metal-stack/metal-hammer/pkg/api"

	"github.com/metal-stack/metal-go/api/models"
	img "github.com/metal-stack/metal-hammer/cmd/image"
	"github.com/metal-stack/metal-hammer/cmd/storage"
	"github.com/metal-stack/metal-hammer/pkg/kernel"
	"gopkg.in/yaml.v3"
)

// Install a given image to the disk by using genuinetools/img
func (h *Hammer) Install(machine *models.V1MachineResponse) (*api.Bootinfo, error) {
	s := storage.New(h.log, h.ChrootPrefix, *h.FilesystemLayout)
	err := s.Run()
	if err != nil {
		return nil, err
	}

	image := machine.Allocation.Image.URL

	err = img.NewImage(h.log).Pull(image, h.OsImageDestination)
	if err != nil {
		return nil, err
	}

	err = img.NewImage(h.log).Burn(h.ChrootPrefix, image, h.OsImageDestination)
	if err != nil {
		return nil, err
	}

	info, err := h.install(h.ChrootPrefix, machine)
	if err != nil {
		return nil, err
	}

	// This is executed after installation to be compatible with images which create fstab by their own
	// TODO can be removed and be done in s.Run() once all images do not create fstab anymore
	err = s.CreateFSTab()
	if err != nil {
		return nil, err
	}

	s.Umount()

	return info, nil
}

// install will execute /install.sh in the pulled docker image which was extracted onto disk
// to finish installation e.g. install mbr, grub, write network and filesystem config
func (h *Hammer) install(prefix string, machine *models.V1MachineResponse) (*api.Bootinfo, error) {
	h.log.Infow("install", "image", machine.Allocation.Image.URL)

	err := h.writeInstallerConfig(machine)
	if err != nil {
		return nil, fmt.Errorf("writing configuration install.yaml failed %w", err)
	}

	err = h.writeUserData(machine)
	if err != nil {
		return nil, fmt.Errorf("writing userdata failed %w", err)
	}

	err = h.writeLVMLocalConf()
	if err != nil {
		return nil, err
	}

	h.log.Infow("running /install.sh on", "prefix", prefix)
	err = os.Chdir(prefix)
	if err != nil {
		return nil, fmt.Errorf("unable to chdir to: %s error %w", prefix, err)
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
		return nil, fmt.Errorf("running install.sh in chroot failed %w", err)
	}

	err = os.Chdir("/")
	if err != nil {
		return nil, fmt.Errorf("unable to chdir to: / error %w", err)
	}
	h.log.Infow("finish running install.sh")

	err = os.Remove(path.Join(prefix, "install.sh"))
	if err != nil {
		h.log.Warnw("unable to remove install.sh, ignoring...", "error")
	}

	info, err := kernel.ReadBootinfo(path.Join(prefix, "etc", "metal", "boot-info.yaml"))
	if err != nil {
		return info, fmt.Errorf("unable to read boot-info.yaml %w", err)
	}

	err = h.EnsureBootOrder(info.BootloaderID)
	if err != nil {
		return info, fmt.Errorf("unable to ensure boot order %w", err)
	}

	tmp := "/tmp"
	_, err = utils.Copy(path.Join(prefix, info.Kernel), path.Join(tmp, filepath.Base(info.Kernel)))
	if err != nil {
		h.log.Errorw("install", "could not copy kernel", "error", err)
		return info, err
	}
	info.Kernel = path.Join(tmp, filepath.Base(info.Kernel))

	if info.Initrd == "" {
		return info, nil
	}

	_, err = utils.Copy(path.Join(prefix, info.Initrd), path.Join(tmp, filepath.Base(info.Initrd)))
	if err != nil {
		h.log.Errorw("install", "could not copy initrd", "error", err)
		return info, err
	}
	info.Initrd = path.Join(tmp, filepath.Base(info.Initrd))

	return info, nil
}

// writeLVMLocalConf to make lvm more compatible with os without udevd
// will only be written if lvm is installed in the target image
func (h *Hammer) writeLVMLocalConf() error {
	srclvmlocal := "/etc/lvm/lvmlocal.conf"
	dstlvm := path.Join(h.ChrootPrefix, "/etc/lvm")
	dstlvmlocal := path.Join(h.ChrootPrefix, srclvmlocal)

	_, err := os.Stat(srclvmlocal)
	if os.IsNotExist(err) {
		h.log.Infow("src lvmlocal.conf not present, not creating lvmlocal.conf")
		return nil
	}
	_, err = os.Stat(dstlvm)
	if os.IsNotExist(err) {
		h.log.Infow("dst /etc/lvm not present, not creating lvmlocal.conf")
		return nil
	}

	input, err := os.ReadFile(srclvmlocal)
	if err != nil {
		return fmt.Errorf("unable to read lvmlocal.conf %w", err)
	}

	err = os.WriteFile(dstlvmlocal, input, 0600)
	if err != nil {
		return fmt.Errorf("unable to write lvmlocal.conf %w", err)
	}
	return nil
}

func (h *Hammer) writeUserData(machine *models.V1MachineResponse) error {
	configdir := path.Join(h.ChrootPrefix, "etc", "metal")
	destination := path.Join(configdir, "userdata")

	base64UserData := machine.Allocation.UserData
	if base64UserData != "" {
		userdata, err := base64.StdEncoding.DecodeString(base64UserData)
		if err != nil {
			h.log.Infow("install", "base64 decode of userdata failed, using plain text", err)
			userdata = []byte(base64UserData)
		}
		return os.WriteFile(destination, userdata, 0600)
	}
	return nil
}

func (h *Hammer) writeInstallerConfig(machine *models.V1MachineResponse) error {
	h.log.Infow("write installation configuration")
	configdir := path.Join(h.ChrootPrefix, "etc", "metal")
	err := os.MkdirAll(configdir, 0755)
	if err != nil {
		return fmt.Errorf("mkdir of %s target os failed %w", configdir, err)
	}
	destination := path.Join(configdir, "install.yaml")

	alloc := machine.Allocation

	sshPubkeys := strings.Join(alloc.SSHPubKeys, "\n")
	cmdline, err := kernel.ParseCmdline()
	if err != nil {
		return fmt.Errorf("unable to get kernel cmdline map %w", err)
	}

	console, ok := cmdline["console"]
	if !ok {
		console = "ttyS0"
	}

	y := &api.InstallerConfig{
		Hostname:     *alloc.Hostname,
		SSHPublicKey: sshPubkeys,
		Networks:     alloc.Networks,
		MachineUUID:  h.Spec.MachineUUID,
		Password:     h.Spec.ConsolePassword,
		Console:      console,
		Timestamp:    time.Now().Format(time.RFC3339),
		Nics:         h.onlyNicsWithNeighbors(machine.Hardware.Nics),
		VPN:          alloc.Vpn,
		Role:         *alloc.Role,
	}

	yamlContent, err := yaml.Marshal(y)
	if err != nil {
		return err
	}

	return os.WriteFile(destination, yamlContent, 0600)
}
func (h *Hammer) onlyNicsWithNeighbors(nics []*models.V1MachineNic) []*models.V1MachineNic {
	noNeighbors := func(neighbors []*models.V1MachineNic) bool {
		if len(neighbors) == 0 {
			return true
		}
		for _, n := range neighbors {
			if n.Mac == nil || *n.Mac == "" {
				return true
			}
		}
		return false
	}

	result := []*models.V1MachineNic{}
	for i := range nics {
		nic := nics[i]
		if noNeighbors(nic.Neighbors) {
			continue
		}
		result = append(result, nic)
	}
	h.log.Infow("onlyNicWithNeighbors add", "result", result)
	return result
}
