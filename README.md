# Discover a machine with metal-hammer

in order to be able to register a new machine, or check whether a machine is already registered, we execute from the pxeboot image a binary which does the hardware discovery and send the output to the metal api.

## Build

```bash
make initrd
```

## Local Testing

```
vagrant destroy -f && make clean initrd && vagrant up && virsh console metal-hammer_pxeclient
```

## Create a PXE boot initrd with u-root

In order to be able to create an initrd image which is suitable to boot a bare metal server with the required tools to discover and install the target os, we use u-root.

### Quickstart

- download u-root:

```
go get -u github.com/u-root/u-root
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
