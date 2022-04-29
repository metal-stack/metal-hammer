package event

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/metal-stack/metal-api/pkg/api/v1"
	"go.uber.org/zap"
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
	log         *zap.SugaredLogger
	eventClient v1.EventServiceClient
	machineID   string
}

func NewEventEmitter(log *zap.SugaredLogger, eventClient v1.EventServiceClient, machineID string) *EventEmitter {
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
	e.log.Infow("event", "event", eventString, "message", message)
	_, err := e.eventClient.Send(context.Background(), &v1.EventServiceSendRequest{
		MachineId: e.machineID,
		Time:      timestamppb.Now(),
		Event:     eventString,
		Message:   message,
	})
	if err != nil {
		e.log.Errorw("event", "cannot send event", eventType, "error", err)
	}
}
