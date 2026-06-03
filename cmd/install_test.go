package cmd

import (
	"log/slog"
	"testing"

	"github.com/google/go-cmp/cmp"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-go/api/models"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestHammer_onlyNicsWithNeighbors(t *testing.T) {

	tests := []struct {
		name string
		nics []*models.V1MachineNic
		want []*apiv2.MachineNic
	}{
		{
			name: "4 interfaces, two with neighbors",
			nics: []*models.V1MachineNic{
				{Name: new("eth0")},
				{Name: new("eth1")},
				{Name: new("eth2"), Mac: new("aa:bb"), Neighbors: []*models.V1MachineNic{{Name: new("swp1"), Mac: new("cc:dd")}}},
				{Name: new("eth3"), Mac: new("aa:bc"), Neighbors: []*models.V1MachineNic{{Name: new("swp2"), Mac: new("cc:de")}}},
			},
			want: []*apiv2.MachineNic{
				{Name: "eth2", Mac: "aa:bb", Neighbors: []*apiv2.MachineNic{{Name: "swp1", Mac: "cc:dd"}}},
				{Name: "eth3", Mac: "aa:bc", Neighbors: []*apiv2.MachineNic{{Name: "swp2", Mac: "cc:de"}}},
			},
		},
		{
			name: "4 interfaces, two with neighbors, one with empty Mac",
			nics: []*models.V1MachineNic{
				{Name: new("eth0")},
				{Name: new("eth1"), Mac: new("aa:bb"), Neighbors: []*models.V1MachineNic{{Name: new("swp1")}}},
				{Name: new("eth2"), Mac: new("aa:bb"), Neighbors: []*models.V1MachineNic{{Name: new("swp1"), Mac: new("cc:dd")}}},
				{Name: new("eth3"), Mac: new("aa:bc"), Neighbors: []*models.V1MachineNic{{Name: new("swp2"), Mac: new("cc:de")}}},
			},
			want: []*apiv2.MachineNic{
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
			got := h.onlyNicsWithNeighbors(tt.nics)
			if diff := cmp.Diff(got, tt.want, protocmp.Transform()); diff != "" {
				t.Errorf("Hammer.onlyNicsWithNeighbors() diff = %s", diff)
			}
		})
	}
}
