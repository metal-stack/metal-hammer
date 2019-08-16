package uuid

import (
	"testing"
)

func TestMachineUUID(t *testing.T) {

	mockedioutilReadFileMock := func(filename string) ([]byte, error) {
		return []byte("4C4C4544-0042-4810-8056-B4C04F395332"), nil
	}

	tests := []struct {
		name string
		want string
	}{
		{
			name: "TestMachineUUID Test 1",
			want: "4C4C4544-0042-4810-8056-B4C04F395332",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := machineUUIDMockable(mockedioutilReadFileMock); got != tt.want {
				t.Errorf("MachineUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_MachineUUID(t *testing.T) {
	// already done in TestMachineUUID
}

func Test_machineUUIDMockable(t *testing.T) {
	// Already done in TestMachineUUID
}
