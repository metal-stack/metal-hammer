package network

import (
	"errors"
	"log/slog"
	"net"
	"os"
	"testing"
	"time"

	"github.com/vishvananda/netlink"
)

func TestLoadTestFrame(t *testing.T) {
	var (
		self  = net.HardwareAddr{0x02, 0, 0, 0, 0, 0x01}
		peer  = net.HardwareAddr{0x02, 0, 0, 0, 0, 0x02}
		other = net.HardwareAddr{0x02, 0, 0, 0, 0, 0x03}
		runID = uint64(4711)
	)

	frame, payloadCRC := newLoadTestFrame(self, peer, loadTestKindLoad, runID, 9014)
	if len(frame) != 9014 {
		t.Fatalf("unexpected frame size %d", len(frame))
	}
	sealLoadTestFrame(frame, 42, payloadCRC)

	status, kind := parseLoadTestFrame(frame, self, peer, runID)
	if status != frameValid || kind != loadTestKindLoad {
		t.Errorf("expected valid load frame, got status:%d kind:%d", status, kind)
	}

	// a small frame is padded by the nic
	small, smallCRC := newLoadTestFrame(self, peer, loadTestKindProbe, runID, 0)
	sealLoadTestFrame(small, 0, smallCRC)
	status, kind = parseLoadTestFrame(append(small, 0, 0, 0, 0), self, peer, runID)
	if status != frameValid || kind != loadTestKindProbe {
		t.Errorf("expected valid probe frame, got status:%d kind:%d", status, kind)
	}

	if status, _ := parseLoadTestFrame(frame, self, peer, runID+1); status != frameForeign {
		t.Errorf("expected frame of other run to be foreign, got %d", status)
	}
	if status, _ := parseLoadTestFrame(frame, self, other, runID); status != frameForeign {
		t.Errorf("expected frame of other sender to be foreign, got %d", status)
	}
	if status, _ := parseLoadTestFrame(frame, peer, self, runID); status != frameForeign {
		t.Errorf("expected own frame to be foreign, got %d", status)
	}
	if status, _ := parseLoadTestFrame(frame[:20], self, peer, runID); status != frameForeign {
		t.Errorf("expected short frame to be foreign, got %d", status)
	}

	if status, _ := parseLoadTestFrame(frame[:5000], self, peer, runID); status != frameCorrupt {
		t.Errorf("expected truncated frame to be corrupt, got %d", status)
	}
	frame[5000] ^= 0x01
	if status, _ := parseLoadTestFrame(frame, self, peer, runID); status != frameCorrupt {
		t.Errorf("expected frame with flipped payload bit to be corrupt, got %d", status)
	}
	frame[5000] ^= 0x01
	frame[37] ^= 0x01
	if status, _ := parseLoadTestFrame(frame, self, peer, runID); status != frameCorrupt {
		t.Errorf("expected frame with flipped sequence bit to be corrupt, got %d", status)
	}
}

func TestDirectionResultEvaluate(t *testing.T) {
	tests := []struct {
		name    string
		result  directionResult
		wantErr bool
	}{
		{
			name:   "clean run",
			result: directionResult{sent: 100000, received: 100000},
		},
		{
			name:   "loss below threshold",
			result: directionResult{sent: 100000, received: 99950},
		},
		{
			name:    "loss above threshold",
			result:  directionResult{sent: 100000, received: 99000},
			wantErr: true,
		},
		{
			name:   "drops of a slow receiver are not a loss",
			result: directionResult{sent: 100000, received: 60000, socketDrops: 30000, rx: nicCounters{rxMissed: 10000}},
		},
		{
			name:    "loss on top of drops of a slow receiver",
			result:  directionResult{sent: 100000, received: 60000, socketDrops: 30000},
			wantErr: true,
		},
		{
			name:    "nothing received",
			result:  directionResult{sent: 100000},
			wantErr: true,
		},
		{
			name:    "nothing sent",
			result:  directionResult{},
			wantErr: true,
		},
		{
			name:    "corrupt frame",
			result:  directionResult{sent: 100000, received: 99999, corrupt: 1},
			wantErr: true,
		},
		{
			name:    "crc errors",
			result:  directionResult{sent: 100000, received: 99999, rx: nicCounters{rxCrc: 1}},
			wantErr: true,
		},
		{
			name:    "transmit errors",
			result:  directionResult{sent: 100000, received: 100000, tx: nicCounters{txErrors: 1}},
			wantErr: true,
		},
		{
			name:   "errors in the other direction are ignored",
			result: directionResult{sent: 100000, received: 100000, tx: nicCounters{rxCrc: 1}, rx: nicCounters{txErrors: 1}},
		},
		{
			name:    "sending aborted",
			result:  directionResult{sent: 100000, received: 100000, sendErr: errors.New("network is down")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.result.from, tt.result.to = "eth0", "eth1"
			err := tt.result.evaluate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("evaluate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !errors.Is(err, ErrUplinkLoadTest) {
				t.Errorf("expected error to wrap ErrUplinkLoadTest, got %v", err)
			}
		})
	}
}

// TestRunLoadTestPair requires privileges to create interfaces, it can be executed without being root with:
//
//	unshare -Urn go test -run TestRunLoadTestPair -v ./cmd/network/
func TestRunLoadTestPair(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("requires root")
	}

	const nameA, nameB = "mhlt0", "mhlt1"
	veth := &netlink.Veth{LinkAttrs: netlink.NewLinkAttrs(), PeerName: nameB}
	veth.Name = nameA
	veth.MTU = MTU
	if err := netlink.LinkAdd(veth); err != nil {
		t.Skipf("unable to create veth pair: %v", err)
	}
	t.Cleanup(func() { _ = netlink.LinkDel(veth) })

	for _, name := range []string{nameA, nameB} {
		if err := linkSetMTU(name, MTU); err != nil {
			t.Fatal(err)
		}
		if err := linkSetUp(name); err != nil {
			t.Fatal(err)
		}
	}

	log := slog.New(slog.NewTextHandler(os.Stderr, nil))

	results, err := runLoadTestPair(log, nameA, nameB, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("expected results for two directions, got %d", len(results))
	}
	for _, r := range results {
		t.Logf("%s -> %s: framesize:%d sent:%d received:%d localdrops:%d lost:%d %.2f gbit/s", r.from, r.to, r.frameSize, r.sent, r.received, r.localDrops(), r.lost(), r.gbitPerSecond())
		if r.frameSize != loadTestEthHeaderLen+MTU {
			t.Errorf("expected jumbo frames, got frame size %d", r.frameSize)
		}
		if err := r.evaluate(); err != nil {
			t.Errorf("unexpected failure: %v", err)
		}
	}

	// now break one direction
	link, err := netlink.LinkByName(nameA)
	if err != nil {
		t.Fatal(err)
	}
	netem := netlink.NewNetem(
		netlink.QdiscAttrs{LinkIndex: link.Attrs().Index, Handle: netlink.MakeHandle(1, 0), Parent: netlink.HANDLE_ROOT},
		netlink.NetemQdiscAttrs{Loss: 5, Limit: 100000},
	)
	if err := netlink.QdiscAdd(netem); err != nil {
		t.Skipf("unable to add netem qdisc: %v", err)
	}

	results, err = runLoadTestPair(log, nameA, nameB, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		t.Logf("%s -> %s: framesize:%d sent:%d received:%d localdrops:%d lost:%d %.2f gbit/s", r.from, r.to, r.frameSize, r.sent, r.received, r.localDrops(), r.lost(), r.gbitPerSecond())
		err := r.evaluate()
		if r.from == nameA && !errors.Is(err, ErrUplinkLoadTest) {
			t.Errorf("expected a failure for the direction with packet loss, got %v", err)
		}
		if r.from == nameB && err != nil {
			t.Errorf("unexpected failure for the direction without packet loss: %v", err)
		}
	}
}
