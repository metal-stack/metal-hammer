package lldp

import (
	"fmt"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	log "github.com/inconshreveable/log15"
	"github.com/pkg/errors"
)

// LinkType can be Interface or Mac
type LinkType string

var (
	// Interface LinkType
	Interface LinkType = "Interface"
	// Mac LinkType
	Mac LinkType = "Mac"
)

// Chassis of a lldp Neighbor
type Chassis struct {
	Type  LinkType
	Value string
}

// Port of a lldp Neighbor
type Port struct {
	Type  LinkType
	Value string
}

// Neighbor is the direct ethernet neighbor
type Neighbor struct {
	Name        string
	Description string
	Interface   string
	Chassis     Chassis
	Port        Port
}

func (c Chassis) String() string {
	return fmt.Sprintf("%s:%s", c.Type, c.Value)
}
func (p Port) String() string {
	return fmt.Sprintf("%s:%s", p.Type, p.Value)
}
func (n Neighbor) String() string {
	return fmt.Sprintf("Name:%s Desc:%s Chassis:%s Port:%s", n.Name, n.Description, n.Chassis, n.Port)
}

// Client consumes lldp messages.
type Client struct {
	Source    *gopacket.PacketSource
	Handle    *pcap.Handle
	Interface *net.Interface
}

// NewClient create a new lldp client.
func NewClient(ifi string) (*Client, error) {
	iface, err := net.InterfaceByName(ifi)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to lookup interface:%s", ifi)
	}
	log.Info("lldp", "listen on interface", iface.Name)

	handle, err := pcap.OpenLive(iface.Name, 65536, true, 5*time.Second)

	if err != nil {
		return nil, errors.Wrapf(err, "unable to open interface:%s in promiscuous mode", iface.Name)
	}
	// Only snoop for LLDP Packets not coming from this interface
	bpfFilter := fmt.Sprintf("ether proto 0x88cc and not ether host %s", iface.HardwareAddr)
	err = handle.SetBPFFilter(bpfFilter)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to filter ethernet traffic 088cc on interface:%s", iface.Name)
	}
	src := gopacket.NewPacketSource(handle, handle.LinkType())
	return &Client{Source: src, Handle: handle, Interface: iface}, nil
}

// Close the lldp client
func (l *Client) Close() {
	l.Handle.Close()
}

// Neighbors search on a interface for neighbors announced via lldp
func (l *Client) Neighbors(neighChan chan Neighbor) {
	for {
		for packet := range l.Source.Packets() {
			switch packet.LinkLayer().LayerType() {
			case layers.LayerTypeEthernet:
				neigh := Neighbor{}
				for _, layer := range packet.Layers() {
					layerType := layer.LayerType()
					switch layerType {
					case layers.LayerTypeLinkLayerDiscovery:
						lldp := layer.(*layers.LinkLayerDiscovery)
						chassis := Chassis{}
						port := Port{}
						var chassismac net.HardwareAddr
						var portmac net.HardwareAddr
						switch lldp.PortID.Subtype {
						case layers.LLDPPortIDSubtypeMACAddr:
							portmac = lldp.PortID.ID
							port.Type = Mac
							port.Value = portmac.String()
						case layers.LLDPPortIDSubtypeIfaceName:
							port.Type = Interface
							port.Value = string(lldp.PortID.ID)
						}
						switch lldp.ChassisID.Subtype {
						case layers.LLDPChassisIDSubTypeMACAddr:
							chassismac = lldp.ChassisID.ID
							chassis.Type = Mac
							chassis.Value = chassismac.String()
						case layers.LLDPChassisIDSubtypeIfaceName:
							chassis.Type = Interface
							chassis.Value = string(lldp.ChassisID.ID)
						}
						neigh.Chassis = chassis
						neigh.Port = port
					case layers.LayerTypeLinkLayerDiscoveryInfo:
						lldpi := layer.(*layers.LinkLayerDiscoveryInfo)
						neigh.Name = lldpi.SysName
						neigh.Description = lldpi.SysDescription
						neigh.Interface = l.Interface.Name
						neighChan <- neigh
					}
				}
			}
		}
	}
}
