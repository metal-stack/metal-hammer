# Discover a device with metal-hammer

in order to be able to register a new device, or check whether a device is already registered, we execute from the pxeboot image a binary which does the hardware discovery. This is done with the *lshw* command which is available on various linux distributions, call it with the `-json` option as root and send the output to the maas api.

## Build

```bash
make
```

## Usage

### Configuration

Is done via Environment variables.

```bash
bin/metal-hammer -h

This application is configured via the environment. The following environment
variables can be used:

KEY                   TYPE             DEFAULT                                  REQUIRED    DESCRIPTION
METAL_HAMMER_DEBUG    True or False    false                                    False       turn on debug log
```

### Execution

```bash
sudo bin/metal-hammer
INFO[09-24|14:09:59] configuration                            debug=false reportURL=http://localhost:8080/device/register
INFO[09-24|14:10:00] device already registered                uuid=4C3CEF61-F536-B211-A85C-B765E03E138F caller=lshw.go:63
```


## Create a PXE boot image with linuxkit

In order to be able to create a kernel and initrd image which is suitable to boot a bare metal server with the required tools to discover and install the target os, we use linuxkit.

### Quickstart

- download linuxkit:

```bash
sudo curl -fSL https://github.com/linuxkit/linuxkit/releases/download/v0.6/linuxkit-linux-amd64 -o /usr/local/bin/linuxkit && sudo chmod +x /usr/local/bin/linuxkit

OR

go get -u github.com/linuxkit/linuxkit/src/cmd/linuxkit

```

- build the kernel and image:

```bash
linuxkit build pxeboot.yaml
```

- check by running it:

```bash
linuxkit run qemu -disk size=4G pxeboot
```

## Create a initial ramdisk with u-root

```
u-root -format=cpio -build=bb -files="bin/metal-hammer:bbin/metal-hammer" -o metal-initrd.cpio
```

executing metal-hammer directly:
```
u-root -format=cpio -build=bb -files="bin/metal-hammer:bbin/metal-hammer" -defaultsh="/bbin/metal-hammer" -o metal-initrd.cpio
```

### check content

```
cpio -itv < metal-initrd.cpio 
```

### start it

```
qemu-system-x86_64 -m 2G -kernel pxeboot-kernel -initrd metal-initrd.cpio -nographic -append console=ttyS0
```
