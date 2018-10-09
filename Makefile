.ONESHELL:
SHA := $(shell git rev-parse --short=8 HEAD)
GITVERSION := $(shell git describe --long --all)
BUILDDATE := $(shell date -Iseconds)
VERSION := $(or ${VERSION},devel)

BINARY := bin/metal-hammer

.PHONY: clean image

all: $(BINARY)

${BINARY}: clean
	CGO_ENABLE=0 \
	GO111MODULE=on \
	go build \
		-tags netgo \
		-ldflags "-X 'main.version=$(VERSION)' \
				  -X 'main.revision=$(GITVERSION)' \
				  -X 'main.gitsha1=$(SHA)' \
				  -X 'main.builddate=$(BUILDDATE)'" \
	-o $@

clean:
	rm -f ${BINARY}

image:
	docker build -t registry.fi-ts.io/metal/metal-hammer .

SGDISK := $(shell which sgdisk)
VFAT := $(shell which mkfs.vfat)
FAT := $(shell which mkfs.fat)
EXT4 := $(shell which mkfs.ext4)
MKFS := $(shell which mke2fs)
RNGD := $(shell which rngd)

uroot: ${BINARY}
	u-root \
		-format=cpio -build=bb \
		-files="bin/metal-hammer:bbin/metal-hammer" \
		-files="${SGDISK}:usr/bin/sgdisk" \
		-files="${VFAT}:sbin/mkfs.vfat" \
		-files="${EXT4}:sbin/mkfs.ext4" \
		-files="${MKFS}:sbin/mke2fs" \
		-files="${FAT}:sbin/mkfs.fat" \
		-files="${RNGD}:usr/sbin/rngd" \
		-files="metal.key:id_rsa" \
		-files="metal.key.pub:authorized_keys" \
		-files="metal-hammer.sh:bbin/uinit" \
	-o metal-hammer-initrd.img

initrd: uroot
	gzip -f metal-hammer-initrd.img
