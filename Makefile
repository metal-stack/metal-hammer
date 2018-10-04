.ONESHELL:
SHA := $(shell git rev-parse --short=8 HEAD)
GITVERSION := $(shell git describe --long --all)
BUILDDATE := $(shell date -Iseconds)
VERSION := $(or ${VERSION},devel)

BINARY := metal-hammer

all: $(BINARY);

%:
	CGO_ENABLE=0 GO111MODULE=on go build -tags netgo -ldflags "-linkmode external -extldflags -static -X 'main.version=$(VERSION)' -X 'main.revision=$(GITVERSION)' -X 'main.gitsha1=$(SHA)' -X 'main.builddate=$(BUILDDATE)'" -o bin/$@

image: docker build -t registry.fi-ts.io/maas/metal-hammer .
