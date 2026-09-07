package storage

import (
	"errors"
	"fmt"
	"log/slog"
	gos "os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"syscall"

	"github.com/u-root/u-root/pkg/mount/block"

	"github.com/metal-stack/api/go/enum"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-hammer/pkg/os"
	"github.com/metal-stack/metal-hammer/pkg/os/command"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/metal-stack/v"
)

type Filesystem struct {
	log *slog.Logger

	config *apiv2.FilesystemLayout
	// chroot defines the root of the mounts
	chroot string
	// mounts are collected to be able to umount all in reverse order
	mounts       []string
	fstabEntries fstabEntries
	RootUUID     string
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

func New(log *slog.Logger, chroot string, config *apiv2.FilesystemLayout) *Filesystem {
	return &Filesystem{
		log:          log,
		config:       config,
		chroot:       chroot,
		fstabEntries: fstabEntries{},
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

	err = f.createLogicalVolumes()
	if err != nil {
		return fmt.Errorf("create logical volumes failed:%w", err)
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

	return nil
}
func (f *Filesystem) Umount() error {
	return f.umountFilesystems()
}

func (f *Filesystem) createPartitions() error {
	if len(f.config.Disks) == 0 {
		return nil
	}
	for _, disk := range f.config.Disks {
		opts := []string{}

		for _, p := range disk.Partitions {
			opts = append(opts, fmt.Sprintf("--new=%d:0:+%dM", p.Number, p.Size))
			opts = append(opts, fmt.Sprintf("--change-name=%d:%s", p.Number, pointer.SafeDeref(p.Label)))
			if p.GptType != nil {
				opts = append(opts, fmt.Sprintf("--typecode=%d:%s", p.Number, *p.GptType))
			}
		}
		f.log.Info("wipe existing partition signatures", "command", command.WIPEFS+" --all"+" "+disk.Device)
		err := os.ExecuteCommand(command.WIPEFS, "--all", disk.Device)
		if err != nil {
			f.log.Error("wipe existing partition signatures failed", "error", err)
			return fmt.Errorf("unable wipe existing partitions on %s %w", disk.Device, err)
		}
		opts = append(opts, disk.Device)
		f.log.Info("sgdisk create partitions", "command", opts)
		err = os.ExecuteCommand(command.SGDisk, opts...)
		if err != nil {
			f.log.Error("sgdisk creating partitions failed", "error", err)
			return fmt.Errorf("unable to create partitions on %s %w", disk.Device, err)
		}

		blkdev, err := block.Device(disk.Device)
		if err != nil {
			return fmt.Errorf("unable to find block device %s: %v", disk.Device, err)
		}

		err = blkdev.ReadPartitionTable()
		if err != nil {
			return fmt.Errorf("unable to re-read the partition table. Kernel still uses old partition table: %v", err)
		}
	}
	return nil
}

func (f *Filesystem) createRaids() error {
	if len(f.config.Raid) == 0 {
		return nil
	}

	for _, raid := range f.config.Raid {
		spares := raid.Spares
		var level string
		switch raid.Level {
		case apiv2.RaidLevel_RAID_LEVEL_0:
			level = "0"
		case apiv2.RaidLevel_RAID_LEVEL_1:
			level = "1"
		default:
			// not supported
		}

		args := []string{
			"--create", raid.ArrayName,
			"--force",
			"--run",
			"--homehost", "any",
			"--level", level,
			"--raid-devices", fmt.Sprintf("%d", len(raid.Devices)-int(spares)),
		}

		switch raid.Level {
		case apiv2.RaidLevel_RAID_LEVEL_0, apiv2.RaidLevel_RAID_LEVEL_1:
			args = append(args, "--assume-clean")
		default:
			// only safe to skip initial sync for raid 0 and 1
			// see https://raid.wiki.kernel.org/index.php/Initial_Array_Creation#raid5
		}

		if spares > 0 {
			args = append(args, "--spare-devices", fmt.Sprintf("%d", spares))
		}

		for _, o := range raid.CreateOptions {
			args = append(args, string(o))
		}

		args = append(args, raid.Devices...)

		f.log.Info("create mdadm raid", "args", args)
		err := os.ExecuteCommand(command.MDADM, args...)
		if err != nil {
			f.log.Error("create mdadm raid", "error", err)
			return fmt.Errorf("unable to create mdadm raid %s %w", raid.ArrayName, err)
		}

		// set sync speed
		err = gos.WriteFile("/proc/sys/dev/raid/speed_limit_min", []byte("200000000"), 0644) // nolint:gosec
		if err != nil {
			f.log.Error("unable to set min sync speed, ignoring...", "error", err)
		}
	}
	return nil
}

func (f *Filesystem) createLogicalVolumes() error {
	if len(f.config.VolumeGroups) == 0 {
		return nil
	}

	pvcount := make(map[string]int)
	for _, vg := range f.config.VolumeGroups {
		if vg.Name == "" {
			continue
		}
		if vgExists(f.log, vg.Name) {
			continue
		}
		args := []string{
			"vgcreate",
			"--verbose",
			vg.Name,
		}
		for _, tag := range vg.Tags {
			args = append(args, "--addtag", tag)
		}
		args = append(args, vg.Devices...)

		pvcount[vg.Name] = len(vg.Devices)
		err := os.ExecuteCommand(command.LVM, args...)
		if err != nil {
			f.log.Error("vgcreate", "error", err)
			return fmt.Errorf("unable to create volume group %s %w", vg.Name, err)
		}
	}

	for _, lv := range f.config.LogicalVolumes {
		if lv.Name == "" || lv.VolumeGroup == "" {
			continue
		}
		if lvExists(f.log, lv.VolumeGroup, lv.Name) {
			continue
		}

		args := []string{
			"lvcreate",
			"--yes",
			"--verbose",
			"--name", lv.Name,
			"--wipesignatures", "y",
		}

		if lv.Size > 0 {
			args = append(args, "--size", fmt.Sprintf("%dm", lv.Size))
		} else {
			args = append(args, "--extents", "100%FREE")
		}

		switch lv.LvmType {
		case apiv2.LVMType_LVM_TYPE_LINEAR:
			if pvcount[lv.VolumeGroup] < 2 {
				f.log.Warn("volumegroup has only 1 pv, only linear is supported", "lv", lv.Name, "vg", lv.VolumeGroup)
			}
		case apiv2.LVMType_LVM_TYPE_STRIPED:
			args = append(args, "--type", "striped", "--stripes", fmt.Sprintf("%d", pvcount[lv.VolumeGroup]))
		case apiv2.LVMType_LVM_TYPE_RAID1:
			args = append(args, "--type", "raid1", "--mirrors", "1", "--nosync")
		default:
			return fmt.Errorf("unsupported lvmtype:%s", lv.LvmType)
		}
		args = append(args, lv.VolumeGroup)

		f.log.Info("lvcreate", "args", args)
		err := os.ExecuteCommand(command.LVM, args...)
		if err != nil {
			f.log.Error("lvcreate", "error", err)
			return fmt.Errorf("unable to create logical volume %s %w", lv.Name, err)
		}
	}

	return nil
}

func (f *Filesystem) createFilesystems() error {
	if len(f.config.Filesystems) == 0 {
		return nil
	}

	for _, fs := range f.config.Filesystems {
		if fs.Format == apiv2.Format_FORMAT_TMPFS {
			continue
		}
		mkfs := ""
		args := []string{}
		args = append(args, fs.CreateOptions...)
		switch fs.Format {
		case apiv2.Format_FORMAT_EXT3:
			mkfs = command.MKFSExt3
			args = append(args, "-F")
			if fs.Label != nil {
				args = append(args, "-L", pointer.SafeDeref(fs.Label))
			}
		case apiv2.Format_FORMAT_EXT4:
			mkfs = command.MKFSExt4
			args = append(args, "-F")
			if fs.Label != nil {
				args = append(args, "-L", pointer.SafeDeref(fs.Label))
			}
		case apiv2.Format_FORMAT_SWAP:
			mkfs = command.MKSwap
			args = append(args, "-f")
			if fs.Label != nil {
				args = append(args, "-L", pointer.SafeDeref(fs.Label))
			}
		case apiv2.Format_FORMAT_VFAT:
			mkfs = command.MKFSVFat
			// There is no force flag for mkfs.vfat, it always destroys any data on
			// the device at which it is pointed.
			// also label is added with -n
			if fs.Label != nil {
				args = append(args, "-n", pointer.SafeDeref(fs.Label))
			}
		case apiv2.Format_FORMAT_NONE:
			//
		default:
			return fmt.Errorf("unsupported filesystem format: %q", fs.Format)
		}
		args = append(args, fs.Device)
		f.log.Info("create filesystem", "args", args)
		err := os.ExecuteCommand(mkfs, args...)
		if err != nil {
			f.log.Error("create filesystem failed", "device", fs.Device, "error", err)
			return fmt.Errorf("unable to create filesystem on %s %w", fs.Device, err)
		}
	}

	return nil
}

func (f *Filesystem) mountFilesystems() error {
	var fss []*apiv2.Filesystem
	for _, fs := range f.config.Filesystems {
		if fs.Path == nil {
			continue
		}
		fss = append(fss, fs)
	}
	sort.Slice(fss, func(i, j int) bool {
		return depth(pointer.SafeDeref(fss[i].Path)) < depth(pointer.SafeDeref(fss[j].Path))
	})
	for _, fs := range fss {
		path, err := mountFs(f.log, f.chroot, fs)
		if err != nil {
			return err
		}
		if path != "" {
			f.mounts = append(f.mounts, path)
		}

		var (
			passno     = uint(2)
			properties = map[string]string{"UUID": ""}
			spec       string
		)
		format, err := enum.GetStringValue(fs.Format)
		if err != nil {
			return err
		}
		if fs.Format == apiv2.Format_FORMAT_TMPFS {
			passno = 0
			spec = *format
		} else {
			properties, err = FetchBlockIDProperties(fs.Device)
			if err != nil {
				return err
			}
			spec = fmt.Sprintf("UUID=%s", properties["UUID"])
		}
		if pointer.SafeDeref(fs.Path) == "/" {
			passno = 1
		}
		mountOpts := []string{"defaults"}
		if len(fs.MountOptions) > 0 {
			mountOpts = fs.MountOptions
		}

		fstabEntry := fstabEntry{
			spec:      spec,
			file:      pointer.SafeDeref(fs.Path),
			vfsType:   *format,
			mountOpts: mountOpts,
			freq:      0,
			passno:    passno,
		}
		f.fstabEntries = append(f.fstabEntries, fstabEntry)
		if pointer.SafeDeref(fs.Label) == "root" {
			f.RootUUID = properties["UUID"]
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
		// /dev and /run are bind mounts, a bind mount must have MS_BIND flags set see man 2 mount
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

		f.log.Info("mount", "source", m.source, "target", mountPoint, "fstype", m.fstype, "flags", m.flags, "data", m.data)
		err := syscall.Mount(m.source, mountPoint, m.fstype, m.flags, m.data)
		if err != nil {
			return fmt.Errorf("mounting %s to %s failed %w", m.source, m.target, err)
		}
	}
	return nil
}

func (f *Filesystem) umountFilesystems() error {
	syscall.Sync()

	for _, specialMount := range slices.Backward(specialMounts) {
		m := filepath.Join(f.chroot, specialMount.target)
		f.log.Info("unmounting", "mountpoint", m)
		if err := syscall.Unmount(m, 0); err != nil {
			f.log.Error("unmount failed, detaching lazily", "path", m, "error", err)
			if err := syscall.Unmount(m, syscall.MNT_DETACH); err != nil { // won't block the data-fs umount below
				f.log.Error("unable to lazy unmount, ignoring", "path", m, "error", err)
			}
		}
	}

	var errs []error
	for _, m := range slices.Backward(f.mounts) {
		if m == "" {
			continue
		}
		f.log.Info("unmounting", "mountpoint", m)
		if err := syscall.Unmount(m, 0); err != nil {
			errs = append(errs, fmt.Errorf("unable to unmount %q %w", m, err))
		}
	}
	return errors.Join(errs...)
}

func (f *Filesystem) CreateFSTab() error {
	return f.fstabEntries.write(f.log, f.chroot)
}

func mountFs(log *slog.Logger, chroot string, fs *apiv2.Filesystem) (string, error) {
	switch fs.Format {
	case apiv2.Format_FORMAT_NONE, apiv2.Format_FORMAT_SWAP, apiv2.Format_FORMAT_TMPFS:
		return "", nil
	}
	path := filepath.Join(chroot, pointer.SafeDeref(fs.Path))

	if _, err := gos.Stat(path); err != nil && gos.IsNotExist(err) {
		if err := gos.MkdirAll(path, 0755); err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}
	opts := optionSliceToString(fs.MountOptions, ",")
	log.Info("mount filesystem", "device", fs.Device, "path", path, "format", fs.Format, "opts", opts)
	var args []string
	if len(opts) > 0 {
		args = append(args, "-o", opts)
	}
	format, err := enum.GetStringValue(fs.Format)
	if err != nil {
		return "", err
	}

	args = append(args, "-t", *format, fs.Device, path)

	if err := os.ExecuteCommand("mount", args...); err != nil {
		log.Error("mount filesystem failed", "device", fs.Device, "path", fs.Path, "opts", opts, "error", err)
		return "", fmt.Errorf("unable to mount filesystem %s on %s opts:%v error:%w", fs.Device, pointer.SafeDeref(fs.Path), opts, err)
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

// from man mount:
// The command mount does not pass the mount options
// unbindable, runbindable, private, rprivate, slave, rslave, shared, rshared, auto, noauto, comment, x-*, loop, offset and sizelimit
// to the mount.<suffix> helpers. All other options are used in a comma-separated list as an argument to the -o option.
// defaults is special and always set.
var impossibleMountOptions = []string{
	"defaults", "unbindable", "runbindable", "private", "rprivate", "slave", "rslave", "shared", "rshared", "auto", "noauto", "comment", "loop", "offset", "sizelimit",
}

func optionSliceToString(opts []string, separator string) string {
	var mountOpts []string
	for _, o := range opts {
		option := string(o)
		if slices.Contains(impossibleMountOptions, option) || strings.HasPrefix(option, "x-") {
			continue
		}
		mountOpts = append(mountOpts, option)
	}
	return strings.Join(mountOpts, separator)
}

// write all fstab entries to /etc/fstab inside chroot
func (fss fstabEntries) write(log *slog.Logger, chroot string) error {
	entries := []string{}
	for _, fs := range fss {
		entries = append(entries, fs.string())
	}
	fstab := strings.Join(entries, "\n")
	header := fmt.Sprintf("# created by metal-hammer: %q\n", v.V)
	content := header + fstab + "\n"
	log.Info("write fstab", "content", content)
	//nolint:gosec
	return gos.WriteFile(path.Join(chroot, "/etc/fstab"), []byte(content), 0644)
}

func (fs fstabEntry) string() string {
	return fmt.Sprintf("%s %s %s %s %d %d", fs.spec, fs.file, fs.vfsType, strings.Join(fs.mountOpts, ","), fs.freq, fs.passno)
}

func lvExists(log *slog.Logger, vg string, name string) bool {
	//nolint:gosec
	cmd := exec.Command("lvm", "lvs", vg+"/"+name, "--noheadings", "-o", "lv_name")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Info("unable to list existing volumes", "lv", name, "error", err)
		return false
	}
	return name == strings.TrimSpace(string(out))
}

func vgExists(log *slog.Logger, vgname string) bool {
	cmd := exec.Command("lvm", "vgs", vgname, "--noheadings", "-o", "vg_name")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Info("unable to list existing volumegroups", "vg", vgname, "error", err)
		return false
	}
	return vgname == strings.TrimSpace(string(out))
}
