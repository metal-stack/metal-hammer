FROM metalstack/builder:latest as builder
# Install Intel Firmware for e800 based network cards
ENV ICE_VERSION=1.8.8
ENV ICE_PKG_VERSION=1.3.28.0
RUN curl -fLsS https://sourceforge.net/projects/e1000/files/ice%20stable/${ICE_VERSION}/ice-${ICE_VERSION}.tar.gz/download -o ice.tar.gz \
 && tar -xf ice.tar.gz ice-${ICE_VERSION}/ddp/ice-${ICE_PKG_VERSION}.pkg \
 && mkdir -p /lib/firmware/intel/ice/ddp/ \
 && mv ice-${ICE_VERSION}/ddp/ice-${ICE_PKG_VERSION}.pkg /work/ice.pkg

FROM r.metal-stack.io/metal/supermicro:2.5.2 as sum

FROM golang:1.14-buster as initrd-builder
ENV UROOT_GIT_SHA_OR_TAG=v0.7.0
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
 && GO111MODULE=off go install
WORKDIR /work
COPY lvmlocal.conf metal.key metal.key.pub passwd varrun Makefile .git /work/
COPY --from=sum /usr/bin/sum /work/
COPY --from=builder /common /common
COPY --from=builder /work/ice.pkg /work/ice.pkg
COPY --from=builder /work/bin/metal-hammer /work/bin/
RUN COMMONDIR=/common make ramdisk

FROM scratch
COPY --from=builder /work/bin/metal-hammer /
COPY --from=initrd-builder /work/metal-hammer-initrd.img.lz4 /
