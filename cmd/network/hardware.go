package network

import (
	log "github.com/inconshreveable/log15"
	"github.com/jaypipes/ghw"
	"github.com/metal-stack/metal-hammer/cmd/storage"
	"github.com/metal-stack/metal-hammer/metal-core/models"
	"github.com/metal-stack/metal-hammer/pkg/bios"
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
	"io/ioutil"
	gonet "net"
	"strings"
	"syscall"
	"unsafe"
)

// ReadHardwareDetails returns the hardware details of the machine
func (n *Network) ReadHardwareDetails() (*models.DomainMetalHammerRegisterMachineRequest, error) {
	err := createSyslog()
	if err != nil {
		return nil, errors.Wrap(err, "unable to write kernel boot message to /var/log/syslog")
	}

	hw := &models.DomainMetalHammerRegisterMachineRequest{}

	memory, err := ghw.Memory()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get system memory")
	}
	hw.Memory = &memory.TotalPhysicalBytes

	// FIXME can be replaced by runtime.NumCPU()
	cpu, err := ghw.CPU()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get system cpu(s)")
	}
	cores := int32(cpu.TotalCores)
	hw.CPUCores = &cores

	nics := []*models.ModelsV1MachineNicExtended{}
	loFound := false
	links, err := netlink.LinkList()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get all links")
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
			n.Eth0Mac = mac
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
		neighbors, err := n.Neighbors(*nic.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to determine neighbors of interface:%s", *nic.Name)
		}
		nic.Neighbors = neighbors
	}

	hw.Nics = nics
	hw.UUID = n.MachineUUID

	blockInfo, err := ghw.Block()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get system block devices")
	}
	for _, disk := range blockInfo.Disks {
		if strings.HasPrefix(disk.Name, storage.DiskPrefixToIgnore) {
			continue
		}
		var parts []*models.ModelsV1MachineDiskPartition
		containsOS := false // not yet allocated
		for _, p := range blockInfo.Partitions {
			size := int64(p.SizeBytes)
			parts = append(parts, &models.ModelsV1MachineDiskPartition{
				Filesystem: &p.Type,
				Device:     &p.Name,
				Label:      &p.Label,
				Mountpoint: &p.MountPoint,
				Containsos: &containsOS,
				Size:       &size,
			})
		}
		primary := false // not allocated yet
		size := int64(disk.SizeBytes)
		blockDevice := &models.ModelsV1MachineBlockDevice{
			Name:       &disk.Name,
			Size:       &size,
			Primary:    &primary,
			Partitions: parts,
		}
		hw.Disks = append(hw.Disks, blockDevice)
	}

	ipmiconfig, err := readIPMIDetails(n.Eth0Mac)
	if err != nil {
		return nil, err
	}
	hw.IPMI = ipmiconfig

	b := bios.Bios()
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

	return ioutil.WriteFile("/var/log/syslog", b[:amt], 0666)
}
