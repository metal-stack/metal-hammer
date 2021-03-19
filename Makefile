BINARY := metal-hammer
INITRD := ${BINARY}-initrd.img
COMPRESSOR := lz4
COMPRESSOR_ARGS := -f -l
INITRD_COMPRESSED := ${INITRD}.${COMPRESSOR}
MAINMODULE := .
COMMONDIR := $(or ${COMMONDIR},../builder)
CGO_ENABLED := 1

in-docker: clean-local-dirs generate-client gofmt test all;

include $(COMMONDIR)/Makefile.inc

release:: generate-client gofmt test all ;

.PHONY: clean
clean::
	rm -f ${INITRD} ${INITRD_COMPRESSED}

${INITRD_COMPRESSED}:
	rm -f ${INITRD_COMPRESSED}
	docker-make --no-push --Lint

.PHONY: initrd
initrd: ${INITRD_COMPRESSED}

# place all binaries in the same directory (/sbin) which is in the PATH of root.
.PHONY: ramdisk
ramdisk:
	GO111MODULE=off u-root \
		-format=cpio -build=bb \
		-defaultsh=/bin/bash \
		-files="bin/metal-hammer:bbin/uinit" \
		-files="/etc/localtime:etc/localtime" \
		-files="/bin/bash:bin/bash" \
		-files="/sbin/blkid:sbin/blkid" \
		-files="/sbin/ethtool:sbin/ethtool" \
		-files="/usr/bin/lspci:bin/lspci" \
		-files="/usr/bin/strace:bin/strace" \
		-files="/usr/share/misc/pci.ids:usr/share/misc/pci.ids" \
		-files="/bin/netstat:bin/netstat" \
		-files="/sbin/hdparm:sbin/hdparm" \
		-files="/usr/bin/ipmitool:usr/bin/ipmitool" \
		-files="/sbin/mkfs.vfat:sbin/mkfs.vfat" \
		-files="/sbin/mkfs.ext3:sbin/mkfs.ext3" \
		-files="/sbin/mkfs.ext4:sbin/mkfs.ext4" \
		-files="/sbin/mke2fs:sbin/mke2fs" \
		-files="/sbin/mkswap:sbin/mkswap" \
		-files="/sbin/mkfs.fat:sbin/mkfs.fat" \
		-files="/usr/sbin/lldptool:sbin/lldptool" \
		-files="/usr/sbin/nvme:sbin/nvme" \
		-files="/sbin/sgdisk:sbin/sgdisk" \
		-files="/etc/ssl/certs/ca-certificates.crt:etc/ssl/certs/ca-certificates.crt" \
		-files="/usr/lib/x86_64-linux-gnu/libnss_files.so:lib/libnss_files.so.2" \
		-files="passwd:etc/passwd" \
		-files="varrun:var/run/keep" \
		-files="ice.pkg:lib/firmware/intel/ice/ddp/ice.pkg" \
		-files="metal.key:id_rsa" \
		-files="metal.key.pub:authorized_keys" \
		-files="sum:sbin/sum" \
	-o ${INITRD} \
	&& ${COMPRESSOR} ${COMPRESSOR_ARGS} ${INITRD} ${INITRD_COMPRESSED} \
	&& rm -f ${INITRD}

clean-local-dirs:
	rm -rf metal-core
	mkdir metal-core

# 'swaggergenerate' generates swagger client with SWAGGERSPEC="swagger.json" SWAGGERTARGET="./".
generate-client: SWAGGERSPEC="metal-core.json"
generate-client: SWAGGERTARGET="metal-core"
generate-client: clean-local-dirs swaggergenerate

vagrant-destroy:
	vagrant destroy -f

vagrant-up: vagrant-destroy
	vagrant up && virsh console metal-hammer_pxeclient

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
		-kernel metal-hammer-kernel \
		-initrd metal-hammer-initrd.img.lz4
