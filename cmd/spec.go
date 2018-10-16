package cmd

import (
	log "github.com/inconshreveable/log15"
)

//Specification defines configuration items of the application
type Specification struct {
	// Debug turn on debug log
	Debug bool
	// ReportURL is the endpoint URL where to report installation success
	ReportURL string
	// RegisterURL is the endpoint where to send device discovery information
	RegisterURL string
	// InstallURL the url where to get the installation from
	InstallURL string
	// ImageURL if given grabs a fixed OS image to install, only suitable in DevMode
	ImageURL string
	// DevMode turn on devmode which prevents failing in some situations
	DevMode bool
}

// Log print configuration options
func (s *Specification) Log() {
	log.Info("configuration",
		"debug", s.Debug,
		"reportURL", s.ReportURL,
		"installURL", s.InstallURL,
		"imageURL", s.ImageURL,
		"devmode", s.DevMode,
	)
}
