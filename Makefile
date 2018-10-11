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