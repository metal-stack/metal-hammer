package cmd

import (
	"log/slog"
	"reflect"
	"testing"

	"github.com/metal-stack/metal-go/api/models"
)

func TestHammer_onlyNicsWithNeighbors(t *testing.T) {

	tests := []struct {
		name string
		nics []*models.V1MachineNic
		want []*models.V1MachineNic
	}{
		{
			name: "4 interfaces, two with neighbors",
			nics: []*models.V1MachineNic{
				{Name: ptr("eth0")},
				{Name: ptr("eth1")},
				{Name: ptr("eth2"), Mac: ptr("aa:bb"), Neighbors: []*models.V1MachineNic{{Name: ptr("swp1"), Mac: ptr("cc:dd")}}},
				{Name: ptr("eth3"), Mac: ptr("aa:bc"), Neighbors: []*models.V1MachineNic{{Name: ptr("swp2"), Mac: ptr("cc:de")}}},
			},
			want: []*models.V1MachineNic{
				{Name: ptr("eth2"), Mac: ptr("aa:bb"), Neighbors: []*models.V1MachineNic{{Name: ptr("swp1"), Mac: ptr("cc:dd")}}},
				{Name: ptr("eth3"), Mac: ptr("aa:bc"), Neighbors: []*models.V1MachineNic{{Name: ptr("swp2"), Mac: ptr("cc:de")}}},
			},
		},
		{
			name: "4 interfaces, two with neighbors, one with empty Mac",
			nics: []*models.V1MachineNic{
				{Name: ptr("eth0")},
				{Name: ptr("eth1"), Mac: ptr("aa:bb"), Neighbors: []*models.V1MachineNic{{Name: ptr("swp1")}}},
				{Name: ptr("eth2"), Mac: ptr("aa:bb"), Neighbors: []*models.V1MachineNic{{Name: ptr("swp1"), Mac: ptr("cc:dd")}}},
				{Name: ptr("eth3"), Mac: ptr("aa:bc"), Neighbors: []*models.V1MachineNic{{Name: ptr("swp2"), Mac: ptr("cc:de")}}},
			},
			want: []*models.V1MachineNic{
				{Name: ptr("eth2"), Mac: ptr("aa:bb"), Neighbors: []*models.V1MachineNic{{Name: ptr("swp1"), Mac: ptr("cc:dd")}}},
				{Name: ptr("eth3"), Mac: ptr("aa:bc"), Neighbors: []*models.V1MachineNic{{Name: ptr("swp2"), Mac: ptr("cc:de")}}},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			h := &Hammer{
				log: slog.Default(),
			}
			if got := h.onlyNicsWithNeighbors(tt.nics); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hammer.onlyNicsWithNeighbors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ptr(s string) *string {
	return &s
}
