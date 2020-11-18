package cmd

import (
	"context"
	log "github.com/inconshreveable/log15"
	v1 "github.com/metal-stack/metal-api/pkg/api/v1"
	"github.com/metal-stack/metal-hammer/cmd/event"
	"io"
	"time"
)

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
			log.Error("failed waiting for allocation, retry in 2sec", "error", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for {
			_, err := stream.Recv()
			if err == io.EOF {
				log.Info("machine has been requested for allocation", "machineID", machineID)
				return nil
			}

			if err != nil {
				log.Error("failed waiting for allocation, retry in 2sec", "error", err)
				time.Sleep(2 * time.Second)
				break
			}

			log.Info("wait for allocation...", "machineID", machineID)
		}
	}
}
