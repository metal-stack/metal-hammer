FROM golang:1.26-alpine AS builder

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
ENV ICE_VERSION=2.4.5
ENV ICE_PKG_VERSION=1.3.53.0
RUN curl -fLsS https://github.com/intel/ethernet-linux-ice/releases/download/v${ICE_VERSION}/ice-${ICE_VERSION}.tar.gz -o ice.tar.gz \
 && tar -xf ice.tar.gz ice-${ICE_VERSION}/ddp/ice-${ICE_PKG_VERSION}.pkg \
 && mkdir -p /lib/firmware/intel/ice/ddp/ \
 && mv ice-${ICE_VERSION}/ddp/ice-${ICE_PKG_VERSION}.pkg /work/ice.pkg

# Install the Intel NVM Update Utility to update the firmware of e810 based network cards.
# If E810_NVM_URL is left empty, the utility is not part of the initrd and metal-hammer
# skips the nic firmware update at runtime.

ARG E810_NVM_URL=
ARG E810_NVM_SHA256=

RUN mkdir -p /work/intel \
 && if [ -z "${E810_NVM_URL}" ]; then \
      echo "E810_NVM_URL is not set, the intel nic firmware update will be disabled" ; \
    else \
      version=$(basename "${E810_NVM_URL}" | sed -n 's/^E810_NVMUpdatePackage_v\([0-9]\{1,\}\)_\([0-9]\{1,\}\)_Linux\.tar\.gz$/\1.\2/p') \
   && if [ -z "${version}" ]; then echo "unable to determine the nvm version from ${E810_NVM_URL}" && exit 1 ; fi \
   && if [ -z "${E810_NVM_SHA256}" ]; then echo "E810_NVM_SHA256 is required if E810_NVM_URL is set" && exit 1 ; fi \
   && curl -fLsS --proto '=https' "${E810_NVM_URL}" -o e810-nvm.tar.gz \
   && echo "$(echo "${E810_NVM_SHA256}" | tr 'A-Z' 'a-z')  e810-nvm.tar.gz" | sha256sum -c - \
   && tar -xf e810-nvm.tar.gz -C /work/intel --strip-components=2 E810/Linux_x64 \
   && chmod +x /work/intel/nvmupdate64e \
   && echo "${version}" > /work/intel/version \
   && rm e810-nvm.tar.gz ; \
    fi

# ipmitool from bookworm is broken and returns with error on most commands, seems fixed
# sgdisk from debian:13 is broken and creates a corrupt GPT partition layout
FROM golang:1.24-bookworm AS initrd-builder
ENV UROOT_GIT_SHA_OR_TAG=v0.14.0
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
	lvm2 \
	lz4 \
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
 && go install
WORKDIR /work
RUN mkdir -p /work/etc/lvm /work/etc/ssl/certs /work/lib/firmware/intel/ice/ddp/ /work/var/run \
 && cp /usr/share/zoneinfo/Etc/UTC /work/etc/localtime
COPY lvmlocal.conf metal.key metal.key.pub passwd varrun Makefile .git /work/
COPY --from=r.metal-stack.io/metal/supermicro:2.14.0 /usr/bin/sum /work/
COPY --from=builder /work/ice.pkg /work/ice.pkg
COPY --from=builder /work/intel /work/intel
COPY --from=builder /work/bin/metal-hammer /work/bin/
RUN make ramdisk

FROM scratch
COPY --from=builder /work/bin/metal-hammer /
COPY --from=initrd-builder /work/metal-hammer-initrd.img.lz4 /
