FROM registry.fi-ts.io/cloud-native/go-builder:latest as builder

FROM golang:1.12-stretch as storcli-builder
# Raidcontroller configuration cli storcli
# check here for latestes releases, there is no public debian repo unfortunately
# given download directory can even not listed.
# https://www.broadcom.com/support/download-search/?pg=&pf=&pn=&pa=&po=&dk=storcli
# FIRMWARE:
# https://www.supermicro.com/wftp/driver/SAS/LSI/3108/Firmware/
ENV STORCLI_VERSION=7.8-007.0813.0000.0000 \
    STORCLI_DOWNLOAD_URL=https://docs.broadcom.com/docs-and-downloads/raid-controllers/raid-controllers-common-files
WORKDIR /work
RUN set -ex \
 && wget -q https://github.com/Microsoft/ethr/releases/download/v0.2.1/ethr_linux.zip -O ethr.zip \
 && wget -q ${STORCLI_DOWNLOAD_URL}/MR_SAS_Unified_StorCLI_${STORCLI_VERSION}.zip -O storcli.zip \
 && apt-get update \
 && apt-get install -y --no-install-recommends unzip \
 && unzip storcli.zip \
 && unzip ethr.zip \
 && dpkg -i Unified_storcli_all_os/Ubuntu/storcli*.deb

FROM golang:1.12-stretch as initrd-builder
ENV UROOT_GIT_SHA_OR_TAG=v4.0.0
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
	mstflint \
	net-tools \
	nvme-cli \
	pciutils \
	util-linux
RUN mkdir -p ${GOPATH}/src/github.com/u-root \
 && cd ${GOPATH}/src/github.com/u-root \
 && git clone https://github.com/u-root/u-root \
 && cd u-root \
 && git checkout ${UROOT_GIT_SHA_OR_TAG} \
 && go install
WORKDIR /work
COPY metal.key /work/
COPY metal.key.pub /work/
COPY Makefile /work/
COPY .git /work/
COPY --from=storcli-builder /opt/MegaRAID/storcli/storcli64 /work/bin/
COPY --from=storcli-builder /work/ethr /work/bin/
COPY --from=builder /common /common
COPY --from=builder /work/bin/metal-hammer /work/bin/
RUN COMMONDIR=/common make ramdisk

FROM scratch
COPY --from=builder /work/bin/metal-hammer /
COPY --from=initrd-builder /work/metal-hammer-initrd.img.lz4 /
