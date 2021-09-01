package register

import (
	"fmt"
	gonet "net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/jaypipes/ghw"
	"github.com/metal-stack/go-hal"
	"github.com/metal-stack/go-hal/pkg/api"
	"github.com/metal-stack/metal-hammer/cmd/network"
	"github.com/metal-stack/metal-hammer/cmd/storage"
	"github.com/metal-stack/metal-hammer/metal-core/client/machine"
	"github.com/metal-stack/metal-hammer/metal-core/models"
	"github.com/vishvananda/netlink"

	log "github.com/inconshreveable/log15"
)

// Register the Machine
type Register struct {
	MachineUUID string
	Client      machine.ClientService
	Network     *network.Network
	Hal         hal.InBand
}

// RegisterMachine register a machine at the metal-api via metal-core
func (r *Register) RegisterMachine(hw *models.DomainMetalHammerRegisterMachineRequest) error {
	params := machine.NewRegisterParams()
	params.SetBody(hw)
	params.ID = hw.UUID
	resp, err := r.Client.Register(params)

	if err != nil {
		return fmt.Errorf("unable to register machine:%#v %w", hw, err)
	}
	if resp == nil {
		return fmt.Errorf("unable to register machine:%#v response payload is nil", hw)
	}

	log.Info("register machine returned", "response", resp.Payload)
	return nil
}

// ReadHardwareDetails returns the hardware details of the machine
func (r *Register) ReadHardwareDetails() (*models.DomainMetalHammerRegisterMachineRequest, error) {
	err := createSyslog()
	if err != nil {
		return nil, fmt.Errorf("unable to write kernel boot message to /var/log/syslog %w", err)
	}

	hw := &models.DomainMetalHammerRegisterMachineRequest{}

	memory, err := ghw.Memory()
	if err != nil {
		return nil, fmt.Errorf("unable to get system memory %w", err)
	}
	hw.Memory = &memory.TotalPhysicalBytes

	// FIXME can be replaced by runtime.NumCPU()
	cpu, err := ghw.CPU()
	if err != nil {
		return nil, fmt.Errorf("unable to get system cpu(s) %w", err)
	}
	cores := int32(cpu.TotalCores)
	hw.CPUCores = &cores

	nics := []*models.ModelsV1MachineNicExtended{}
	loFound := false
	links, err := netlink.LinkList()
	if err != nil {
		return nil, fmt.Errorf("unable to get all links %w", err)
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
			r.Network.Eth0Mac = mac
		}

		nic := &models.ModelsV1MachineNicExtended{
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
		lo := &models.ModelsV1MachineNicExtended{
			Mac:  &mac,
			Name: &name,
		}
		nics = append(nics, lo)
	}

	// now attach neighbors, this will wait up to 2*tx-intervall
	// if during this timeout not all required neighbors where found abort and reboot.
	for _, nic := range nics {
		neighbors, err := r.Network.Neighbors(*nic.Name)
		if err != nil {
			return nil, fmt.Errorf("unable to determine neighbors of interface:%s %w", *nic.Name, err)
		}
		nic.Neighbors = neighbors
	}

	hw.Nics = nics
	hw.UUID = r.MachineUUID

	blockInfo, err := ghw.Block()
	if err != nil {
		return nil, fmt.Errorf("unable to get system block devices %w", err)
	}
	for _, disk := range blockInfo.Disks {
		if strings.HasPrefix(disk.Name, storage.DiskPrefixToIgnore) {
			continue
		}
		size := int64(disk.SizeBytes)
		diskName := disk.Name
		if !strings.HasPrefix(diskName, "/dev/") {
			diskName = fmt.Sprintf("/dev/%s", disk.Name)
		}
		blockDevice := &models.ModelsV1MachineBlockDevice{
			Name: &diskName,
			Size: &size,
		}
		hw.Disks = append(hw.Disks, blockDevice)
	}

	ipmiconfig, err := readIPMIDetails(r.Network.Eth0Mac, r.Hal)
	if err != nil {
		return nil, err
	}
	hw.Ipmi = ipmiconfig

	board := r.Hal.Board()
	b := board.BIOS
	if b == nil {
		return nil, fmt.Errorf("unable to read bios informations from bmc")
	}
	hw.Bios = &models.ModelsV1MachineBIOS{
		Version: &b.Version,
		Vendor:  &b.Vendor,
		Date:    &b.Date,
	}

	return hw, nil
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

	//nolint:gosec
	return os.WriteFile("/var/log/syslog", b[:amt], 0666)
}

// IPMI configuration and
func readIPMIDetails(eth0Mac string, hal hal.InBand) (*models.ModelsV1MachineIPMI, error) {
	var pw string
	intf := "lanplus"
	details := &models.ModelsV1MachineIPMI{
		Interface: &intf,
	}
	defaultIPMIPort := "623"
	bmcVersion := "unknown"
	bmcConn := hal.BMCConnection()
	if bmcConn.Present() {
		log.Info("ipmi details from bmc")
		board := hal.Board()
		bmc := board.BMC
		if bmc == nil {
			return nil, fmt.Errorf("unable to read ipmi bmc info configuration")
		}

		// FIXME userid should be verified if available
		pw, err := bmcConn.CreateUserAndPassword(bmcConn.User(), api.AdministratorPrivilege)
		if err != nil {
			return nil, fmt.Errorf("ipmi create user failed %w", err)
		}

		bmcUser := bmcConn.User().Name
		bmcVersion = bmc.FirmwareRevision
		fru := models.ModelsV1MachineFru{
			ChassisPartNumber:   bmc.ChassisPartNumber,
			ChassisPartSerial:   bmc.ChassisPartSerial,
			BoardMfg:            bmc.BoardMfg,
			BoardMfgSerial:      bmc.BoardMfgSerial,
			BoardPartNumber:     bmc.BoardPartNumber,
			ProductManufacturer: bmc.ProductManufacturer,
			ProductPartNumber:   bmc.ProductPartNumber,
			ProductSerial:       bmc.ProductSerial,
		}
		bmc.IP = bmc.IP + ":" + defaultIPMIPort
		details.Address = &bmc.IP
		details.Mac = &bmc.MAC
		details.User = &bmcUser
		details.Password = &pw
		details.Bmcversion = &bmcVersion
		details.Fru = &fru
		return details, nil
	}

	log.Info("ipmi details faked")
	if len(eth0Mac) == 0 {
		eth0Mac = "00:00:00:00:00:00"
	}

	macParts := strings.Split(eth0Mac, ":")
	lastOctet := macParts[len(macParts)-1]
	port, err := strconv.ParseUint(lastOctet, 16, 32)
	if err != nil {
		return nil, fmt.Errorf("unable to parse last octet of eth0 mac to a integer %w", err)
	}

	const baseIPMIPort = 6230
	// Fixed IP of vagrant environment gateway
	bmcIP := fmt.Sprintf("192.168.121.1:%d", baseIPMIPort+port)
	bmcMAC := "00:00:00:00:00:00"
	pw = "vagrant"
	user := "vagrant"
	details.Address = &bmcIP
	details.Mac = &bmcMAC
	details.User = &user
	details.Password = &pw
	details.Bmcversion = &bmcVersion
	return details, nil
}
