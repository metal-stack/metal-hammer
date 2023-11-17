# FIXME this points to the last go-1.20 image of the builder, 
# go-1.21 cant be used actually because bookworm mangled libsystemd-dev which breaks the build
# and go-1.21 is not available with bullseye.
# maybe we should switch away from depending on the builder image
FROM metalstack/builder@sha256:d2050a3bef9bbd9d9ea769a71a4a70b9ff4b24c537d29d5870b83fc652bb67f8 as builder
# Install Intel Firmware for e800 based network cards
ENV ICE_VERSION=1.9.11
ENV ICE_PKG_VERSION=1.3.30.0
RUN curl -fLsS https://sourceforge.net/projects/e1000/files/ice%20stable/${ICE_VERSION}/ice-${ICE_VERSION}.tar.gz/download -o ice.tar.gz \
 && tar -xf ice.tar.gz ice-${ICE_VERSION}/ddp/ice-${ICE_PKG_VERSION}.pkg \
 && mkdir -p /lib/firmware/intel/ice/ddp/ \
 && mv ice-${ICE_VERSION}/ddp/ice-${ICE_PKG_VERSION}.pkg /work/ice.pkg

# ipmitool from bookworm is broken and returns with error on most commands
FROM golang:1.20-bullseye as initrd-builder
ENV UROOT_GIT_SHA_OR_TAG=v0.11.0
RUN apt-get update \
 && apt-get install -y --no-install-recommends \
	ca-certificates \
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
	util-linux \
 # this is required, otherwise uroot complains that these files already exist
 && rm -f /etc/passwd /etc/lvm/lvmlocal.conf
RUN mkdir -p ${GOPATH}/src/github.com/u-root \
 && cd ${GOPATH}/src/github.com/u-root \
 && git clone https://github.com/u-root/u-root \
 && cd u-root \
 && git checkout ${UROOT_GIT_SHA_OR_TAG} \
 && GO111MODULE=off go install
WORKDIR /work
RUN mkdir -p /work/etc/lvm /work/etc/ssl/certs /work/lib/firmware/intel/ice/ddp/ /work/var/run
COPY enterprise-numbers.txt lvmlocal.conf metal.key metal.key.pub passwd varrun Makefile .git /work/
COPY --from=r.metal-stack.io/metal/supermicro:2.12.0 /usr/bin/sum /work/
COPY --from=builder /common /common
COPY --from=builder /work/ice.pkg /work/ice.pkg
COPY --from=builder /work/bin/metal-hammer /work/bin/
RUN COMMONDIR=/common make ramdisk

FROM scratch
COPY --from=builder /work/bin/metal-hammer /
COPY --from=initrd-builder /work/metal-hammer-initrd.img.lz4 /
