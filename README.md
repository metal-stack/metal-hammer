# metal-stack.io | metal-hammer

![Go version](https://img.shields.io/github/go-mod/go-version/metal-stack/metal-hammer)
[![Go Report Card](https://goreportcard.com/badge/github.com/metal-stack/metal-hammer)](https://goreportcard.com/report/github.com/metal-stack/metal-hammer)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/metal-stack/metal-hammer)
[![Build](https://github.com/metal-stack/metal-hammer/actions/workflows/build.yml/badge.svg?branch=master)](https://github.com/metal-stack/metal-hammer/actions)
[![Slack](https://img.shields.io/badge/slack-metal--stack-brightgreen.svg?logo=slack)](https://metal-stack.slack.com/)

The `metal-hammer` is a component of metal-stack for performing bare-metal provisioning of servers. It runs as an `initrd` that contains a small Go binary as `init` process. The `metal-hammer` is loaded during the PXE booting process together with the `metal-kernel`.

This component performs the following actions:

- Ensures a reboot after 24 hours, if the machine was not allocated yet
- Creates or updates the [BMC users](https://metal-stack.io/docs/next/security-principles/#bmc-user-management), that are used for administrative tasks
- Registers the machine at the `metal-api` using the following hardware specifications:
  - CPUs (vendor, model, cores, threads)
  - GPUs (vendor, model)
  - NICs (MAC, interface name, interface neighbors, MAC of switch chassis)
  - Disks (size, device path)
  - IPMI interface details (IP address, MAC, BMC firmware version)
  - IPMI FRU details (board manufacturer, board type)
  - BIOS (vendor, version, date)
- Ensures all interfaces are up for link-local neighbor discovery
- Wipes disks:
  - Using `mkfs.ext4` (and `-E discard` option) for rotational disks
  - Using nvme cli (and `--format` option) for NVMe disks
  - Using `dd` as fallback option
- Ensures BIOS uses UEFI mode
- Ensures `Waiting` loop until machine gets requested by the `metal-api`

If a machine is requested, `metal-hammer` initiates the installation including the following steps:

- Installs a filesystem layout (defined in the `metal-api`):
  - Cleanup disks using `wipefs` to remove signatures
  - Create partitions using `sgdisk` (and `--zap-all` option, to destroy partition table structure)
  - Create RAID using `mdadm`
  - Create logical volumes using `lvm`
  - Create filesystems. Supported are: `ext3`, `ext4`, `swap`, `vfat`
  - Mount filesystems
- Downloads the requested OS image as a `.tar` and unpacks it into a directory on the disk
- Writes the requested allocation configuration for the OS (e.g. hostname, networks, DNS servers, NTP servers...) into `/etc/metal/install.yaml`
- Writes the user specific data into `/etc/metal/userdata`
- Writes the LVM configuration into `/etc/lvm/lvmlocal.conf` for compatibility with the new OS
- Executes the Go binary `install.go` inside the downloaded OS metal-image
- Reads the kernel's boot info from `/etc/metal/boot-info.yaml`
- Configures the boot order to boot the installed OS after rebooting
- Writes all fstab entries to `/etc/fstab` inside chroot
- Reports the installation to the `metal-api`
- Boots into the new kernel from the installed OS

## Local Development

Use the following command for local development:

```bash
make clean initrd vagrant-up
```

## Create your own PXE boot initrd with u-root

We use `u-root` to create an `initrd` image which is suitable to boot a bare metal server with the required tools to discover and install a target OS.
Follow these steps to create your custom one:

```bash
# Download u-root
go get -u github.com/u-root/u-root

# Build the initrd
make initrd

# Verify content
cpio -itv < metal-hammer-initrd.img

# Test it
make vagrant-up
```
