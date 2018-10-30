FROM registry.fi-ts.io/cloud-native/go-builder:latest as builder

FROM golang:1.11-stretch as initrd-builder
ENV UROOT_GIT_SHA=edd248adfa09bfe392ba0f552f6574dbd37e1747
RUN apt-get update \
 && apt-get install -y \
	curl \
	dosfstools \
	e2fsprogs \
	gcc \
	gdisk \
	hdparm \
	liblz4-tool \
	nvme-cli
RUN mkdir -p ${GOPATH}/src/github.com/u-root \
 && cd ${GOPATH}/src/github.com/u-root \
 && git clone https://github.com/u-root/u-root \
 && cd u-root \
 && git checkout ${UROOT_GIT_SHA} \
 && go install
WORKDIR /work
COPY metal.key /work/
COPY metal.key.pub /work/
COPY Makefile /work/
COPY .git /work/
COPY --from=builder /common /common
COPY --from=builder /work/bin/metal-hammer /work/bin/
RUN COMMONDIR=/common make ramdisk

FROM scratch
COPY --from=builder /work/bin/metal-hammer /
COPY --from=initrd-builder /work/metal-hammer-initrd.img.lz4 /
