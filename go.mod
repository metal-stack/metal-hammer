module github.com/metal-stack/metal-hammer

go 1.15

require (
	github.com/beevik/ntp v0.3.0
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/fatih/color v1.10.0 // indirect
	github.com/frankban/quicktest v1.11.3 // indirect
	github.com/go-openapi/errors v0.20.0
	github.com/go-openapi/runtime v0.19.26
	github.com/go-openapi/strfmt v0.20.0
	github.com/go-openapi/swag v0.19.14
	github.com/go-openapi/validate v0.20.2
	github.com/google/gopacket v1.1.19
	github.com/google/uuid v1.2.0
	github.com/inconshreveable/log15 v0.0.0-20201112154412-8562bdadbbac
	github.com/jaypipes/ghw v0.7.0
	github.com/mattn/go-runewidth v0.0.10 // indirect
	github.com/mdlayher/ethernet v0.0.0-20190606142754-0394541c37b7
	github.com/mdlayher/lldp v0.0.0-20150915211757-afd9f83164c5
	github.com/mdlayher/raw v0.0.0-20191009151244-50f2db8cc065
	github.com/metal-stack/go-hal v0.3.0
	github.com/metal-stack/metal-api v0.13.0
	github.com/metal-stack/metal-lib v0.7.0
	github.com/metal-stack/v v1.0.3
	// archiver must stay in version v2.1.0, see replace below
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/pierrec/lz4 v2.6.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	github.com/u-root/u-root v7.0.0+incompatible
	github.com/vishvananda/netlink v1.1.0
	github.com/vishvananda/netns v0.0.0-20210104183010-2eb08e3e575f // indirect
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/sys v0.0.0-20210319071255-635bc2c9138d
	google.golang.org/grpc v1.36.0
	gopkg.in/cheggaaa/pb.v1 v1.0.28
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/mholt/archiver => github.com/mholt/archiver v2.1.0+incompatible
