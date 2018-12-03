package cmd

import (
	"fmt"
	"strings"

	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/models"

	"github.com/jaypipes/ghw"
	"github.com/vishvananda/netlink"
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
func Neighbors(name string) ([]*models.ModelsMetalNic, error) {
	neighbors := make([]*models.ModelsMetalNic, 0)

	link, err := netlink.LinkByName(name)
	if err != nil {
		return neighbors, err
	}

	// TODO: Maybe we can use FAMILY_ALL as well for both v4 and v6,
	// but we need an environment with IPv6 neighbors to check if it's working
	v4, err := netlink.NeighList(link.Attrs().Index, netlink.FAMILY_V4)
	if err != nil {
		return neighbors, err
	}
	v6, err := netlink.NeighList(link.Attrs().Index, netlink.FAMILY_V6)
	if err != nil {
		return neighbors, err
	}

	macs := map[string]bool{}

	for _, n := range v4 {
		macs[n.HardwareAddr.String()] = true
	}
	for _, n := range v6 {
		macs[n.HardwareAddr.String()] = true
	}

	for mac := range macs {
		macAddress := mac
		neighbors = append(neighbors, &models.ModelsMetalNic{Mac: &macAddress})
	}

	return neighbors, nil
}
