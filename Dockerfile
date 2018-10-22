FROM golang:1.11-stretch as metal-hammer-builder
RUN apt-get update \
 && apt-get install -y make git
WORKDIR /work
COPY .git /work/
COPY go.mod /work/
COPY go.sum /work/
COPY main.go /work/
COPY cmd /work/cmd/
COPY pkg /work/pkg/
COPY Makefile /work/
RUN make bin/metal-hammer

FROM golang:1.11-stretch as initrd-builder
ENV UROOT_GIT_SHA=5909da7ef93be40da573f61005189e5270078bb7
RUN apt-get update \
 && apt-get install -y \
	curl \
	dosfstools \
	e2fsprogs \
	gcc \
	gdisk \
	hdparm \
	liblz4-tool \
	nvme-cli \
	rng-tools
RUN mkdir -p ${GOPATH}/src/github.com/u-root \
 && cd ${GOPATH}/src/github.com/u-root \
 && git clone https://github.com/u-root/u-root \
 && cd u-root \
 && git checkout ${UROOT_GIT_SHA} \
 && go install
WORKDIR /work
COPY metal.key /work/
COPY metal.key.pub /work/
COPY metal-hammer.sh /work/
COPY Makefile /work/
COPY --from=metal-hammer-builder /work/bin/metal-hammer /work/bin/
RUN make ramdisk

FROM scratch
COPY --from=metal-hammer-builder /work/bin/metal-hammer /
COPY --from=initrd-builder /work/metal-hammer-initrd.img.lz4 /
