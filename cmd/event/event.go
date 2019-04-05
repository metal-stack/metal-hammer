package event

import (
	"fmt"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/client/machine"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/models"
	log "github.com/inconshreveable/log15"
	"time"
)

// ProvisioningEventType indicates an event emitted by a machine during the provisioning sequence
// FIXME factor out to metal-lib
type ProvisioningEventType string

// The enums for the machine provisioning events.
const (
	ProvisioningEventAlive            ProvisioningEventType = "Alive"
	ProvisioningEventPlannedReboot    ProvisioningEventType = "Planned Reboot"
	ProvisioningEventPreparing        ProvisioningEventType = "Preparing"
	ProvisioningEventRegistering      ProvisioningEventType = "Registering"
	ProvisioningEventWaiting          ProvisioningEventType = "Waiting"
	ProvisioningEventInstalling       ProvisioningEventType = "Installing"
	ProvisioningEventBootingNewKernel ProvisioningEventType = "Booting New Kernel"
)

type EventEmitter struct {
	client    *machine.Client
	machineID string
}

func NewEventEmitter(client *machine.Client, machineID string) *EventEmitter {
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
	event := &models.ModelsMetalProvisioningEvent{
		Event:   &eventString,
		Message: message,
	}
	params := machine.NewAddProvisioningEventParams()
	params.ID = e.machineID
	params.Body = event

	log.Info("event", "event", event.Event, "message", event.Message)
	_, err := e.client.AddProvisioningEvent(params)
	if err != nil {
		log.Error("event", "cannot sent event", eventType, "error", err)
	}
}
