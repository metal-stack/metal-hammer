package report

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	v1 "github.com/metal-stack/metal-api/pkg/api/v1"
)

type Report struct {
	Client          v1.BootServiceClient
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
	report := &v1.BootServiceReportRequest{
		Uuid:            r.MachineUUID,
		Success:         true,
		ConsolePassword: r.ConsolePassword,
	}
	report.BootInfo = &v1.BootInfo{
		Initrd:       r.Initrd,
		Cmdline:      r.Cmdline,
		Kernel:       r.Kernel,
		BootloaderId: r.BootloaderID,
	}
	if r.InstallError != nil {
		message := r.InstallError.Error()
		report.Success = false
		report.Message = message
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.Client.Report(ctx, report)
	if err != nil {
		r.Log.Error("report", "error", err)
		return fmt.Errorf("unable to report image installation %w", err)
	}
	r.Log.Info("report image installation was successful")
	return nil
}
