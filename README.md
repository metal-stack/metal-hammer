# Metal Stack Hammer

Hammer is used to boot a bare metal server via PXE together with the Metal Stack kernel. Hammer is a initrd which runs a small golang binary as init process. This does the following actions:

- Ensures all interfaces are up
- Check if the server was booted in UEFI, if not modify the bios tu uefi and reboots
- Wipes as existing disks by either:
  - run secure erase if possible by using the mechanism in modern disks, this is true for most SSD´s and NVME disks.
  - If not possible run mkfs.ext4 --discard on the disks.
- Gather HW informations and report them back to metal-api:
  - CPU Core count
  - Memory count
  - Disks with their size and device path
  - Network adapters which have an active uplink with their interface name, own mac address and mac address of the switch chassis where this network card is connected to. 2 distinct switch chassis are required.
  - IPMI interface with mac and ipaddress.
  - create a metal user on IPMI with a strong password
- Set BIOS boot order to contain only PXE and Hard Disk as possible options.
- Wait until a `machine create` command was issued from metal-api

## Local Testing

```bash
make clean initrd vagrant-up
```

## Create a PXE boot initrd with u-root

In order to be able to create an initrd image which is suitable to boot a bare metal server with the required tools to discover and install the target os, we use u-root.

### Quickstart

- download u-root:

```bash
go get -u github.com/u-root/u-root
```

- build the initrd

```bash
make initrd
```

### check content

```bash
cpio -itv < metal-hammer-initrd.img
```

### start it

```bash
make vagrant-up
```
