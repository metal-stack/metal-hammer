FROM golang:1.22-alpine as builder

RUN apk add \
	binutils \
	coreutils \
	curl \
	gcc \
	git \
	make \
	musl-dev \
	libpcap-dev
WORKDIR /work
COPY . .
RUN make all
# Install Intel Firmware for e800 based network cards
ENV ICE_VERSION=1.13.7
ENV ICE_PKG_VERSION=1.3.35.0
RUN curl -fLsS https://sourceforge.net/projects/e1000/files/ice%20stable/${ICE_VERSION}/ice-${ICE_VERSION}.tar.gz/download -o ice.tar.gz \
 && tar -xf ice.tar.gz ice-${ICE_VERSION}/ddp/ice-${ICE_PKG_VERSION}.pkg \
 && mkdir -p /lib/firmware/intel/ice/ddp/ \
 && mv ice-${ICE_VERSION}/ddp/ice-${ICE_PKG_VERSION}.pkg /work/ice.pkg

# ipmitool from bookworm is broken and returns with error on most commands
FROM golang:1.22-bullseye as initrd-builder
ENV UROOT_GIT_SHA_OR_TAG=update-pci-ids
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
 && git clone https://github.com/majst01/u-root \
 && cd u-root \
 && git checkout ${UROOT_GIT_SHA_OR_TAG} \
 && go install
WORKDIR /work
RUN mkdir -p /work/etc/lvm /work/etc/ssl/certs /work/lib/firmware/intel/ice/ddp/ /work/var/run \
 && cp /usr/share/zoneinfo/Etc/UTC /work/etc/localtime
COPY lvmlocal.conf metal.key metal.key.pub passwd varrun Makefile .git /work/
COPY --from=r.metal-stack.io/metal/supermicro:2.13.0 /usr/bin/sum /work/
COPY --from=builder /work/ice.pkg /work/ice.pkg
COPY --from=builder /work/bin/metal-hammer /work/bin/
RUN make ramdisk

FROM scratch
COPY --from=builder /work/bin/metal-hammer /
COPY --from=initrd-builder /work/metal-hammer-initrd.img.lz4 /
