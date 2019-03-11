package network

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
	mutex             sync.RWMutex
	neighbors         map[string][]*lldp.Neighbor
	interfaces        []string
	start             time.Time
	done              bool
	timeout           time.Duration
	minimumInterfaces int
	minimumNeighbors  int
}

const (
	// LLDPTxInterval is set to 10 seconds in /etc/lldpd.d/tx-interval.conf on each leaf.
	// FIXME set to 5 minutes until we have a workin setup
	LLDPTxInterval = 150 * time.Second

	// LLDPTxIntervalTimeout is set to double of tx-interval of lldpd on the switch side.
	// This ensures we get all lldp pdu`s.
	// We add 2 seconds to be on the save side.
	LLDPTxIntervalTimeout = (2 * LLDPTxInterval) + 2
)

// NewLLDPClient create a lldp client.
func NewLLDPClient(interfaces []string, minimumInterfaces, minimumNeighbors int, timeout time.Duration) *LLDPClient {
	if timeout == 0 {
		timeout = LLDPTxIntervalTimeout
	}
	return &LLDPClient{
		Host: &Host{
			mutex:             sync.RWMutex{},
			neighbors:         make(map[string][]*lldp.Neighbor),
			interfaces:        interfaces,
			start:             time.Now(),
			done:              false,
			timeout:           timeout,
			minimumInterfaces: minimumInterfaces,
			minimumNeighbors:  minimumNeighbors,
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
			l.Host.mutex.RLock()
			for _, existingNeigh := range l.Host.neighbors {
				for _, en := range existingNeigh {
					if en.Chassis.Value == neigh.Chassis.Value &&
						en.Port.Value == neigh.Port.Value {
						neighExists = true
						break
					}
				}
				if neighExists {
					break
				}
			}
			l.Host.mutex.RUnlock()
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

func (l *LLDPClient) requirementsMet() bool {
	// First check we have at least neighbors for 2 interfaces found
	if l.Host.minimumInterfaces == 0 && l.Host.minimumNeighbors == 0 {
		return true
	}
	if len(l.Host.neighbors) < l.Host.minimumInterfaces {
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
	if len(neighMap) >= l.Host.minimumNeighbors {
		return true
	}

	// Requirements are not met
	return false
}
