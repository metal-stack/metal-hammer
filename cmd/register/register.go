package register

import (
	"context"
	"fmt"
	"log/slog"
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
	"github.com/metal-stack/metal-hammer/cmd/event"
	"github.com/metal-stack/metal-hammer/cmd/network"
	"github.com/metal-stack/metal-hammer/cmd/storage"
	"github.com/metal-stack/v"
	"github.com/u-root/u-root/pkg/pci"
	"github.com/vishvananda/netlink"
)

// Register the Machine
type Register struct {
	machineUUID string
	client      v1.BootServiceClient
	emitter     *event.EventEmitter
	network     *network.Network
	inband      hal.InBand
	log         *slog.Logger
}

func New(log *slog.Logger, machineID string, bootClient v1.BootServiceClient, emitter *event.EventEmitter, network *network.Network, inband hal.InBand) *Register {
	return &Register{
		machineUUID: machineID,
		client:      bootClient,
		emitter:     emitter,
		network:     network,
		inband:      inband,
		log:         log,
	}
}

// RegisterMachine register a machine at the metal-api via metal-api
func (r *Register) RegisterMachine() error {
	r.emitter.Emit(event.ProvisioningEventRegistering, "start registering")
	req, err := r.readHardwareDetails()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := r.client.Register(ctx, req)

	if err != nil {
		return fmt.Errorf("unable to register machine:%#v %w", req, err)
	}
	if resp == nil {
		return fmt.Errorf("unable to register machine:%#v response payload is nil", req)
	}

	r.log.Info("machine registered", "response", resp)
	return nil
}

// ReadHardwareDetails returns the hardware details of the machine
func (r *Register) readHardwareDetails() (*v1.BootServiceRegisterRequest, error) {
	err := createSyslog()
	if err != nil {
		return nil, fmt.Errorf("unable to write kernel boot message to /var/log/syslog %w", err)
	}
	memory, err := ghw.Memory()
	if err != nil {
		return nil, fmt.Errorf("unable to get system memory %w", err)
	}
	// FIXME can be replaced by runtime.NumCPU()
	cpu, err := ghw.CPU()
	if err != nil {
		return nil, fmt.Errorf("unable to get system cpu(s) %w", err)
	}
	r.log.Info("cpu", "processors", cpu.String())
	var metalCPUs []*v1.MachineCPU
	for _, cpu := range cpu.Processors {
		metalCPUs = append(metalCPUs, &v1.MachineCPU{
			Vendor:  cpu.Vendor,
			Model:   cpu.Model,
			Cores:   cpu.NumCores,
			Threads: cpu.NumThreads,
		})
	}

	// 0000:bd:00.0: DisplayVGA: NVIDIA Corporation AD102GL [RTX 6000 Ada Generation]

	gpus, err := r.detectGPUs()
	if err != nil {
		return nil, fmt.Errorf("unable to get system gpu(s) %w", err)
	}

	var metalGPUs []*v1.MachineGPU
	for _, g := range gpus {
		r.log.Info("found gpu", "gpu", g.String())
		metalGPUs = append(metalGPUs, &v1.MachineGPU{
			Vendor: g.VendorName,
			Model:  g.DeviceName,
		})
	}

	// Nics
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
			r.log.Debug("skip interface with invalid mac", "interface", name, "mac", mac)
			continue
		}
		// check if after mac validation loopback is still present
		if name == "lo" {
			loFound = true
		}
		if name == "eth0" {
			r.network.Eth0Mac = mac
		}

		nic := &v1.MachineNic{
			Mac:  mac,
			Name: name,
		}
		r.log.Info("register", "nic", name, "mac", mac)
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

	// now attach neighbors, this will wait up to 2*tx-interval
	// if during this timeout not all required neighbors where found abort and reboot.
	for _, nic := range nics {
		r.log.Info("register search neighbor for", "nic", nic.Name)
		neighbors, err := r.network.Neighbors(nic.Name)
		if err != nil {
			return nil, fmt.Errorf("unable to determine neighbors of interface:%s %w", nic.Name, err)
		}
		r.log.Info("register found neighbor for", "nic", nic.Name, "neighbors", neighbors)
		nic.Neighbors = neighbors
	}

	// Disks
	blockInfo, err := ghw.Block()
	if err != nil {
		return nil, fmt.Errorf("unable to get system block devices %w", err)
	}
	disks := []*v1.MachineBlockDevice{}
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
		disks = append(disks, blockDevice)
	}

	hardware := &v1.MachineHardware{
		Memory: uint64(memory.TotalPhysicalBytes),
		Nics:   nics,
		Disks:  disks,
		Cpus:   metalCPUs,
		Gpus:   metalGPUs,
	}

	// IPMI
	ipmi, err := r.readIPMIDetails()
	if err != nil {
		return nil, err
	}

	// Bios
	board := r.inband.Board()
	b := board.BIOS
	if b == nil {
		return nil, fmt.Errorf("unable to read bios information from bmc")
	}
	bios := &v1.MachineBIOS{
		Version: b.Version,
		Vendor:  b.Vendor,
		Date:    b.Date,
	}

	request := &v1.BootServiceRegisterRequest{
		Uuid:               r.machineUUID,
		Hardware:           hardware,
		Bios:               bios,
		Ipmi:               ipmi,
		MetalHammerVersion: v.Version,
	}

	r.log.Info("register", "request", request)
	return request, nil
}

func (r *Register) detectGPUs() (pci.Devices, error) {
	pciReader, err := pci.NewBusReader("*")
	if err != nil {
		return nil, err
	}

	var devices pci.Devices
	if devices, err = pciReader.Read(); err != nil {
		return nil, err
	}

	devices.SetVendorDeviceName()

	var result pci.Devices
	for _, device := range devices {
		// "vendor":"NVIDIA Corporation","device":"AD102GL [RTX 6000 Ada Generation]"}
		if !strings.Contains(strings.ToLower(device.VendorName), "nvidia") {
			continue
		}

		r.log.Info("add gpu", "vendor", device.VendorName, "device", device.DeviceName)
		result = append(result, device)
	}

	return result, nil
}

// save the content of kernel ring buffer to /var/log/syslog
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
func (r *Register) readIPMIDetails() (*v1.MachineIPMI, error) {
	var pw string
	intf := "lanplus"
	details := &v1.MachineIPMI{
		Interface: intf,
	}
	defaultIPMIPort := "623"
	bmcVersion := "unknown"
	bmcConn := r.inband.BMCConnection()
	if bmcConn.Present() {
		r.log.Info("ipmi details from bmc")
		board := r.inband.Board()
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

	r.log.Info("ipmi details faked")
	eth0Mac := r.network.Eth0Mac
	if len(r.network.Eth0Mac) == 0 {
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
