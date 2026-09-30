package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	infrav2 "github.com/metal-stack/api/go/metalstack/infra/v2"
	"github.com/metal-stack/api/go/metalstack/infra/v2/infrav2connect"
)

// WaitForAllocation can be used to call the wait method continuously until an allocation was made.
// This is made for the metal-hammer and located here for better testability.
func WaitForAllocation(ctx context.Context, log *slog.Logger, c infrav2connect.BootServiceClient, machineID string) (*apiv2.MachineAllocation, error) {
	msgs, errs := client.ReconnectingStreamRead(ctx, func(ctx context.Context) (*connect.ServerStreamForClient[infrav2.BootServiceWaitResponse], error) {
		return c.Wait(ctx, &infrav2.BootServiceWaitRequest{Uuid: machineID})
	}, client.WithStreamBackoff(10*time.Second), client.WithStreamLogger(log))

	for {
		log.Info("wait for allocation...")

		select {
		case message := <-msgs:
			log.Info("received allocation")

			return message.Allocation, nil
		case err := <-errs:
			log.Error("error waiting for allocation", "error", err)

			continue
		case <-ctx.Done():
			log.Info("context cancelled, stop waiting for allocation")

			return nil, fmt.Errorf("stopped waiting for allocation")
		}
	}
}
