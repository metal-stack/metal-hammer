package cmd

import (
	"fmt"
	"io/ioutil"
	gonet "net"
	"os"
	"strconv"
	"strings"

	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/client/device"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/models"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/ipmi"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/password"

	"github.com/google/uuid"
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

// this mac is used to calculate the IPMI Port offset in the metal-lab environment.
var eth0Mac = ""

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
		if n.Name == "eth0" {
			eth0Mac = n.MacAddress
		}
		neighbors, err := Neighbors(n.Name)
		if err != nil {
			log.Error("unable to determine neighbors of interface", "interface", n.Name)
		}

		nic := &models.ModelsMetalNic{
			Mac:       &n.MacAddress,
			Name:      &n.Name,
			Neighbors: neighbors,
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

	uuid := getServerUUID()
	hw.UUID = uuid

	ipmiconfig, err := h.readIPMIDetails(eth0Mac)
	if err != nil {
		return nil, err
	}
	hw.IPMI = ipmiconfig

	return hw, nil
}

const dmiUUID = "/sys/class/dmi/id/product_uuid"
const dmiSerial = "/sys/class/dmi/id/product_serial"

func getServerUUID() string {
	if _, err := os.Stat(dmiUUID); !os.IsNotExist(err) {
		productUUID, err := ioutil.ReadFile(dmiUUID)
		if err != nil {
			log.Error("error getting product_uuid", "error", err)
		} else {
			log.Info("create UUID from", "source", dmiUUID)
			return strings.TrimSpace(string(productUUID))
		}
	}

	if _, err := os.Stat(dmiSerial); !os.IsNotExist(err) {
		productSerial, err := ioutil.ReadFile(dmiSerial)
		if err != nil {
			log.Error("error getting product_serial", "error", err)
		} else {
			productSerialBytes, err := uuid.FromBytes([]byte(fmt.Sprintf("%16s", string(productSerial))))
			if err != nil {
				log.Error("error getting converting product_serial to uuid", "error", err)
			} else {
				log.Info("create UUID from", "source", dmiSerial)
				return strings.TrimSpace(productSerialBytes.String())
			}
		}
	}
	log.Error("no valid UUID found", "return uuid", "00000000-0000-0000-0000-000000000000")
	return "00000000-0000-0000-0000-000000000000"
}

const defaultIpmiPort = "623"

const defaultIpmiUser = "metal"

// IPMI configuration and
func (h *Hammer) readIPMIDetails(eth0Mac string) (*models.ModelsMetalIPMI, error) {
	config := ipmi.LanConfig{}
	i := ipmi.New()
	var pw string
	var user string
	if i.DevicePresent() {
		log.Info("ipmi details from bmc")
		pw = password.Generate(10)
		user = defaultIpmiUser
		// FIXME userid should be verified if available
		err := i.CreateUser(user, pw, 2, ipmi.Administrator)
		if err != nil {
			return nil, fmt.Errorf("ipmi error: %v", err)
		}
		config, err = i.GetLanConfig()
		if err != nil {
			return nil, fmt.Errorf("unable to read ipmi lan configuration, info:%v", err)
		}
		config.IP = config.IP + ":" + defaultIpmiPort
	} else {
		log.Info("ipmi details faked")

		if len(eth0Mac) == 0 {
			eth0Mac = "00:00:00:00:00:00"
		}

		macParts := strings.Split(eth0Mac, ":")
		lastOctet := macParts[len(macParts)-1]
		port, err := strconv.ParseUint(lastOctet, 16, 32)
		if err != nil {
			return nil, fmt.Errorf("unable to parse last octet of eth0 mac to a integer: %v", err)
		}

		const baseIPMIPort = 6230
		// Fixed IP of vagrant environment gateway
		config.IP = fmt.Sprintf("192.168.121.1:%d", baseIPMIPort+port)
		config.Mac = "00:00:00:00:00:00"
		pw = "vagrant"
		user = "vagrant"
	}

	intf := "lanplus"
	details := &models.ModelsMetalIPMI{
		Address:   &config.IP,
		Mac:       &config.Mac,
		Password:  &pw,
		User:      &user,
		Interface: &intf,
	}

	return details, nil
}
