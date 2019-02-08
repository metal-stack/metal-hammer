package storage

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"syscall"

	log "github.com/inconshreveable/log15"
)

type mount struct {
	source string
	target string
	fstype string
	flags  uintptr
	data   string
}

var (
	// Order is important and must be preserved.
	specialMounts = []mount{
		{source: "proc", target: "/proc", fstype: "proc", flags: 0, data: ""},
		{source: "sys", target: "/sys", fstype: "sysfs", flags: 0, data: ""},
		{source: "efivarfs", target: "/sys/firmware/efi/efivars", fstype: "efivarfs", flags: 0, data: ""},
		{source: "tmpfs", target: "/tmp", fstype: "tmpfs", flags: 0, data: ""},
		// /dev is a bind mount, a bind mount must have MS_BIND flags set see man 2 mount
		{source: "/dev", target: "/dev", fstype: "", flags: syscall.MS_BIND, data: ""},
	}
	// This slice is filled by MountPartitions to be able to unmount at the end.
	diskMounts = []mount{}
)

// MountPartitions mounts all partitions under prefix
func (disk Disk) MountPartitions(prefix string) error {
	log.Info("mount", "disk", disk)
	// "/" must be mounted first
	partitions := disk.SortByMountPoint()

	for _, p := range partitions {

		err := p.MkFS()
		if err != nil {
			log.Error("create filesystem failed", "error", err)
			return errors.Wrap(err, "create filesystem failed")
		}

		err = p.fetchBlockIDProperties()
		if err != nil {
			log.Error("reading blkid properties failed", "error", err)
			return errors.Wrap(err, "reading blkid properties failed")
		}
		log.Info("partition properties set", "device", p.Device, "properties", p.Properties)

		if p.MountPoint == "" {
			continue
		}

		mountPoint := filepath.Join(prefix, p.MountPoint)
		err = os.MkdirAll(mountPoint, os.ModePerm)
		if err != nil {
			log.Error("create directory failed", "error", err)
			return errors.Wrap(err, "create directory failed")
		}
		log.Info("mount", "source", p.Device, "target", mountPoint, "fstype", p.Filesystem)
		// see man 2 mount
		err = syscall.Mount(p.Device, mountPoint, string(p.Filesystem), 0, "")
		if err != nil {
			log.Error("mount failed", "partition", p.Device, "mountPoint", mountPoint, "error", err)
			return errors.Wrapf(err, "mount: %s to:%s failed", p.Device, mountPoint)
		}
		diskMounts = append(diskMounts, mount{target: p.MountPoint})
	}

	return nil
}

// SortByMountPoint ensures that "/" is the first, which is required for mounting
func (disk *Disk) SortByMountPoint() []*Partition {
	ordered := make([]*Partition, 0)
	for _, p := range disk.Partitions {
		if p.MountPoint == "/" {
			ordered = append(ordered, p)
		}
	}
	for _, p := range disk.Partitions {
		if p.MountPoint != "/" {
			ordered = append(ordered, p)
		}
	}
	return ordered
}

// MountSpecialFilesystems mounts all special filesystems needed by a chroot
func MountSpecialFilesystems(prefix string) error {
	for _, m := range specialMounts {
		mountPoint := filepath.Join(prefix, m.target)
		log.Info("mount", "source", m.source, "target", mountPoint, "fstype", m.fstype, "flags", m.flags, "data", m.data)
		err := syscall.Mount(m.source, mountPoint, m.fstype, m.flags, m.data)
		if err != nil {
			return errors.Wrapf(err, "mounting %s to %s failed", m.source, m.target)
		}
	}
	return nil
}

// UnMountAll will unmount all filesystems
func UnMountAll(prefix string) {
	allmounts := [][]mount{specialMounts, diskMounts}
	for _, mounts := range allmounts {
		for index := len(mounts) - 1; index >= 0; index-- {
			m := filepath.Join(prefix, mounts[index].target)
			log.Info("unmounting", "mountpoint", m)
			err := syscall.Unmount(m, syscall.MNT_FORCE)
			if err != nil {
				log.Error("unable to unmount", "path", m, "error", err)
			}
		}
	}
}
