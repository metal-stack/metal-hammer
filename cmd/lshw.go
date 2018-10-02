package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

	log "github.com/inconshreveable/log15"
)

var (
	lshwCommand = "/usr/sbin/lshw"
)

//RegisterDevice register a device at the maas api
func RegisterDevice(spec *Specification) (string, error) {
	lshw, err := executeCommand()
	if err != nil {
		return "", fmt.Errorf("error reading lshw output %v", err.Error())
	}
	log.Debug("lshw output", "raw", lshw)
	return register(spec.ReportURL, lshw)
}

func executeCommand() (string, error) {
	lshwOutput, err := exec.Command(lshwCommand, "-quiet", "-json").Output()
	if err != nil {
		return "", err
	}
	return string(lshwOutput), nil
}

func register(url, lshw string) (string, error) {
	var jsonStr = []byte(lshw)

	fakeUUID := "1234-1234-1234"
	e := fmt.Sprintf("%v/%v", url,fakeUUID)
	req, err := http.NewRequest(http.MethodPost, e, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	log.Info("registering device", "uuid", fakeUUID)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot POST lshw json struct to register endpoint: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response from register call %v", err)
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("POST of lshw to register endpoint did not succeed %v", resp.Status)
	}

	result := make(map[string]interface{})
	var uuid interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		uuid = "unknown"
	} else {
		uuid = result["id"]
	}

	if resp.StatusCode == 200 {
		log.Info("device already registered", "uuid", uuid)
	} else if resp.StatusCode == 201 {
		log.Info("device registered", "uuid", uuid)
	}
	return fakeUUID, nil
}
