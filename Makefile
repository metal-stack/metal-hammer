.ONESHELL:
SHA := $(shell git rev-parse --short=8 HEAD)
GITVERSION := $(shell git describe --long --all)
BUILDDATE := $(shell date -Iseconds)
VERSION := $(or ${VERSION},devel)

BINARY := bin/metal-hammer
INITRD := metal-hammer-initrd.img.gz

.PHONY: clean initrd

all: $(BINARY)

${BINARY}: clean
	CGO_ENABLE=0 \
	GO111MODULE=on \
	go test -v -race -cover $(shell go list ./...) && \
	go build \
		-tags netgo \
		-ldflags "-X 'main.version=$(VERSION)' \
				  -X 'main.revision=$(GITVERSION)' \
				  -X 'main.gitsha1=$(SHA)' \
				  -X 'main.builddate=$(BUILDDATE)'" \
	-o $@

clean:
	rm -f ${BINARY} ${INITRD}

${INITRD}:
	rm -f ${INITRD}
	docker-make --no-push --Lint

initrd: ${INITRD}

fast: clean ${BINARY}
	u-root \
		-format=cpio -build=bb \
		-files="bin/metal-hammer:bbin/metal-hammer" \
		-files="/sbin/sgdisk:usr/bin/sgdisk" \
		-files="/sbin/mkfs.vfat:sbin/mkfs.vfat" \
		-files="/sbin/mkfs.ext4:sbin/mkfs.ext4" \
		-files="/sbin/mke2fs:sbin/mke2fs" \
		-files="/sbin/mkfs.fat:sbin/mkfs.fat" \
		-files="/usr/sbin/rngd:usr/sbin/rngd" \
		-files="/etc/ssl/certs/ca-certificates.crt:etc/ssl/certs/ca-certificates.crt" \
		-files="metal.key:id_rsa" \
		-files="metal.key.pub:authorized_keys" \
		-files="metal-hammer.sh:bbin/uinit" \
	-o metal-hammer-initrd.img \
	&& gzip -f metal-hammer-initrd.img \
	&& rm -f metal-hammer-initrd.img
