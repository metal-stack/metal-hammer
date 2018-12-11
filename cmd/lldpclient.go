package cmd

import (
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/lldp"
	log "github.com/inconshreveable/log15"
	"sync"
	"time"
)

type (
	neighbor map[string][]*lldp.Neighbor

	Host struct {
		mutex     sync.Mutex
		neighbors neighbor
		done      bool
	}
)

var (
	host = Host{
		mutex:     sync.Mutex{},
		neighbors: make(map[string][]*lldp.Neighbor),
		done:      false,
	}
)

const (
	// LLDPTxInterval is set to 10 seconds in /etc/lldpd.d/tx-interval.conf on each leaf.
	LLDPTxInterval = 10 * time.Second

	// LLDPTxIntervalTimeout is set to double of tx-interval of lldpd on the switch side.
	// This ensures we get all lldp pdu`s.
	// We add 2 seconds to be on the save side.
	LLDPTxIntervalTimeout = (2 * LLDPTxInterval) + 2
)

// StartLLDPDClient will start lldpd for neighbor discovery.
func (h *Hammer) StartLLDPDClient(interfaces []string) {
	log.Info("lldp start discovery")
	neighChan := make(chan lldp.Neighbor)
	for _, ifi := range interfaces {
		lldpcli, err := lldp.NewLLDPClient(ifi)
		if err != nil {
			log.Error("lldp", "unable to start client on", ifi, "error", err)
			continue
		}
		go lldpcli.Neighbors(neighChan)
	}

	go func(timeout time.Duration) {
		log.Info("lldp", "wait", timeout)
		time.Sleep(timeout)
		host.done = true
	}(LLDPTxIntervalTimeout)

	for {
		select {
		case neigh := <-neighChan:
			log.Debug("lldp", "neigh", neigh)
			found := false
			for _, value := range host.neighbors {
				for _, v := range value {
					if v.Chassis.Value == neigh.Chassis.Value &&
						v.Port.Value == neigh.Port.Value {
						found = true
						log.Debug("lldp", "neigh known", neigh)
					}
				}
			}
			if found {
				break
			}
			host.mutex.Lock()
			host.neighbors[neigh.Interface] = append(host.neighbors[neigh.Interface], &neigh)
			host.mutex.Unlock()
			log.Info("lldp", "neighbors", host.neighbors)
		}
	}
}
