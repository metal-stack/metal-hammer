module github.com/metal-stack/metal-hammer

go 1.15

require (
	github.com/beevik/ntp v0.3.0
	github.com/cheggaaa/pb/v3 v3.0.8
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/fatih/color v1.12.0 // indirect
	github.com/frankban/quicktest v1.13.0 // indirect
	github.com/go-openapi/errors v0.20.0
	github.com/go-openapi/runtime v0.19.29
	github.com/go-openapi/strfmt v0.20.1
	github.com/go-openapi/swag v0.19.15
	github.com/go-openapi/validate v0.20.2
	github.com/google/gopacket v1.1.19
	github.com/google/uuid v1.2.0
	github.com/inconshreveable/log15 v0.0.0-20201112154412-8562bdadbbac
	github.com/jaypipes/ghw v0.8.0
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/metal-stack/go-hal v0.3.3
	github.com/metal-stack/go-lldpd v0.3.1
	github.com/metal-stack/metal-api v0.15.0
	github.com/metal-stack/metal-lib v0.8.0
	github.com/metal-stack/v v1.0.3
	// archiver must stay in version v2.1.0, see replace below
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/u-root/u-root v7.0.0+incompatible
	github.com/vishvananda/netlink v1.1.0
	github.com/vishvananda/netns v0.0.0-20210104183010-2eb08e3e575f // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210611083646-a4fc73990273
	google.golang.org/grpc v1.38.0
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/mholt/archiver => github.com/mholt/archiver v2.1.0+incompatible
