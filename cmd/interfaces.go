package cmd

import (
	"fmt"
	log "github.com/inconshreveable/log15"
	"strings"
	"time"

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

	interfaces := make([]string, 0)
	for _, nic := range net.NICs {
		if !strings.HasPrefix(nic.Name, "eth") {
			continue
		}
		interfaces = append(interfaces, nic.Name)

		err := linkSetUp(nic.Name)
		if err != nil {
			return fmt.Errorf("Error set link %s up: %v", nic.Name, err)
		}
	}

	go h.StartLLDPDClient(interfaces)

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

	for !host.done {
		log.Info("not all lldp pdu's are received, waiting...", "interface", name)
		time.Sleep(1 * time.Second)
	}
	log.Info("all lldp pdu's received", "interface", name)

	neighs, _ := host.neighbors[name]
	for _, neigh := range neighs {
		macAddress := neigh.Chassis.Value
		neighbors = append(neighbors, &models.ModelsMetalNic{Mac: &macAddress})
	}
	return neighbors, nil
}
