package storage

import (
	"fmt"
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

// GPTType is the GUID Partition table type
type GPTType string

// GPTGuid is the UID of the GPT partition to create
type GPTGuid string

// FSType defines the Filesystem of a Partition
type FSType string

// MountOption a option given to a mountpoint
type MountOption string

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
	GPTGuid    GPTGuid

	// Properties from blkid
	Properties map[string]string
}

func (p *Partition) String() string {
	return fmt.Sprintf("%s", p.Device)
}

// Disk is a physical Disk
type Disk struct {
	Device string
	// Partitions to create on this disk, order is preserved
	Partitions []*Partition
}

// DefaultDisk returns a default partitioning scheme for our devices
func DefaultDisk() Disk {
	defaultDisk := Disk{
		Device: "/dev/sda",
		Partitions: []*Partition{
			&Partition{
				Label:      "efi",
				Number:     1,
				MountPoint: "/boot/efi",
				Filesystem: VFAT,
				GPTType:    GPTBoot,
				GPTGuid:    EFISystemPartition,
				Size:       300,
				Properties: make(map[string]string),
			},
			&Partition{
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

	for _, p := range defaultDisk.Partitions {
		p.Device = fmt.Sprintf("%s%d", defaultDisk.Device, p.Number)
	}
	return defaultDisk
}
