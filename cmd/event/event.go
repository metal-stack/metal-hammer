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

type EventEmitter struct {
	log        *slog.Logger
	bootClient infrav2connect.BootServiceClient
	machineID  string
}

func NewEventEmitter(log *slog.Logger, bootClient infrav2connect.BootServiceClient, machineID string) *EventEmitter {
	emitter := &EventEmitter{
		bootClient: bootClient,
		machineID:  machineID,
		log:        log,
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

	_, err := e.bootClient.SendEvent(ctx, &infrav2.BootServiceSendEventRequest{
		Uuid: e.machineID,
		Event: &apiv2.MachineProvisioningEvent{
			Time:    timestamppb.Now(),
			Event:   eventType,
			Message: message,
		},
	})
	if err != nil {
		e.log.Error("event", "cannot send event", eventType, "error", err)
	}
}
