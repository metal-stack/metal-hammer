package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	gos "os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	log "github.com/inconshreveable/log15"
	"github.com/metal-stack/metal-hammer/metal-core/models"
	"github.com/metal-stack/metal-hammer/pkg/os"
	"github.com/metal-stack/metal-hammer/pkg/os/command"
	"github.com/metal-stack/v"
	"github.com/pkg/errors"
)

type Filesystem struct {
	config models.ModelsV1FilesystemLayoutResponse
	// chroot defines the root of the mounts
	chroot string
	// mounts are collected to be able to umount all in reverse order
	mounts       []string
	fstabEntries fstabEntries
	// disk is the legacy disk.json representatio
	// TODO remove once old images are gone
	disk Disk
}

type fstabEntries []fstabEntry

// fstabEntry see man fstab for reference
type fstabEntry struct {
	spec      string
	file      string
	vfsType   string
	mountOpts []string
	freq      uint
	passno    uint
}

func New(chroot string, config models.ModelsV1FilesystemLayoutResponse) *Filesystem {
	return &Filesystem{
		config:       config,
		chroot:       chroot,
		fstabEntries: fstabEntries{},
		disk:         Disk{Device: "legacy", Partitions: []Partition{}},
	}
}

func (f *Filesystem) Run() error {

	err := f.createPartitions()
	if err != nil {
		return fmt.Errorf("create partitions failed:%w", err)
	}

	err = f.createRaids()
	if err != nil {
		return fmt.Errorf("create raids failed:%w", err)
	}

	err = f.createFilesystems()
	if err != nil {
		return fmt.Errorf("create filesystems failed:%w", err)
	}

	err = f.mountFilesystems()
	if err != nil {
		return fmt.Errorf("mount filesystems failed:%w", err)
	}
	err = f.mountSpecialFilesystems()
	if err != nil {
		return fmt.Errorf("mount special filesystems failed:%w", err)
	}

	err = f.createDiskJSON()
	if err != nil {
		return fmt.Errorf("disk.json creation failed:%w", err)
	}
	return nil
}
func (f *Filesystem) Umount() error {
	err := f.umountFilesystems()
	if err != nil {
		return fmt.Errorf("umount filesystems failed:%w", err)
	}
	return nil
}

func (f *Filesystem) createPartitions() error {
	if len(f.config.Disks) == 0 {
		return nil
	}
	for _, disk := range f.config.Disks {
		opts := []string{}

		if disk.Wipeonreinstall != nil && *disk.Wipeonreinstall {
			opts = append(opts, "--zap-all")
		}
		for _, p := range disk.Partitions {
			opts = append(opts, fmt.Sprintf("--new=%d:0:+%dM", *p.Number, *p.Size))
			if p.Label != nil {
				opts = append(opts, fmt.Sprintf("--change-name=%d:%s", *p.Number, *p.Label))
			}
			if p.Gpttype != nil {
				opts = append(opts, fmt.Sprintf("--typecode=%d:%s", *p.Number, *p.Gpttype))
			}
		}
		if disk.Device != nil {
			opts = append(opts, *disk.Device)
			log.Info("sgdisk create partitions", "command", opts)
			err := os.ExecuteCommand(command.SGDisk, opts...)
			if err != nil {
				log.Error("sgdisk creating partitions failed", "error", err)
				return errors.Wrapf(err, "unable to create partitions on %s", *disk.Device)
			}
		}
	}
	return nil
}

func (f *Filesystem) createRaids() error {
	if len(f.config.Raid) == 0 {
		return nil
	}

	for _, raid := range f.config.Raid {
		if raid.Arrayname == nil {
			continue
		}
		spares := int32(0)
		if raid.Spares != nil {
			spares = *raid.Spares
		}
		level := "1"
		if raid.Level != nil {
			level = *raid.Level
		}
		args := []string{
			"--create", *raid.Arrayname,
			"--force",
			"--run",
			"--homehost", "any",
			"--level", level,
			"--raid-devices", fmt.Sprintf("%d", int32(len(raid.Devices))-spares),
		}

		if spares > 0 {
			args = append(args, "--spare-devices", fmt.Sprintf("%d", spares))
		}

		for _, o := range raid.Createoptions {
			args = append(args, string(o))
		}

		args = append(args, raid.Devices...)

		log.Info("create mdadm raid", "args", args)
		err := os.ExecuteCommand(command.MDADM, args...)
		if err != nil {
			log.Error("create mdadm raid", "error", err)
			return errors.Wrapf(err, "unable to create mdadm raid %s", *raid.Arrayname)
		}
	}
	return nil
}

func (f *Filesystem) createFilesystems() error {
	if len(f.config.Filesystems) == 0 {
		return nil
	}

	for _, fs := range f.config.Filesystems {
		if fs.Format == nil || *fs.Format == "tmpfs" {
			continue
		}
		mkfs := ""
		args := []string{}
		args = append(args, fs.Createoptions...)
		switch *fs.Format {
		case "ext3":
			mkfs = command.MKFSExt3
			args = append(args, "-F")
			if fs.Label != nil {
				args = append(args, "-L", *fs.Label)
			}
		case "ext4":
			mkfs = command.MKFSExt4
			args = append(args, "-F")
			if fs.Label != nil {
				args = append(args, "-L", *fs.Label)
			}
		case "swap":
			mkfs = command.MKSwap
			args = append(args, "-f")
			if fs.Label != nil {
				args = append(args, "-L", *fs.Label)
			}
		case "vfat":
			mkfs = command.MKFSVFat
			// There is no force flag for mkfs.vfat, it always destroys any data on
			// the device at which it is pointed.
			if fs.Label != nil {
				args = append(args, "-n", *fs.Label)
			}
		case "none":
			//
		default:
			return fmt.Errorf("unsupported filesystem format: %q", *fs.Format)
		}
		args = append(args, *fs.Device)
		log.Info("create filesystem", "args", args)
		err := os.ExecuteCommand(mkfs, args...)
		if err != nil {
			log.Error("create filesystem failed", "device", *fs.Device, "error", err)
			return errors.Wrapf(err, "unable to create filesystem on %s", *fs.Device)
		}
	}

	return nil
}

func (f *Filesystem) mountFilesystems() error {
	fss := []models.ModelsV1Filesystem{}
	for _, fs := range f.config.Filesystems {
		if fs.Path == nil || *fs.Path == "" {
			continue
		}
		fss = append(fss, *fs)
	}
	sort.Slice(fss, func(i, j int) bool { return depth(*fss[i].Path) < depth(*fss[j].Path) })
	for _, fs := range fss {
		path, err := mountFs(f.chroot, fs)
		if err != nil {
			return err
		}
		f.mounts = append(f.mounts, path)

		passno := uint(2)
		spec := ""
		properties := map[string]string{"UUID": ""}
		if *fs.Format == "tmpfs" {
			spec = *fs.Format
			passno = 0
		} else {
			properties, err = FetchBlockIDProperties(*fs.Device)
			if err != nil {
				return err
			}
			spec = fmt.Sprintf("UUID=%s", properties["UUID"])
		}
		if *fs.Path == "/" {
			passno = 1
		}
		mountOpts := []string{"defaults"}
		if len(fs.Mountoptions) > 0 {
			mountOpts = fs.Mountoptions
		}
		fstabEntry := fstabEntry{
			spec:      spec,
			file:      *fs.Path,
			vfsType:   *fs.Format,
			mountOpts: mountOpts,
			freq:      0,
			passno:    passno,
		}
		f.fstabEntries = append(f.fstabEntries, fstabEntry)
		if fs.Label == nil {
			continue
		}
		// create legacy disk.json
		switch *fs.Label {
		case "root", "efi", "varlib":
			part := Partition{
				Label:      *fs.Label,
				Filesystem: *fs.Format,
				Properties: map[string]string{"UUID": properties["UUID"]},
			}
			f.disk.Partitions = append(f.disk.Partitions, part)
		}
	}
	return nil
}

type mount struct {
	source string
	target string
	fstype string
	flags  uintptr
	data   string
}

var (
	specialMounts = []mount{
		{source: "proc", target: "/proc", fstype: "proc", flags: 0, data: ""},
		{source: "sys", target: "/sys", fstype: "sysfs", flags: 0, data: ""},
		{source: "efivarfs", target: "/sys/firmware/efi/efivars", fstype: "efivarfs", flags: 0, data: ""},
		{source: "tmpfs", target: "/tmp", fstype: "tmpfs", flags: 0, data: ""},
		// /dev is a bind mount, a bind mount must have MS_BIND flags set see man 2 mount
		{source: "/dev", target: "/dev", fstype: "", flags: syscall.MS_BIND, data: ""},
	}
)

func (f *Filesystem) mountSpecialFilesystems() error {
	// Order is important and must be preserved.
	for _, m := range specialMounts {
		mountPoint := filepath.Join(f.chroot, m.target)

		if _, err := gos.Stat(mountPoint); err != nil && gos.IsNotExist(err) {
			if err := gos.MkdirAll(mountPoint, 0755); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		log.Info("mount", "source", m.source, "target", mountPoint, "fstype", m.fstype, "flags", m.flags, "data", m.data)
		err := syscall.Mount(m.source, mountPoint, m.fstype, m.flags, m.data)
		if err != nil {
			return errors.Wrapf(err, "mounting %s to %s failed", m.source, m.target)
		}
	}
	return nil
}

func (f *Filesystem) umountFilesystems() error {
	for index := len(specialMounts) - 1; index >= 0; index-- {
		m := filepath.Join(f.chroot, specialMounts[index].target)
		log.Info("unmounting", "mountpoint", m)
		err := syscall.Unmount(m, syscall.MNT_FORCE)
		if err != nil {
			log.Error("unable to unmount", "path", m, "error", err)
		}
	}
	for index := len(f.mounts) - 1; index >= 0; index-- {
		m := f.mounts[index]
		log.Info("unmounting", "mountpoint", m)
		err := syscall.Unmount(m, syscall.MNT_FORCE)
		if err != nil {
			log.Error("unable to unmount", "path", m, "error", err)
		}
	}
	return nil
}

func (f *Filesystem) CreateFSTab() error {
	return f.fstabEntries.write(f.chroot)
}

func (f *Filesystem) createDiskJSON() error {
	configdir := path.Join(f.chroot, "etc", "metal")
	destination := path.Join(configdir, "disk.json")

	if _, err := gos.Stat(configdir); err != nil && gos.IsNotExist(err) {
		if err := gos.MkdirAll(configdir, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	j, err := json.MarshalIndent(f.disk, "", "  ")
	if err != nil {
		return errors.Wrap(err, "unable to marshal to json")
	}
	log.Info("create legacy disk.json", "content", string(j))
	return ioutil.WriteFile(destination, j, 0600)
}

func mountFs(chroot string, fs models.ModelsV1Filesystem) (string, error) {
	if fs.Format == nil || *fs.Format == "swap" || *fs.Format == "" || *fs.Format == "tmpfs" {
		return "", nil
	}
	path := filepath.Join(chroot, *fs.Path)

	if _, err := gos.Stat(path); err != nil && gos.IsNotExist(err) {
		if err := gos.MkdirAll(path, 0755); err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}
	opts := optionSliceToString(fs.Mountoptions, ",")
	log.Info("mount filesystem", "device", *fs.Device, "path", path, "format", fs.Format, "opts", opts)
	err := os.ExecuteCommand("mount", "-o", opts, "-t", *fs.Format, *fs.Device, path)
	if err != nil {
		log.Error("mount filesystem failed", "device", *fs.Device, "path", fs.Path, "opts", opts, "error", err)
		return "", errors.Wrapf(err, "unable to create filesystem %s on %s", *fs.Device, *fs.Path)
	}
	return path, nil
}

func depth(path string) uint {
	var count uint = 0
	for p := filepath.Clean(path); p != "/"; count++ {
		p = filepath.Dir(p)
	}
	return count
}

func optionSliceToString(opts []string, separator string) string {
	mountOpts := make([]string, len(opts))
	for i, o := range opts {
		mountOpts[i] = string(o)
	}
	return strings.Join(mountOpts, separator)
}

// write all fstab entries to /etc/fstab inside chroot
func (fss fstabEntries) write(chroot string) error {
	entries := []string{}
	for _, fs := range fss {
		entries = append(entries, fs.string())
	}
	fstab := strings.Join(entries, "\n")
	header := fmt.Sprintf("# created by metal-hammer: %q\n", v.V)
	content := header + fstab + "\n"
	log.Info("write fstab", "content", content)
	return ioutil.WriteFile(path.Join(chroot, "/etc/fstab"), []byte(content), 0644)
}

func (fs fstabEntry) string() string {
	return fmt.Sprintf("%s %s %s %s %d %d", fs.spec, fs.file, fs.vfsType, strings.Join(fs.mountOpts, ","), fs.freq, fs.passno)
}
