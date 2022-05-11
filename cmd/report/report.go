package report

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/metal-stack/metal-api/pkg/api/v1"
	"go.uber.org/zap"
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
	Log             *zap.SugaredLogger
}

// ReportInstallation will tell metal-core the result of the installation
func (r *Report) ReportInstallation() error {
	report := &v1.BootServiceReportRequest{
		Uuid:            r.MachineUUID,
		Success:         true,
		ConsolePassword: r.ConsolePassword,
	}
	report.BootInfo = &v1.BootInfo{
		// FIXME other fields.
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
		r.Log.Errorw("report", "error", err)
		return fmt.Errorf("unable to report image installation %w", err)
	}
	r.Log.Infow("report image installation was successful")
	return nil
}
