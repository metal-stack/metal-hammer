package cmd

import (
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"git.f-i-ts.de/cloud-native/maas/metal-hammer/pkg"
	log "github.com/inconshreveable/log15"
	"github.com/mholt/archiver"
	pb "gopkg.in/cheggaaa/pb.v1"
	"gopkg.in/yaml.v2"
)

var (
	sgdiskCommand    = "/usr/bin/sgdisk"
	ext4MkFsCommand  = "/sbin/mkfs.ext4"
	ext3MkFsCommand  = "/sbin/mkfs.ext3"
	fat32MkFsCommand = "/sbin/mkfs.vfat"
	mkswapCommand    = "/sbin/mkswap"
	defaultDisk      = Disk{
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

const (
	prefix             = "/rootfs"
	osImageDestination = "/tmp/os.tgz"
)

// GPTType is the GUID Partition table type
type GPTType string

// GPTGuid is the UID of the GPT partition to create
type GPTGuid string

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
	GPTGuid    GPTGuid
}

func (p *Partition) String() string {
	return fmt.Sprintf("%s", p.Device)
}

// MountOption a option given to a mountpoint
type MountOption string

// Disk is a physical Disk
type Disk struct {
	Device string
	// Partitions to create on this disk, order is preserved
	Partitions []*Partition
}

// InstallerConfig contains configuration items which are
// consumed by the install.sh of the individual target OS.
type InstallerConfig struct {
	Hostname     string `yaml:"hostname"`
	SSHPublicKey string `yaml:"sshpublickey"`
}

// Wait until a device create request was fired
func Wait(url, uuid string) (*Device, error) {
	log.Info("waiting for install, long polling", "uuid", uuid)
	e := fmt.Sprintf("%v/%v", url, uuid)

	var resp *http.Response
	var err error
	for {
		resp, err = http.Get(e)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Warn("wait for install failed, retrying...", "error", err)
		} else {
			break
		}
		time.Sleep(2 * time.Second)
	}

	deviceJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("wait for install reading response failed with: %v", err)
	}

	var device Device
	err = json.Unmarshal(deviceJSON, &device)
	if err != nil {
		return nil, fmt.Errorf("wait for install could not unmarshal response with error: %v", err)
	}
	log.Info("stopped waiting got", "device", device)

	return &device, nil
}

// Install a given image to the disk by using genuinetools/img
func Install(device *Device) (*pkg.Bootinfo, error) {
	image := device.Image.Url
	err := partition(defaultDisk)
	if err != nil {
		return nil, err
	}

	err = mountPartitions(prefix, defaultDisk)
	if err != nil {
		return nil, err
	}

	err = pull(image)
	if err != nil {
		return nil, err
	}
	err = burn(prefix, image)
	if err != nil {
		return nil, err
	}

	info, err := install(prefix, device)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func partition(disk Disk) error {
	log.Info("partition disk", "disk", disk)

	err := executeCommand(sgdiskCommand, "-Z", disk.Device)
	if err != nil {
		log.Error("sgdisk zapping existing partitions failed, ignoring...", "error", err)
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
		if p.GPTGuid != "" {
			args = append(args, fmt.Sprintf("-u=%d:%s", p.Number, p.GPTGuid))
		}

		// TODO format must not have the side effect to change incoming data
		p.Device = fmt.Sprintf("%s%d", disk.Device, p.Number)
	}

	args = append(args, disk.Device)
	log.Info("sgdisk create partitions", "command", args)
	err = executeCommand(sgdiskCommand, args...)
	if err != nil {
		log.Error("sgdisk creating partitions failed", "error", err)
		return fmt.Errorf("unable to create partitions on %s error:%v", disk, err)
	}

	return nil
}

func mountPartitions(prefix string, disk Disk) error {
	log.Info("mount disk", "disk", disk)
	// "/" must be mounted first
	partitions := disk.SortByMountPoint()

	// FIXME error handling
	for _, p := range partitions {
		err := createFilesystem(p)
		if err != nil {
			log.Error("mount partition create filesystem failed", "error", err)
			return fmt.Errorf("mount partitions create fs failed: %v", err)
		}

		if p.MountPoint == "" {
			continue
		}

		mountPoint := filepath.Join(prefix, p.MountPoint)
		err = os.MkdirAll(mountPoint, os.ModePerm)
		if err != nil {
			log.Error("mount partition create directory", "error", err)
			return fmt.Errorf("mount partitions create directory failed: %v", err)
		}
		log.Info("mount partition", "partition", p.Device, "mountPoint", mountPoint)
		// see man 2 mount
		err = syscall.Mount(p.Device, mountPoint, string(p.Filesystem), 0, "")
		if err != nil {
			log.Error("unable to mount", "partition", p.Device, "mountPoint", mountPoint, "error", err)
			return fmt.Errorf("mount partitions mount: %s to:%s failed: %v", p.Device, mountPoint, err)

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
		args = append(args, "-v", "-F")
		if p.Label != "" {
			args = append(args, "-L", p.Label)
		}
	case EXT3:
		mkfs = ext3MkFsCommand
		args = append(args, "-v", "-F")
		if p.Label != "" {
			args = append(args, "-L", p.Label)
		}
	case FAT32, VFAT:
		mkfs = fat32MkFsCommand
		args = append(args, "-v", "-F", "32")
		if p.Label != "" {
			args = append(args, "-n", strings.ToUpper(p.Label))
		}
	case SWAP:
		mkfs = mkswapCommand
		args = append(args, "-f")
		if p.Label != "" {
			args = append(args, "-L", p.Label)
		}
	default:
		return fmt.Errorf("unsupported filesystem format: %q", p.Filesystem)
	}
	args = append(args, p.Device)
	err := executeCommand(mkfs, args...)
	if err != nil {
		return fmt.Errorf("mkfs failed with error:%v", err)
	}

	return nil
}

// SortByMountPoint ensures that "/" is the first, which is required for mounting
func (d *Disk) SortByMountPoint() []*Partition {
	ordered := make([]*Partition, 0)
	for _, p := range d.Partitions {
		if p.MountPoint == "/" {
			ordered = append(ordered, p)
		}
	}
	for _, p := range d.Partitions {
		if p.MountPoint != "/" {
			ordered = append(ordered, p)
		}
	}
	return ordered
}

// pull a image from s3
func pull(image string) error {
	log.Info("pull image", "image", image)
	destination := osImageDestination
	md5destination := destination + ".md5"
	md5file := image + ".md5"
	err := download(image, destination)
	if err != nil {
		return fmt.Errorf("unable to pull image %s error: %v", image, err)
	}
	err = download(md5file, md5destination)
	defer os.Remove(md5destination)
	if err != nil {
		return fmt.Errorf("unable to pull md5 %s error: %v", md5file, err)
	}
	log.Info("check md5")
	matches, err := checkMD5(destination, md5destination)
	if err != nil || !matches {
		return fmt.Errorf("md5sum mismatch %v", err)
	}

	log.Info("pull image done", "image", image)
	return nil
}

// burn a image pulling a tarball and unpack to a specific directory
func burn(prefix, image string) error {
	log.Info("burn image", "image", image)

	source := osImageDestination

	file, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("%s: failed to open archive: %v", source, err)

	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("unable to stat %s error: %v", source, err)
	}

	// last four bytes of the gzip contain the uncompressed file size
	buf := make([]byte, 4)
	start := stat.Size() - 4
	_, err = file.ReadAt(buf, start)
	if err != nil {
		return fmt.Errorf("cannot read uncompressed file size of gzip: %v", err)
	}

	bar := pb.New64(int64(binary.LittleEndian.Uint32(buf))).SetUnits(pb.U_BYTES)
	bar.Start()
	bar.SetWidth(80)
	bar.ShowSpeed = true

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("error decompressing: %v", err)
	}
	defer gzr.Close()

	reader := bar.NewProxyReader(gzr)

	err = archiver.Tar.Read(reader, prefix)
	if err != nil {
		return fmt.Errorf("unable to burn image %s error: %v", source, err)
	}

	bar.Finish()

	err = os.Remove(source)
	if err != nil {
		log.Warn("burn image unable to remove image source", "error", err)
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

// install will execute /install.sh in the pulled docker image which was extracted onto disk
// to finish installation e.g. install mbr, grub, write network and filesystem config
func install(prefix string, device *Device) (*pkg.Bootinfo, error) {
	log.Info("install image", "image", device.Image.Url)
	mounts := []mount{
		mount{source: "proc", target: "/proc", fstype: "proc", flags: 0, data: ""},
		mount{source: "sys", target: "/sys", fstype: "sysfs", flags: 0, data: ""},
		mount{source: "tmpfs", target: "/tmp", fstype: "tmpfs", flags: 0, data: ""},
		// /dev is a bind mount, a bind mount must have MS_BIND flags set see man 2 mount
		mount{source: "/dev", target: "/dev", fstype: "", flags: syscall.MS_BIND, data: ""},
	}

	for _, m := range mounts {
		err := syscall.Mount(m.source, prefix+m.target, m.fstype, m.flags, m.data)
		if err != nil {
			log.Error("mounting failed", "source", m.source, "target", m.target, "error", err)
		}
	}

	err := writeInstallerConfig(device)
	if err != nil {
		return nil, fmt.Errorf("writing configuration install.yaml failed:%v", err)
	}

	log.Info("running /install.sh on", "prefix", prefix)
	err = os.Chdir(prefix)
	if err != nil {
		return nil, fmt.Errorf("unable to chdir to: %s error:%v", prefix, err)
	}
	cmd := exec.Command("/install.sh")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid:    uint32(0),
			Gid:    uint32(0),
			Groups: []uint32{0},
		},
		Chroot: prefix,
	}
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("running install.sh in chroot failed: %v", err)
	}
	err = os.Chdir("/")
	if err != nil {
		return nil, fmt.Errorf("unable to chdir to: / error:%v", err)
	}
	log.Info("finish running /install.sh")

	log.Info("read /etc/metal/boot-info.yaml")
	bi, err := ioutil.ReadFile(path.Join(prefix, "etc", "metal", "boot-info.yaml"))
	if err != nil {
		log.Error("could not read boot-info.yaml", "error", err)
		return nil, err
	}

	var info pkg.Bootinfo
	err = yaml.Unmarshal(bi, &info)
	if err != nil {
		log.Error("could not unmarshal boot-info.yaml", "error", err)
		return nil, err
	}

	files := []string{info.Kernel, info.Initrd}
	tmp := "/tmp"
	for _, f := range files {
		src := path.Join(prefix, f)
		dest := path.Join(tmp, filepath.Base(f))
		_, err := copy(src, dest)
		if err != nil {
			log.Error("could not copy", "src", src, "dest", dest, "error", err)
			return nil, err
		}
	}
	info.Kernel = path.Join(tmp, filepath.Base(info.Kernel))
	info.Initrd = path.Join(tmp, filepath.Base(info.Initrd))

	umounts := [6]string{"/boot/efi", "/proc", "/sys", "/dev", "/tmp", "/"}
	for _, m := range umounts {
		p := prefix + m
		log.Info("unmounting", "mountpoint", p)
		err := syscall.Unmount(p, syscall.MNT_FORCE)
		if err != nil {
			log.Error("unable to umount", "path", p, "error", err)
		}
	}

	return &info, nil
}

func writeInstallerConfig(device *Device) error {
	log.Info("write installation configuration")
	configdir := path.Join(prefix, "etc", "metal")
	err := os.MkdirAll(configdir, 0755)
	if err != nil {
		return fmt.Errorf("mkdir of %s target os failed: %v", configdir, err)
	}
	destination := path.Join(configdir, "install.yaml")

	y := &InstallerConfig{
		Hostname:     device.Hostname,
		SSHPublicKey: device.SSHPubKey,
	}
	yamlContent, err := yaml.Marshal(y)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(destination, yamlContent, 0600)
}
