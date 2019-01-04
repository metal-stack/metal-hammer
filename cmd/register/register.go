package register

import (
	"fmt"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/cmd/network"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/client/device"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/models"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/ipmi"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/password"
	"io/ioutil"
	gonet "net"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	log "github.com/inconshreveable/log15"
	"github.com/jaypipes/ghw"
	"github.com/vishvananda/netlink"
)

// Register the Device
type Register struct {
	DeviceUUID string
	Client     *device.Client
	Network    *network.Network
}

// RegisterDevice register a device at the metal-api via metal-core
func (r *Register) RegisterDevice() (string, error) {
	hw, err := r.readHardwareDetails()
	if err != nil {
		return "", fmt.Errorf("unable to read all hardware details error:%v", err)
	}
	params := device.NewRegisterParams()
	params.SetBody(hw)
	params.ID = hw.UUID
	resp, err := r.Client.Register(params)

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

func (r *Register) readHardwareDetails() (*models.DomainMetalHammerRegisterDeviceRequest, error) {
	err := createSyslog()
	if err != nil {
		return nil, fmt.Errorf("unable to write kernel boot message to /var/log/syslog, info:%v", err)
	}

	hw := &models.DomainMetalHammerRegisterDeviceRequest{}

	memory, err := ghw.Memory()
	if err != nil {
		return nil, fmt.Errorf("unable to get system memory, info:%v", err)
	}
	hw.Memory = &memory.TotalPhysicalBytes

	// FIXME can be replaced by runtime.NumCPU()
	cpu, err := ghw.CPU()
	if err != nil {
		return nil, fmt.Errorf("unable to get system cpu(s), info:%v", err)
	}
	cores := int32(cpu.TotalCores)
	hw.CPUCores = &cores

	nics := []*models.ModelsMetalNic{}
	loFound := false
	links, err := netlink.LinkList()
	if err != nil {
		return nil, fmt.Errorf("unable to get all links:%v", err)
	}
	for _, l := range links {
		attrs := l.Attrs()
		name := attrs.Name
		mac := attrs.HardwareAddr.String()
		_, err := gonet.ParseMAC(mac)

		if err != nil {
			log.Debug("skip interface with invalid mac", "interface", name, "mac", mac)
			continue
		}
		// check if after mac validation loopback is still present
		if name == "lo" {
			loFound = true
		}
		if name == "eth0" {
			eth0Mac = mac
		}

		nic := &models.ModelsMetalNic{
			Mac:  &mac,
			Name: &name,
		}
		log.Info("register", "nic", name, "mac", mac)
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

	// now attach neighbors, this will wait up to 2*tx-intervall
	// if during this timeout not all required neighbors where found abort and reboot.
	for _, n := range nics {
		neighbors, err := r.Network.Neighbors(*n.Name)
		if err != nil {
			return nil, fmt.Errorf("unable to determine neighbors of interface:%s error:%v", *n.Name, err)
		}
		n.Neighbors = neighbors
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
	hw.UUID = r.DeviceUUID

	ipmiconfig, err := readIPMIDetails(eth0Mac)
	if err != nil {
		return nil, err
	}
	hw.IPMI = ipmiconfig

	return hw, nil
}

const defaultIpmiPort = "623"

const defaultIpmiUser = "metal"

// IPMI configuration and
func readIPMIDetails(eth0Mac string) (*models.ModelsMetalIPMI, error) {
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

// save the content of kernel ringbuffer to /var/log/syslog
// by calling the appropriate syscall.
// Only required if Memory is gathered by ghw.Memory()
// FIXME consider different implementation
func createSyslog() error {
	const SyslogActionReadAll = 3
	level := uintptr(SyslogActionReadAll)

	b := make([]byte, 256*1024)
	amt, _, err := syscall.Syscall(syscall.SYS_SYSLOG, level, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))
	if err != 0 {
		return err
	}

	return ioutil.WriteFile("/var/log/syslog", b[:amt], 0666)
}
