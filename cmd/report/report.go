package report

import (
	"github.com/metal-stack/metal-hammer/metal-core/client/machine"
	"github.com/metal-stack/metal-hammer/metal-core/models"

	log "github.com/inconshreveable/log15"
	"github.com/pkg/errors"
)

type Report struct {
	Client          *machine.Client
	ConsolePassword string
	MachineUUID     string
	InstallError    error
	PrimaryDisk     string
	OSPartition     string
}

// ReportInstallation will tell metal-core the result of the installation
func (r *Report) ReportInstallation() error {
	report := &models.DomainReport{
		Success:         true,
		ConsolePassword: &r.ConsolePassword,
		PrimaryDisk:     &r.PrimaryDisk,
		OsPartition:     &r.OSPartition,
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
		log.Error("report", "error", err)
		return errors.Wrap(err, "unable to report image installation")
	}
	log.Info("report image installation was successful")
	return nil
}
