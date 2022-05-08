package register

import (
	"context"
	"fmt"
	gonet "net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/jaypipes/ghw"
	"github.com/metal-stack/go-hal"
	"github.com/metal-stack/go-hal/pkg/api"
	v1 "github.com/metal-stack/metal-api/pkg/api/v1"
	"github.com/metal-stack/metal-hammer/cmd/network"
	"github.com/metal-stack/metal-hammer/cmd/storage"
	"github.com/vishvananda/netlink"
	"go.uber.org/zap"
)

// Register the Machine
type Register struct {
	MachineUUID string
	Client      v1.BootServiceClient
	Network     *network.Network
	Hal         hal.InBand
	Log         *zap.SugaredLogger
}

// RegisterMachine register a machine at the metal-api via metal-core
func (r *Register) RegisterMachine(hw *v1.BootServiceRegisterRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := r.Client.Register(ctx, hw)

	if err != nil {
		return fmt.Errorf("unable to register machine:%#v %w", hw, err)
	}
	if resp == nil {
		return fmt.Errorf("unable to register machine:%#v response payload is nil", hw)
	}

	r.Log.Infow("register machine returned", "response", resp)
	return nil
}

// ReadHardwareDetails returns the hardware details of the machine
func (r *Register) ReadHardwareDetails() (*v1.BootServiceRegisterRequest, error) {
	err := createSyslog()
	if err != nil {
		return nil, fmt.Errorf("unable to write kernel boot message to /var/log/syslog %w", err)
	}

	res := &v1.BootServiceRegisterRequest{}
	hw := &v1.MachineHardware{}

	memory, err := ghw.Memory()
	if err != nil {
		return nil, fmt.Errorf("unable to get system memory %w", err)
	}
	hw.Memory = uint64(memory.TotalPhysicalBytes)

	// FIXME can be replaced by runtime.NumCPU()
	cpu, err := ghw.CPU()
	if err != nil {
		return nil, fmt.Errorf("unable to get system cpu(s) %w", err)
	}
	hw.CpuCores = uint32(cpu.TotalCores)

	nics := []*v1.MachineNic{}
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
			r.Log.Debugw("skip interface with invalid mac", "interface", name, "mac", mac)
			continue
		}
		// check if after mac validation loopback is still present
		if name == "lo" {
			loFound = true
		}
		if name == "eth0" {
			r.Network.Eth0Mac = mac
		}

		nic := &v1.MachineNic{
			Mac:  mac,
			Name: name,
		}
		r.Log.Infow("register", "nic", name, "mac", mac)
		nics = append(nics, nic)
	}
	// add a lo interface if not present
	// this is required to have this interface present
	// in our DCIM management to add a ip later.
	if !loFound {
		mac := "00:00:00:00:00:00"
		name := "lo"
		lo := &v1.MachineNic{
			Mac:  mac,
			Name: name,
		}
		nics = append(nics, lo)
	}

	// now attach neighbors, this will wait up to 2*tx-intervall
	// if during this timeout not all required neighbors where found abort and reboot.
	for _, nic := range nics {
		neighbors, err := r.Network.Neighbors(nic.Name)
		if err != nil {
			return nil, fmt.Errorf("unable to determine neighbors of interface:%s %w", nic.Name, err)
		}
		ns := []*v1.MachineNic{}
		for i := range neighbors {

			ns = append(ns, &v1.MachineNic{
				Mac:  *neighbors[i].Mac,
				Name: *neighbors[i].Name,
			})
		}

		nic.Neighbor = ns
	}

	hw.Nics = nics
	res.Uuid = r.MachineUUID

	blockInfo, err := ghw.Block()
	if err != nil {
		return nil, fmt.Errorf("unable to get system block devices %w", err)
	}
	for _, disk := range blockInfo.Disks {
		if strings.HasPrefix(disk.Name, storage.DiskPrefixToIgnore) {
			continue
		}
		size := uint64(disk.SizeBytes)
		diskName := disk.Name
		if !strings.HasPrefix(diskName, "/dev/") {
			diskName = fmt.Sprintf("/dev/%s", disk.Name)
		}
		blockDevice := &v1.MachineBlockDevice{
			Name: diskName,
			Size: size,
		}
		hw.Disks = append(hw.Disks, blockDevice)
	}

	ipmiconfig, err := readIPMIDetails(r.Log, r.Network.Eth0Mac, r.Hal)
	if err != nil {
		return nil, err
	}
	res.Ipmi = ipmiconfig

	board := r.Hal.Board()
	b := board.BIOS
	if b == nil {
		return nil, fmt.Errorf("unable to read bios informations from bmc")
	}
	res.Bios = &v1.MachineBIOS{
		Version: b.Version,
		Vendor:  b.Vendor,
		Date:    b.Date,
	}

	return res, nil
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
func readIPMIDetails(log *zap.SugaredLogger, eth0Mac string, hal hal.InBand) (*v1.MachineIPMI, error) {
	var pw string
	intf := "lanplus"
	details := &v1.MachineIPMI{
		Interface: intf,
	}
	defaultIPMIPort := "623"
	bmcVersion := "unknown"
	bmcConn := hal.BMCConnection()
	if bmcConn.Present() {
		log.Infow("ipmi details from bmc")
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
		fru := &v1.MachineFRU{
			ChassisPartNumber:   &bmc.ChassisPartNumber,
			ChassisPartSerial:   &bmc.ChassisPartSerial,
			BoardMfg:            &bmc.BoardMfg,
			BoardMfgSerial:      &bmc.BoardMfgSerial,
			BoardPartNumber:     &bmc.BoardPartNumber,
			ProductManufacturer: &bmc.ProductManufacturer,
			ProductPartNumber:   &bmc.ProductPartNumber,
			ProductSerial:       &bmc.ProductSerial,
		}
		bmc.IP = bmc.IP + ":" + defaultIPMIPort
		details.Address = bmc.IP
		details.Mac = bmc.MAC
		details.User = bmcUser
		details.Password = pw
		details.BmcVersion = bmcVersion
		details.Fru = fru
		return details, nil
	}

	log.Infow("ipmi details faked")
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
	details.Address = bmcIP
	details.Mac = bmcMAC
	details.User = user
	details.Password = pw
	details.BmcVersion = bmcVersion
	return details, nil
}
