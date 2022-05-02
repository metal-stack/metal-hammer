package cmd

import (
	"context"
	"errors"
	"io"
	"time"

	v1 "github.com/metal-stack/metal-api/pkg/api/v1"
	"github.com/metal-stack/metal-hammer/cmd/event"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const defaultWaitTimeOut = 2 * time.Second

func (c *GrpcClient) WaitForAllocation(machineID string) error {
	client, closer, err := c.NewWaitClient()
	if err != nil {
		return err
	}
	defer closer.Close()

	e, eventCloser, err := c.NewEventClient()
	if err != nil {
		return err
	}
	defer eventCloser.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s, err := e.Send(ctx, &v1.EventServiceSendRequest{
		Events: map[string]*v1.MachineProvisioningEvent{
			machineID: {
				Time:    timestamppb.Now(),
				Event:   string(event.ProvisioningEventWaiting),
				Message: "waiting for allocation",
			},
		},
	})
	if err != nil {
		c.log.Errorw("wait unable to send event", "error", err)
		return err
	}
	if s != nil {
		c.log.Infow("event", "send", s.Events, "failed", s.Failed)
	}

	req := &v1.WaitRequest{
		MachineID: machineID,
	}
	for {
		stream, err := client.Wait(context.Background(), req)
		if err != nil {
			c.log.Errorw("failed waiting for allocation", "retry after", defaultWaitTimeOut, "error", err)
			time.Sleep(defaultWaitTimeOut)
			continue
		}

		for {
			_, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				c.log.Infow("machine has been requested for allocation", "machineID", machineID)
				return nil
			}

			if err != nil {
				c.log.Errorw("failed stream receiving during waiting for allocation", "retry after", defaultWaitTimeOut, "error", err)
				time.Sleep(defaultWaitTimeOut)
				break
			}

			c.log.Infow("wait for allocation...", "machineID", machineID)
		}
	}
}
