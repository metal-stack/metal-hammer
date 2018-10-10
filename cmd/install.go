package cmd

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/inconshreveable/log15"
	"github.com/jaypipes/ghw"
	"github.com/mholt/archiver"
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
	prefix = "/rootfs"
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

// MountOption a option given to a mountpoint
type MountOption string

// Disk is a physical Disk
type Disk struct {
	Device string
	// Partitions to create on this disk, order is preserved
	Partitions []*Partition
}

type InstallerConfig struct {
	Hostname     string `yaml:"hostname"`
	SSHPublicKey string `yaml:"sshpublickey"`
}

func (p *Partition) String() string {
	return fmt.Sprintf("%s", p.Device)
}

// Install a given image to the disk by using genuinetools/img
func Install(device *Device) error {
	image := device.Image.Url
	err := partition(defaultDisk)
	if err != nil {
		return err
	}

	err = mountPartitions(prefix, defaultDisk)
	if err != nil {
		return err
	}

	err = pull(image)
	if err != nil {
		return err
	}
	err = burn(prefix, image)
	if err != nil {
		return err
	}

	err = install(prefix, device)
	if err != nil {
		return err
	}
	return nil
}

// WipeDisks will erase all content and partitions of all existing Disks
func WipeDisks() error {
	log.Info("wipe all disks")
	block, err := ghw.Block()
	if err != nil {
		return fmt.Errorf("unable to gather disks: %v", err)
	}
	for _, disk := range block.Disks {
		log.Info("TODO wipe disk", "disk", disk)

		diskDevice := fmt.Sprintf("/dev/%s", disk.Name)
		log.Info("sgdisk zap all existing partitions", "disk", diskDevice)
		err := executeCommand(sgdiskCommand, "-Z", diskDevice)
		if err != nil {
			log.Error("sgdisk zap all existing partitions failed", "error", err)
		}
	}
	return nil
}

func partition(disk Disk) error {
	log.Info("partition disk", "disk", disk)

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
	err := executeCommand(sgdiskCommand, args...)
	// FIXME sgdisk return 0 in case of failure, and > 0 if succeed
	// TODO still the case ?
	if err != nil {
		log.Error("sgdisk creating partitions failed", "error", err)
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
	md5file := image + ".md5"
	log.Info("pull image", "image", image)
	err := downloadFile("/tmp/os.tgz", image)
	if err != nil {
		return fmt.Errorf("unable to pull image %s error: %v", image, err)
	}
	err = downloadFile("/tmp/os.tgz.md5", md5file)
	if err != nil {
		return fmt.Errorf("unable to pull md5 %s error: %v", md5file, err)
	}
	log.Info("check md5")
	matches, err := checkMD5("/tmp/os.tgz", "/tmp/os.tgz.md5")
	if err != nil {
		return fmt.Errorf("unable to check md5sum error: %v", err)
	}
	if !matches {
		log.Error("md5sum mismatch")
	}

	log.Debug("pull image done", "image", image)
	return nil
}

// downloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func downloadFile(filepath string, url string) error {
	log.Info("download", "from", url, "to", filepath)
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

func checkMD5(file, md5file string) (bool, error) {
	md5fileContent, err := ioutil.ReadFile(md5file)
	if err != nil {
		return false, fmt.Errorf("unable to read md5sum file: %v", err)
	}
	expectedMD5 := strings.Split(string(md5fileContent), " ")[0]

	f, err := os.Open(file)
	if err != nil {
		return false, fmt.Errorf("unable to read file: %v", err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return false, fmt.Errorf("unable to calculate md5sum of file: %v", err)
	}
	sourceMD5 := fmt.Sprintf("%x", h.Sum(nil))
	log.Info("checkMD5", "source md5", sourceMD5, "expected md5", expectedMD5)
	if sourceMD5 != expectedMD5 {
		return false, nil
	}
	return true, nil
}

// burn a image by calling genuinetools/img unpack to a specific directory
func burn(prefix, image string) error {
	log.Info("burn image", "image", image)

	err := archiver.TarGz.Open("/tmp/os.tgz", prefix)
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

type mount struct {
	source string
	target string
	fstype string
	flags  uintptr
	data   string
}

// install will execute /install.sh in the pulled docker image which was extracted onto disk
// to finish installation e.g. install mbr, grub, write network and filesystem config
func install(prefix string, device *Device) error {
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

	log.Info("write installation configuration")
	configdir := path.Join(prefix, "etc", "metal")
	err := os.MkdirAll(configdir, 0755)
	if err != nil {
		log.Error("mkdir of configuration in target os failed", "error", err)
	}
	configpath := path.Join(configdir, "install.yaml")
	err = writeInstallerConfig(device, configpath)
	if err != nil {
		log.Error("writing configuration in target os failed", "configpath", configpath, "error", err)
	}

	log.Info("running /install.sh on", "prefix", prefix)

	err = os.Chdir(prefix)
	if err != nil {
		log.Error("unable to chdir", "chroot", prefix, "error", err)
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
		log.Error("running install.sh in chroot failed", "error", err)
		return fmt.Errorf("running install.sh in chroot failed: %v", err)
	}
	err = os.Chdir("/")
	if err != nil {
		log.Error("unable to chdir to /", "error", err)
	}
	log.Info("finish running /install.sh")

	umounts := [6]string{"/boot/efi", "/proc", "/sys", "/dev", "/tmp", "/"}
	for _, m := range umounts {
		p := prefix + m
		err := syscall.Unmount(p, syscall.MNT_FORCE)
		if err != nil {
			log.Error("unable to umount", "path", p, "error", err)
		}
	}

	return nil
}

func writeInstallerConfig(device *Device, destination string) error {
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

// small helper to execute a command, redirect stdout/stderr.
func executeCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}
