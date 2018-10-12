package cmd

import (
	log "github.com/inconshreveable/log15"
)

//Specification defines configuration items which can be configured vi env variables
type Specification struct {
	Debug       bool   `default:"true" desc:"turn on debug log" required:"False"`
	ReportURL   string `default:"http://localhost:4242/device/report" desc:"Report endpoint url" required:"False"`
	RegisterURL string `default:"http://localhost:4242/device/register" desc:"Register endpoint url" required:"False"`
	InstallURL  string `default:"http://localhost:4242/device/install" desc:"Get Image url of OS to install" required:"False"`
	ImageURL    string `default:"" desc:"Use a fixed Image url of OS to install" required:"False"`
	DevMode     bool   `default:"false" desc:"turn on devmode which prevents failing in some situations" required:"False"`
}

// Log print configuration options
func (s *Specification) Log() {
	log.Info("configuration", "debug", s.Debug, "reportURL", s.ReportURL, "installURL", s.InstallURL, "imageURL", s.ImageURL)
}
