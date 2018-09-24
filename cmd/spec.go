package cmd

import (
	log "github.com/inconshreveable/log15"
)

//Specification defines configuration items which can be configured vi env variables
type Specification struct {
	Debug     bool   `default:"false" desc:"turn on debug log" required:"False"`
	ReportURL string `default:"http://localhost:8080/device/register" desc:"Register endpoint url" required:"False"`
}

// Log print configuration options
func (s *Specification) Log() {
	log.Info("configuration", "debug", s.Debug, "reportURL", s.ReportURL)
}
