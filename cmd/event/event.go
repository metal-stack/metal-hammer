package event

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	v1 "github.com/metal-stack/metal-api/pkg/api/v1"
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
	eventClient v1.EventServiceClient
	machineID   string
}

func NewEventEmitter(log *slog.Logger, eventClient v1.EventServiceClient, machineID string) *EventEmitter {
	emitter := &EventEmitter{
		eventClient: eventClient,
		machineID:   machineID,
		log:         log,
	}

	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for t := range ticker.C {
			emitter.Emit(ProvisioningEventAlive, fmt.Sprintf("still alive at: %s", t))
		}
	}()
	return emitter
}

func (e *EventEmitter) Emit(eventType ProvisioningEventType, message string) {
	eventString := string(eventType)
	e.log.Info("event", "event", eventString, "message", message)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s, err := e.eventClient.Send(ctx, &v1.EventServiceSendRequest{
		Events: map[string]*v1.MachineProvisioningEvent{
			e.machineID: {
				Time:    timestamppb.Now(),
				Event:   eventString,
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
