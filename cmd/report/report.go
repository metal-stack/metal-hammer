package report

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	infrav2 "github.com/metal-stack/api/go/metalstack/infra/v2"
	"github.com/metal-stack/api/go/metalstack/infra/v2/infrav2connect"
)

type Report struct {
	ConsolePassword string
	MachineUUID     string
	InstallError    error
	Log             *slog.Logger
}

// ReportInstallation will tell metal-apiserver that the installation succeeded.
//
// A failed installation is not reported here: InstallationSucceeded has no way to carry a failure,
// and calling it makes metal-apiserver move the machine into the tenant vrf and boot it from disk.
// The caller returns the installation error instead, which ends up in a crashed provisioning event.
func (r *Report) ReportInstallation(bootClient infrav2connect.BootServiceClient) error {
	if r.InstallError != nil {
		r.Log.Error("installation failed, not reporting it as succeeded", "error", r.InstallError)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := bootClient.InstallationSucceeded(ctx, &infrav2.BootServiceInstallationSucceededRequest{
		Uuid:            r.MachineUUID,
		ConsolePassword: r.ConsolePassword,
	})
	if err != nil {
		r.Log.Error("report", "error", err)
		return fmt.Errorf("unable to report image installation: %w", err)
	}

	r.Log.Info("report image installation was successful")

	return nil
}
