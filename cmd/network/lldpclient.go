package network

import (
	"sync"
	"time"

	"github.com/metal-stack/metal-hammer/pkg/lldp"
	"go.uber.org/zap"
)

// LLDPClient act as a small wrapper about low level lldp primitives.
type LLDPClient struct {
	Host *Host
	log  *zap.SugaredLogger
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
	LLDPTxInterval = 15 * time.Second

	// LLDPTxIntervalTimeout is set to double of tx-interval of lldpd on the switch side.
	// This ensures we get all lldp pdu`s.
	// We add 2 seconds to be on the save side.
	LLDPTxIntervalTimeout = (2 * LLDPTxInterval) + 2
)

// NewLLDPClient create a lldp client.
func NewLLDPClient(log *zap.SugaredLogger, interfaces []string, minimumInterfaces, minimumNeighbors int, timeout time.Duration) *LLDPClient {
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
		log: log,
	}
}

// Start starts lldpd for neighbor discovery.
func (l *LLDPClient) Start() {
	l.log.Infow("lldp start discovery")
	neighChan := make(chan lldp.Neighbor)
	for _, ifi := range l.Host.interfaces {
		lldpcli, err := lldp.NewClient(l.log, ifi)
		if err != nil {
			l.log.Errorw("lldp", "unable to start client on", ifi, "error", err)
			continue
		}
		go lldpcli.Neighbors(neighChan)
	}

	for detectedNeighbor := range neighChan {
		l.log.Debugw("lldp", "detectedNeighbor", detectedNeighbor)
		if l.neighborKnown(detectedNeighbor) {
			continue
		}

		l.addNeighbor(detectedNeighbor)
		l.log.Infow("lldp", "neighbors", l.Host.neighbors)
	}
}

// neighborKnown returns if the given neighbor is already known
func (l *LLDPClient) neighborKnown(neighbor lldp.Neighbor) bool {
	l.Host.mutex.RLock()
	defer l.Host.mutex.RUnlock()

	for _, knownNeighbors := range l.Host.neighbors {
		for _, kn := range knownNeighbors {
			if kn.Chassis.Value == neighbor.Chassis.Value && kn.Port.Value == neighbor.Port.Value {
				return true
			}
		}
	}
	return false
}

// addNeighbor adds the neighbor to the known neighbors
func (l *LLDPClient) addNeighbor(neighbor lldp.Neighbor) {
	l.Host.mutex.Lock()
	defer l.Host.mutex.Unlock()

	l.Host.neighbors[neighbor.Interface] = append(l.Host.neighbors[neighbor.Interface], &neighbor)
	l.Host.done = l.requirementsMet()
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
	// and every port type of a interface on the switch is set to mac
	neighMap := make(map[string]string)
	for iface, neighs := range l.Host.neighbors {
		for _, neigh := range neighs {
			if neigh.Chassis.Type == lldp.Mac && neigh.Port.Type == lldp.Mac {
				neighMap[neigh.Chassis.Value] = iface
			}
		}
	}
	// Requirements are met if we found at least 2 distinct chassis mac's
	return len(neighMap) >= l.Host.minimumNeighbors
}
