package event

import (
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/metal-stack/metal-hammer/metal-core/client/machine"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/stretchr/testify/assert"
)

func TestNewEventEmitter(t *testing.T) {
	transport := httptransport.New("MetalCoreURL", "", nil)
	client := machine.New(transport, strfmt.Default)

	type args struct {
		client    machine.ClientService
		machineID string
	}
	tests := []struct {
		name string
		args args
		want *EventEmitter
	}{
		{
			name: "TestNewEventEmitter Test 1",
			args: args{
				client:    client,
				machineID: "machineID",
			},
			want: &EventEmitter{
				machineID: "machineID",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.args.machineID, tt.want.machineID, "check machine ID")
		})
	}
}
