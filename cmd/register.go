package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/inconshreveable/log15"

	"github.com/jaypipes/ghw"
)

type Nic struct {
	MacAddress string   `json:"mac"`
	Name       string   `json:"name"`
	Vendor     string   `json:"vendor"`
	Features   []string `json:"features"`
}

type BlockDevice struct {
	Name string `json:"name"`
	Size uint64 `json:"size"`
}

type registerRequest struct {
	UUID     string        `json:"uuid" description:"the uuid of the device to register"`
	Memory   int64         `json:"memory" description:"the memory in bytes of the device to register"`
	CPUCores uint32        `json:"cpucores" description:"the cpu core of the device to register"`
	Nics     []Nic         `json:"nics"`
	Disks    []BlockDevice `json:"disks"`
}

// RegisterDevice register a device at the metal-api via metal-core
func RegisterDevice(spec *Specification) (string, error) {
	hw := registerRequest{}

	memory, err := ghw.Memory()
	if err != nil {
		return "", fmt.Errorf("unable to get system memory, info:%v", err)
	}
	hw.Memory = memory.TotalPhysicalBytes

	cpu, err := ghw.CPU()
	if err != nil {
		return "", fmt.Errorf("unable to get system cpu(s), info:%v", err)
	}
	hw.CPUCores = cpu.TotalCores

	net, err := ghw.Network()
	if err != nil {
		return "", fmt.Errorf("unable to get system nic(s), info:%v", err)
	}
	nics := []Nic{}
	for _, n := range net.NICs {
		features := []string{}
		if n.EnabledFeatures != nil {
			features = n.EnabledFeatures
		}
		nic := Nic{
			MacAddress: n.MacAddress,
			Name:       n.Name,
			Vendor:     n.Vendor,
			Features:   features,
		}
		nics = append(nics, nic)
	}
	hw.Nics = nics

	blockInfo, err := ghw.Block()
	if err != nil {
		return "", fmt.Errorf("unable to get system block devices, info:%v", err)
	}
	disks := []BlockDevice{}
	for _, disk := range blockInfo.Disks {
		blockDevice := BlockDevice{
			Name: disk.Name,
			Size: disk.SizeBytes,
		}
		disks = append(disks, blockDevice)
	}
	hw.Disks = disks

	productUUID, err := ioutil.ReadFile("/sys/class/dmi/id/product_uuid")
	if err != nil {
		return "", fmt.Errorf("error getting product_uuid info: %v", err)
	}

	hw.UUID = strings.TrimSpace(string(productUUID))
	return register(spec.RegisterURL, hw)
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
		return "", fmt.Errorf("cannot POST hw %s to register endpoint:%s %v", string(hwJSON), url, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response from register call %v", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return "", fmt.Errorf("POST of hw %s to register endpoint:%s did not succeed %v response body:%s", string(hwJSON), url, resp.Status, body)
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
