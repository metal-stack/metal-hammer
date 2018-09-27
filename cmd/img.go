package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	log "github.com/inconshreveable/log15"
	"github.com/jaypipes/ghw"
)

var (
	imgCommand      = "/bin/img"
	sgdiskCommand   = "/usr/bin/sgdisk"
	ext4MkFsCommand = "/sbin/mkfs.ext4"
	ext3MkFsCommand = "/sbin/mkfs.ext3"
	mkswapCommand   = "/sbin/mkswap"
	defaultDisk     = Disk{
		Device: "/dev/sda",
		Partitions: []*Partition{

			&Partition{
				Label:      "boot",
				Number:     1,
				MountPoint: "/boot",
				Filesystem: EXT3,
				GPTType:    GPTBoot,
				Size:       100,
			},
			&Partition{
				Label:      "root",
				Number:     2,
				MountPoint: "/",
				Filesystem: EXT4,
				GPTType:    GPTLinux,
				Size:       -1,
			},
		},
	}
)

const (
	// EXT3 is usually only used for /boot
	EXT3 = FSType("ext3")
	// EXT4 is the default fs
	EXT4 = FSType("ext4")
	// SWAP is for the swap partition
	SWAP = FSType("swap")

	// GPTBoot EFI Boot Partition
	GPTBoot = GPTType("ef02")
	// GPTLinux Linux Partition
	GPTLinux = GPTType("8300")
)

// GPTType is the GUID Partition table type
type GPTType string

// FSType defines the Filesystem of a Partition
type FSType string

// Partition defines a disk partition
type Partition struct {
	Label        string
	Device       string
	Number       uint
	MountPoint   string
	MountOptions []*MountOption

	// Size in mebiBytes. If negative all available space is used.
	Size       int64
	Filesystem FSType
	GPTType    GPTType
}

// MountOption a option given to a mountpoint
type MountOption string

// Disk is a physical Disk
type Disk struct {
	Device     string
	SectorSize int64
	// Partitions to create on this disk, order is preserved
	Partitions []*Partition
}

// Install a given image to the disk by using genuinetools/img
func Install(image string) error {
	err := wipeDisks()
	if err != nil {
		return err
	}
	err = format(defaultDisk)
	if err != nil {
		return err
	}

	err = mountPartitions("rootfs", defaultDisk)
	if err != nil {
		return err
	}

	err = pull(image)
	if err != nil {
		return err
	}
	err = burn(image)
	if err != nil {
		return err
	}

	err = install(image)
	if err != nil {
		return err
	}
	return nil
}

func wipeDisks() error {
	log.Info("wipe all disks")
	block, err := ghw.Block()
	if err != nil {
		return fmt.Errorf("unable to gather disks: %v", err)
	}
	for _, disk := range block.Disks {
		log.Info("TODO wipe disk", "disk", disk)
	}
	return nil
}

func format(disk Disk) error {
	log.Info("format disk", "disk", disk)

	log.Info("sgdisk zap all existing partitions", "disk", disk)
	cmd := exec.Command(sgdiskCommand, "--zap-all")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("sgdisk zap all existing partitions failed", "error", err, "output", output)
	}

	newPartitionsCommands := make([]string, 0)
	for _, p := range disk.Partitions {
		size := fmt.Sprintf("%dM", p.Size)
		if p.Size == -1 {
			size = "0"
		}
		newPartitionCommand := fmt.Sprintf("--new %d:0:%s --change-name %d:\"%s\" --type %d:%s", p.Number, size, p.Number, p.Label, p.Number, p.GPTType)
		newPartitionsCommands = append(newPartitionsCommands, newPartitionCommand)

		p.Device = fmt.Sprintf("%s%d", disk.Device, p.Number)
	}

	newPartitionsCommands = append(newPartitionsCommands, disk.Device)
	log.Info("sgdisk create partitions", "command", newPartitionsCommands)
	cmd = exec.Command(sgdiskCommand, newPartitionsCommands...)
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Error("sgdisk creating partitions failed", "error", err, "output", output)
	}

	return nil
}

func mountPartitions(prefix string, disk Disk) error {
	log.Info("mount disk", "disk", disk)
	// "/" must be mounted first
	partitions := orderPartitions(disk.Partitions)

	// FIXME error handling
	for _, p := range partitions {
		err := createFilesystem(p)
		if err != nil {
			log.Error("mount partition create filesystem failed", "error", err)
		}

		mountPoint := filepath.Join(prefix, p.MountPoint)
		err = os.MkdirAll(mountPoint, os.ModePerm)
		if err != nil {
			log.Error("mount partition create directory", "error", err)
		}
		log.Info("mount partition", "partition", p.Device, "mountPoint", mountPoint)
		// see man 2 mount
		err = syscall.Mount(p.Device, mountPoint, string(p.Filesystem), 0, "")
		if err != nil {
			log.Error("unable to mount", "partition", p.Device, "mountPoint", mountPoint, "error", err)
		}
	}

	return nil
}

func createFilesystem(p *Partition) error {
	log.Info("create filesystem", "device", p.Device, "filesystem", p.Filesystem)
	mkfs := ""
	var args []string
	switch p.Filesystem {
	case "ext4":
		mkfs = ext4MkFsCommand
		args = append(args, "-F")
		if p.Label != "" {
			args = append(args, "-L", p.Label)
		}
	case "ext3":
		mkfs = ext3MkFsCommand
		args = append(args, "-F")
		if p.Label != "" {
			args = append(args, "-L", p.Label)
		}
	case "swap":
		mkfs = ext3MkFsCommand
		args = append(args, "-f")
		if p.Label != "" {
			args = append(args, "-L", p.Label)
		}
	default:
		return fmt.Errorf("unsupported filesystem format: %q", p.Filesystem)
	}
	args = append(args, p.Device)
	cmd := exec.Command(mkfs, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mkfs failed: %s error:%v", output, err)
	}

	return nil
}

// orderPartitions ensures that "/" is the first, which is required for mounting
func orderPartitions(partitions []*Partition) []*Partition {
	ordered := make([]*Partition, 0)
	for _, p := range partitions {
		if p.MountPoint == "/" {
			ordered = append(ordered, p)
		}
	}
	for _, p := range partitions {
		if p.MountPoint != "/" {
			ordered = append(ordered, p)
		}
	}
	return ordered
}

// pull a image by calling genuinetools/img pull
func pull(image string) error {
	log.Info("pull image", "image", image)
	cmd := exec.Command(imgCommand, "pull", image)
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("unable to pull image %s error message: %v error: %v", image, string(output), err)
	}
	log.Debug("pull image", "output", output, "image", image)
	return nil
}

// burn a image by calling genuinetools/img unpack to a specific directory
func burn(image string) error {
	log.Info("burn image", "image", image)
	cmd := exec.Command(imgCommand, "unpack", image)
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("unable to burn image %s error message: %v error: %v", image, string(output), err)
	}
	log.Debug("burn image", "output", output, "image", image)
	return nil
}

// install will execute /install.sh in the pulled docker image which was extracted onto disk
// to finish installation e.g. install mbr, grub, write network and filesystem config
func install(image string) error {
	log.Info("TODO: install image", "image", image)
	return nil
}