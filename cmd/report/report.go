package report

import (
	log "github.com/inconshreveable/log15"
	"github.com/metal-stack/go-hal"
	"github.com/metal-stack/metal-hammer/metal-core/client/machine"
	"github.com/metal-stack/metal-hammer/metal-core/models"
	"github.com/pkg/errors"
)

type Report struct {
	Client          *machine.Client
	ConsolePassword string
	MachineUUID     string
	InstallError    error
	PrimaryDisk     string
	OSPartition     string
	Initrd          string
	Cmdline         string
	Kernel          string
	BootloaderID    string
	Hal             hal.InBand
}

// ReportInstallation will tell metal-core the result of the installation
func (r *Report) ReportInstallation() error {
	var fru *models.ModelsV1MachineFru
	bmc, err := r.Hal.BMC()
	if err != nil {
		log.Error("failed to fetch BMC details", "error", err)
	} else if bmc == nil {
		log.Debug("no BMC details")
	} else {
		fru = &models.ModelsV1MachineFru{
			BoardMfg:            bmc.BoardMfg,
			BoardMfgSerial:      bmc.BoardMfgSerial,
			BoardPartNumber:     bmc.BoardPartNumber,
			ChassisPartNumber:   bmc.ChassisPartNumber,
			ChassisPartSerial:   bmc.ChassisPartSerial,
			ProductManufacturer: bmc.ProductManufacturer,
			ProductPartNumber:   bmc.ProductPartNumber,
			ProductSerial:       bmc.ProductSerial,
		}
	}

	report := &models.DomainReport{
		Success:         true,
		ConsolePassword: &r.ConsolePassword,
		PrimaryDisk:     &r.PrimaryDisk,
		OsPartition:     &r.OSPartition,
		Initrd:          &r.Initrd,
		Cmdline:         &r.Cmdline,
		Kernel:          &r.Kernel,
		Bootloaderid:    &r.BootloaderID,
		Fru:             fru,
	}
	if r.InstallError != nil {
		message := r.InstallError.Error()
		report.Success = false
		report.Message = &message
	}

	params := machine.NewReportParams()
	params.SetBody(report)
	params.ID = r.MachineUUID
	_, err = r.Client.Report(params)
	if err != nil {
		log.Error("report", "error", err)
		return errors.Wrap(err, "unable to report image installation")
	}
	log.Info("report image installation was successful")
	return nil
}
