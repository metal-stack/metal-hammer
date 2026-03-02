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
				{Name: new("eth0")},
				{Name: new("eth1")},
				{Name: new("eth2"), Mac: new("aa:bb"), Neighbors: []*models.V1MachineNic{{Name: new("swp1"), Mac: new("cc:dd")}}},
				{Name: new("eth3"), Mac: new("aa:bc"), Neighbors: []*models.V1MachineNic{{Name: new("swp2"), Mac: new("cc:de")}}},
			},
			want: []*models.V1MachineNic{
				{Name: new("eth2"), Mac: new("aa:bb"), Neighbors: []*models.V1MachineNic{{Name: new("swp1"), Mac: new("cc:dd")}}},
				{Name: new("eth3"), Mac: new("aa:bc"), Neighbors: []*models.V1MachineNic{{Name: new("swp2"), Mac: new("cc:de")}}},
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
			want: []*models.V1MachineNic{
				{Name: new("eth2"), Mac: new("aa:bb"), Neighbors: []*models.V1MachineNic{{Name: new("swp1"), Mac: new("cc:dd")}}},
				{Name: new("eth3"), Mac: new("aa:bc"), Neighbors: []*models.V1MachineNic{{Name: new("swp2"), Mac: new("cc:de")}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &hammer{
				log: slog.Default(),
			}
			if got := h.onlyNicsWithNeighbors(tt.nics); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hammer.onlyNicsWithNeighbors() = %v, want %v", got, tt.want)
			}
		})
	}
}

//go:fix inline
func ptr(s string) *string {
	return new(s)
}
