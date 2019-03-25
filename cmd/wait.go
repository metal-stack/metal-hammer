package cmd

import (
	"encoding/json"
	"fmt"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/models"
	log "github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

// Wait until a machine create request was fired
func (h *Hammer) Wait(uuid string) (*models.ModelsMetalMachineWithPhoneHomeToken, error) {
	// We do not use the swagger client because this has no ability to specify a timeout.
	e := fmt.Sprintf("http://%v/machine/install/%v", h.Spec.MetalCoreURL, uuid)
	log.Info("waiting for install, long polling", "url", e, "uuid", uuid)

	var resp *http.Response
	var err error
	// Create a http client with a specific timeout to prevent a infinite wait
	// which could lead to a situation where e.g network outages would never be
	// detected and we will never recover from this situation.
	client := http.Client{
		Timeout: time.Duration(5 * time.Minute),
	}
	for {
		resp, err = client.Get(e)
		if err != nil {
			log.Error("wait for install failed, retrying in 30sec...", "error", err)
			time.Sleep(30 * time.Second)
		} else if resp.StatusCode == http.StatusGatewayTimeout || resp.StatusCode == http.StatusNotModified {
			log.Info("wait for install timeout retrying...", "statuscode", resp.StatusCode)
			continue
		} else if resp.StatusCode != http.StatusOK {
			log.Warn("wait for install timeout with unexpected returncode retrying in 5sec", "statuscode", resp.StatusCode)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

	machineJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "wait for install reading response failed")
	}
	log.Info("wait finished", "statuscode", resp.StatusCode, "response", machineJSON)

	var machineWithToken models.ModelsMetalMachineWithPhoneHomeToken
	err = json.Unmarshal(machineJSON, &machineWithToken)
	if err != nil {
		return nil, errors.Wrap(err, "wait for install could not unmarshal response")
	}
	log.Info("stopped waiting got", "machineWithToken", machineWithToken)

	return &machineWithToken, nil
}
