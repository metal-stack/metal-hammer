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
func (h *hammer) Install(machine *models.V1MachineResponse) (*api.Bootinfo, error) {
	s := storage.New(h.log, h.chrootPrefix, *h.filesystemLayout)
	err := s.Run()
	if err != nil {
		return nil, err
	}

	image := machine.Allocation.Image.URL

	err = img.NewImage(h.log).Pull(image, h.osImageDestination)
	if err != nil {
		return nil, err
	}

	err = img.NewImage(h.log).Burn(h.chrootPrefix, image, h.osImageDestination)
	if err != nil {
		return nil, err
	}

	info, err := h.install(h.chrootPrefix, machine, s.RootUUID)
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
func (h *hammer) install(prefix string, machine *models.V1MachineResponse, rootUUID string) (*api.Bootinfo, error) {
	h.log.Info("install", "image", machine.Allocation.Image.URL)

	err := h.writeInstallerConfig(machine, rootUUID)
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

	installBinary := "/install.sh"
	if fileExists(path.Join(prefix, "install-go")) {
		installBinary = "/install-go"
	}

	h.log.Info("running install", "binary", installBinary, "prefix", prefix)
	err = os.Chdir(prefix)
	if err != nil {
		return nil, fmt.Errorf("unable to chdir to: %s error %w", prefix, err)
	}
	cmd := exec.Command(installBinary)
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
		return nil, fmt.Errorf("running %q in chroot failed %w", installBinary, err)
	}

	err = os.Chdir("/")
	if err != nil {
		return nil, fmt.Errorf("unable to chdir to: / error %w", err)
	}
	h.log.Info("finish running", "binary", installBinary)

	err = os.Remove(path.Join(prefix, installBinary))
	if err != nil {
		h.log.Warn("unable to remove, ignoring", "binary", installBinary, "error", err)
	}

	info, err := kernel.ReadBootinfo(path.Join(prefix, "etc", "metal", "boot-info.yaml"))
	if err != nil {
		return info, fmt.Errorf("unable to read boot-info.yaml %w", err)
	}

	h.log.Info("checking for boot order call", "vendor", h.hal.Board().Vendor.String())
	if h.hal.Board().Vendor.String() != "Giga Computing" {
		h.log.Info("metal-hammer need to ensure boot order", "vendor", h.hal.Board().Vendor.String(), "bootLoaderID", info.BootloaderID)
		err = h.EnsureBootOrder(info.BootloaderID)
		if err != nil {
			return info, fmt.Errorf("unable to ensure boot order %w", err)
		}
	}

	tmp := "/tmp"
	_, err = utils.Copy(path.Join(prefix, info.Kernel), path.Join(tmp, filepath.Base(info.Kernel)))
	if err != nil {
		h.log.Error("could not copy kernel", "error", err)
		return info, err
	}
	info.Kernel = path.Join(tmp, filepath.Base(info.Kernel))

	if info.Initrd == "" {
		return info, nil
	}

	_, err = utils.Copy(path.Join(prefix, info.Initrd), path.Join(tmp, filepath.Base(info.Initrd)))
	if err != nil {
		h.log.Error("could not copy initrd", "error", err)
		return info, err
	}
	info.Initrd = path.Join(tmp, filepath.Base(info.Initrd))

	return info, nil
}

// writeLVMLocalConf to make lvm more compatible with os without udevd
// will only be written if lvm is installed in the target image
func (h *hammer) writeLVMLocalConf() error {
	srclvmlocal := "/etc/lvm/lvmlocal.conf"
	dstlvm := path.Join(h.chrootPrefix, "/etc/lvm")
	dstlvmlocal := path.Join(h.chrootPrefix, srclvmlocal)

	_, err := os.Stat(srclvmlocal) // FIXME use fileExists below
	if os.IsNotExist(err) {
		h.log.Info("src lvmlocal.conf not present, not creating lvmlocal.conf")
		return nil
	}
	_, err = os.Stat(dstlvm) // FIXME use fileExists below
	if os.IsNotExist(err) {
		h.log.Info("dst /etc/lvm not present, not creating lvmlocal.conf")
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

func (h *hammer) writeUserData(machine *models.V1MachineResponse) error {
	configdir := path.Join(h.chrootPrefix, "etc", "metal")
	destination := path.Join(configdir, "userdata")

	base64UserData := machine.Allocation.UserData
	if base64UserData != "" {
		userdata, err := base64.StdEncoding.DecodeString(base64UserData)
		if err != nil {
			h.log.Info("install", "base64 decode of userdata failed, using plain text", err)
			userdata = []byte(base64UserData)
		}
		return os.WriteFile(destination, userdata, 0600)
	}
	return nil
}

func (h *hammer) writeInstallerConfig(machine *models.V1MachineResponse, rootUUiD string) error {
	h.log.Info("write installation configuration")
	configdir := path.Join(h.chrootPrefix, "etc", "metal")
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

	h.log.Info("check command line options", "console", console)

	raidEnabled := false
	if alloc != nil && alloc.Filesystemlayout != nil && len(alloc.Filesystemlayout.Raid) > 0 {
		raidEnabled = true
	}

	h.log.Info("check if raid is enabled", "raid", raidEnabled)

	y := &api.InstallerConfig{
		Hostname:      *alloc.Hostname,
		SSHPublicKey:  sshPubkeys,
		Networks:      alloc.Networks,
		MachineUUID:   h.spec.MachineUUID,
		Password:      h.spec.ConsolePassword,
		Console:       console,
		Timestamp:     time.Now().Format(time.RFC3339),
		Nics:          h.onlyNicsWithNeighbors(machine.Hardware.Nics),
		VPN:           alloc.Vpn,
		Role:          *alloc.Role,
		RaidEnabled:   raidEnabled,
		RootUUID:      rootUUiD,
		FirewallRules: alloc.FirewallRules,
		DNSServers:    alloc.DNSServers,
		NTPServers:    alloc.NtpServers,
	}

	yamlContent, err := yaml.Marshal(y)
	if err != nil {
		return err
	}

	return os.WriteFile(destination, yamlContent, 0600)
}
func (h *hammer) onlyNicsWithNeighbors(nics []*models.V1MachineNic) []*models.V1MachineNic {
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
	h.log.Info("onlyNicWithNeighbors add", "result", result)
	return result
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
