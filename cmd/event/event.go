package event

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/metal-stack/metal-hammer/pkg/kernel"

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
	log                  *zap.SugaredLogger
	eventClient          v1.EventServiceClient
	machineID            string
	consecutiveErrors    atomic.Uint32
	maxConsecutiveErrors uint32
}

func NewEventEmitter(log *zap.SugaredLogger, eventClient v1.EventServiceClient, machineID string, maxErrors uint32) *EventEmitter {
	emitter := &EventEmitter{
		eventClient:          eventClient,
		machineID:            machineID,
		log:                  log,
		consecutiveErrors:    atomic.Uint32{},
		maxConsecutiveErrors: maxErrors,
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
	e.log.Infow("event", "event", eventString, "message", message, "errorCount", e.consecutiveErrors.Load())
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
		e.log.Errorw("event", "cannot send event", eventType, "errorCount", e.consecutiveErrors.Load(), "error", err)
		e.consecutiveErrors.Add(1)
		if e.consecutiveErrors.Load() > e.maxConsecutiveErrors {
			err = kernel.Reboot()
			if err != nil {
				e.log.Errorw("event, unable to reboot because of too many consecutive errors", "error", err)
			}
		}
	}
	if s != nil {
		e.log.Infow("event", "send", s.Events, "failed", s.Failed)
		e.consecutiveErrors.Store(0)
	}
}
