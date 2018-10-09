# Discover a device with metal-hammer

in order to be able to register a new device, or check whether a device is already registered, we execute from the pxeboot image a binary which does the hardware discovery and send the output to the maas api.

## Build

```bash
make initrd
```

## Local Testing

```
vagrant destroy -f && make initrd && vagrant up && virsh console metal-hammer_pxeclient
```


## Create a PXE boot image with linuxkit and u-root

In order to be able to create a kernel and initrd image which is suitable to boot a bare metal server with the required tools to discover and install the target os, we use linuxkit and u-root.

### Quickstart

- download linuxkit:

```bash
sudo curl -fSL https://github.com/linuxkit/linuxkit/releases/download/v0.6/linuxkit-linux-amd64 -o /usr/local/bin/linuxkit && sudo chmod +x /usr/local/bin/linuxkit
```

- download u-root:

```
go get -u github.com/u-root/u-root
```

- build the kernel:

```bash
linuxkit build metal-hammer.yaml
```

- build the initrd

```bash
make initrd
```

### check content

```
cpio -itv < metal-hammer-initrd.img
```

### start it

```
vagrant up
```
