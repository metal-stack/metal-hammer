package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	infrav2 "github.com/metal-stack/api/go/metalstack/infra/v2"
	"github.com/metal-stack/api/go/metalstack/infra/v2/infrav2connect"
)

// WaitForAllocation can be used to call the wait method continuously until an allocation was made.
// This is made for the metal-hammer and located here for better testability.
func WaitForAllocation(ctx context.Context, log *slog.Logger, service infrav2connect.BootServiceClient, machineID string, timeout time.Duration) (*apiv2.MachineAllocation, error) {
	req := &infrav2.BootServiceWaitRequest{
		Uuid: machineID,
	}

	for {
		stream, err := service.Wait(ctx, req)
		defer func() {
			_ = stream.Close()
		}()
		if err != nil {
			log.Error("failed waiting for allocation", "retry after", timeout, "error", err)

			if strings.Contains(err.Error(), "failed to verify certificate") {
				return nil, fmt.Errorf("certificate changed, rebooting")
			}

			time.Sleep(timeout)
			continue
		}

		log.Info("wait for allocation...")
		for stream.Receive() {
			resp := stream.Msg()
			return resp.Allocation, nil
		}
	}
}
