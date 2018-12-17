package cmd

import (
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/lldp"
	log "github.com/inconshreveable/log15"
	"sync"
	"time"
)

// LLDPClient act as a small wrapper about low level lldp primitives.
type LLDPClient struct {
	Host *Host
}

// Host collects lldp neighbor information's.
type Host struct {
	mutex      sync.Mutex
	neighbors  map[string][]*lldp.Neighbor
	interfaces []string
	start      time.Time
	done       bool
}

const (
	// LLDPTxInterval is set to 10 seconds in /etc/lldpd.d/tx-interval.conf on each leaf.
	LLDPTxInterval = 10 * time.Second

	// LLDPTxIntervalTimeout is set to double of tx-interval of lldpd on the switch side.
	// This ensures we get all lldp pdu`s.
	// We add 2 seconds to be on the save side.
	LLDPTxIntervalTimeout = (2 * LLDPTxInterval) + 2
)

// NewLLDPClient create a lldp client.
func NewLLDPClient(interfaces []string) *LLDPClient {
	return &LLDPClient{
		Host: &Host{
			mutex:      sync.Mutex{},
			neighbors:  make(map[string][]*lldp.Neighbor),
			interfaces: interfaces,
			start:      time.Now(),
			done:       false,
		},
	}
}

// Start will start lldpd for neighbor discovery.
func (l *LLDPClient) Start() {
	log.Info("lldp start discovery")
	neighChan := make(chan lldp.Neighbor)
	for _, ifi := range l.Host.interfaces {
		lldpcli, err := lldp.NewClient(ifi)
		if err != nil {
			log.Error("lldp", "unable to start client on", ifi, "error", err)
			continue
		}
		go lldpcli.Neighbors(neighChan)
	}

	for {
		select {
		case neigh := <-neighChan:
			log.Debug("lldp", "neigh", neigh)
			neighExists := false
			for _, existingNeigh := range l.Host.neighbors {
				for _, en := range existingNeigh {
					if en.Chassis.Value == neigh.Chassis.Value &&
						en.Port.Value == neigh.Port.Value {
						neighExists = true
					}
				}
			}
			if neighExists {
				break
			}
			l.Host.mutex.Lock()
			l.Host.neighbors[neigh.Interface] = append(l.Host.neighbors[neigh.Interface], &neigh)
			l.Host.done = l.requirementsMet()
			l.Host.mutex.Unlock()
			log.Info("lldp", "neighbors", l.Host.neighbors)
		}
	}
}

const minimumInterfaces = 2
const minimumNeighbors = 2

func (l *LLDPClient) requirementsMet() bool {
	// First check we have at least neighbors for 2 interfaces found
	if len(l.Host.neighbors) < minimumInterfaces {
		return false
	}
	// Then check if 2 distinct Chassis neighbors where found
	neighMap := make(map[string]string)
	for iface, neighs := range l.Host.neighbors {
		for _, neigh := range neighs {
			if neigh.Chassis.Type == lldp.Mac {
				neighMap[neigh.Chassis.Value] = iface
			}
		}
	}
	// OK we found 2 distinct chassis mac's
	if len(neighMap) >= minimumNeighbors {
		return true
	}

	// Requirements are not met
	return false
}
