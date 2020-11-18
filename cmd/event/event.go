package event

import (
	"fmt"
	"time"

	log "github.com/inconshreveable/log15"
	"github.com/metal-stack/metal-hammer/metal-core/client/machine"
	"github.com/metal-stack/metal-hammer/metal-core/models"
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
	client    machine.ClientService
	machineID string
}

func NewEventEmitter(client machine.ClientService, machineID string) *EventEmitter {
	emitter := &EventEmitter{
		client:    client,
		machineID: machineID,
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

	log.Info("event", "event", eventString, "message", event.Message)
	_, err := e.client.AddProvisioningEvent(params)
	if err != nil {
		log.Error("event", "cannot sent event", eventType, "error", err)
	}
}
