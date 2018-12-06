package lldpd

import (
	"bytes"
	"net"
	"sync"
	"time"

	log "github.com/inconshreveable/log15"
	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/lldp"
	"github.com/mdlayher/raw"
)

// LLDPD is the server for LLDP PDU's
// It will always listen passively. This means, it will
// only send LLDP PDU's in response to a received PDU.
type LLDPD struct {
	filterFn      InterfaceFilterFn
	portLookupFn  PortLookupFn
	replyUnicast  bool
	sourceAddress net.HardwareAddr

	recvChannel chan *message
	sendChannel chan *message

	listenersLock sync.RWMutex
	listeners     map[int]*packetConn
}

type packetConn struct {
	conn   *raw.Conn
	packet []byte
}

// New will return a new LLDPD server with the optional
// options configured.
func New(opts ...Option) *LLDPD {
	l := &LLDPD{
		filterFn:      defaultInterfaceFilterFn,
		portLookupFn:  defaultPortLookupFn,
		replyUnicast:  false,
		sourceAddress: []byte{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad},
		recvChannel:   make(chan *message, 64),
		sendChannel:   make(chan *message, 64),
		listeners:     make(map[int]*packetConn),
	}

	for _, opt := range opts {
		l.SetOption(opt)
	}

	return l
}

func (l *LLDPD) startNLLoop() {
	nl := NewNLListener()
	nl.Start()

	go func() {
		for {
			select {
			case info := <-nl.Messages:
				switch info.op {
				case IF_ADD:
					if l.filterFn(info.ifi) {
						l.ListenOn(info.ifi)
					}
				case IF_DEL:
					l.CancelListenOn(info.ifi)
				}
			}
		}
	}()
}

// ListenOn will listen on the specified interface for
// LLDP PDU's
func (l *LLDPD) ListenOn(ifi *net.Interface) {
	l.listenersLock.Lock()
	defer l.listenersLock.Unlock()
	if _, ok := l.listeners[ifi.Index]; !ok {
		conn, err := raw.ListenPacket(ifi, uint16(lldp.EtherType), nil)
		if err != nil {
			log.Error("error listening on interface", "ifname", ifi.Name, "ifindex", ifi.Index, "error", err)
			return
		}

		l.listeners[ifi.Index] = &packetConn{
			conn: conn,
		}

		go func() {
			log.Info("started listener on interface", "ifname", ifi.Name, "ifindex", ifi.Index)
			b := make([]byte, ifi.MTU)

			for {
				_, src, err := conn.ReadFrom(b)
				if err != nil {
					log.Error("error ReadFrom interface", "ifname", ifi.Name, "ifindex", ifi.Index, "error", err)
					continue
				}
				//spew.Dump(src, err, b[:n])
				l.recvChannel <- &message{
					ifi:  ifi,
					addr: src.(*raw.Addr),
				}
			}
		}()
	}
}

// CancelListenOn will stop listening on the interface
func (l *LLDPD) CancelListenOn(ifi *net.Interface) {
	l.listenersLock.Lock()
	defer l.listenersLock.Unlock()
	if pconn, ok := l.listeners[ifi.Index]; ok {
		pconn.conn.Close()
		delete(l.listeners, ifi.Index)
		log.Info("closed listener on interface", "ifname", ifi.Name, "ifindex", ifi.Index)
	}
}

// Listen will start the main listener loop
func (l *LLDPD) Listen() error {
	l.startNLLoop()

	go func() {
		for {
			select {
			case msg := <-l.sendChannel:
				l.listenersLock.RLock()
				if _, ok := l.listeners[msg.ifi.Index]; !ok {
					l.listenersLock.RUnlock()
					continue
				}
				pconn := l.listeners[msg.ifi.Index]
				l.listenersLock.RUnlock()

				b := l.packetFor(msg)

				_, err := pconn.conn.WriteTo(b, msg.addr)
				if err != nil {
					log.Error("error sending pdu out on interface", "name", msg.ifi.Name, "index", msg.ifi.Index, "error", err)
					continue
				}

				log.Info("send pdu out on interface", "name", msg.ifi.Name, "index", msg.ifi.Index)
				continue
			}
			break
		}
	}()

	for {
		select {
		case msg := <-l.recvChannel:
			log.Info("incoming pdu on interface", "name", msg.ifi.Name, "index", msg.ifi.Index)
			l.sendChannel <- &message{
				ifi:  msg.ifi,
				addr: msg.addr,
			}
			continue
		}
		break
	}
	close(l.sendChannel)

	return nil
}

func (l *LLDPD) packetFor(msg *message) []byte {
	l.listenersLock.RLock()
	if packet, ok := l.listeners[msg.ifi.Index]; ok {
		if packet.packet != nil {
			l.listenersLock.RUnlock()
			return packet.packet
		}
	}
	l.listenersLock.RUnlock()

	pDescr := l.portLookupFn(msg.ifi)
	var portDescr bytes.Buffer
	portDescr.WriteString(pDescr)

	systemName := "metal-hammer"
	systemDesc := "waiting for allocation"
	lf := lldp.Frame{
		ChassisID: &lldp.ChassisID{
			Subtype: lldp.ChassisIDSubtypeMACAddress,
			ID:      l.sourceAddress,
		},
		PortID: &lldp.PortID{
			Subtype: lldp.PortIDSubtypeAgentCircuitID,
			ID:      []byte{'1'},
		},
		TTL: 10 * time.Second,
		Optional: []*lldp.TLV{
			{
				Type:   lldp.TLVTypePortDescription,
				Value:  portDescr.Bytes(),
				Length: uint16(portDescr.Len()),
			},
			{
				Type:   lldp.TLVTypeSystemName,
				Value:  []byte(systemName),
				Length: uint16(len(systemName)),
			},
			{
				Type:   lldp.TLVTypeSystemDescription,
				Value:  []byte(systemDesc),
				Length: uint16(len(systemDesc)),
			},
		},
	}

	b, err := lf.MarshalBinary()
	if err != nil {
		log.Error("error marshalling lldp frame", "error", err)
		return nil
	}

	dest := net.HardwareAddr{0x01, 0x80, 0xc2, 0x00, 0x00, 0x0e}
	if l.replyUnicast {
		dest = msg.addr.HardwareAddr
	}
	f := &ethernet.Frame{
		Destination: dest,
		Source:      l.sourceAddress,
		EtherType:   lldp.EtherType,
		Payload:     b,
	}
	frame, err := f.MarshalBinary()

	if err != nil {
		log.Error("error marshalling ethernet frame", "error", err)
		return nil
	}

	l.listenersLock.Lock()
	l.listeners[msg.ifi.Index].packet = frame
	l.listenersLock.Unlock()

	return frame
}

type message struct {
	addr *raw.Addr
	ifi  *net.Interface
}
