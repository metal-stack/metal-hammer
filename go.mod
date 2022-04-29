module github.com/metal-stack/metal-hammer

go 1.16

require (
	github.com/beevik/ntp v0.3.0
	github.com/cheggaaa/pb/v3 v3.0.8
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/frankban/quicktest v1.14.3 // indirect
	github.com/go-openapi/errors v0.20.2
	github.com/go-openapi/runtime v0.23.3
	github.com/go-openapi/strfmt v0.21.2
	github.com/go-openapi/swag v0.21.1
	github.com/go-openapi/validate v0.21.0
	github.com/google/gopacket v1.1.19
	github.com/google/uuid v1.3.0
	github.com/jaypipes/ghw v0.9.0
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/metal-stack/go-hal v0.4.0
	github.com/metal-stack/go-lldpd v0.4.0
	github.com/metal-stack/metal-api v0.16.7-0.20220429103547-d2c35c9a011f
	github.com/metal-stack/metal-lib v0.9.0
	github.com/metal-stack/v v1.0.3
	// archiver must stay in version v2.1.0, see replace below
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/nwaples/rardecode v1.1.2 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.14
	github.com/stretchr/testify v1.7.1
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/u-root/u-root v0.8.0
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/vishvananda/netlink v1.1.1-0.20211118161826-650dca95af54
	go.uber.org/zap v1.21.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20220422013727-9388b58f7150
	google.golang.org/grpc v1.46.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/mholt/archiver => github.com/mholt/archiver v2.1.0+incompatible
