package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	gonet "net"
	"net/http"
	"strings"
	"time"

	log "github.com/inconshreveable/log15"
	"github.com/jaypipes/ghw"
)

type Facility struct {
	ID          string    `json:"id" description:"a unique ID" unique:"true" modelDescription:"A Facility describes the location where a device is placed."  rethinkdb:"id,omitempty"`
	Name        string    `json:"name" description:"the readable name" rethinkdb:"name"`
	Description string    `json:"description,omitempty" description:"a description for this facility" optional:"true" rethinkdb:"description"`
	Created     time.Time `json:"created" description:"the creation time of this facility" optional:"true" readOnly:"true" rethinkdb:"created"`
	Changed     time.Time `json:"changed" description:"the last changed timestamp" optional:"true" readOnly:"true" rethinkdb:"changed"`
}

type Size struct {
	ID          string `json:"id" description:"a unique ID" unique:"true" modelDescription:"An image that can be put on a device." rethinkdb:"id,omitempty"`
	Name        string `json:"name" description:"the readable name" rethinkdb:"name"`
	Description string `json:"description,omitempty" description:"a description for this image" optional:"true" rethinkdb:"description"`
	// Constraints []*Constraint `json:"constraints" description:"a list of constraints that defines this size" optional:"true"`
	Created time.Time `json:"created" description:"the creation time of this image" optional:"true" readOnly:"true" rethinkdb:"created"`
	Changed time.Time `json:"changed" description:"the last changed timestamp" optional:"true" readOnly:"true" rethinkdb:"changed"`
}

type Image struct {
	ID          string    `json:"id" description:"a unique ID" unique:"true" modelDescription:"An image that can be put on a device." rethinkdb:"id,omitempty"`
	Name        string    `json:"name" description:"the readable name" rethinkdb:"name"`
	Description string    `json:"description,omitempty" description:"a description for this image" optional:"true" rethinkdb:"description"`
	Url         string    `json:"url" description:"the url to this image" rethinkdb:"url"`
	Created     time.Time `json:"created" description:"the creation time of this image" optional:"true" readOnly:"true" rethinkdb:"created"`
	Changed     time.Time `json:"changed" description:"the last changed timestamp" optional:"true" readOnly:"true" rethinkdb:"changed"`
}

type Device struct {
	ID          string         `json:"id" description:"a unique ID" unique:"true" readOnly:"true" modelDescription:"A device representing a bare metal machine." rethinkdb:"id,omitempty"`
	Name        string         `json:"name" description:"the name of the device" rethinkdb:"name"`
	Description string         `json:"description,omitempty" description:"a description for this machine" optional:"true" rethinkdb:"description"`
	Created     time.Time      `json:"created" description:"the creation time of this machine" optional:"true" readOnly:"true" rethinkdb:"created"`
	Changed     time.Time      `json:"changed" description:"the last changed timestamp" optional:"true" readOnly:"true" rethinkdb:"changed"`
	Project     string         `json:"project" description:"the project that this device is assigned to" rethinkdb:"project"`
	Facility    Facility       `json:"facility" description:"the facility assigned to this device" readOnly:"true" rethinkdb:"-"`
	FacilityID  string         `json:"-" rethinkdb:"facilityid"`
	Image       *Image         `json:"image" description:"the image assigned to this device" readOnly:"true"  rethinkdb:"-"`
	ImageID     string         `json:"-" rethinkdb:"imageid"`
	Size        *Size          `json:"size" description:"the size of this device" readOnly:"true" rethinkdb:"-"`
	SizeID      string         `json:"-" rethinkdb:"sizeid"`
	Hardware    DeviceHardware `json:"hardware" description:"the hardware of this device" rethinkdb:"hardware"`
	IP          string         `json:"ip" description:"the ip address of the allocated device" rethinkdb:"ip"`
	Hostname    string         `json:"hostname" description:"the hostname of the device" rethinkdb:"hostname"`
	SSHPubKey   string         `json:"ssh_pub_key" description:"the public ssh key to access the device with" rethinkdb:"sshPubKey"`
}

type DeviceHardware struct {
	Memory   int64         `json:"memory" description:"the total memory of the device" rethinkdb:"memory"`
	CPUCores uint32        `json:"cpu_cores" description:"the total memory of the device" rethinkdb:"cpu_cores"`
	Nics     []Nic         `json:"nics" description:"the list of network interfaces of this device" rethinkdb:"network_interfaces"`
	Disks    []BlockDevice `json:"disks" description:"the list of block devices of this device" rethinkdb:"block_devices"`
}

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
	loFound := false
	for _, n := range net.NICs {
		_, err := gonet.ParseMAC(n.MacAddress)
		if err != nil {
			log.Debug("skip interface with invalid mac", "interface", n.Name, "mac", n.MacAddress)
			continue
		}
		// check if after mac validation loopback is still present
		if n.Name == "lo" {
			loFound = true
		}
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
	// add a lo interface if not present
	if !loFound {
		lo := Nic{
			MacAddress: "00:00:00:00:00:00",
			Name:       "lo",
		}
		nics = append(nics, lo)
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
		log.Error("error getting product_uuid, use default uuid", "error", err)
		productUUID = []byte("00000000-0000-0000-0000-000000000000")
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

	log.Info("registering device", "details", hw)
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

	log.Info("register device returned", "response", string(body))
	result := make(map[string]interface{})
	var uuid string
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Error("unable to unmarshal register response", "error", err)
		uuid = "unknown"
	} else {
		uuid = result["id"].(string)
	}

	if resp.StatusCode == http.StatusOK {
		log.Info("device already registered", "uuid", uuid)
	} else if resp.StatusCode == http.StatusCreated {
		log.Info("device registered", "uuid", uuid)
	}
	return uuid, nil
}
