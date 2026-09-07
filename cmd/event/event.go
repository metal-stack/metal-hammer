package event

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	infrav2 "github.com/metal-stack/api/go/metalstack/infra/v2"
	"github.com/metal-stack/api/go/metalstack/infra/v2/infrav2connect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProvisioningEventType indicates an event emitted by a machine during the provisioning sequence
// FIXME factor out to metal-lib
type ProvisioningEventType string

// The enums for the machine provisioning events.
const (
	ProvisioningEventAlive            ProvisioningEventType = "Alive"
	ProvisioningEventCrashed          ProvisioningEventType = "Crashed"
	ProvisioningEventResetFailCount   ProvisioningEventType = "Reset Fail Count"
	ProvisioningEventPXEBooting       ProvisioningEventType = "PXE Booting"
	ProvisioningEventPlannedReboot    ProvisioningEventType = "Planned Reboot"
	ProvisioningEventPreparing        ProvisioningEventType = "Preparing"
	ProvisioningEventRegistering      ProvisioningEventType = "Registering"
	ProvisioningEventWaiting          ProvisioningEventType = "Waiting"
	ProvisioningEventInstalling       ProvisioningEventType = "Installing"
	ProvisioningEventBootingNewKernel ProvisioningEventType = "Booting New Kernel"
	ProvisioningEventPhonedHome       ProvisioningEventType = "Phoned Home"
)

type EventEmitter struct {
	log         *slog.Logger
	eventClient infrav2connect.EventServiceClient
	machineID   string
}

func NewEventEmitter(log *slog.Logger, eventClient infrav2connect.EventServiceClient, machineID string) *EventEmitter {
	emitter := &EventEmitter{
		eventClient: eventClient,
		machineID:   machineID,
		log:         log,
	}

	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for t := range ticker.C {
			emitter.Emit(apiv2.MachineProvisioningEventType_MACHINE_PROVISIONING_EVENT_TYPE_ALIVE, fmt.Sprintf("still alive at: %s", t))
		}
	}()
	return emitter
}

func (e *EventEmitter) Emit(eventType apiv2.MachineProvisioningEventType, message string) {
	e.log.Info("event", "event", eventType, "message", message)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s, err := e.eventClient.Send(ctx, &infrav2.EventServiceSendRequest{
		Events: map[string]*apiv2.MachineProvisioningEvent{
			e.machineID: {
				Time:    timestamppb.Now(),
				Event:   eventType,
				Message: message,
			},
		},
	})
	if err != nil {
		e.log.Error("event", "cannot send event", eventType, "error", err)
	}
	if s != nil {
		e.log.Info("event", "send", s.Events, "failed", s.Failed)
	}
}
