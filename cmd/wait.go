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

// Wait until a device create request was fired
func (h *Hammer) Wait(uuid string) (*models.ModelsMetalDeviceWithPhoneHomeToken, error) {
	// We do not use the swagger client because this has no ability to specify a timeout.
	e := fmt.Sprintf("http://%v/device/install/%v", h.Spec.MetalCoreURL, uuid)
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
			log.Warn("wait for install failed, retrying...", "error", err)
		} else if resp.StatusCode != http.StatusOK {
			log.Warn("wait for install failed, retrying...", "statuscode", resp.StatusCode)
		} else {
			break
		}
	}

	deviceJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "wait for install reading response failed")
	}

	var deviceWithToken models.ModelsMetalDeviceWithPhoneHomeToken
	err = json.Unmarshal(deviceJSON, &deviceWithToken)
	if err != nil {
		return nil, errors.Wrap(err, "wait for install could not unmarshal response")
	}
	log.Info("stopped waiting got", "deviceWithToken", deviceWithToken)

	return &deviceWithToken, nil
}
