package event

import (
	"reflect"
	"testing"

	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/client/machine"
)

func TestNewEventEmitter(t *testing.T) {
	type args struct {
		client    *machine.Client
		machineID string
	}
	tests := []struct {
		name string
		args args
		want *EventEmitter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEventEmitter(tt.args.client, tt.args.machineID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEventEmitter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventEmitter_Emit(t *testing.T) {
	type args struct {
		eventType ProvisioningEventType
		message   string
	}
	tests := []struct {
		name string
		e    *EventEmitter
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.Emit(tt.args.eventType, tt.args.message)
		})
	}
}
