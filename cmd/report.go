package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/inconshreveable/log15"
)

// Report is send back to metal-core after installation finished
type Report struct {
	Success bool   `json:"success" description:"true if installation succeeded"`
	Message string `json:"message" description:"if installation failed, the error message"`
}

func (r *Report) String() string {
	return fmt.Sprintf("success:%v message:%s", r.Success, r.Message)
}

// ReportInstallation will tell metal-core the result of the installation
func ReportInstallation(url, uuid string, installError error) error {
	e := fmt.Sprintf("%v/%v", url, uuid)
	report := &Report{}
	report.Success = true
	if installError != nil {
		report.Success = false
		report.Message = installError.Error()
	}

	reportJSON, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("unable to serialize report to json %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, e, bytes.NewBuffer(reportJSON))
	req.Header.Set("Content-Type", "application/json")

	log.Info("report device", "uuid", uuid, "report", report)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot POST hw %s to report endpoint:%s %v", string(reportJSON), url, err)
	}
	defer resp.Body.Close()
	if !report.Success {
		log.Error("report image installation was not successful, rebooting")
		return nil
	}

	log.Info("report image installation was successful")
	return nil
}
