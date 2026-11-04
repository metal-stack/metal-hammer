package network

import (
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/metal-stack/metal-hammer/pkg/lldp"
)

// LLDPClient act as a small wrapper about low level lldp primitives.
type LLDPClient struct {
	Host *Host
	log  *slog.Logger
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

const (
	// CmdlineLLDPMinInterfaces is the kernel commandline key that lowers LLDPRequirements.MinInterfaces.
	CmdlineLLDPMinInterfaces = "LLDP_MIN_INTERFACES"
	// CmdlineLLDPMinNeighbors is the kernel commandline key that lowers LLDPRequirements.MinNeighbors.
	CmdlineLLDPMinNeighbors = "LLDP_MIN_NEIGHBORS"
)

// LLDPRequirements define when neighbor discovery is complete.
type LLDPRequirements struct {
	// MinInterfaces is the number of interfaces that must have received at least one LLDP PDU.
	MinInterfaces int
	// MinNeighbors is the number of distinct switches (chassis MACs) that must have been seen.
	MinNeighbors int
}

// DefaultLLDPRequirements is the production setup: every machine is dual-homed to two leaf switches.
var DefaultLLDPRequirements = LLDPRequirements{MinInterfaces: 2, MinNeighbors: 2}

// ParseLLDPRequirements reads LLDP_MIN_INTERFACES and LLDP_MIN_NEIGHBORS from the kernel commandline.
// Missing keys keep the default. Values below 1 are rejected, so neighbor discovery can be relaxed
// for single-switch setups (e.g. hardware labs) but never skipped.
func ParseLLDPRequirements(cmdline map[string]string) (LLDPRequirements, error) {
	r := DefaultLLDPRequirements
	for key, target := range map[string]*int{
		CmdlineLLDPMinInterfaces: &r.MinInterfaces,
		CmdlineLLDPMinNeighbors:  &r.MinNeighbors,
	} {
		raw, ok := cmdline[key]
		if !ok {
			continue
		}
		n, err := strconv.Atoi(raw)
		if err != nil {
			return DefaultLLDPRequirements, fmt.Errorf("%s=%q is not a number", key, raw)
		}
		if n < 1 {
			return DefaultLLDPRequirements, fmt.Errorf("%s must be at least 1, got %d", key, n)
		}
		*target = n
	}
	return r, nil
}

// NewLLDPClient create a lldp client.
func NewLLDPClient(log *slog.Logger, interfaces []string, minimumInterfaces, minimumNeighbors int, timeout time.Duration) *LLDPClient {
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
	l.log.Info("lldp start discovery")
	neighChan := make(chan lldp.Neighbor)
	for _, ifi := range l.Host.interfaces {
		lldpcli, err := lldp.NewClient(l.log, ifi)
		if err != nil {
			l.log.Error("lldp", "unable to start client on", ifi, "error", err)
			continue
		}
		go lldpcli.Neighbors(neighChan)
	}

	for detectedNeighbor := range neighChan {
		l.log.Debug("lldp", "detectedNeighbor", detectedNeighbor)
		if l.neighborKnown(detectedNeighbor) {
			continue
		}

		l.addNeighbor(detectedNeighbor)
		l.log.Info("lldp", "neighbors", l.Host.neighbors)
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
	// First check we have neighbors on at least minimumInterfaces interfaces
	if l.Host.minimumInterfaces == 0 && l.Host.minimumNeighbors == 0 {
		return true
	}
	if len(l.Host.neighbors) < l.Host.minimumInterfaces {
		return false
	}

	// Then check if minimumNeighbors distinct Chassis neighbors were found
	neighMap := make(map[string]string)
	for iface, neighs := range l.Host.neighbors {
		for _, neigh := range neighs {
			if neigh.Chassis.Type == lldp.Mac && (neigh.Port.Type == lldp.Mac || neigh.Port.Type == lldp.Local) {
				neighMap[neigh.Chassis.Value] = iface
			}
		}
	}
	// Requirements are met if we found at least minimumNeighbors distinct chassis mac's
	return len(neighMap) >= l.Host.minimumNeighbors
}
