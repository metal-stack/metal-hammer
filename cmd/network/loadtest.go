package network

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"log/slog"
	"math/rand/v2"
	"net"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

// ErrUplinkLoadTest is returned (wrapped) if the load test between the uplinks ran,
// but its result indicates a faulty link, e.g. a broken GBIC.
// Errors which prevent the test from running at all do not wrap this error.
var ErrUplinkLoadTest = errors.New("uplink load test failed")

const (
	// UplinkLoadTestFailureHold is the duration to wait before rebooting after a failed uplink load test.
	// This prevents a machine with a broken link from generating load with a endless and fast boot cycle.
	UplinkLoadTestFailureHold = 4 * time.Hour

	// loadTestEtherType is one of the ethertypes reserved by IEEE for local experimental use.
	loadTestEtherType uint16 = 0x88b5
	// loadTestMagic identifies a frame of the load test: "mhlt"
	loadTestMagic uint32 = 0x6d686c74

	loadTestKindProbe byte = 1
	loadTestKindLoad  byte = 2

	// frame layout:
	//  0- 5 destination mac
	//  6-11 source mac
	// 12-13 ethertype
	// 14-17 magic
	// 18    kind
	// 19    unused
	// 20-21 payload length
	// 22-29 run id
	// 30-37 sequence number
	// 38-41 crc32 of bytes 14-37 xor crc32 of the payload
	// 42-   payload
	loadTestEthHeaderLen  = 14
	loadTestPayloadOffset = 42
	loadTestMinPayload    = 18
	loadTestStandardMTU   = 1500

	loadTestProbeTimeout  = 5 * time.Second
	loadTestProbeInterval = 50 * time.Millisecond
	// loadTestJumboGrace is the time to wait for full size probes once small probes went through.
	loadTestJumboGrace = time.Second
	// loadTestDrain is the time given to frames still on the wire after the senders stopped.
	loadTestDrain      = 500 * time.Millisecond
	loadTestRcvBuf     = 64 * 1024 * 1024
	loadTestRcvTimeout = 100 * time.Millisecond
	loadTestSndTimeout = time.Second

	// loadTestMaxLossRatio is the maximum tolerated ratio of frames lost on the wire.
	loadTestMaxLossRatio = 0.001
)

// nicCounters are the kernel counters of a interface which are relevant for the load test.
type nicCounters struct {
	rxErrors  uint64
	rxCrc     uint64
	rxFrame   uint64
	rxLength  uint64
	rxDropped uint64
	rxMissed  uint64
	rxOver    uint64
	rxFifo    uint64
	txErrors  uint64
	txCarrier uint64
	txDropped uint64
}

func (c nicCounters) sub(o nicCounters) nicCounters {
	return nicCounters{
		rxErrors:  c.rxErrors - o.rxErrors,
		rxCrc:     c.rxCrc - o.rxCrc,
		rxFrame:   c.rxFrame - o.rxFrame,
		rxLength:  c.rxLength - o.rxLength,
		rxDropped: c.rxDropped - o.rxDropped,
		rxMissed:  c.rxMissed - o.rxMissed,
		rxOver:    c.rxOver - o.rxOver,
		rxFifo:    c.rxFifo - o.rxFifo,
		txErrors:  c.txErrors - o.txErrors,
		txCarrier: c.txCarrier - o.txCarrier,
		txDropped: c.txDropped - o.txDropped,
	}
}

func readNicCounters(name string) (nicCounters, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return nicCounters{}, fmt.Errorf("unable to read counters of %s %w", name, err)
	}
	s := link.Attrs().Statistics
	if s == nil {
		return nicCounters{}, nil
	}
	return nicCounters{
		rxErrors:  s.RxErrors,
		rxCrc:     s.RxCrcErrors,
		rxFrame:   s.RxFrameErrors,
		rxLength:  s.RxLengthErrors,
		rxDropped: s.RxDropped,
		rxMissed:  s.RxMissedErrors,
		rxOver:    s.RxOverErrors,
		rxFifo:    s.RxFifoErrors,
		txErrors:  s.TxErrors,
		txCarrier: s.TxCarrierErrors,
		txDropped: s.TxDropped,
	}, nil
}

// directionResult is the outcome of the load test in one direction.
type directionResult struct {
	from, to  string
	frameSize int
	duration  time.Duration
	// sent is the number of frames handed over to the sending interface
	sent uint64
	// received is the number of valid frames seen on the receiving interface
	received uint64
	// corrupt is the number of frames which arrived with a wrong checksum
	corrupt uint64
	// socketDrops is the number of frames dropped because the receiver was too slow
	socketDrops uint64
	// sendErr is set if sending was aborted, e.g. because the link went down
	sendErr error
	// tx are the counter changes of the sending, rx of the receiving interface
	tx, rx nicCounters
}

// localDrops are frames which made it over the wire but were dropped on this host because it could not keep up.
func (r directionResult) localDrops() uint64 {
	return r.socketDrops + r.rx.rxDropped + r.rx.rxMissed + r.rx.rxOver + r.rx.rxFifo
}

// lost returns the number of frames which got lost on the wire.
func (r directionResult) lost() uint64 {
	arrived := r.received + r.corrupt + r.localDrops()
	if arrived >= r.sent {
		return 0
	}
	return r.sent - arrived
}

func (r directionResult) lossRatio() float64 {
	if r.sent == 0 {
		return 0
	}
	return float64(r.lost()) / float64(r.sent)
}

func (r directionResult) gbitPerSecond() float64 {
	if r.duration <= 0 {
		return 0
	}
	return float64(r.received) * float64(r.frameSize) * 8 / r.duration.Seconds() / 1e9
}

// evaluate returns an error wrapping ErrUplinkLoadTest if the result indicates a faulty link.
func (r directionResult) evaluate() error {
	var reasons []string

	if r.sendErr != nil {
		reasons = append(reasons, fmt.Sprintf("sending aborted: %s", r.sendErr))
	}
	if r.sent == 0 {
		reasons = append(reasons, "no frames sent")
	} else if r.received == 0 {
		reasons = append(reasons, fmt.Sprintf("none of %d frames received", r.sent))
	} else if r.lossRatio() > loadTestMaxLossRatio {
		reasons = append(reasons, fmt.Sprintf("%d of %d frames lost (%.3f%%)", r.lost(), r.sent, r.lossRatio()*100))
	}
	if r.corrupt > 0 {
		reasons = append(reasons, fmt.Sprintf("%d corrupt frames received", r.corrupt))
	}
	if n := r.rx.rxCrc + r.rx.rxFrame + r.rx.rxLength; n > 0 {
		reasons = append(reasons, fmt.Sprintf("%d receive errors on %s (crc:%d frame:%d length:%d)", n, r.to, r.rx.rxCrc, r.rx.rxFrame, r.rx.rxLength))
	}
	if r.tx.txErrors > 0 || r.tx.txCarrier > 0 {
		reasons = append(reasons, fmt.Sprintf("transmit errors on %s (errors:%d carrier:%d)", r.from, r.tx.txErrors, r.tx.txCarrier))
	}

	if len(reasons) == 0 {
		return nil
	}
	return fmt.Errorf("%w: %s -> %s: %s", ErrUplinkLoadTest, r.from, r.to, strings.Join(reasons, ", "))
}

// newLoadTestFrame creates a frame with a total size of frameSize bytes.
// The returned checksum of the payload is required to seal the frame.
func newLoadTestFrame(dst, src net.HardwareAddr, kind byte, runID uint64, frameSize int) ([]byte, uint32) {
	if frameSize < loadTestPayloadOffset+loadTestMinPayload {
		frameSize = loadTestPayloadOffset + loadTestMinPayload
	}
	frame := make([]byte, frameSize)
	copy(frame[0:6], dst)
	copy(frame[6:12], src)
	binary.BigEndian.PutUint16(frame[12:14], loadTestEtherType)
	binary.BigEndian.PutUint32(frame[14:18], loadTestMagic)
	frame[18] = kind
	payload := frame[loadTestPayloadOffset:]
	binary.BigEndian.PutUint16(frame[20:22], uint16(len(payload))) // nolint:gosec
	binary.BigEndian.PutUint64(frame[22:30], runID)

	// random data, because a constant pattern is to easy to transmit
	var seed [32]byte
	binary.BigEndian.PutUint64(seed[:], runID)
	_, _ = rand.NewChaCha8(seed).Read(payload)

	return frame, crc32.ChecksumIEEE(payload)
}

// sealLoadTestFrame sets the sequence number and the checksum of the frame.
func sealLoadTestFrame(frame []byte, seq uint64, payloadCRC uint32) {
	binary.BigEndian.PutUint64(frame[30:38], seq)
	binary.BigEndian.PutUint32(frame[38:42], crc32.ChecksumIEEE(frame[loadTestEthHeaderLen:38])^payloadCRC)
}

type frameStatus int

const (
	// frameForeign is a frame which does not belong to this load test run
	frameForeign frameStatus = iota
	frameCorrupt
	frameValid
)

// parseLoadTestFrame checks if the given frame was sent from peer to self during the run with the given id.
func parseLoadTestFrame(frame []byte, self, peer net.HardwareAddr, runID uint64) (frameStatus, byte) {
	if len(frame) < loadTestPayloadOffset {
		return frameForeign, 0
	}
	if !bytes.Equal(frame[0:6], self) || !bytes.Equal(frame[6:12], peer) {
		return frameForeign, 0
	}
	if binary.BigEndian.Uint16(frame[12:14]) != loadTestEtherType || binary.BigEndian.Uint32(frame[14:18]) != loadTestMagic {
		return frameForeign, 0
	}
	if binary.BigEndian.Uint64(frame[22:30]) != runID {
		return frameForeign, 0
	}
	kind := frame[18]
	end := loadTestPayloadOffset + int(binary.BigEndian.Uint16(frame[20:22]))
	if end > len(frame) {
		return frameCorrupt, kind
	}
	crc := crc32.ChecksumIEEE(frame[loadTestEthHeaderLen:38]) ^ crc32.ChecksumIEEE(frame[loadTestPayloadOffset:end])
	if crc != binary.BigEndian.Uint32(frame[38:42]) {
		return frameCorrupt, kind
	}
	return frameValid, kind
}

// loadTestPort is one side of the load test.
type loadTestPort struct {
	name string
	mac  net.HardwareAddr
	mtu  int
	// separate sockets are used, because a socket created without protocol never receives anything.
	tx, rx int

	// set by the receiver
	smallProbe atomic.Bool
	fullProbe  atomic.Bool
	received   atomic.Uint64
	corrupt    atomic.Uint64
}

func htons(v uint16) uint16 {
	return v<<8 | v>>8
}

func openLoadTestPort(name string) (*loadTestPort, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return nil, fmt.Errorf("unable to lookup interface %s %w", name, err)
	}
	attrs := link.Attrs()
	if len(attrs.HardwareAddr) != 6 {
		return nil, fmt.Errorf("interface %s has no ethernet address", name)
	}

	p := &loadTestPort{name: name, mac: attrs.HardwareAddr, mtu: attrs.MTU, tx: -1, rx: -1}

	p.tx, err = unix.Socket(unix.AF_PACKET, unix.SOCK_RAW|unix.SOCK_CLOEXEC, 0)
	if err != nil {
		return nil, fmt.Errorf("unable to create send socket for %s %w", name, err)
	}
	err = unix.Bind(p.tx, &unix.SockaddrLinklayer{Ifindex: attrs.Index})
	if err != nil {
		p.close()
		return nil, fmt.Errorf("unable to bind send socket to %s %w", name, err)
	}
	snd := unix.NsecToTimeval(loadTestSndTimeout.Nanoseconds())
	err = unix.SetsockoptTimeval(p.tx, unix.SOL_SOCKET, unix.SO_SNDTIMEO, &snd)
	if err != nil {
		p.close()
		return nil, fmt.Errorf("unable to set send timeout on %s %w", name, err)
	}

	// created without protocol and bound to interface and protocol afterwards,
	// otherwise frames from all interfaces are received until bind is called.
	p.rx, err = unix.Socket(unix.AF_PACKET, unix.SOCK_RAW|unix.SOCK_CLOEXEC, 0)
	if err != nil {
		p.close()
		return nil, fmt.Errorf("unable to create receive socket for %s %w", name, err)
	}
	// best effort, the receiver ignores frames sent by this interface anyway.
	_ = unix.SetsockoptInt(p.rx, unix.SOL_PACKET, unix.PACKET_IGNORE_OUTGOING, 1)
	err = unix.SetsockoptInt(p.rx, unix.SOL_SOCKET, unix.SO_RCVBUFFORCE, loadTestRcvBuf)
	if err != nil {
		_ = unix.SetsockoptInt(p.rx, unix.SOL_SOCKET, unix.SO_RCVBUF, loadTestRcvBuf)
	}
	rcv := unix.NsecToTimeval(loadTestRcvTimeout.Nanoseconds())
	err = unix.SetsockoptTimeval(p.rx, unix.SOL_SOCKET, unix.SO_RCVTIMEO, &rcv)
	if err != nil {
		p.close()
		return nil, fmt.Errorf("unable to set receive timeout on %s %w", name, err)
	}
	err = unix.Bind(p.rx, &unix.SockaddrLinklayer{Protocol: htons(loadTestEtherType), Ifindex: attrs.Index})
	if err != nil {
		p.close()
		return nil, fmt.Errorf("unable to bind receive socket to %s %w", name, err)
	}

	return p, nil
}

func (p *loadTestPort) close() {
	if p.tx >= 0 {
		_ = unix.Close(p.tx)
	}
	if p.rx >= 0 {
		_ = unix.Close(p.rx)
	}
}

// socketDrops returns the number of frames dropped since the last call because the receive buffer was full.
func (p *loadTestPort) socketDrops() uint64 {
	stats, err := unix.GetsockoptTpacketStats(p.rx, unix.SOL_PACKET, unix.PACKET_STATISTICS)
	if err != nil {
		return 0
	}
	return uint64(stats.Drops)
}

// receive consumes frames sent by peer until stop is set.
func (p *loadTestPort) receive(peer net.HardwareAddr, runID uint64, stop *atomic.Bool) {
	buf := make([]byte, loadTestEthHeaderLen+p.mtu+64)
	for !stop.Load() {
		n, err := unix.Read(p.rx, buf)
		if err != nil || n <= 0 {
			// timeout or interrupted
			continue
		}
		status, kind := parseLoadTestFrame(buf[:n], p.mac, peer, runID)
		switch {
		case status == frameForeign:
		case kind == loadTestKindProbe && status == frameValid:
			if n > loadTestEthHeaderLen+loadTestStandardMTU {
				p.fullProbe.Store(true)
			} else {
				p.smallProbe.Store(true)
			}
		case kind == loadTestKindLoad && status == frameValid:
			p.received.Add(1)
		case kind == loadTestKindLoad:
			p.corrupt.Add(1)
		}
	}
}

func isTemporarySendError(err error) bool {
	return errors.Is(err, unix.ENOBUFS) || errors.Is(err, unix.EAGAIN) || errors.Is(err, unix.EINTR)
}

// send sends frames as fast as possible to peer until stop is set and returns the number of frames sent.
func (p *loadTestPort) send(peer net.HardwareAddr, runID uint64, frameSize int, stop *atomic.Bool) (uint64, error) {
	frame, payloadCRC := newLoadTestFrame(peer, p.mac, loadTestKindLoad, runID, frameSize)
	var (
		sent     uint64
		failing  time.Time
		lastErr  error
		maxRetry = 2 * loadTestSndTimeout
	)
	for !stop.Load() {
		sealLoadTestFrame(frame, sent, payloadCRC)
		_, err := unix.Write(p.tx, frame)
		if err == nil {
			sent++
			failing = time.Time{}
			continue
		}
		if !isTemporarySendError(err) {
			return sent, err
		}
		// the transmit queue is full
		if failing.IsZero() {
			failing = time.Now()
		} else if time.Since(failing) > maxRetry {
			lastErr = fmt.Errorf("transmit queue of %s is stuck %w", p.name, err)
			break
		}
		time.Sleep(50 * time.Microsecond)
	}
	return sent, lastErr
}

func (p *loadTestPort) sendProbe(peer net.HardwareAddr, runID uint64, frameSize int) {
	frame, payloadCRC := newLoadTestFrame(peer, p.mac, loadTestKindProbe, runID, frameSize)
	sealLoadTestFrame(frame, 0, payloadCRC)
	_, _ = unix.Write(p.tx, frame)
}

// probe checks if frames get through in both directions, which also lets the switches learn the mac addresses.
// It returns the frame size to use for the load test.
func probeLoadTestPair(log *slog.Logger, a, b *loadTestPort, runID uint64) (int, error) {
	mtu := min(a.mtu, b.mtu)
	fullSize := loadTestEthHeaderLen + mtu
	jumbo := mtu > loadTestStandardMTU

	var (
		start    = time.Now()
		smallAt  time.Time
		smallOK  bool
		fullOK   bool
		received = func(get func(p *loadTestPort) bool) bool { return get(a) && get(b) }
	)
	for time.Since(start) < loadTestProbeTimeout {
		for _, d := range [][2]*loadTestPort{{a, b}, {b, a}} {
			d[0].sendProbe(d[1].mac, runID, 0)
			if jumbo {
				d[0].sendProbe(d[1].mac, runID, fullSize)
			}
		}
		time.Sleep(loadTestProbeInterval)

		smallOK = received(func(p *loadTestPort) bool { return p.smallProbe.Load() })
		fullOK = received(func(p *loadTestPort) bool { return p.fullProbe.Load() })
		if smallOK && smallAt.IsZero() {
			smallAt = time.Now()
		}
		if smallOK && (!jumbo || fullOK || time.Since(smallAt) > loadTestJumboGrace) {
			break
		}
	}

	if !smallOK {
		return 0, fmt.Errorf("%w: no layer 2 connectivity between %s and %s within %s (%s received:%t, %s received:%t)",
			ErrUplinkLoadTest, a.name, b.name, loadTestProbeTimeout, a.name, a.smallProbe.Load(), b.name, b.smallProbe.Load())
	}
	if jumbo && !fullOK {
		// this is not a problem of the link of this machine, but of the network in between.
		log.Warn("uplink load test: frames with full mtu do not get through, falling back to standard mtu", "from", a.name, "to", b.name, "mtu", mtu)
		return loadTestEthHeaderLen + loadTestStandardMTU, nil
	}
	return fullSize, nil
}

// runLoadTestPair sends frames between the two given interfaces in both directions at the same time.
func runLoadTestPair(log *slog.Logger, nameA, nameB string, duration time.Duration) ([]directionResult, error) {
	a, err := openLoadTestPort(nameA)
	if err != nil {
		return nil, err
	}
	defer a.close()
	b, err := openLoadTestPort(nameB)
	if err != nil {
		return nil, err
	}
	defer b.close()

	runID := rand.Uint64()

	var (
		stopReceive atomic.Bool
		receivers   sync.WaitGroup
	)
	receivers.Go(func() { a.receive(b.mac, runID, &stopReceive) })
	receivers.Go(func() { b.receive(a.mac, runID, &stopReceive) })
	stopReceivers := func() {
		stopReceive.Store(true)
		receivers.Wait()
	}

	frameSize, err := probeLoadTestPair(log, a, b, runID)
	if err != nil {
		stopReceivers()
		return nil, err
	}

	// reading resets the drop counters
	_ = a.socketDrops()
	_ = b.socketDrops()
	beforeA, errA := readNicCounters(a.name)
	beforeB, errB := readNicCounters(b.name)
	if err := errors.Join(errA, errB); err != nil {
		stopReceivers()
		return nil, err
	}

	var (
		stopSend       atomic.Bool
		senders        sync.WaitGroup
		sentA, sentB   uint64
		sendErrA       error
		sendErrB       error
		elapsedA       time.Duration
		elapsedB       time.Duration
		start          = time.Now()
		sendAndMeasure = func(from, to *loadTestPort, sent *uint64, sendErr *error, elapsed *time.Duration) {
			*sent, *sendErr = from.send(to.mac, runID, frameSize, &stopSend)
			*elapsed = time.Since(start)
		}
	)
	senders.Go(func() { sendAndMeasure(a, b, &sentA, &sendErrA, &elapsedA) })
	senders.Go(func() { sendAndMeasure(b, a, &sentB, &sendErrB, &elapsedB) })

	time.Sleep(duration)
	stopSend.Store(true)
	senders.Wait()
	time.Sleep(loadTestDrain)
	stopReceivers()

	afterA, errA := readNicCounters(a.name)
	afterB, errB := readNicCounters(b.name)
	if err := errors.Join(errA, errB); err != nil {
		return nil, err
	}
	deltaA := afterA.sub(beforeA)
	deltaB := afterB.sub(beforeB)

	return []directionResult{
		{
			from: a.name, to: b.name, frameSize: frameSize, duration: elapsedA,
			sent: sentA, received: b.received.Load(), corrupt: b.corrupt.Load(), socketDrops: b.socketDrops(),
			sendErr: sendErrA, tx: deltaA, rx: deltaB,
		},
		{
			from: b.name, to: a.name, frameSize: frameSize, duration: elapsedB,
			sent: sentB, received: a.received.Load(), corrupt: a.corrupt.Load(), socketDrops: a.socketDrops(),
			sendErr: sendErrB, tx: deltaB, rx: deltaA,
		},
	}, nil
}

// uplinks returns all interfaces with a lldp neighbor together with a description of this neighbor.
func (n *Network) uplinks() map[string]string {
	host := n.LLDPClient.Host
	host.mutex.RLock()
	defer host.mutex.RUnlock()

	uplinks := make(map[string]string)
	for name, neighbors := range host.neighbors {
		if len(neighbors) == 0 {
			continue
		}
		uplinks[name] = fmt.Sprintf("%s connected to %s port %s", name, neighbors[0].Name, neighbors[0].PortDescription)
	}
	return uplinks
}

// UplinkLoadTest transfers data between the uplinks of this machine for the given duration to detect faulty links.
// As all uplinks are in the same layer 2 network, frames sent from one uplink to the mac address of another uplink
// have to pass both links, their GBICs and the switches in between.
//
// It must be called after the lldp neighbors were detected.
// The returned error wraps ErrUplinkLoadTest if a faulty link was detected.
func (n *Network) UplinkLoadTest(duration time.Duration) error {
	uplinks := n.uplinks()
	names := make([]string, 0, len(uplinks))
	for name := range uplinks {
		names = append(names, name)
	}
	slices.Sort(names)
	if len(names) < 2 {
		return fmt.Errorf("uplink load test requires at least two uplinks, found:%v", names)
	}

	var errs []error
	for i := 0; i+1 < len(names); i++ {
		a, b := names[i], names[i+1]
		n.Log.Info("uplink load test start", "from", uplinks[a], "to", uplinks[b], "duration", duration.String())

		results, err := runLoadTestPair(n.Log, a, b, duration)
		if err != nil {
			errs = append(errs, fmt.Errorf("%w (%s, %s)", err, uplinks[a], uplinks[b]))
			continue
		}

		for _, r := range results {
			n.Log.Info("uplink load test result",
				"from", uplinks[r.from], "to", uplinks[r.to],
				"framesize", r.frameSize, "sent", r.sent, "received", r.received, "corrupt", r.corrupt,
				"lost", r.lost(), "localdrops", r.localDrops(), "gbit/s", fmt.Sprintf("%.2f", r.gbitPerSecond()),
				"rx-errors", r.rx.rxErrors, "rx-crc-errors", r.rx.rxCrc, "rx-frame-errors", r.rx.rxFrame, "rx-length-errors", r.rx.rxLength,
				"tx-errors", r.tx.txErrors, "tx-carrier-errors", r.tx.txCarrier, "tx-dropped", r.tx.txDropped,
			)
			err := r.evaluate()
			if err != nil {
				errs = append(errs, fmt.Errorf("%w (%s, %s)", err, uplinks[r.from], uplinks[r.to]))
			}
		}
	}

	return errors.Join(errs...)
}
