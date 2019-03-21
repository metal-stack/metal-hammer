package uuid

import "testing"

func TestMachineUUID(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "TestMachineUUID Test 1",
			want: "00000000-0000-0000-0000-000000000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MachineUUID(); got != tt.want {
				t.Errorf("MachineUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}
