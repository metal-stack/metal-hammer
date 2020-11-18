module github.com/metal-stack/metal-hammer

go 1.15

require (
	github.com/beevik/ntp v0.3.0
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/fatih/color v1.10.0 // indirect
	github.com/frankban/quicktest v1.11.2 // indirect
	github.com/go-openapi/errors v0.19.8
	github.com/go-openapi/runtime v0.19.24
	github.com/go-openapi/strfmt v0.19.11
	github.com/go-openapi/swag v0.19.12
	github.com/go-openapi/validate v0.19.14
	github.com/google/gopacket v1.1.19
	github.com/google/uuid v1.1.2
	github.com/inconshreveable/log15 v0.0.0-20201112154412-8562bdadbbac
	github.com/jaypipes/ghw v0.6.1
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mdlayher/ethernet v0.0.0-20190606142754-0394541c37b7
	github.com/mdlayher/lldp v0.0.0-20150915211757-afd9f83164c5
	github.com/mdlayher/raw v0.0.0-20191009151244-50f2db8cc065
	github.com/metal-stack/go-hal v0.3.0
	github.com/metal-stack/metal-api v0.11.1
	github.com/metal-stack/metal-lib v0.6.6
	github.com/metal-stack/v v1.0.2
	// archiver must stay in version v2.1.0, see replace below
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/pierrec/lz4 v2.5.2+incompatible
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	github.com/u-root/u-root v7.0.0+incompatible
	github.com/vishvananda/netlink v1.1.0
	github.com/vishvananda/netns v0.0.0-20200728191858-db3c7e526aae // indirect
	golang.org/x/sys v0.0.0-20201117222635-ba5294a509c7
	google.golang.org/grpc v1.33.2
	gopkg.in/cheggaaa/pb.v1 v1.0.28
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/mholt/archiver => github.com/mholt/archiver v2.1.0+incompatible
