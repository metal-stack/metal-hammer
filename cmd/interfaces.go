package cmd

import (
	"fmt"
	"github.com/jaypipes/ghw"
	"github.com/vishvananda/netlink"
	"strings"
)

// UpAllInterfaces set all available eth* interfaces up
// to ensure they do ipv6 link local autoconfiguration and
// therefore neighbor discovery,
// which is required to make all local mac's visible on the switch side.
func (h *Hammer) UpAllInterfaces() error {
	net, err := ghw.Network()
	if err != nil {
		return fmt.Errorf("Error getting network info: %v", err)
	}

	for _, nic := range net.NICs {
		if !strings.HasPrefix(nic.Name, "eth") {
			continue
		}

		err := linkSetUp(nic.Name)
		if err != nil {
			return fmt.Errorf("Error set link %s up: %v", nic.Name, err)
		}
	}
	return nil
}

func linkSetUp(name string) error {
	iface, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}
	err = netlink.LinkSetUp(iface)
	if err != nil {
		return err
	}
	return nil
}

// Neighbors of a interface, detected via ip neighbor detection
func Neighbors(name string) ([]string, error) {
	iface, err := netlink.LinkByName(name)
	if err != nil {
		return nil, err
	}
	neighbors := make([]string, 0)
	neigh, err := netlink.NeighList(iface.Attrs().Index, 4)
	if err != nil {
		return nil, err
	}

	for _, n := range neigh {
		neighbors = append(neighbors, string(n.HardwareAddr))
	}

	return neighbors, nil
}
