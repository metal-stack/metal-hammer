BINARY := metal-hammer
INITRD := ${BINARY}-initrd.img
COMPRESSOR := lz4
COMPRESSOR_ARGS := -f -l
INITRD_COMPRESSED := ${INITRD}.${COMPRESSOR}
MAINMODULE := .
COMMONDIR := $(or ${COMMONDIR},../../common)

in-docker: genrate-client all;

include $(COMMONDIR)/Makefile.inc

.PHONY: clean
clean::
	rm -f ${INITRD} ${INITRD_COMPRESSED}

${INITRD_COMPRESSED}:
	rm -f ${INITRD_COMPRESSED}
	docker-make --no-push --Lint

.PHONY: initrd
initrd: ${INITRD_COMPRESSED}

.PHONY: ramdisk
ramdisk:
	u-root \
		-format=cpio -build=bb \
		-files="bin/metal-hammer:bbin/uinit" \
		-files="/sbin/sgdisk:usr/bin/sgdisk" \
		-files="/sbin/mkfs.vfat:sbin/mkfs.vfat" \
		-files="/sbin/mkfs.ext4:sbin/mkfs.ext4" \
		-files="/sbin/mke2fs:sbin/mke2fs" \
		-files="/sbin/mkfs.fat:sbin/mkfs.fat" \
		-files="/sbin/hdparm:sbin/hdparm" \
		-files="/usr/sbin/nvme:usr/sbin/nvme" \
		-files="/etc/ssl/certs/ca-certificates.crt:etc/ssl/certs/ca-certificates.crt" \
		-files="metal.key:id_rsa" \
		-files="metal.key.pub:authorized_keys" \
	-o ${INITRD} \
	&& ${COMPRESSOR} ${COMPRESSOR_ARGS} ${INITRD} ${INITRD_COMPRESSED} \
	&& rm -f ${INITRD}

generate-client:
	rm -rf metal-core \
	&& mkdir metal-core \
	&& GO111MODULE=off swagger generate client -f metal-core.json --skip-validation --target metal-core
