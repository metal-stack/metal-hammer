package cmd

import (
	log "github.com/inconshreveable/log15"
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
	// DevMode turn on devmode which prevents failing in some situations
	DevMode bool
	// BGPEnabled if set to true real bgp configuration is configured, otherwise dhcp will be used
	BGPEnabled bool
	// Cidr of BGP interface in DEV Mode
	Cidr string
	// ConsolePassword of the metal user valid for one day.
	ConsolePassword string
	// DeviceUUID is the unique identifier of this device
	DeviceUUID string
	// Ip of this instance
	Ip string
}

// Log print configuration options
func (s *Specification) Log() {
	log.Info("configuration",
		"debug", s.Debug,
		"metalCoreURL", s.MetalCoreURL,
		"imageURL", s.ImageURL,
		"imageID", s.ImageID,
		"devmode", s.DevMode,
		"bgpenabled", s.BGPEnabled,
		"cidr", s.Cidr,
	)
}
