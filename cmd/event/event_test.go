package event

import (
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/metal-stack/metal-hammer/metal-core/client/machine"
	"go.uber.org/zap/zaptest"

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
				client:    nil,
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

func TestEventEmitter_Emit(t *testing.T) {
	transport := httptransport.New("metalcoreURL", "", nil)
	client := machine.New(transport, strfmt.Default)

	type args struct {
		eventType ProvisioningEventType
		message   string
	}
	tests := []struct {
		name string
		e    *EventEmitter
		args args
	}{
		{
			name: "TestEventEmitter_Emit Test 1",
			e:    NewEventEmitter(zaptest.NewLogger(t).Sugar(), client, "machineID"),
			args: args{
				eventType: "ProvisioningEventType",
				message:   "message",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.e.Emit(tt.args.eventType, tt.args.message)
		})
	}
}
