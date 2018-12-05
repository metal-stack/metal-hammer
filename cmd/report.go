package cmd

import (
	"fmt"

	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/client/device"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/models"

	log "github.com/inconshreveable/log15"
)

// ReportInstallation will tell metal-core the result of the installation
func (h *Hammer) ReportInstallation(uuid string, installError error) error {
	report := &models.DomainReport{
		Success:         true,
		ConsolePassword: &h.Spec.ConsolePassword,
	}
	if installError != nil {
		message := installError.Error()
		report.Success = false
		report.Message = &message
	}

	params := device.NewReportParams()
	params.SetBody(report)
	params.ID = uuid
	resp, err := h.Client.Report(params)
	if err != nil {
		return fmt.Errorf("unable to report image installation error:%v", err)
	}
	log.Info("report image installation was successful", "response", resp.Payload)
	return nil
}
