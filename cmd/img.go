package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	log "github.com/inconshreveable/log15"
	"github.com/jaypipes/ghw"
	"github.com/mholt/archiver"
)

var (
	imgCommand       = "/bin/img"
	sgdiskCommand    = "/usr/bin/sgdisk"
	ext4MkFsCommand  = "/sbin/mkfs.ext4"
	ext3MkFsCommand  = "/sbin/mkfs.ext3"
	fat32MkFsCommand = "/sbin/mkfs.vfat"
	mkswapCommand    = "/sbin/mkswap"
	defaultDisk      = Disk{
		Device: "/dev/sda",
		Partitions: []*Partition{

			&Partition{
				Label:      "boot",
				Number:     1,
				MountPoint: "",
				Filesystem: FAT32,
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
	// FAT32 is ised for the UEFI boot partition
	FAT32 = FSType("fat32")
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
	cmd := exec.Command(sgdiskCommand, "-Z")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("sgdisk zap all existing partitions failed", "error", err, "output", output)
	}

	args := make([]string, 0)
	for _, p := range disk.Partitions {
		size := fmt.Sprintf("%dM", p.Size)
		if p.Size == -1 {
			size = "0"
		}
		args = append(args, fmt.Sprintf("-n=%d:0:%s", p.Number, size))
		args = append(args, fmt.Sprintf(`-c=%d:"%s"`, p.Number, p.Label))
		args = append(args, fmt.Sprintf("-t=%d:%s", p.Number, p.GPTType))

		p.Device = fmt.Sprintf("%s%d", disk.Device, p.Number)
	}

	args = append(args, disk.Device)
	log.Info("sgdisk create partitions", "command", args)
	cmd = exec.Command(sgdiskCommand, args...)
	output, err = cmd.Output()
	// FIXME sgdisk return 0 in case of failure, and > 0 if succeed
	if err != nil {
		log.Error("sgdisk creating partitions failed", "error", err, "output", string(output))
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

		if p.MountPoint == "" {
			continue
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
	case EXT4:
		mkfs = ext4MkFsCommand
		args = append(args, "-F")
		if p.Label != "" {
			args = append(args, "-L", p.Label)
		}
	case EXT3:
		mkfs = ext3MkFsCommand
		args = append(args, "-F")
		if p.Label != "" {
			args = append(args, "-L", p.Label)
		}
	case FAT32:
		mkfs = fat32MkFsCommand
		if p.Label != "" {
			args = append(args, "-L", p.Label)
		}
	case SWAP:
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
		return fmt.Errorf("mkfs failed: %s error:%v", string(output), err)
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
	err := downloadFile("/tmp/os.tgz", image)
	if err != nil {
		return fmt.Errorf("unable to pull image %s error: %v", image, err)
	}
	log.Debug("pull image", "image", image)
	return nil
}

// downloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func downloadFile(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// burn a image by calling genuinetools/img unpack to a specific directory
func burn(image string) error {
	log.Info("burn image", "image", image)

	err := archiver.TarGz.Open("/tmp/os.tgz", "rootfs")
	if err != nil {
		return fmt.Errorf("unable to burn image %s error: %v", image, err)
	}
	log.Debug("burn image", "image", image)
	err = os.Remove("/tmp/os.tgz")
	if err != nil {
		log.Warn("burn image unable to remove image source", "error", err)
	}
	return nil
}

// install will execute /install.sh in the pulled docker image which was extracted onto disk
// to finish installation e.g. install mbr, grub, write network and filesystem config
func install(image string) error {
	log.Info("TODO: install image", "image", image)
	return nil
}
