package storage

import (
	"fmt"

	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/models"
	log "github.com/inconshreveable/log15"
)

const (
	// FAT32 is used for the UEFI boot partition
	FAT32 = FSType("fat32")
	// VFAT is used for the UEFI boot partition
	VFAT = FSType("vfat")
	// EXT3 is usually only used for /boot
	EXT3 = FSType("ext3")
	// EXT4 is the default fs
	EXT4 = FSType("ext4")
	// SWAP is for the swap partition
	SWAP = FSType("swap")

	// GPTBoot EFI Boot Partition
	GPTBoot = GPTType("ef00")
	// GPTLinux Linux Partition
	GPTLinux = GPTType("8300")
	// EFISystemPartition see https://en.wikipedia.org/wiki/EFI_system_partition
	EFISystemPartition = "C12A7328-F81F-11D2-BA4B-00A0C93EC93B"
)

type (
	// GPTType is the GUID Partition table type
	GPTType string

	// GPTGuid is the UID of the GPT partition to create
	GPTGuid string

	// FSType defines the Filesystem of a Partition
	FSType string

	// MountOption a option given to a mountpoint
	MountOption string

	// Partition defines a disk partition
	Partition struct {
		Label        string
		Device       string
		Number       uint
		MountPoint   string
		MountOptions []*MountOption

		// Size in mebiBytes. If negative all available space is used.
		Size       int64
		Filesystem FSType
		GPTType    GPTType
		GPTGuid    GPTGuid

		// Properties from blkid
		Properties map[string]string
	}

	// PrimaryDevice is the device where the installation happens.
	PrimaryDevice struct {
		DeviceName      string
		PartitionPrefix string
	}

	// Disk is a physical Disk
	Disk struct {
		Device string
		// Partitions to create on this disk, order is preserved
		Partitions []*Partition
	}
)

var (
	defaultDisk = Disk{
		Partitions: []*Partition{
			{
				Label:      "efi",
				Number:     1,
				MountPoint: "/boot/efi",
				Filesystem: VFAT,
				GPTType:    GPTBoot,
				GPTGuid:    EFISystemPartition,
				Size:       300,
				Properties: make(map[string]string),
			},
			{
				Label:      "root",
				Number:     2,
				MountPoint: "/",
				Filesystem: EXT4,
				GPTType:    GPTLinux,
				Size:       -1,
				Properties: make(map[string]string),
			},
		},
	}

	clearlinuxDisk = Disk{
		Partitions: []*Partition{
			{
				Label:      "efi",
				Number:     1,
				MountPoint: "/boot",
				Filesystem: VFAT,
				GPTType:    GPTBoot,
				GPTGuid:    EFISystemPartition,
				Size:       300,
				Properties: make(map[string]string),
			},
			{
				Label:      "root",
				Number:     2,
				MountPoint: "/",
				Filesystem: EXT4,
				GPTType:    GPTLinux,
				Size:       -1,
				Properties: make(map[string]string),
			},
		},
	}

	diskByImage = map[string]Disk{
		"default":      defaultDisk,
		"ubuntu-18.04": defaultDisk,
		"ubuntu-18.10": defaultDisk,
		"alpine-3.8":   defaultDisk,
		"alpine-3.9":   defaultDisk,
		"clearlinux":   clearlinuxDisk,
	}

	primaryDeviceBySize = map[string]PrimaryDevice{
		"v1-small-x86": PrimaryDevice{
			DeviceName:      "/dev/sda",
			PartitionPrefix: "",
		},
		"t1-small-x86": PrimaryDevice{
			DeviceName:      "/dev/sda",
			PartitionPrefix: "",
		},
		"s1-large-x86": PrimaryDevice{
			DeviceName:      "/dev/nvme0n1",
			PartitionPrefix: "p",
		},
		"c1-medium-x86": PrimaryDevice{
			DeviceName:      "/dev/sda",
			PartitionPrefix: "",
		},
		"c1-large-x86": PrimaryDevice{
			DeviceName:      "/dev/nvme0n1",
			PartitionPrefix: "p",
		},
		"c1-xlarge-x86": PrimaryDevice{
			DeviceName:      "/dev/nvme0n1",
			PartitionPrefix: "p",
		},
	}
)

// String for a Partition
func (p *Partition) String() string {
	return fmt.Sprintf("%s", p.Device)
}

// GetDisk returns a partitioning scheme for the given image, if image.ID is unknown default is used.
func GetDisk(image *models.ModelsMetalImage, size *models.ModelsMetalSize) Disk {
	log.Info("getdisk", "imageID", *image.ID)
	disk, ok := diskByImage[*image.ID]
	if !ok {
		log.Warn("getdisk", "imageID unknown, using default", *image.ID)
		disk = defaultDisk
	}

	primaryDevice, ok := primaryDeviceBySize[*size.ID]
	if !ok {
		log.Warn("getdisk", "sizeID unknown, using default", *size.ID)
		primaryDevice = PrimaryDevice{
			DeviceName:      "/dev/sda",
			PartitionPrefix: "",
		}
	}
	disk.Device = primaryDevice.DeviceName

	for _, p := range disk.Partitions {
		p.Device = fmt.Sprintf("%s%s%d", disk.Device, primaryDevice.PartitionPrefix, p.Number)
	}
	return disk
}
