package network

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/metal-stack/go-lldpd/pkg/lldp"
	"github.com/metal-stack/metal-hammer/metal-core/models"
	"github.com/metal-stack/v"

	log "github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
)

// Network provides networking operations.
type Network struct {
	IPAddress   string
	Started     time.Time
	MachineUUID string
	LLDPClient  *LLDPClient
	Eth0Mac     string // this mac is used to calculate the IPMI Port offset in the metal-lab environment.
}

// We expect to have storage and MTU of 9000 supports efficient transmission.
// In our clos topology MTU 9000 (non vxlan)/9216 (vxlan) is status quo.
const MTU = 9000

// UpAllInterfaces set all available eth* interfaces up
// to ensure they do ipv6 link local autoconfiguration and
// therefore neighbor discovery,
// which is required to make all local mac's visible on the switch side.
func (n *Network) UpAllInterfaces() error {
	description := fmt.Sprintf("metal-hammer IP:%s version:%s waiting since %s for installation", n.IPAddress, v.V, n.Started)
	interfaces := make([]string, 0)
	ethtool := NewEthtool()
	for _, name := range Interfaces() {
		if !strings.HasPrefix(name, "eth") {
			continue
		}
		interfaces = append(interfaces, name)

		err := linkSetMTU(name, MTU)
		if err != nil {
			return errors.Wrapf(err, "Error set link %s mtu", name)
		}

		err = linkSetUp(name)
		if err != nil {
			return errors.Wrapf(err, "Error set link %s up", name)
		}
		// This will take time
		// if !linkIsUp(name) {
		// 	continue
		// }
		ethtool.disableFirmwareLLDP(name)

		lldpd, err := lldp.NewDaemon(n.MachineUUID, description, name, 5*time.Second)
		if err != nil {
			return errors.Wrapf(err, "Error start lldpd on %s", name)
		}
		lldpd.Start()
	}

	lc := NewLLDPClient(interfaces, 2, 2, 0)
	n.LLDPClient = lc
	go lc.Start()

	return nil
}

func linkSetMTU(name string, mtu int) error {
	iface, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}

	err = netlink.LinkSetMTU(iface, mtu)
	if err != nil {
		return err
	}
	return err
}

// Kept for documentary purpose
// func linkIsUp(name string) bool {
// 	iface, err := netlink.LinkByName(name)
// 	if err != nil {
// 		return false
// 	}

// 	return iface.Attrs().OperState == netlink.OperUp
// }

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
func (n *Network) Neighbors(name string) ([]*models.ModelsV1MachineNicExtended, error) {
	neighbors := make([]*models.ModelsV1MachineNicExtended, 0)

	host := n.LLDPClient.Host

	for !host.done {
		actualNeigh := len(host.neighbors)
		minimumNeigh := host.minimumNeighbors
		log.Info("waiting for lldp neighbors", "interface", name, "actual", actualNeigh, "minimum", minimumNeigh)
		time.Sleep(1 * time.Second)

		duration := time.Since(host.start)
		if duration > host.timeout {
			return nil, errors.Errorf("not all neighbor requirements where met within: %s, exiting", host.timeout)
		}
	}
	log.Info("all lldp pdu's received", "interface", name)

	neighs := host.neighbors[name]
	for _, neigh := range neighs {
		macAddress := neigh.Port.Value
		neighbors = append(neighbors, &models.ModelsV1MachineNicExtended{Mac: &macAddress})
	}
	return neighbors, nil
}

// InternalIP returns the first ipv4 ip of a eth* interface.
func InternalIP() string {
	for _, name := range Interfaces() {
		if !strings.HasPrefix(name, "eth") {
			continue
		}
		itf, _ := net.InterfaceByName(name)
		item, _ := itf.Addrs()
		for _, addr := range item {
			switch v := addr.(type) {
			case *net.IPNet:
				if !v.IP.IsLoopback() {
					if v.IP.To4() != nil {
						return v.IP.String()
					}
				}
			}
		}
	}
	return ""
}

// Interfaces return a list of all known interfaces.
func Interfaces() []string {
	var interfaces []string
	links, err := netlink.LinkList()
	if err != nil {
		return interfaces
	}
	for _, nic := range links {
		name := nic.Attrs().Name
		interfaces = append(interfaces, name)
	}
	return interfaces
}
