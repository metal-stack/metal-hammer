package cmd

import (
	"fmt"
	"io/ioutil"
	"net"
	gonet "net"
	"strings"

	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/client/device"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/models"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/ipmi"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/password"

	log "github.com/inconshreveable/log15"

	"github.com/jaypipes/ghw"
)

// RegisterDevice register a device at the metal-api via metal-core
func (h *Hammer) RegisterDevice() (string, error) {
	hw, err := h.readHardwareDetails()
	if err != nil {
		return "", fmt.Errorf("unable to read all hardware details error:%v", err)
	}
	params := device.NewRegisterParams()
	params.SetBody(hw)
	params.ID = hw.UUID
	resp, err := h.Client.Register(params)

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

func (h *Hammer) readHardwareDetails() (*models.DomainMetalHammerRegisterDeviceRequest, error) {
	hw := &models.DomainMetalHammerRegisterDeviceRequest{}

	memory, err := ghw.Memory()
	if err != nil {
		return nil, fmt.Errorf("unable to get system memory, info:%v", err)
	}
	hw.Memory = &memory.TotalPhysicalBytes

	cpu, err := ghw.CPU()
	if err != nil {
		return nil, fmt.Errorf("unable to get system cpu(s), info:%v", err)
	}
	cores := int32(cpu.TotalCores)
	hw.CPUCores = &cores

	net, err := ghw.Network()
	if err != nil {
		return nil, fmt.Errorf("unable to get system nic(s), info:%v", err)
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
		nic := &models.ModelsMetalNic{
			Mac:  &n.MacAddress,
			Name: &n.Name,
		}
		nics = append(nics, nic)
	}
	// add a lo interface if not present
	// this is required to have this interface present
	// in our DCIM management to add a ip later.
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
		return nil, fmt.Errorf("unable to get system block devices, info:%v", err)
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
	hw.UUID = uuid

	ipmiconfig, err := h.readIPMIDetails()
	if err != nil {
		return nil, err
	}
	hw.IPMI = ipmiconfig

	return hw, nil
}

const defaultIpmiPort = "623"

var defaultIpmiUser = "metal"

// IPMI configuration and
func (h *Hammer) readIPMIDetails() (*models.ModelsMetalIPMI, error) {
	config := ipmi.LanConfig{}
	pw := password.Generate(10)
	var i ipmi.Ipmi
	if h.Spec.IPMIPort != defaultIpmiPort {
		// Wild guess, set the last octet to 1 to get the gateway
		gwip := net.ParseIP(h.IPAddress)
		gwip = gwip.To4()
		gwip[3] = 1

		config.IP = fmt.Sprintf("%s:%s", gwip, h.Spec.IPMIPort)
		config.Mac = "00:00:00:00:00:00"
	} else {
		var err error
		i = ipmi.New()
		// FIXME userid should be verified if available
		err = i.CreateUser("metal", pw, 2, ipmi.Administrator)
		if err != nil {
			return nil, fmt.Errorf("ipmi error: %v", err)
		}
		config, err = i.GetLanConfig()
		if err != nil {
			return nil, fmt.Errorf("unable to read ipmi lan configuration, info:%v", err)
		}
		config.IP = config.IP + ":" + defaultIpmiPort
	}

	intf := "lanplus"
	details := &models.ModelsMetalIPMI{
		Address:   &config.IP,
		Mac:       &config.Mac,
		Password:  &pw,
		User:      &defaultIpmiUser,
		Interface: &intf,
	}

	return details, nil
}
