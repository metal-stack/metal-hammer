FROM registry.fi-ts.io/cloud-native/go-builder:latest as builder

FROM golang:1.12-stretch as storcli-builder
# Check here for latestes releases, there is no public deb repo unfortunately
# given download directory can even not listed.
# https://www.broadcom.com/support/download-search/?pg=&pf=&pn=&pa=&po=&dk=storcli
# FIRMWARE:
# https://www.supermicro.com/wftp/driver/SAS/LSI/3108/Firmware/
ENV STORCLI_VERSION=7.8-007.0813.0000.0000 \
    STORCLI_DOWNLOAD_URL=https://docs.broadcom.com/docs-and-downloads/raid-controllers/raid-controllers-common-files \
	RAID_FIRMWARE_VERSION=4.680.00-8290 \
	RAID_FIRMWARE_DOWNLOAD_URL=https://www.supermicro.com/wftp/driver/SAS/LSI/3108/Firmware/
WORKDIR /work
RUN set -ex \
 && wget -q ${STORCLI_DOWNLOAD_URL}/MR_SAS_Unified_StorCLI_${STORCLI_VERSION}.zip -O storcli.zip \
 && wget -q ${RAID_FIRMWARE_DOWNLOAD_URL}/${RAID_FIRMWARE_VERSION}.zip -O firmware.zip \
 && apt-get update \
 && apt-get install -y --no-install-recommends unzip \
 && unzip storcli.zip \
 && unzip firmware.zip \
 && dpkg -i Unified_storcli_all_os/Ubuntu/storcli*.deb

FROM golang:1.12-stretch as initrd-builder
ENV UROOT_GIT_SHA_OR_TAG=v4.0.0
RUN apt-get update \
 && apt-get install -y --no-install-recommends \
	curl \
	dosfstools \
	e2fsprogs \
	gcc \
	gdisk \
	hdparm \
	ipmitool \
	liblz4-tool \
	nvme-cli \
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
# firmware update via
# storcli /cX download file=smc3108.rom
COPY --from=storcli-builder /work/smc3108.rom /work/bin/
COPY --from=builder /common /common
COPY --from=builder /work/bin/metal-hammer /work/bin/
RUN COMMONDIR=/common make ramdisk

FROM scratch
COPY --from=builder /work/bin/metal-hammer /
COPY --from=initrd-builder /work/metal-hammer-initrd.img.lz4 /
