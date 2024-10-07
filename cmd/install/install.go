package install

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strconv"
	"strings"
	"time"

	config "github.com/flatcar/ignition/config/v2_4"
	v1 "github.com/metal-stack/metal-hammer/cmd/install/v1"
	"github.com/metal-stack/metal-hammer/pkg/api"
	"github.com/metal-stack/metal-networker/pkg/netconf"
	"github.com/metal-stack/v"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

const (
	installYAML = "/etc/metal/install.yaml"
	userdata    = "/etc/metal/userdata"
)

func parseInstallYAML(fs afero.Fs) (*api.InstallerConfig, error) {
	var config api.InstallerConfig
	content, err := afero.ReadFile(fs, installYAML)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func runFromCI() bool {
	ciEnv := os.Getenv("INSTALL_FROM_CI")

	ci, err := strconv.ParseBool(ciEnv)
	if err != nil {
		return false
	}

	return ci
}

type installer struct {
	log    *slog.Logger
	fs     afero.Fs
	oss    operatingsystem
	config *api.InstallerConfig
	exec   *cmdexec
}

func Run() error {
	start := time.Now()
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
	log := slog.New(jsonHandler)

	log.Info("running install", "version", v.V.String())

	fs := afero.OsFs{}

	oss, err := detectOS(fs)
	if err != nil {
		return err
	}

	config, err := parseInstallYAML(fs)
	if err != nil {
		return err
	}

	i := installer{
		log:    log.WithGroup("install-go"),
		fs:     fs,
		oss:    oss,
		config: config,
		exec: &cmdexec{
			log: log.WithGroup("cmdexec"),
			c:   exec.CommandContext,
		},
	}

	err = i.do()
	if err != nil {
		i.log.Error("installation failed", "duration", time.Since(start))
		panic(err)
	}

	i.log.Info("installation succeeded", "duration", time.Since(start))
	return nil
}

func (i *installer) do() error {
	err := i.detectFirmware()
	if err != nil {
		i.log.Warn("no efi detected", "error", err)
		return err
	}

	if !i.fileExists(installYAML) {
		return fmt.Errorf("no install.yaml found")
	}

	// remove .dockerenv, otherwise systemd-detect-virt guesses docker which modifies the behavior of many services.
	if i.fileExists("/.dockerenv") {
		err := os.Remove("/.dockerenv")
		if err != nil {
			return fmt.Errorf("unable to delete .dockerenv")
		}
	}

	err = i.writeResolvConf()
	if err != nil {
		i.log.Warn("writing resolv.conf failed", "error", err)
		return err
	}

	err = i.createMetalUser()
	if err != nil {
		return err
	}
	err = i.configureNetwork()
	if err != nil {
		return err
	}

	err = i.copySSHKeys()
	if err != nil {
		return err
	}

	err = i.fixPermissions()
	if err != nil {
		return err
	}

	err = i.processUserdata()
	if err != nil {
		return err
	}

	cmdLine := i.buildCMDLine()

	err = i.writeBootInfo(cmdLine)
	if err != nil {
		return err
	}

	err = i.grubInstall(cmdLine)
	if err != nil {
		return err
	}

	err = i.unsetMachineID()
	if err != nil {
		return err
	}

	err = i.writeBuildMeta()
	if err != nil {
		return err
	}

	return nil
}

func (i *installer) detectFirmware() error {
	i.log.Info("detect firmware")

	if !i.isVirtual() && !i.fileExists("/sys/firmware/efi") {
		return fmt.Errorf("not running efi mode")
	}
	return nil
}

func (i *installer) isVirtual() bool {
	return !i.fileExists("/sys/class/dmi")
}

func (i *installer) unsetMachineID() error {
	i.log.Info("unset machine-id")
	for _, p := range []string{"/etc/machine-id", "/var/lib/dbus/machine-id"} {
		if !i.fileExists(p) {
			continue
		}
		f, err := i.fs.Create(p)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

func (i *installer) fileExists(filename string) bool {
	info, err := i.fs.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (i *installer) writeResolvConf() error {
	i.log.Info("write /etc/resolv.conf")
	// Must be written here because during docker build this file is synthetic
	// FIXME enable systemd-resolved based approach again once we figured out why it does not work on the firewall
	// most probably because the resolved must be running in the internet facing vrf.
	// ln -sf /run/systemd/resolve/stub-resolv.conf /etc/resolv.conf
	// in ignite this file is a symlink to /proc/net/pnp, to pass integration test, remove this first
	err := i.fs.Remove("/etc/resolv.conf")
	if err != nil {
		i.log.Info("no /etc/resolv.conf present")
	}

	// FIXME migrate to dns0.eu resolvers
	content := []byte(
		`nameserver 8.8.8.8
nameserver 8.8.4.4
`)
	return afero.WriteFile(i.fs, "/etc/resolv.conf", content, 0644)
}

func (i *installer) buildCMDLine() string {
	i.log.Info("build kernel cmdline")

	rootUUID := i.config.RootUUID

	parts := []string{
		fmt.Sprintf("console=%s", i.config.Console),
		fmt.Sprintf("root=UUID=%s", rootUUID),
		"init=/sbin/init",
		"net.ifnames=0",
		"biosdevname=0",
		"nvme_core.io_timeout=300", // 300 sec should be enough for firewalls to be replaced
		"systemd.unified_cgroup_hierarchy=0",
	}

	mdUUID, found := i.findMDUUID()
	if found {
		mdParts := []string{
			"rdloaddriver=raid0",
			"rdloaddriver=raid1",
			fmt.Sprintf("rd.md.uuid=%s", mdUUID),
		}
		parts = append(parts, mdParts...)
	}

	return strings.Join(parts, " ")
}

func (i *installer) findMDUUID() (mdUUID string, found bool) {
	i.log.Info("detect software raid uuid")
	if !i.config.RaidEnabled {
		return "", false
	}

	blkidOut, err := i.exec.command(&cmdParams{
		name:    "blkid",
		timeout: 10 * time.Second,
	})
	if err != nil {
		i.log.Error("unable to run blkid", "error", err)
		return "", false
	}
	rootUUID := i.config.RootUUID

	var rootDisk string
	for _, line := range strings.Split(string(blkidOut), "\n") {
		if strings.Contains(line, rootUUID) {
			rd, _, found := strings.Cut(line, ":")
			if found {
				rootDisk = strings.TrimSpace(rd)
				break
			}
		}
	}
	if rootDisk == "" {
		i.log.Error("unable to detect rootdisk")
		return "", false
	}

	mdadmOut, err := i.exec.command(&cmdParams{
		name:    "mdadm",
		args:    []string{"--detail", "--export", rootDisk},
		timeout: 10 * time.Second,
	})
	if err != nil {
		i.log.Error("unable to run mdadm", "error", err)
		return "", false
	}

	for _, line := range strings.Split(string(mdadmOut), "\n") {
		_, md, found := strings.Cut(line, "MD_UUID=")
		if found {
			mdUUID = md
			break
		}
	}

	if mdUUID == "" {
		i.log.Error("unable to detect md root disk")
		return "", false
	}

	return mdUUID, true
}

func (i *installer) createMetalUser() error {
	i.log.Info("create user", "user", "metal")

	u, err := user.Lookup("metal")
	if err != nil {
		if err.Error() != user.UnknownUserError("metal").Error() {
			return err
		}
	}
	if u != nil {
		i.log.Info("user already exists, recreating")
		_, err = i.exec.command(&cmdParams{
			name:    "userdel",
			args:    []string{"metal"},
			timeout: 10 * time.Second,
		})
		if err != nil {
			return err
		}
	}

	_, err = i.exec.command(&cmdParams{
		name:    "useradd",
		args:    []string{"--create-home", "--uid", "1000", "--gid", i.oss.SudoGroup(), "--shell", "/bin/bash", "metal"},
		timeout: 10 * time.Second,
	})
	if err != nil {
		return err
	}

	_, err = i.exec.command(&cmdParams{
		name:    "passwd",
		args:    []string{"metal"},
		timeout: 10 * time.Second,
		stdin:   i.config.Password + "\n" + i.config.Password + "\n",
	})
	if err != nil {
		return err
	}

	if i.oss == osAlmalinux {
		// otherwise in rescue mode the root account is locked
		_, err = i.exec.command(&cmdParams{
			name:    "passwd",
			args:    []string{"root"},
			timeout: 10 * time.Second,
			stdin:   i.config.Password + "\n" + i.config.Password + "\n",
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *installer) configureNetwork() error {
	i.log.Info("configure network")
	kb, err := netconf.New(i.log.WithGroup("networker"), installYAML)
	if err != nil {
		return err
	}

	var kind netconf.BareMetalType
	switch i.config.Role {
	case "firewall":
		kind = netconf.Firewall
	case "machine":
		kind = netconf.Machine
	default:
		return fmt.Errorf("unknown role:%s", i.config.Role)
	}

	err = kb.Validate(kind)
	if err != nil {
		return err
	}

	c, err := netconf.NewConfigurator(kind, *kb, false)
	if err != nil {
		return err
	}
	c.Configure(netconf.ForwardPolicyDrop)
	return nil
}

func (i *installer) copySSHKeys() error {
	i.log.Info("copy ssh keys")
	err := i.fs.MkdirAll("/home/metal/.ssh", 0700)
	if err != nil {
		return err
	}

	u, err := user.Lookup("metal")
	if err != nil {
		return err
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		return err
	}

	err = i.fs.Chown("/home/metal/.ssh", uid, gid)
	if err != nil {
		return err
	}

	err = afero.WriteFile(i.fs, "/home/metal/.ssh/authorized_keys", []byte(i.config.SSHPublicKey), 0600)
	if err != nil {
		return err
	}
	return i.fs.Chown("/home/metal/.ssh/authorized_keys", uid, gid)
}

func (i *installer) fixPermissions() error {
	i.log.Info("fix permissions")
	for p, perm := range map[string]fs.FileMode{
		"/var/tmp":   01777,
		"/etc/hosts": 0644,
	} {
		err := i.fs.Chmod(p, perm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *installer) processUserdata() error {
	i.log.Info("process userdata")
	if ok := i.fileExists(userdata); !ok {
		i.log.Info("no userdata present, not processing userdata", "path", userdata)
		return nil
	}

	content, err := afero.ReadFile(i.fs, userdata)
	if err != nil {
		return err
	}

	defer func() {
		out, err := i.exec.command(&cmdParams{
			name: "systemctl",
			args: []string{"preset-all"},
		})
		if err != nil {
			i.log.Error("error when running systemctl preset-all, continuing anyway", "error", err, "output", string(out))
		}
	}()

	if isCloudInitFile(content) {
		_, err := i.exec.command(&cmdParams{
			name: "cloud-init",
			args: []string{"devel", "schema", "--config-file", userdata},
		})
		if err != nil {
			i.log.Error("error when running cloud-init userdata, continuing anyway", "error", err)
		}

		return nil
	}

	err = i.fs.Rename(userdata, "/etc/metal/config.ign")
	if err != nil {
		return err
	}

	rawConfig, err := afero.ReadFile(i.fs, "/etc/metal/config.ign")
	if err != nil {
		return err
	}
	_, report, err := config.Parse(rawConfig)
	if err != nil {
		i.log.Error("error when validating ignition userdata, continuing anyway", "error", err)
	}

	i.log.Info("executing ignition")
	_, err = i.exec.command(&cmdParams{
		name: "ignition",
		args: []string{"-oem", "file", "-stage", "files", "-log-to-stdout"},
		dir:  "/etc/metal",
	})
	if err != nil {
		i.log.Error("error when running ignition, continuing anyway", "report", report.Entries, "error", err)
	}

	return nil
}

func isCloudInitFile(content []byte) bool {
	for i, line := range strings.Split(string(content), "\n") {
		if strings.Contains(line, "#cloud-config") {
			return true
		}
		if i > 1 {
			return false
		}
	}
	return false
}

func (i *installer) writeBootInfo(cmdLine string) error {
	i.log.Info("write boot-info")

	kern, initrd, err := i.kernelAndInitrdPath()
	if err != nil {
		return err
	}

	content, err := yaml.Marshal(api.Bootinfo{
		Initrd:       initrd,
		Cmdline:      cmdLine,
		Kernel:       kern,
		BootloaderID: i.oss.BootloaderID(),
	})
	if err != nil {
		return fmt.Errorf("unable to write boot-info.yaml %w", err)
	}

	return afero.WriteFile(i.fs, "/etc/metal/boot-info.yaml", content, 0700)
}

func (i *installer) kernelAndInitrdPath() (kern string, initrd string, err error) {
	// Debian 10
	// root@1f223b59051bcb12:/boot# ls -l
	// total 83500
	// -rw-r--r-- 1 root root       83 Aug 13 15:25 System.map-5.10.0-17-amd64
	// -rw-r--r-- 1 root root   236286 Aug 13 15:25 config-5.10.0-17-amd64
	// -rw-r--r-- 1 root root    93842 Jul 19  2021 config-5.10.51
	// drwxr-xr-x 2 root root     4096 Oct  3 11:21 grub
	// -rw-r--r-- 1 root root 34665690 Oct  3 11:22 initrd.img-5.10.0-17-amd64
	// lrwxrwxrwx 1 root root       21 Jul 19  2021 vmlinux -> /boot/vmlinux-5.10.51
	// -rwxr-xr-x 1 root root 43526368 Jul 19  2021 vmlinux-5.10.51
	// -rw-r--r-- 1 root root  6962816 Aug 13 15:25 vmlinuz-5.10.0-17-amd64

	// Ubuntu 20.04
	// root@568551f94559b121:~# ls -l /boot/
	// total 83500
	// -rw-r--r-- 1 root root       83 Aug 13 15:25 System.map-5.10.0-17-amd64
	// -rw-r--r-- 1 root root   236286 Aug 13 15:25 config-5.10.0-17-amd64
	// -rw-r--r-- 1 root root    93842 Jul 19  2021 config-5.10.51
	// drwxr-xr-x 2 root root     4096 Oct  3 11:21 grub
	// -rw-r--r-- 1 root root 34665690 Oct  3 11:22 initrd.img-5.10.0-17-amd64
	// lrwxrwxrwx 1 root root       21 Jul 19  2021 vmlinux -> /boot/vmlinux-5.10.51
	// -rwxr-xr-x 1 root root 43526368 Jul 19  2021 vmlinux-5.10.51
	// -rw-r--r-- 1 root root  6962816 Aug 13 15:25 vmlinuz-5.10.0-17-amd64

	var (
		bootPartition   = "/boot"
		systemMapPrefix = "/boot/System.map-"
	)

	systemMaps, err := afero.Glob(i.fs, systemMapPrefix+"*")
	if err != nil {
		return "", "", fmt.Errorf("unable to find a System.map, probably no kernel installed %w", err)
	}
	if len(systemMaps) != 1 {
		return "", "", fmt.Errorf("more or less than a single System.map found(%v), probably no kernel or more than one kernel installed", systemMaps)
	}

	systemMap := systemMaps[0]
	_, kernelVersion, found := strings.Cut(systemMap, systemMapPrefix)
	if !found {
		return "", "", fmt.Errorf("unable to detect kernel version in System.map :%q", systemMap)
	}

	kern = path.Join(bootPartition, "vmlinuz"+"-"+kernelVersion)
	if !i.fileExists(kern) {
		return "", "", fmt.Errorf("kernel image %q not found", kern)
	}
	initrd = path.Join(bootPartition, i.oss.Initramdisk(kernelVersion))
	if !i.fileExists(initrd) {
		return "", "", fmt.Errorf("ramdisk %q not found", initrd)
	}

	i.log.Info("detect kernel and initrd", "kernel", kern, "initrd", initrd)

	return
}

func (i *installer) grubInstall(cmdLine string) error {
	i.log.Info("install grub")
	// ttyS1,115200n8
	serialPort, serialSpeed, found := strings.Cut(i.config.Console, ",")
	if !found {
		return fmt.Errorf("serial console could not be split into port and speed")
	}

	_, serialPort, found = strings.Cut(serialPort, "ttyS")
	if !found {
		return fmt.Errorf("serial port could not be split")
	}

	serialSpeed, _, found = strings.Cut(serialSpeed, "n8")
	if !found {
		return fmt.Errorf("serial speed could not be split")
	}

	defaultGrub := fmt.Sprintf(`GRUB_DEFAULT=0
GRUB_TIMEOUT=5
GRUB_DISTRIBUTOR=%s
GRUB_CMDLINE_LINUX_DEFAULT=""
GRUB_CMDLINE_LINUX="%s"
GRUB_TERMINAL=serial
GRUB_SERIAL_COMMAND="serial --speed=%s --unit=%s --word=8"
`, i.oss.BootloaderID(), cmdLine, serialSpeed, serialPort)

	if i.oss == osAlmalinux {
		defaultGrub += fmt.Sprintf("GRUB_DEVICE=UUID=%s\n", i.config.RootUUID)
		defaultGrub += "GRUB_ENABLE_BLSCFG=false\n"
	}

	err := afero.WriteFile(i.fs, "/etc/default/grub", []byte(defaultGrub), 0755)
	if err != nil {
		return err
	}

	grubInstallArgs := []string{
		"--target=x86_64-efi",
		"--efi-directory=/boot/efi",
		"--boot-directory=/boot",
		"--bootloader-id=" + i.oss.BootloaderID(),
	}
	if i.config.RaidEnabled {
		grubInstallArgs = append(grubInstallArgs, "--no-nvram")
	}

	if i.oss == osAlmalinux {
		path := "/boot/grub2/grub.cfg"
		if i.oss == osAlmalinux {
			path = "/boot/efi/EFI/almalinux/grub.cfg"
		}
		_, err := i.exec.command(&cmdParams{
			name: "grub2-mkconfig",
			args: []string{"-o", path},
		})
		if err != nil {
			return err
		}

		grubInstallArgs = append(grubInstallArgs, fmt.Sprintf("UUID=%s", i.config.RootUUID))
	} else {
		grubInstallArgs = append(grubInstallArgs, "--removable")
	}

	if i.config.RaidEnabled {
		out, err := i.exec.command(&cmdParams{
			name:    "mdadm",
			args:    []string{"--examine", "--scan"},
			timeout: 10 * time.Second,
		})
		if err != nil {
			return err
		}

		out += "\nMAILADDR root\n"

		err = afero.WriteFile(i.fs, "/etc/mdadm.conf", []byte(out), 0700)
		if err != nil {
			return err
		}

		if i.oss.NeedUpdateInitRamfs() {
			err = i.fs.MkdirAll("/var/lib/initramfs-tools", 0755)
			if err != nil {
				return err
			}

			_, err = i.exec.command(&cmdParams{
				name: "update-initramfs",
				args: []string{"-u"},
			})
			if err != nil {
				return err
			}
		}

		out, err = i.exec.command(&cmdParams{
			name: "blkid",
		})
		if err != nil {
			return err
		}

		for _, line := range strings.Split(string(out), "\n") {
			if strings.Contains(line, `PARTLABEL="efi"`) {
				disk, _, found := strings.Cut(line, ":")
				if !found {
					return fmt.Errorf("unable to process blkid output lines")
				}
				shim := fmt.Sprintf(`\\EFI\\%s\\grubx64.efi`, i.oss.BootloaderID())
				if i.oss == osAlmalinux {
					shim = fmt.Sprintf(`\\EFI\\%s\\shimx64.efi`, i.oss.BootloaderID())
				}

				_, err = i.exec.command(&cmdParams{
					name: "efibootmgr",
					args: []string{"-c", "-d", disk, "-p1", "-l", shim, "-L", i.oss.BootloaderID()},
				})
				if err != nil {
					return err
				}
			}
		}
	}

	if i.oss.GrubInstallCmd() != "" && !runFromCI() {
		_, err = i.exec.command(&cmdParams{
			name: i.oss.GrubInstallCmd(),
			args: grubInstallArgs,
		})
		if err != nil {
			return err
		}
	}

	if i.oss == osAlmalinux {
		if !i.config.RaidEnabled {
			return nil
		}

		v, err := i.getKernelVersion()
		if err != nil {
			return err
		}

		_, err = i.exec.command(&cmdParams{
			name: "dracut",
			args: []string{
				"--mdadmconf",
				"--kver", v,
				"--kmoddir", "/lib/modules/" + v,
				"--include", "/lib/modules/" + v, "/lib/modules/" + v,
				"--fstab",
				"--add=dm mdraid",
				"--add-drivers=raid0 raid1",
				"--hostonly",
				"--force",
			},
		})
		if err != nil {
			return err
		}

		return nil
	}

	_, err = i.exec.command(&cmdParams{
		name: "update-grub2",
	})
	if err != nil {
		return err
	}

	_, err = i.exec.command(&cmdParams{
		name: "dpkg-reconfigure",
		args: []string{"grub-efi-amd64-bin"},
		env: []string{
			"DEBCONF_NONINTERACTIVE_SEEN=true",
			"DEBIAN_FRONTEND=noninteractive",
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (i *installer) writeBuildMeta() error {
	i.log.Info("writing build meta file", "path", "/etc/metal/build-meta.yaml")

	meta := &v1.BuildMeta{
		Version:  v.Version,
		Date:     v.BuildDate,
		SHA:      v.GitSHA1,
		Revision: v.Revision,
	}

	out, err := i.exec.command(&cmdParams{
		name: "ignition",
		args: []string{"-version"},
	})
	if err != nil {
		i.log.Error("error detecting ignition version for build meta, continuing anyway", "error", err)
	} else {
		meta.IgnitionVersion = strings.TrimSpace(out)
	}

	content, err := yaml.Marshal(meta)
	if err != nil {
		return err
	}

	content = append([]byte("---\n"), content...)

	return afero.WriteFile(i.fs, "/etc/metal/build-meta.yaml", content, 0644)
}

func (i *installer) getKernelVersion() (string, error) {
	kern, _, err := i.kernelAndInitrdPath()
	if err != nil {
		return "", err
	}

	_, version, found := strings.Cut(kern, "vmlinuz-")
	if !found {
		return "", fmt.Errorf("unable to determine kernel version from: %s", kern)
	}

	return version, nil
}
