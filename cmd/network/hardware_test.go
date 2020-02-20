package network

import (
	"github.com/metal-stack/metal-hammer/metal-core/models"
	"reflect"
	"testing"
)

func Test_readHardwareDetails(t *testing.T) {
	// FIXME test is disabled
	t.Skip()
	tests := []struct {
		name    string
		want    *models.DomainMetalHammerRegisterMachineRequest
		wantErr bool
	}{
		{
			name: "simple",
			want: &models.DomainMetalHammerRegisterMachineRequest{
				UUID: "00000000-0000-0000-0000-000000000000",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Network{
				MachineUUID: "00000000-0000-0000-0000-000000000000",
				Eth0Mac:     "00:00:00:00:00:01",
			}

			got, err := n.ReadHardwareDetails()
			if (err != nil) != tt.wantErr {
				t.Errorf("readHardwareDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got.UUID) != len(tt.want.UUID) {
				t.Errorf("readHardwareDetails() expected uuid: %s got %s", tt.want.UUID, got.UUID)
			}
			if *got.CPUCores == 0 {
				t.Errorf("readHardwareDetails() expected cpucores: %d", got.CPUCores)
			}
			if *got.Memory == 0 {
				t.Errorf("readHardwareDetails() expected memory: %d", got.Memory)
			}
		})
	}
}

func TestHammer_readIPMIDetails(t *testing.T) {
	tests := []struct {
		name    string
		want    *models.ModelsV1MachineIPMI
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readIPMIDetails("00:00:00:00:00:01")
			if (err != nil) != tt.wantErr {
				t.Errorf("Hammer.readIPMIDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hammer.readIPMIDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}
