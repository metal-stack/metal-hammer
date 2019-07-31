module git.f-i-ts.de/cloud-native/metal/metal-hammer

require (
	github.com/beevik/ntp v0.2.0
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/fatih/color v1.7.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/runtime v0.19.2
	github.com/go-openapi/strfmt v0.19.0
	github.com/go-openapi/swag v0.19.2
	github.com/go-openapi/validate v0.19.2
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/gopacket v1.1.17
	github.com/google/uuid v1.1.1
	github.com/inconshreveable/log15 v0.0.0-20180818164646-67afb5ed74ec
	github.com/jaypipes/ghw v0.0.0-20190529200504-f4fa514ff57b
	github.com/jaypipes/pcidb v0.0.0-20190603110524-92fa8668ddf5 // indirect
	github.com/mailru/easyjson v0.0.0-20190620125010-da37f6c1e481 // indirect
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/mdlayher/ethernet v0.0.0-20190606142754-0394541c37b7
	github.com/mdlayher/lldp v0.0.0-20150915211757-afd9f83164c5
	github.com/mdlayher/raw v0.0.0-20190606144222-a54781e5f38f
	github.com/metal-pod/v v0.0.2
	// archiver must stay in version v2.1.0, see replace below
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/nwaples/rardecode v1.0.0 // indirect
	github.com/pierrec/lz4 v2.0.5+incompatible
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.3.0
	github.com/u-root/u-root v4.0.0+incompatible
	github.com/vishvananda/netlink v1.0.0
	github.com/vishvananda/netns v0.0.0-20180720170159-13995c7128cc // indirect
	golang.org/x/net v0.0.0-20190619014844-b5b0513f8c1b // indirect
	golang.org/x/sys v0.0.0-20190620070143-6f217b454f45
	gopkg.in/cheggaaa/pb.v1 v1.0.28
	gopkg.in/yaml.v2 v2.2.2
	howett.net/plist v0.0.0-20181124034731-591f970eefbb // indirect
)

replace github.com/mholt/archiver => github.com/mholt/archiver v2.1.0+incompatible
