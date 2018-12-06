package cmd

import (
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/lldpd"
	log "github.com/inconshreveable/log15"
	"net"
)

// StartLLDPD will start lldpd for neighbor discovery.
func (h *Hammer) StartLLDPD(name string) {
	ifi, err := net.InterfaceByName(name)
	if err != nil {
		log.Error("lldpd", "unable get nic", err)
	}
	l := lldpd.New(lldpd.SourceAddress(ifi.HardwareAddr))
	err = l.Listen()
	if err != nil {
		log.Error("lldpd", "listen", err)
	}
}
