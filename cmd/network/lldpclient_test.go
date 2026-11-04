package network

import (
	"log/slog"
	"testing"

	"github.com/metal-stack/metal-hammer/pkg/lldp"
)

func TestParseLLDPRequirements(t *testing.T) {
	tests := []struct {
		name    string
		cmdline map[string]string
		want    LLDPRequirements
		wantErr bool
	}{
		{name: "default is two interfaces and two switches", want: LLDPRequirements{MinInterfaces: 2, MinNeighbors: 2}},
		{name: "unrelated keys are ignored", cmdline: map[string]string{"BGP": "1", "console": "ttyS0"}, want: DefaultLLDPRequirements},
		{
			name:    "single switch lab",
			cmdline: map[string]string{CmdlineLLDPMinInterfaces: "1", CmdlineLLDPMinNeighbors: "1"},
			want:    LLDPRequirements{MinInterfaces: 1, MinNeighbors: 1},
		},
		{
			name:    "two cables to one switch",
			cmdline: map[string]string{CmdlineLLDPMinNeighbors: "1"},
			want:    LLDPRequirements{MinInterfaces: 2, MinNeighbors: 1},
		},
		{name: "zero would skip discovery", cmdline: map[string]string{CmdlineLLDPMinInterfaces: "0"}, wantErr: true},
		{name: "negative", cmdline: map[string]string{CmdlineLLDPMinNeighbors: "-1"}, wantErr: true},
		{name: "not a number", cmdline: map[string]string{CmdlineLLDPMinNeighbors: "one"}, wantErr: true},
		{name: "empty value", cmdline: map[string]string{CmdlineLLDPMinInterfaces: ""}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLLDPRequirements(tt.cmdline)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if got != DefaultLLDPRequirements {
					t.Errorf("on error got %+v, want the default", got)
				}
				return
			}
			if got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestRequirementsMet(t *testing.T) {
	neigh := func(iface, chassis, port string) lldp.Neighbor {
		return lldp.Neighbor{
			Interface: iface,
			Chassis:   lldp.Chassis{Type: lldp.Mac, Value: chassis},
			Port:      lldp.Port{Type: lldp.Local, Value: port},
		}
	}
	const leaf1, leaf2 = "aa:aa:aa:aa:aa:01", "aa:aa:aa:aa:aa:02"

	tests := []struct {
		name      string
		req       LLDPRequirements
		neighbors []lldp.Neighbor
		want      bool
	}{
		{name: "default: nothing seen", req: DefaultLLDPRequirements, want: false},
		{
			name:      "default: dual-homed to two leaves",
			req:       DefaultLLDPRequirements,
			neighbors: []lldp.Neighbor{neigh("eth0", leaf1, "swp1"), neigh("eth1", leaf2, "swp1")},
			want:      true,
		},
		{
			name:      "default: one interface only",
			req:       DefaultLLDPRequirements,
			neighbors: []lldp.Neighbor{neigh("eth0", leaf1, "swp1")},
			want:      false,
		},
		{
			name:      "default: two cables to the same leaf",
			req:       DefaultLLDPRequirements,
			neighbors: []lldp.Neighbor{neigh("eth0", leaf1, "swp1"), neigh("eth1", leaf1, "swp2")},
			want:      false,
		},
		{
			name: "single switch: still waits for a neighbor",
			req:  LLDPRequirements{MinInterfaces: 1, MinNeighbors: 1},
			want: false,
		},
		{
			name:      "single switch: one neighbor",
			req:       LLDPRequirements{MinInterfaces: 1, MinNeighbors: 1},
			neighbors: []lldp.Neighbor{neigh("eth0", leaf1, "swp1")},
			want:      true,
		},
		{
			name:      "two cables to one switch: met",
			req:       LLDPRequirements{MinInterfaces: 2, MinNeighbors: 1},
			neighbors: []lldp.Neighbor{neigh("eth0", leaf1, "swp1"), neigh("eth1", leaf1, "swp2")},
			want:      true,
		},
		{
			name: "single switch: neighbor without mac chassis does not count",
			req:  LLDPRequirements{MinInterfaces: 1, MinNeighbors: 1},
			neighbors: []lldp.Neighbor{{
				Interface: "eth0",
				Chassis:   lldp.Chassis{Type: lldp.Interface, Value: "leaf1"},
				Port:      lldp.Port{Type: lldp.Local, Value: "swp1"},
			}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewLLDPClient(slog.Default(), []string{"eth0", "eth1"}, tt.req.MinInterfaces, tt.req.MinNeighbors, 0)
			for _, n := range tt.neighbors {
				c.addNeighbor(n)
			}
			if got := c.requirementsMet(); got != tt.want {
				t.Errorf("requirementsMet() = %v, want %v", got, tt.want)
			}
			if c.Host.done != tt.want && len(tt.neighbors) > 0 {
				t.Errorf("Host.done = %v, want %v", c.Host.done, tt.want)
			}
		})
	}
}
