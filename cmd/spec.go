package cmd

import (
	"log/slog"
	"strconv"
	"strings"

	"os"

	"github.com/metal-stack/metal-hammer/pkg/kernel"
	pixiecore "github.com/metal-stack/pixie/api"
)

// Specification defines configuration items of the application
type Specification struct {
	// Debug turn on debug log
	Debug bool
	// PixieAPIUrl is the endpoint URL where the pixie reside
	PixieAPIUrl string
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
	// MetalConfig is fetched from pixiecore to get the certs for the metal-api and logging config
	MetalConfig *pixiecore.MetalConfig

	log *slog.Logger
}

// NewSpec fills Specification with configuration made by kernel commandline
func NewSpec(log *slog.Logger) *Specification {
	spec := &Specification{}
	// Grab metal-hammer configuration from kernel commandline
	envmap, err := kernel.ParseCmdline()
	if err != nil {
		log.Error("parse cmdline", "error", err)
		os.Exit(1)
	}

	if d, ok := envmap["DEBUG"]; ok && (d == "1" || strings.ToLower(d) == "true") {
		spec.Debug = true
		os.Setenv("DEBUG", "1")
	}

	// PIXIE_API_URL must be in the form http://ip-of-pixie:4242
	if url, ok := envmap["PIXIE_API_URL"]; ok {
		spec.PixieAPIUrl = url
	}

	metalConfig, err := fetchMetalConfig(spec.PixieAPIUrl)
	if err != nil {
		log.Error("unable to fetch configuration from pixiecore", "error", err)
		os.Exit(1)
	}

	spec.MetalConfig = metalConfig

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
	s.log.Info("configuration",
		"debug", s.Debug,
		"pixieAPIUrl", s.PixieAPIUrl,
		"bgpenabled", s.BGPEnabled,
		"cidr", s.Cidr,
		"machineUUID", s.MachineUUID,
		"ip", s.IP,
	)
}
