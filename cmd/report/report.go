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
	Client          infrav2connect.BootServiceClient
	ConsolePassword string
	MachineUUID     string
	InstallError    error
	Initrd          string
	Cmdline         string
	Kernel          string
	BootloaderID    string
	Log             *slog.Logger
}

// ReportInstallation will tell metal-api the result of the installation
func (r *Report) ReportInstallation() error {
	report := &infrav2.BootServiceInstallationSucceededRequest{
		Uuid:            r.MachineUUID,
		ConsolePassword: r.ConsolePassword,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.Client.InstallationSucceeded(ctx, report)
	if err != nil {
		r.Log.Error("report", "error", err)
		return fmt.Errorf("unable to report image installation %w", err)
	}
	r.Log.Info("report image installation was successful")
	return nil
}
