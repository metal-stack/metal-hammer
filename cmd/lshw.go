package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/inconshreveable/log15"
)

type registerRequest struct {
	UUID string `json:"uuid" description:"the uuid of the device to register"`
}

//RegisterDevice register a device at the maas api
func RegisterDevice(spec *Specification) (string, error) {
	result := registerRequest{}
	uuid, err := ioutil.ReadFile("/sys/class/dmi/id/product_uuid")
	if err != nil {
		return "", fmt.Errorf("error getting product_uuid info: %v", err)
	}

	result.UUID = string(uuid)
	return register(spec.ReportURL, result)
}

func register(url string, hw registerRequest) (string, error) {
	e := fmt.Sprintf("%v/%v", url, hw.UUID)
	hwJSON, err := json.Marshal(hw)
	if err != nil {
		return "", fmt.Errorf("unable to serialize hw to json %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, e, bytes.NewBuffer(hwJSON))
	req.Header.Set("Content-Type", "application/json")

	log.Info("registering device", "uuid", hw.UUID)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot POST hw json struct to register endpoint: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response from register call %v", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return "", fmt.Errorf("POST of hw to register endpoint did not succeed %v", resp.Status)
	}

	result := make(map[string]interface{})
	var uuid interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		uuid = "unknown"
	} else {
		uuid = result["id"]
	}

	if resp.StatusCode == http.StatusOK {
		log.Info("device already registered", "uuid", uuid)
	} else if resp.StatusCode == http.StatusCreated {
		log.Info("device registered", "uuid", uuid)
	}
	return hw.UUID, nil
}
