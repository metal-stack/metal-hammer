package cmd

import (
	"fmt"
	"os"
	"os/exec"

	log "github.com/inconshreveable/log15"
	"github.com/jaypipes/ghw"
)

var (
	imgCommand = "/bin/img"
)

const (
	// EXT3 is usually only used for /boot
	EXT3 = FSType("ext3")
	// EXT4 is the default fs
	EXT4 = FSType("ext4")
	// SWAP is for the swap partition
	SWAP = FSType("swap")
)

const (
	// GB GigaByte
	GB = 1024 * 1024 * 1024
)

// FSType defines the Filesystem of a Partition
type FSType string

// Partition defines a disk partition
type Partition struct {
	Label      string
	Number     int
	MountPoint string

	// Size in Bytes. If negative all available space is used.
	Size       int64
	Filesystem FSType
}

// Disk is a physical Disk
type Disk struct {
	Device string
	// Partitions to create on this disk, order is important when mounting
	Partitions []*Partition
}

// Install a given image to the disk by using genuinetools/img
func Install(image string) error {
	err := wipeDisks()
	if err != nil {
		return err
	}

	boot := &Partition{
		Label:      "boot",
		Number:     1,
		MountPoint: "/boot",
		Filesystem: EXT3,
		Size:       1 * GB,
	}
	root := &Partition{
		Label:      "root",
		Number:     2,
		MountPoint: "/",
		Filesystem: EXT4,
		Size:       -1,
	}
	partitions := make([]*Partition, 0)
	partitions = append(partitions, root)
	partitions = append(partitions, boot)
	disk := Disk{
		Device:     "/dev/sda",
		Partitions: partitions,
	}
	err = formatDisk(disk)
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
	return nil
}

func wipeDisks() error {
	log.Info("wipe all disks")
	block, err := ghw.Block()
	if err != nil {
		return fmt.Errorf("unable to gather disks: %v", err)
	}
	for _, disk := range block.Disks {
		log.Info("wipe disk", "disk", disk)
	}
	return nil
}

func formatDisk(disk Disk) error {
	log.Info("format disk", "disk", disk)
	return nil
}

// pull a image by calling genuinetools/img pull
func pull(image string) error {
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
	cmd := exec.Command(imgCommand, "unpack", image)
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("unable to burn image %s error message: %v error: %v", image, string(output), err)
	}
	log.Debug("burn image", "output", output, "image", image)
	return nil
}
