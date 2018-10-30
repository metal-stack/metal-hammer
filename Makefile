INITRD := metal-hammer-initrd.img.lz4
BINARY := metal-hammer
MAINMODULE := .
COMMONDIR := $(or ${COMMONDIR},../../common)

include $(COMMONDIR)/Makefile.inc

.PHONY: clean
clean::
	rm ${INITRD}

${INITRD}:
	rm -f ${INITRD}
	docker-make --no-push --Lint

.PHONY: initrd
initrd: ${INITRD}

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
	-o metal-hammer-initrd.img \
	&& lz4 -f -l metal-hammer-initrd.img metal-hammer-initrd.img.lz4 \
	&& rm -f metal-hammer-initrd.img
