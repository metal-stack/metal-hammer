BINARY := metal-hammer
INITRD := ${BINARY}-initrd.img
COMPRESSOR := lz4
COMPRESSOR_ARGS := -f -l
INITRD_COMPRESSED := ${INITRD}.${COMPRESSOR}
MAINMODULE := .
COMMONDIR := $(or ${COMMONDIR},../common)
CGO_ENABLED := 1

in-docker: clean-client generate-client test all;

include $(COMMONDIR)/Makefile.inc

release:: generate-client test all ;

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
	u-root \
		-format=cpio -build=bb \
		-files="bin/metal-hammer:bbin/uinit" \
		-files="/bin/bash:bin/bash" \
		-files="/sbin/blkid:sbin/blkid" \
		-files="/sbin/ethtool:sbin/ethtool" \
		-files="/bin/netstat:bin/netstat" \
		-files="/sbin/hdparm:sbin/hdparm" \
		-files="/usr/bin/ipmitool:usr/bin/ipmitool" \
		-files="/sbin/mkfs.vfat:sbin/mkfs.vfat" \
		-files="/sbin/mkfs.ext4:sbin/mkfs.ext4" \
		-files="/sbin/mke2fs:sbin/mke2fs" \
		-files="/sbin/mkfs.fat:sbin/mkfs.fat" \
		-files="/usr/sbin/nvme:sbin/nvme" \
		-files="/sbin/sgdisk:sbin/sgdisk" \
		-files="bin/storcli64:sbin/storcli" \
		-files="bin/smc3108.rom:bin/smc3108.rom" \
		-files="/etc/ssl/certs/ca-certificates.crt:etc/ssl/certs/ca-certificates.crt" \
		-files="metal.key:id_rsa" \
		-files="metal.key.pub:authorized_keys" \
	-o ${INITRD} \
	&& ${COMPRESSOR} ${COMPRESSOR_ARGS} ${INITRD} ${INITRD_COMPRESSED} \
	&& rm -f ${INITRD}

clean-local-dirs:
	rm -rf metal-core
	mkdir metal-core
clean-client: clean-local-dirs
	cp ../metal-core/spec/metal-core.json metal-core.json

generate-client: clean-local-dirs
	GO111MODULE=off swagger generate client -f metal-core.json --skip-validation --target metal-core

vagrant-destroy:
	vagrant destroy -f

vagrant-up: vagrant-destroy
	vagrant up && virsh console metal-hammer_pxeclient
