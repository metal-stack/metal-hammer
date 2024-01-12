package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	v1 "github.com/metal-stack/metal-api/pkg/api/v1"
	"github.com/metal-stack/metal-hammer/cmd/event"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const defaultWaitTimeOut = 2 * time.Second

func (c *MetalAPIClient) WaitForAllocation(e *event.EventEmitter, machineID string) error {
	e.Emit(event.ProvisioningEventWaiting, "waiting for allocation")

	req := &v1.BootServiceWaitRequest{
		MachineId: machineID,
	}
	for {
		stream, err := c.BootService().Wait(context.Background(), req)
		if err != nil {
			c.log.Error("failed waiting for allocation", "retry after", defaultWaitTimeOut, "error", err)

			time.Sleep(defaultWaitTimeOut)
			continue
		}

		for {
			_, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				c.log.Info("machine has been requested for allocation", "machineID", machineID)
				return nil
			}

			if err != nil {
				if e, ok := status.FromError(err); ok {
					c.log.Error("got error from wait call", "code", e.Code(), "message", e.Message(), "details", e.Details())
					switch e.Code() { // nolint:exhaustive
					case codes.Unimplemented:
						return fmt.Errorf("metal-api breaking change detected, rebooting: %w", err)
					}
				}

				c.log.Error("failed stream receiving during waiting for allocation", "retry after", defaultWaitTimeOut, "error", err)

				time.Sleep(defaultWaitTimeOut)
				break
			}

			c.log.Info("wait for allocation...", "machineID", machineID)
		}
	}
}
