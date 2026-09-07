package cmd

import (
	"log/slog"
	"testing"

	"github.com/google/go-cmp/cmp"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	installerv1 "github.com/metal-stack/os-installer/api/v1"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestHammer_onlyNicsWithNeighbors(t *testing.T) {
	t.Skip()

	tests := []struct {
		name string
		nics []*apiv2.MachineNic
		want []*installerv1.V1MachineNic
	}{
		{
			name: "4 interfaces, two with neighbors",
			want: []*installerv1.V1MachineNic{
				{Name: new("eth0")},
				{Name: new("eth1")},
				{Name: new("eth2"), Mac: new("aa:bb"), Neighbors: []*installerv1.V1MachineNic{{Name: new("swp1"), Mac: new("cc:dd")}}},
				{Name: new("eth3"), Mac: new("aa:bc"), Neighbors: []*installerv1.V1MachineNic{{Name: new("swp2"), Mac: new("cc:de")}}},
			},
			nics: []*apiv2.MachineNic{
				{Name: "eth2", Mac: "aa:bb", Neighbors: []*apiv2.MachineNic{{Name: "swp1", Mac: "cc:dd"}}},
				{Name: "eth3", Mac: "aa:bc", Neighbors: []*apiv2.MachineNic{{Name: "swp2", Mac: "cc:de"}}},
			},
		},
		{
			name: "4 interfaces, two with neighbors, one with empty Mac",
			want: []*installerv1.V1MachineNic{
				{Name: new("eth0")},
				{Name: new("eth1"), Mac: new("aa:bb"), Neighbors: []*installerv1.V1MachineNic{{Name: new("swp1")}}},
				{Name: new("eth2"), Mac: new("aa:bb"), Neighbors: []*installerv1.V1MachineNic{{Name: new("swp1"), Mac: new("cc:dd")}}},
				{Name: new("eth3"), Mac: new("aa:bc"), Neighbors: []*installerv1.V1MachineNic{{Name: new("swp2"), Mac: new("cc:de")}}},
			},
			nics: []*apiv2.MachineNic{
				{Name: "eth2", Mac: "aa:bb", Neighbors: []*apiv2.MachineNic{{Name: "swp1", Mac: "cc:dd"}}},
				{Name: "eth3", Mac: "aa:bc", Neighbors: []*apiv2.MachineNic{{Name: "swp2", Mac: "cc:de"}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &hammer{
				log: slog.Default(),
			}
			got := h.onlyNicsWithNeighborsLegacy(tt.nics)
			if diff := cmp.Diff(got, tt.nics, protocmp.Transform()); diff != "" {
				t.Errorf("Hammer.onlyNicsWithNeighbors() diff = %s", diff)
			}
		})
	}
}
