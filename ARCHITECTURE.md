# Metal Hammer Architecture

Metal Hammer is used to install a virtual or physical server with a configurable operating system as fast as possible and in a reproducible manner.

## Goals

Installation of the Operating System is done in a fixed way with not much configurable knobs.
These operating systems (OS) is almost always used in a cloud native environment where general purpose is not required.
Instead it focuses on container only workloads, orchestrated with kubernetes.

Therefore only a very minimal OS installation is required with only a minimal set of enabled and installed features, daemons.
There is also a fixed disk configuration which does not fancy stuff like raid, lvm, swap etc.

The installation procedure is fully automatic with no manual interaction possible.

Choice of a small set of different Linux distributions.

## Non Goals

Configuration of the OS installation.

## High level design

The installation of a OS is based on a PXE boot approach. A linux kernel and a initrd is loaded from a central location.
The initrd contains the `metal-hammer` binary which acts as *init* process.

`metal-hammer` requires a control plane which delivers the minimal basic configuration items, this is called `metal-core`.
Full API documentation for this communication follows below.

All Disks are wiped and the first disk available is formatted and a tarball of the target OS is unpacked onto this disk.
After that a kexec of the new installed kernel with the contents on that disk is executed, 
the target OS installation is now usable.

On every error which happens during installation, the process dies and triggers a reboot.

## Theorie of operation
