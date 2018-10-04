package cmd

import (
	log "github.com/inconshreveable/log15"
)

//Specification defines configuration items which can be configured vi env variables
type Specification struct {
	Debug      bool   `default:"false" desc:"turn on debug log" required:"False"`
	ReportURL  string `default:"http://localhost:4242/device/register" desc:"Register endpoint url" required:"False"`
	InstallURL string `default:"http://localhost:4242/device/install" desc:"Get Image url of OS to install" required:"False"`
}

// Log print configuration options
func (s *Specification) Log() {
	log.Info("configuration", "debug", s.Debug, "reportURL", s.ReportURL)
}
