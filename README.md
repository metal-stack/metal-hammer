# Discover a machine with metal-hammer

in order to be able to register a new machine, or check whether a machine is already registered, we execute from the pxeboot image a binary which does the hardware discovery and send the output to the metal api.

# Build

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
