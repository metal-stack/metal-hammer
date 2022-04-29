package cmd

import (
	"strconv"
	"strings"

	"os"

	"github.com/metal-stack/metal-hammer/pkg/kernel"
	"go.uber.org/zap"
)

//Specification defines configuration items of the application
type Specification struct {
	// Debug turn on debug log
	Debug bool
	// MetalCoreURL is the endpoint URL where the metalcore reside
	MetalCoreURL string
	// ImageURL if given grabs a fixed OS image to install, only suitable in DevMode
	ImageURL string
	// ImageID if given defines the image.ID which normally comes from a allocation
	// can be something like ubuntu-18.04, alpine-3.9 or "default"
	// only suitable in DevMode
	ImageID string
	// SizeID if given defines the size.ID which normally comes from a allocation
	// can be something like v1-small-x86
	// only suitable in DevMode
	SizeID string
	// DevMode turn on devmode which prevents failing in some situations
	DevMode bool
	// BGPEnabled if set to true real bgp configuration is configured, otherwise dhcp will be used
	BGPEnabled bool
	// Cidr of BGP interface in DEV Mode
	Cidr string
	// ConsolePassword of the metal user valid for one day.
	ConsolePassword string
	// MachineUUID is the unique identifier of this machine
	MachineUUID string
	// IP of this instance
	IP string

	log *zap.SugaredLogger
}

// NewSpec fills Specification with configuration made by kernel commandline
func NewSpec(log *zap.SugaredLogger) *Specification {
	spec := &Specification{}
	// Grab metal-hammer configuration from kernel commandline
	envmap, err := kernel.ParseCmdline()
	if err != nil {
		log.Errorw("parse cmdline", "error", err)
		os.Exit(1)
	}

	if d, ok := envmap["DEBUG"]; ok && (d == "1" || strings.ToLower(d) == "true") {
		spec.Debug = true
		os.Setenv("DEBUG", "1")
	}

	// METAL_CORE_URL must be in the form http://metal-core:4242
	if url, ok := envmap["METAL_CORE_ADDRESS"]; ok {
		spec.MetalCoreURL = url
	}

	if i, ok := envmap["IMAGE_URL"]; ok {
		spec.ImageURL = i
		spec.DevMode = true
	}

	if i, ok := envmap["IMAGE_ID"]; ok {
		spec.ImageID = i
		spec.DevMode = true
	}

	if s, ok := envmap["SIZE_ID"]; ok {
		spec.SizeID = s
		spec.DevMode = true
	}

	if c, ok := envmap["CIDR"]; ok {
		spec.Cidr = c
		spec.DevMode = true
	}

	if bgp, ok := envmap["BGP"]; ok {
		enabled, err := strconv.ParseBool(bgp)
		if err == nil {
			spec.BGPEnabled = enabled
		}
	}
	spec.log = log

	return spec
}

// Log print configuration options
func (s *Specification) Log() {
	s.log.Infow("configuration",
		"debug", s.Debug,
		"metalCoreURL", s.MetalCoreURL,
		"imageURL", s.ImageURL,
		"imageID", s.ImageID,
		"sizeID", s.SizeID,
		"devmode", s.DevMode,
		"bgpenabled", s.BGPEnabled,
		"cidr", s.Cidr,
		"machineUUID", s.MachineUUID,
		"ip", s.IP,
	)
}
