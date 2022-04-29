package report

import (
	"fmt"

	"github.com/metal-stack/metal-hammer/metal-core/client/machine"
	"github.com/metal-stack/metal-hammer/metal-core/models"
	"go.uber.org/zap"
)

type Report struct {
	Client          machine.ClientService
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
	report := &models.DomainReport{
		Success:         true,
		ConsolePassword: &r.ConsolePassword,
		Initrd:          &r.Initrd,
		Cmdline:         &r.Cmdline,
		Kernel:          &r.Kernel,
		Bootloaderid:    &r.BootloaderID,
	}
	if r.InstallError != nil {
		message := r.InstallError.Error()
		report.Success = false
		report.Message = &message
	}

	params := machine.NewReportParams()
	params.SetBody(report)
	params.ID = r.MachineUUID
	_, err := r.Client.Report(params)
	if err != nil {
		r.Log.Errorw("report", "error", err)
		return fmt.Errorf("unable to report image installation %w", err)
	}
	r.Log.Infow("report image installation was successful")
	return nil
}
