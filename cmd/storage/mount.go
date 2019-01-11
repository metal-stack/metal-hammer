package storage

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"syscall"

	log "github.com/inconshreveable/log15"
)

// MountPartitions mounts all partitions under prefix
func (disk Disk) MountPartitions(prefix string) error {
	log.Info("mount disk", "disk", disk)
	// "/" must be mounted first
	partitions := disk.SortByMountPoint()

	for _, p := range partitions {

		err := p.MkFS()
		if err != nil {
			log.Error("mount partition create filesystem failed", "error", err)
			return errors.Wrap(err, "mount partitions create fs failed")
		}

		err = p.fetchBlockIDProperties()
		if err != nil {
			log.Error("reading blkid properties failed", "error", err)
			return errors.Wrap(err, "reading blkid properties failed")
		}
		log.Info("set partition properties", "device", p.Device, "properties", p.Properties)

		if p.MountPoint == "" {
			continue
		}

		mountPoint := filepath.Join(prefix, p.MountPoint)
		err = os.MkdirAll(mountPoint, os.ModePerm)
		if err != nil {
			log.Error("mount partition create directory", "error", err)
			return errors.Wrap(err, "mount partitions create directory failed")
		}
		log.Info("mount partition", "partition", p.Device, "mountPoint", mountPoint)
		// see man 2 mount
		err = syscall.Mount(p.Device, mountPoint, string(p.Filesystem), 0, "")
		if err != nil {
			log.Error("unable to mount", "partition", p.Device, "mountPoint", mountPoint, "error", err)
			return errors.Wrapf(err, "mount partitions mount: %s to:%s failed", p.Device, mountPoint)
		}
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

type mount struct {
	source string
	target string
	fstype string
	flags  uintptr
	data   string
}

// MountSpecialFilesystems mounts all special filesystems needed by a chroot
func MountSpecialFilesystems(prefix string) error {
	mounts := []mount{
		{source: "proc", target: "/proc", fstype: "proc", flags: 0, data: ""},
		{source: "sys", target: "/sys", fstype: "sysfs", flags: 0, data: ""},
		{source: "tmpfs", target: "/tmp", fstype: "tmpfs", flags: 0, data: ""},
		// /dev is a bind mount, a bind mount must have MS_BIND flags set see man 2 mount
		{source: "/dev", target: "/dev", fstype: "", flags: syscall.MS_BIND, data: ""},
	}

	for _, m := range mounts {
		err := syscall.Mount(m.source, prefix+m.target, m.fstype, m.flags, m.data)
		if err != nil {
			return errors.Wrapf(err, "mounting %s to %s failed", m.source, m.target)
		}
	}
	return nil
}

// UnMountAll will unmount all filesystems
func UnMountAll(prefix string) {
	umounts := [6]string{"/boot/efi", "/proc", "/sys", "/dev", "/tmp", "/"}
	for _, m := range umounts {
		p := prefix + m
		log.Info("unmounting", "mountpoint", p)
		err := syscall.Unmount(p, syscall.MNT_FORCE)
		if err != nil {
			log.Error("unable to unmount", "path", p, "error", err)
		}
	}
}
