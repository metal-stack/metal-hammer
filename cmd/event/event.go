package event

import (
	"fmt"
	"time"

	"github.com/metal-stack/metal-hammer/metal-core/client/machine"
	"github.com/metal-stack/metal-hammer/metal-core/models"
	"go.uber.org/zap"
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
	log       *zap.SugaredLogger
	client    machine.ClientService
	machineID string
}

func NewEventEmitter(log *zap.SugaredLogger, client machine.ClientService, machineID string) *EventEmitter {
	emitter := &EventEmitter{
		client:    client,
		machineID: machineID,
		log:       log,
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
	event := &models.ModelsV1MachineProvisioningEvent{
		Event:   &eventString,
		Message: message,
	}
	params := machine.NewAddProvisioningEventParams()
	params.ID = e.machineID
	params.Body = event

	e.log.Infow("event", "event", eventString, "message", event.Message)
	_, err := e.client.AddProvisioningEvent(params)
	if err != nil {
		e.log.Errorw("event", "cannot send event", eventType, "error", err)
	}
}
