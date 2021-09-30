module github.com/metal-stack/metal-hammer

go 1.16

require (
	github.com/beevik/ntp v0.3.0
	github.com/cheggaaa/pb/v3 v3.0.8
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/frankban/quicktest v1.13.1 // indirect
	github.com/go-openapi/errors v0.20.1
	github.com/go-openapi/runtime v0.19.31
	github.com/go-openapi/strfmt v0.20.2
	github.com/go-openapi/swag v0.19.15
	github.com/go-openapi/validate v0.20.2
	github.com/google/gopacket v1.1.19
	github.com/google/uuid v1.3.0
	github.com/inconshreveable/log15 v0.0.0-20201112154412-8562bdadbbac
	github.com/jaypipes/ghw v0.8.0
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/metal-stack/go-hal v0.3.6
	github.com/metal-stack/go-lldpd v0.3.5
	github.com/metal-stack/metal-api v0.15.7
	github.com/metal-stack/metal-lib v0.8.2
	github.com/metal-stack/v v1.0.3
	// archiver must stay in version v2.1.0, see replace below
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/nwaples/rardecode v1.1.2 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/pretty v1.2.0 // indirect
	// keep u-root sha in sync with Dockerfile
	github.com/u-root/u-root v0.0.0-20210920205541-c0a6cbaae564
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/vishvananda/netlink v1.1.1-0.20200221165523-c79a4b7b4066
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6
	google.golang.org/grpc v1.41.0
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/mholt/archiver => github.com/mholt/archiver v2.1.0+incompatible
