SHA := $(shell git rev-parse --short=8 HEAD)
GITVERSION := $(shell git describe --long --all)
BUILDDATE := $(shell date --iso-8601=seconds)
VERSION := $(or ${VERSION},$(shell git describe --tags --exact-match 2> /dev/null || git symbolic-ref -q --short HEAD || git rev-parse --short HEAD))
GO := go
GOSRC = $(shell find . -not \( -path vendor -prune \) -type f -name '*.go')

BINARY := metal-hammer
INITRD := ${BINARY}-initrd.img
COMPRESSOR := lz4
COMPRESSOR_ARGS := -f -l
INITRD_COMPRESSED := ${INITRD}.${COMPRESSOR}
MAINMODULE := .
CGO_ENABLED := 1
# export CGO_LDFLAGS := "-lsystemd" "-lpcap" "-ldbus-1"
BINARY_TAR := ${BINARY}.tar
HASH := md5
INITRD_HASH := ${INITRD_COMPRESSED}.${HASH}


.PHONY: all
all:: bin/$(BINARY);

in-docker: all;

LINKMODE := -linkmode external -extldflags '-static -s -w' \
		 -X 'github.com/metal-stack/v.Version=$(VERSION)' \
		 -X 'github.com/metal-stack/v.Revision=$(GITVERSION)' \
		 -X 'github.com/metal-stack/v.GitSHA1=$(SHA)' \
		 -X 'github.com/metal-stack/v.BuildDate=$(BUILDDATE)'

bin/$(BINARY): test $(GOSRC)
	$(info CGO_ENABLED="$(CGO_ENABLED)")
	$(GO) build \
		-tags netgo \
		-ldflags \
		"$(LINKMODE)" \
		-o bin/$(BINARY) \
		$(MAINMODULE)
	strip bin/$(BINARY)

.PHONY: test
test:
	CGO_ENABLED=1 $(GO) test -cover ./...

.PHONY: clean
clean::
	rm -f ${INITRD} ${INITRD_COMPRESSED}

${INITRD_COMPRESSED}:
	rm -f ${INITRD_COMPRESSED}
	docker buildx build -o - . > ${BINARY_TAR}
	tar -xf ${BINARY_TAR} ${INITRD_COMPRESSED}
	rm -f ${BINARY_TAR}
	md5sum ${INITRD_COMPRESSED} > ${INITRD_HASH}

.PHONY: initrd
initrd: ${INITRD_COMPRESSED}

# place all binaries in the same directory (/sbin) which is in the PATH of root.
# keep them alphabetically sorted
.PHONY: ramdisk
ramdisk:
	GO111MODULE=off u-root \
		-format=cpio -build=bb \
		-defaultsh=/bin/bash \
		-files="bin/metal-hammer:bbin/uinit" \
		-files="/bin/bash:bin/bash" \
		-files="/bin/netstat:bin/netstat" \
		-files="/etc/localtime:etc/localtime" \
		-files="/etc/lvm/lvm.conf:etc/lvm/lvm.conf" \
		-files="/etc/ssl/certs/ca-certificates.crt:etc/ssl/certs/ca-certificates.crt" \
		-files="/lib/x86_64-linux-gnu/libnss_files-2.31.so:lib/x86_64-linux-gnu/libnss_files-2.31.so" \
		-files="/lib/x86_64-linux-gnu/libnss_files.so.2:lib/x86_64-linux-gnu/libnss_files.so.2" \
		-files="/sbin/blkid:sbin/blkid" \
		-files="/sbin/ethtool:sbin/ethtool" \
		-files="/sbin/hdparm:sbin/hdparm" \
		-files="/sbin/lvm:sbin/lvm" \
		-files="/sbin/mdadm:sbin/mdadm" \
		-files="/sbin/mdmon:sbin/mdmon" \
		-files="/sbin/mke2fs:sbin/mke2fs" \
		-files="/sbin/mkfs.ext3:sbin/mkfs.ext3" \
		-files="/sbin/mkfs.ext4:sbin/mkfs.ext4" \
		-files="/sbin/mkfs.fat:sbin/mkfs.fat" \
		-files="/sbin/mkfs.vfat:sbin/mkfs.vfat" \
		-files="/sbin/mkswap:sbin/mkswap" \
		-files="/sbin/sgdisk:sbin/sgdisk" \
		-files="/sbin/wipefs:sbin/wipefs" \
		-files="/usr/bin/ipmitool:usr/bin/ipmitool" \
		-files="/usr/bin/efibootmgr:/usr/bin/efibootmgr" \
		-files="/usr/bin/lspci:bin/lspci" \
		-files="/usr/bin/strace:bin/strace" \
		-files="/usr/sbin/nvme:sbin/nvme" \
		-files="/usr/share/misc/pci.ids:usr/share/misc/pci.ids" \
		-files="lvmlocal.conf:etc/lvm/lvmlocal.conf" \
		-files="passwd:etc/passwd" \
		-files="varrun:var/run/keep" \
		-files="ice.pkg:lib/firmware/intel/ice/ddp/ice.pkg" \
		-files="metal.key:id_rsa" \
		-files="metal.key.pub:authorized_keys" \
		-files="sum:sbin/sum" \
	-o ${INITRD} \
	&& ${COMPRESSOR} ${COMPRESSOR_ARGS} ${INITRD} ${INITRD_COMPRESSED} \
	&& rm -f ${INITRD}

vagrant-destroy:
	vagrant destroy -f

vagrant-up: vagrant-destroy
	vagrant up && virsh console metal-hammerpxeclient

# TODO make this work as with vagrant as a lightweight alternative.
# networking is not working atm.
# http://nickdesaulniers.github.io/blog/2018/10/24/booting-a-custom-linux-kernel-in-qemu-and-debugging-it-with-gdb/
qemu-up:
	qemu-system-x86_64 \
		--enable-kvm \
		-m 512 \
		-nographic \
		-object rng-random,filename=/dev/urandom,id=rng0 \
		-device virtio-rng-pci,rng=rng0 \
		-append "console=ttyS0 ip=dhcp \
          METAL_CORE_ADDRESS=192.168.121.1:4712 \
          IMAGE_ID=default  \
          SIZE_ID=v1-small-x86  \
          IMAGE_URL=http://192.168.121.1:4711/images/ubuntu/19.04/img.tar.lz4  \
          DEBUG=1  \
          BGP=1" \
		-kernel metal-kernel \
		-initrd metal-hammer-initrd.img.lz4

start:
	# sudo setcap cap_net_admin+ep ~/bin/cloud-hypervisor
	# /usr/src/linux-headers-6.5.0-9/scripts/extract-vmlinux metal-kernel-6.6.2 > vmlinux
	cloud-hypervisor \
		--kernel ./vmlinux \
		--console off \
		--serial tty \
		--initramfs=metal-hammer-initrd.img.lz4 \
		--cmdline "console=ttyS0" \
		--cpus boot=4 \
		--memory size=1024M \
		--net "tap=,mac=,ip=,mask="