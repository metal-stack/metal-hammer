package cmd

import (
	"context"
	"errors"
	"io"
	"time"

	v1 "github.com/metal-stack/metal-api/pkg/api/v1"
	"github.com/metal-stack/metal-hammer/cmd/event"
)

const defaultWaitTimeOut = 2 * time.Second

func (c *GrpcClient) NewWaitClient() (v1.WaitClient, io.Closer, error) {
	conn, err := c.newConnection()
	if err != nil {
		return nil, nil, err
	}
	return v1.NewWaitClient(conn), conn, nil
}

func (c *GrpcClient) WaitForAllocation(machineID string) error {
	client, closer, err := c.NewWaitClient()
	if err != nil {
		return err
	}
	defer closer.Close()

	c.Emit(event.ProvisioningEventWaiting, "waiting for allocation")

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
