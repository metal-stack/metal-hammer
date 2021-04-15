package storage

import (
	"reflect"
	"testing"

	"github.com/metal-stack/metal-hammer/metal-core/models"
)

func TestGuessDisk(t *testing.T) {
	sde := "sde"
	sda := "sda"
	nvme := "nvme0n"
	gib64 := 64 * GIB
	tib10 := 10 * TIB
	tests := []struct {
		name  string
		disks []*models.ModelsV1MachineBlockDevice
		want  string
	}{
		{
			name: "working",
			disks: []*models.ModelsV1MachineBlockDevice{
				&models.ModelsV1MachineBlockDevice{
					Name: &sda,
					Size: &tib10,
				},
				&models.ModelsV1MachineBlockDevice{
					Name: &sde,
					Size: &gib64,
				},
				&models.ModelsV1MachineBlockDevice{
					Name: &nvme,
					Size: &gib64,
				},
			},
			want: "/dev/sde",
		},
		{
			name: "no guess possible",
			disks: []*models.ModelsV1MachineBlockDevice{
				&models.ModelsV1MachineBlockDevice{
					Name: &sda,
					Size: &tib10,
				},
				&models.ModelsV1MachineBlockDevice{
					Name: &nvme,
					Size: &gib64,
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := guessDisk(tt.disks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("guessDisk() = %v, want %v", got, tt.want)
			}
		})
	}
}
