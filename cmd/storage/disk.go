package storage

import (
	"fmt"
	"strings"

	log "github.com/inconshreveable/log15"
	"github.com/metal-stack/metal-hammer/metal-core/models"
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
	// None
	NONE = FSType("none")

	// GPTBoot EFI Boot Partition
	GPTBoot = GPTType("ef00")
	// GPTLinux Linux Partition
	GPTLinux = GPTType("8300")
	// GPTLinux Linux Partition
	GPTLinuxLVM = GPTType("8e00")
	// EFISystemPartition see https://en.wikipedia.org/wiki/EFI_system_partition
	EFISystemPartition = "C12A7328-F81F-11D2-BA4B-00A0C93EC93B"
	// GIB bytes of a Gigabyte
	GIB = int64(1024 * 1024 * 1024)
	// TIB bytes of a Terabyte
	TIB = int64(1024 * GIB)
	// DiskPrefixToIgnore disks with this prefix will not be reported and wiped.
	DiskPrefixToIgnore = "ram"
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
		// Device the name of the disk device visible from kernel side, e.g. sda
		Device string
		// Partitions to create on this disk, order is preserved
		Partitions []*Partition
	}
)

var (
	// FIXME this whole struct should be part of the size, comming with the allocate response.
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
				Size:       5000,
				Properties: make(map[string]string),
			},
			{
				Label:      "varlib",
				Number:     3,
				MountPoint: "/var/lib",
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
	s3LargeDisk = Disk{
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
				Size:       50000,
				Properties: make(map[string]string),
			},
			// Keep room for a additional Partition to be used by LVM
			{
				Label:      "vgroot",
				Number:     3,
				MountPoint: "",
				Filesystem: NONE,
				GPTType:    GPTLinuxLVM,
				Size:       -1,
				Properties: make(map[string]string),
			},
		},
	}
)

// String for a Partition
func (p *Partition) String() string {
	return p.Device
}

// primaryDeviceBySize will configure the disk device where the OS gets installed.
func primaryDeviceBySize(sizeID string, disks []*models.ModelsV1MachineBlockDevice) PrimaryDevice {
	switch sizeID {
	case "v1-small-x86", "t1-small-x86", "s1-large-x86", "c1-medium-x86", "c1-large-x86", "c1-xlarge-x86":
		return PrimaryDevice{DeviceName: "/dev/sda", PartitionPrefix: ""}
	case "nvm-size-x86":
		// Example how to specify disk partitioning on NVME disks if they need be be used as root disk.
		return PrimaryDevice{DeviceName: "/dev/nvme0n1", PartitionPrefix: "p"}
	case "y1-medium-x86":
		return PrimaryDevice{DeviceName: "/dev/nvme0n1", PartitionPrefix: "p"}
	case "s3-large-x86":
		return PrimaryDevice{DeviceName: "/dev/sda", PartitionPrefix: ""}
	default:
		log.Info("getdisk", "sizeID unknown, try to guess disk", sizeID)
		deviceName := guessDisk(disks)
		if deviceName == "" {
			deviceName = "/dev/sda"
		}
		log.Warn("getdisk", "using for OS device", deviceName)
		primaryDevice := PrimaryDevice{
			DeviceName:      deviceName,
			PartitionPrefix: "",
		}
		return primaryDevice
	}
}

// diskByImage based on the distribution choose a partition layout
func diskByImage(imageID string) Disk {
	if strings.HasPrefix(imageID, "ubuntu") {
		return defaultDisk
	}
	if strings.HasPrefix(imageID, "firewall") {
		return defaultDisk
	}
	if strings.HasPrefix(imageID, "alpine") {
		return defaultDisk
	}
	if strings.HasPrefix(imageID, "clearlinux") {
		return clearlinuxDisk
	}
	return defaultDisk
}

// Guess the disk for OS installation
func guessDisk(disks []*models.ModelsV1MachineBlockDevice) string {
	guess := ""
	// skip nvmes and large devices (> 1 TiB) and try to find a single best guess
	for _, d := range disks {
		if strings.Contains(*d.Name, "nvme") {
			continue
		} else if (*d.Size) > TIB {
			continue
		} else if guess != "" {
			log.Warn("getdisk", "guess for OS is ambiguous: %s, %s", guess, *d.Name)
			return ""
		} else {
			guess = fmt.Sprintf("/dev/%s", *d.Name)
		}
	}
	return guess
}

// GetDisk returns a partitioning scheme for the given image, if image.ID is unknown default is used.
func GetDisk(imageID string, size *models.ModelsV1SizeResponse, disks []*models.ModelsV1MachineBlockDevice) Disk {
	log.Info("getdisk", "imageID", imageID)
	disk := diskByImage(imageID)
	// TODO hack, must be moved to metal-api as all other code here as well
	if *size.ID == "s3-large-x86" {
		disk = s3LargeDisk
	}

	primaryDevice := primaryDeviceBySize(*size.ID, disks)
	disk.Device = primaryDevice.DeviceName

	for _, p := range disk.Partitions {
		p.Device = fmt.Sprintf("%s%s%d", disk.Device, primaryDevice.PartitionPrefix, p.Number)
	}
	return disk
}
