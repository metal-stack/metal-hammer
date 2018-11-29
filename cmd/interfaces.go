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
