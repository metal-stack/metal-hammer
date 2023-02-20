package network

import (
	"fmt"
	"net"
	"strings"
	"time"

	v1 "github.com/metal-stack/metal-api/pkg/api/v1"

	"github.com/metal-stack/go-lldpd/pkg/lldp"
	"github.com/metal-stack/v"
	"go.uber.org/zap"

	"github.com/vishvananda/netlink"
)

// Network provides networking operations.
type Network struct {
	IPAddress   string
	Started     time.Time
	MachineUUID string
	LLDPClient  *LLDPClient
	Eth0Mac     string // this mac is used to calculate the IPMI Port offset in the metal-lab environment.
	Log         *zap.SugaredLogger
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
	ethtool := NewEthtool(n.Log)
	for _, name := range Interfaces() {
		if !strings.HasPrefix(name, "eth") {
			continue
		}
		interfaces = append(interfaces, name)

		err := linkSetMTU(name, MTU)
		if err != nil {
			return fmt.Errorf("error set link %s mtu %w", name, err)
		}

		err = linkSetUp(name)
		if err != nil {
			return fmt.Errorf("error set link %s up %w", name, err)
		}

		ethtool.disableFirmwareLLDP(name)

		lldpd, err := lldp.NewDaemon(n.Log, n.MachineUUID, description, name, 5*time.Second)
		if err != nil {
			return fmt.Errorf("error start lldpd on %s %w", name, err)
		}
		lldpd.Start()
	}

	lc := NewLLDPClient(n.Log, interfaces, 2, 2, 0)
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
func (n *Network) Neighbors(name string) (neighbors []*v1.MachineNic, err error) {
	host := n.LLDPClient.Host

	for !host.done {
		actualNeigh := len(host.neighbors)
		minimumNeigh := host.minimumNeighbors
		n.Log.Infow("waiting for lldp neighbors", "interface", name, "actual", actualNeigh, "minimum", minimumNeigh)
		time.Sleep(1 * time.Second)

		duration := time.Since(host.start)
		if duration > host.timeout {
			return nil, fmt.Errorf("not all neighbor requirements where met within: %s, exiting", host.timeout)
		}
	}
	n.Log.Infow("all lldp pdu's received", "interface", name)

	neighs := host.neighbors[name]
	for _, neigh := range neighs {
		identifier := neigh.Port.Value
		n.Log.Infow("register add neighbor", "nic", name, "identifier", identifier)
		neighbors = append(neighbors, &v1.MachineNic{
			Mac:        identifier,
			Identifier: identifier,
			Name:       name,
			Hostname:   neigh.Name,
		})
	}
	return neighbors, nil
}

// InternalIP returns the first ipv4 ip of an eth* interface.
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
