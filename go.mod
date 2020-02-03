module git.f-i-ts.de/cloud-native/metal/metal-hammer

require (
	github.com/avast/retry-go v2.4.3+incompatible
	github.com/beevik/ntp v0.2.0
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/fatih/color v1.7.0 // indirect
	github.com/frankban/quicktest v1.5.0 // indirect
	github.com/go-openapi/analysis v0.19.7 // indirect
	github.com/go-openapi/errors v0.19.3
	github.com/go-openapi/runtime v0.19.9
	github.com/go-openapi/spec v0.19.5 // indirect
	github.com/go-openapi/strfmt v0.19.4
	github.com/go-openapi/swag v0.19.6
	github.com/go-openapi/validate v0.19.5
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/gopacket v1.1.17
	github.com/google/uuid v1.1.1
	github.com/inconshreveable/log15 v0.0.0-20180818164646-67afb5ed74ec
	github.com/jaypipes/ghw v0.5.0
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.11 // indirect
	github.com/mattn/go-runewidth v0.0.7 // indirect
	github.com/mdlayher/ethernet v0.0.0-20190606142754-0394541c37b7
	github.com/mdlayher/lldp v0.0.0-20150915211757-afd9f83164c5
	github.com/mdlayher/raw v0.0.0-20191009151244-50f2db8cc065
	github.com/metal-pod/v v0.0.2
	// archiver must stay in version v2.1.0, see replace below
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/nwaples/rardecode v1.0.0 // indirect
	github.com/pierrec/lz4 v2.4.0+incompatible
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.4.0
	github.com/u-root/u-root v6.0.0+incompatible
	github.com/vishvananda/netlink v1.0.0
	github.com/vishvananda/netns v0.0.0-20191106174202-0a2b9b5464df // indirect
	go.mongodb.org/mongo-driver v1.2.0 // indirect
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553
	golang.org/x/sys v0.0.0-20191218084908-4a24b4065292
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/cheggaaa/pb.v1 v1.0.28
	gopkg.in/yaml.v2 v2.2.7
)

replace github.com/mholt/archiver => github.com/mholt/archiver v2.1.0+incompatible

go 1.13
