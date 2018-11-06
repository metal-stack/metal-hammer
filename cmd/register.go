package cmd

import (
	"fmt"
	"io/ioutil"
	gonet "net"
	"strings"

	"git.f-i-ts.de/cloud-native/maas/metal-hammer/metal-core/client/device"
	"git.f-i-ts.de/cloud-native/maas/metal-hammer/metal-core/models"

	log "github.com/inconshreveable/log15"

	"github.com/jaypipes/ghw"
)

// RegisterDevice register a device at the metal-api via metal-core
func (h *Hammer) RegisterDevice() (string, error) {
	hw := &models.DomainMetalHammerRegisterDeviceRequest{}

	memory, err := ghw.Memory()
	if err != nil {
		return "", fmt.Errorf("unable to get system memory, info:%v", err)
	}
	hw.Memory = &memory.TotalPhysicalBytes

	cpu, err := ghw.CPU()
	if err != nil {
		return "", fmt.Errorf("unable to get system cpu(s), info:%v", err)
	}
	cores := int64(cpu.TotalCores)
	hw.CPUCores = &cores

	net, err := ghw.Network()
	if err != nil {
		return "", fmt.Errorf("unable to get system nic(s), info:%v", err)
	}
	nics := []*models.ModelsMetalNic{}
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
		nic := &models.ModelsMetalNic{
			Mac:      &n.MacAddress,
			Name:     &n.Name,
			Vendor:   &n.Vendor,
			Features: features,
		}
		nics = append(nics, nic)
	}
	// add a lo interface if not present
	if !loFound {
		mac := "00:00:00:00:00:00"
		name := "lo"
		lo := &models.ModelsMetalNic{
			Mac:  &mac,
			Name: &name,
		}
		nics = append(nics, lo)
	}

	hw.Nics = nics

	blockInfo, err := ghw.Block()
	if err != nil {
		return "", fmt.Errorf("unable to get system block devices, info:%v", err)
	}
	disks := []*models.ModelsMetalBlockDevice{}
	for _, disk := range blockInfo.Disks {
		size := int64(disk.SizeBytes)
		blockDevice := &models.ModelsMetalBlockDevice{
			Name: &disk.Name,
			Size: &size,
		}
		disks = append(disks, blockDevice)
	}
	hw.Disks = disks

	productUUID, err := ioutil.ReadFile("/sys/class/dmi/id/product_uuid")
	if err != nil {
		log.Error("error getting product_uuid, use default uuid", "error", err)
		productUUID = []byte("00000000-0000-0000-0000-000000000000")
	}

	uuid := strings.TrimSpace(string(productUUID))
	hw.UUID = &uuid

	params := device.NewRegisterEndpointParams()
	params.SetBody(hw)
	params.ID = uuid
	resp, err := h.Client.RegisterEndpoint(params)

	if err != nil {
		return "", fmt.Errorf("unable to register device:%#v error:%#v", hw, err.Error())
	}
	if resp == nil {
		return "", fmt.Errorf("unable to register device:%#v response payload is nil", hw)
	}

	log.Info("register device returned", "response", resp.Payload)
	// FIXME add different logging based on created/already registered
	// if resp.StatusCode() == http.StatusOK {
	//	log.Info("device already registered", "uuid", uuid)
	//} else if resp.StatusCode == http.StatusCreated {
	//	log.Info("device registered", "uuid", uuid)
	//}
	return *resp.Payload.ID, nil
}
