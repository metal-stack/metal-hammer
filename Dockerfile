FROM metalstack/builder:latest as builder

FROM r.metal-stack.io/metal/supermicro:2.5.2 as sum

FROM golang:1.14-buster as initrd-builder
ENV UROOT_GIT_SHA_OR_TAG=v7.0.0
RUN apt-get update \
 && apt-get install -y --no-install-recommends \
	curl \
	dosfstools \
	e2fsprogs \
	ethtool \
	gcc \
	gdisk \
	hdparm \
	ipmitool \
	liblz4-tool \
	lvm2 \
	mdadm \
	net-tools \
	nvme-cli \
	pciutils \
	strace \
	util-linux
RUN mkdir -p ${GOPATH}/src/github.com/u-root \
 && cd ${GOPATH}/src/github.com/u-root \
 && git clone https://github.com/u-root/u-root \
 && cd u-root \
 && git checkout ${UROOT_GIT_SHA_OR_TAG} \
 && go install
WORKDIR /work
COPY lvmlocal.conf ice.pkg metal.key metal.key.pub passwd varrun Makefile .git /work/
COPY --from=sum /usr/bin/sum /work/
COPY --from=builder /common /common
COPY --from=builder /work/bin/metal-hammer /work/bin/
RUN COMMONDIR=/common make ramdisk

FROM scratch
COPY --from=builder /work/bin/metal-hammer /
COPY --from=initrd-builder /work/metal-hammer-initrd.img.lz4 /
