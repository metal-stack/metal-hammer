package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-hammer/cmd/utils"
	"github.com/metal-stack/metal-lib/pkg/pointer"

	installerv1 "github.com/metal-stack/os-installer/api/v1"
	"github.com/metal-stack/os-installer/pkg/installer"

	"github.com/metal-stack/metal-go/api/models"
	img "github.com/metal-stack/metal-hammer/cmd/image"
	"github.com/metal-stack/metal-hammer/cmd/storage"
	"github.com/metal-stack/metal-hammer/pkg/chroot"
	"github.com/metal-stack/metal-hammer/pkg/kernel"
)

// Install a given image to the disk by using genuinetools/img
func (h *hammer) Install(machine *models.V1MachineResponse) (*installerv1.Bootinfo, error) {
	s := storage.New(h.log, h.chrootPrefix, *h.filesystemLayout)
	err := s.Run()
	if err != nil {
		return nil, err
	}

	image := machine.Allocation.Image.URL

	err = img.NewImage(h.log).Pull(image, h.osImageDestination)
	if err != nil {
		return nil, err
	}

	err = img.NewImage(h.log).Burn(h.chrootPrefix, image, h.osImageDestination)
	if err != nil {
		return nil, err
	}

	info, err := h.install(h.chrootPrefix, machine, s.RootUUID)
	if err != nil {
		return nil, err
	}

	// This is executed after installation to be compatible with images which create fstab by their own
	// TODO can be removed and be done in s.Run() once all images do not create fstab anymore
	err = s.CreateFSTab()
	if err != nil {
		return nil, err
	}

	s.Umount()

	return info, nil
}

// install will execute Install from os-installer in the chroot where the os-image was extracted
// to finish installation e.g. install mbr, grub, write network and filesystem config
func (h *hammer) install(prefix string, machine *models.V1MachineResponse, rootUUID string) (*installerv1.Bootinfo, error) {
	h.log.Info("install", "image", machine.Allocation.Image.URL)

	machineDetails, machineAllocation, err := h.convertConfigs(machine, rootUUID)
	if err != nil {
		return nil, fmt.Errorf("error converting configuration: %w", err)
	}

	legacyConfig, err := h.generateLegacyConfig(machine, machineDetails, rootUUID)
	if err != nil {
		return nil, fmt.Errorf("error converting legacy configuration: %w", err)
	}

	err = h.writeUserData(machine)
	if err != nil {
		return nil, fmt.Errorf("writing userdata failed %w", err)
	}

	err = h.writeLVMLocalConf()
	if err != nil {
		return nil, err
	}

	// Write configuration to /etc/metal and execute the installer in chroot
	if err := chroot.RunInChroot(h.log, prefix, func() error {
		h.log.Debug("write configs in chroot")
		err = h.writeConfigs(legacyConfig, machineDetails, machineAllocation)
		if err != nil {
			return fmt.Errorf("error writing configuration: %w", err)
		}

		h.log.Debug("start install in chroot", "details", machineDetails, "allocation", machineAllocation)
		i := installer.New(h.log, machineDetails, machineAllocation)
		err = i.Install(context.TODO())
		if err != nil {
			h.log.Error("error during install", "error", err)
			return fmt.Errorf("error during install: %w", err)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("unable to run the installer %w", err)
	}

	h.log.Info("finish running the installer")

	info, err := kernel.ReadBootinfo(path.Join(prefix, "etc", "metal", "boot-info.yaml"))
	if err != nil {
		return info, fmt.Errorf("unable to read boot-info.yaml %w", err)
	}

	h.log.Info("bootinfo", "info", info)

	err = h.EnsureBootOrder(info.BootloaderID)
	if err != nil {
		return info, fmt.Errorf("unable to ensure boot order %w", err)
	}

	tmp := "/tmp"
	_, err = utils.Copy(path.Join(prefix, info.Kernel), path.Join(tmp, filepath.Base(info.Kernel)))
	if err != nil {
		h.log.Error("could not copy kernel", "error", err)
		return info, err
	}
	info.Kernel = path.Join(tmp, filepath.Base(info.Kernel))

	if info.Initrd == "" {
		return info, nil
	}

	_, err = utils.Copy(path.Join(prefix, info.Initrd), path.Join(tmp, filepath.Base(info.Initrd)))
	if err != nil {
		h.log.Error("could not copy initrd", "error", err)
		return info, err
	}
	info.Initrd = path.Join(tmp, filepath.Base(info.Initrd))

	return info, nil
}

// writeLVMLocalConf to make lvm more compatible with os without udevd
// will only be written if lvm is installed in the target image
func (h *hammer) writeLVMLocalConf() error {
	srclvmlocal := "/etc/lvm/lvmlocal.conf"
	dstlvm := path.Join(h.chrootPrefix, "/etc/lvm")
	dstlvmlocal := path.Join(h.chrootPrefix, srclvmlocal)

	_, err := os.Stat(srclvmlocal) // FIXME use fileExists below
	if os.IsNotExist(err) {
		h.log.Info("src lvmlocal.conf not present, not creating lvmlocal.conf")
		return nil
	}
	_, err = os.Stat(dstlvm) // FIXME use fileExists below
	if os.IsNotExist(err) {
		h.log.Info("dst /etc/lvm not present, not creating lvmlocal.conf")
		return nil
	}

	input, err := os.ReadFile(srclvmlocal)
	if err != nil {
		return fmt.Errorf("unable to read lvmlocal.conf %w", err)
	}

	err = os.WriteFile(dstlvmlocal, input, 0600)
	if err != nil {
		return fmt.Errorf("unable to write lvmlocal.conf %w", err)
	}
	return nil
}

func (h *hammer) writeUserData(machine *models.V1MachineResponse) error {
	configdir := path.Join(h.chrootPrefix, "etc", "metal")
	destination := path.Join(configdir, "userdata")

	base64UserData := machine.Allocation.UserData
	if base64UserData != "" {
		userdata, err := base64.StdEncoding.DecodeString(base64UserData)
		if err != nil {
			h.log.Info("install", "base64 decode of userdata failed, using plain text", err)
			userdata = []byte(base64UserData)
		}
		return os.WriteFile(destination, userdata, 0600)
	}
	return nil
}

func (h *hammer) writeConfigs(legacyConfig *installerv1.InstallerConfig, details *installerv1.MachineDetails, allocation *apiv2.MachineAllocation) error {
	h.log.Info("write installation configuration")
	configdir := path.Join("etc", "metal")
	err := os.MkdirAll(configdir, 0755)
	if err != nil {
		return fmt.Errorf("mkdir of %s target os failed %w", configdir, err)
	}

	i := installer.New(h.log, details, allocation)

	err = i.PersistConfigurations()
	if err != nil {
		return fmt.Errorf("unable to persist configuration: %w", err)
	}

	err = i.PersistLegacyInstallYaml(legacyConfig)
	if err != nil {
		return fmt.Errorf("unable to persist configuration: %w", err)
	}

	return nil
}

func (h *hammer) generateLegacyConfig(machine *models.V1MachineResponse, details *installerv1.MachineDetails, rootUUiD string) (*installerv1.InstallerConfig, error) {
	var vpn *installerv1.V1MachineVPN
	if machine.Allocation.Vpn != nil {
		vpn = &installerv1.V1MachineVPN{
			Address: machine.Allocation.Vpn.Address,
			AuthKey: machine.Allocation.Vpn.AuthKey,
		}
	}

	var (
		dnsServers []*installerv1.V1DNSServer
		ntpServers []*installerv1.V1NTPServer
	)
	for _, dns := range machine.Allocation.DNSServers {
		dnsServers = append(dnsServers, &installerv1.V1DNSServer{
			IP: dns.IP,
		})
	}
	for _, ntp := range machine.Allocation.NtpServers {
		ntpServers = append(ntpServers, &installerv1.V1NTPServer{
			Address: ntp.Address,
		})
	}

	var firewallRules *installerv1.V1FirewallRules
	if machine.Allocation.FirewallRules != nil {
		var (
			egress  []*installerv1.V1FirewallEgressRule
			ingress []*installerv1.V1FirewallIngressRule
		)
		for _, e := range machine.Allocation.FirewallRules.Egress {
			egress = append(egress, &installerv1.V1FirewallEgressRule{
				Comment:  e.Comment,
				Protocol: e.Protocol,
				Ports:    e.Ports,
				To:       e.To,
			})
		}

		for _, i := range machine.Allocation.FirewallRules.Ingress {
			ingress = append(ingress, &installerv1.V1FirewallIngressRule{
				Comment:  i.Comment,
				Protocol: i.Protocol,
				Ports:    i.Ports,
				To:       i.To,
				From:     i.From,
			})
		}

		firewallRules = &installerv1.V1FirewallRules{
			Egress:  egress,
			Ingress: ingress,
		}
	}

	var networks []*installerv1.V1MachineNetwork
	for _, nw := range machine.Allocation.Networks {
		networks = append(networks, &installerv1.V1MachineNetwork{
			Asn:                 nw.Asn,
			Destinationprefixes: nw.Destinationprefixes,
			Ips:                 nw.Ips,
			Nat:                 nw.Nat,
			Nattypev2:           nw.Nattypev2,
			Networkid:           nw.Networkid,
			Networktype:         nw.Networktype,
			Networktypev2:       nw.Nattypev2,
			Prefixes:            nw.Prefixes,
			Private:             nw.Private,
			Projectid:           nw.Projectid,
			Underlay:            nw.Underlay,
			Vrf:                 nw.Vrf,
		})
	}

	legacyConfig := &installerv1.InstallerConfig{
		Hostname:      pointer.SafeDeref(machine.Allocation.Hostname),
		Password:      details.Password,
		Console:       details.Console,
		RaidEnabled:   details.RaidEnabled,
		RootUUID:      details.RootUUID,
		SSHPublicKey:  strings.Join(machine.Allocation.SSHPubKeys, "\n"),
		MachineUUID:   h.spec.MachineUUID,
		Timestamp:     time.Now().Format(time.RFC3339),
		Networks:      networks,
		Nics:          h.onlyNicsWithNeighborsLegacy(machine.Hardware.Nics),
		VPN:           vpn,
		Role:          pointer.SafeDeref(machine.Allocation.Role),
		FirewallRules: firewallRules,
		DNSServers:    dnsServers,
		NTPServers:    ntpServers,
	}

	return legacyConfig, nil
}

func (h *hammer) convertConfigs(machine *models.V1MachineResponse, rootUUiD string) (*installerv1.MachineDetails, *apiv2.MachineAllocation, error) {
	alloc := machine.Allocation

	cmdline, err := kernel.ParseCmdline()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get kernel cmdline map %w", err)
	}

	console, ok := cmdline["console"]
	if !ok {
		console = "ttyS0"
	}

	var raidEnabled bool
	if alloc != nil && alloc.Filesystemlayout != nil && len(alloc.Filesystemlayout.Raid) > 0 {
		raidEnabled = true
	}

	machineDetails := &installerv1.MachineDetails{
		ID:          h.spec.MachineUUID,
		Password:    h.spec.ConsolePassword,
		Console:     console,
		RaidEnabled: raidEnabled,
		RootUUID:    rootUUiD,
		Nics:        h.onlyNicsWithNeighbors(machine.Hardware.Nics),
	}

	h.log.Info("generated apiv2 machinedetails", "details", machineDetails)

	var vpn *apiv2.MachineVPN
	if alloc.Vpn != nil {
		vpn = &apiv2.MachineVPN{
			ControlPlaneAddress: pointer.SafeDeref(alloc.Vpn.Address),
			AuthKey:             pointer.SafeDeref(alloc.Vpn.AuthKey),
			Connected:           pointer.SafeDeref(alloc.Vpn.Connected),
		}
	}

	allocationType := apiv2.MachineAllocationType_MACHINE_ALLOCATION_TYPE_MACHINE
	if alloc.Role != nil && *alloc.Role == "firewall" {
		allocationType = apiv2.MachineAllocationType_MACHINE_ALLOCATION_TYPE_FIREWALL
	}

	var dnsservers []*apiv2.DNSServer
	for _, dns := range alloc.DNSServers {
		dnsservers = append(dnsservers, &apiv2.DNSServer{
			Ip: pointer.SafeDeref(dns.IP),
		})
	}
	var ntpservers []*apiv2.NTPServer
	for _, ntp := range alloc.NtpServers {
		ntpservers = append(ntpservers, &apiv2.NTPServer{
			Address: pointer.SafeDeref(ntp.Address),
		})
	}

	var firewallRules *apiv2.FirewallRules
	if alloc.FirewallRules != nil {
		var egressrules []*apiv2.FirewallEgressRule

		for _, egress := range alloc.FirewallRules.Egress {
			var proto apiv2.IPProtocol
			if egress.Protocol == "tcp" {
				proto = apiv2.IPProtocol_IP_PROTOCOL_TCP
			}
			if egress.Protocol == "udp" {
				proto = apiv2.IPProtocol_IP_PROTOCOL_UDP
			}
			var ports []uint32
			for _, port := range egress.Ports {
				ports = append(ports, uint32(port))
			}

			egressrules = append(egressrules, &apiv2.FirewallEgressRule{
				Comment:  egress.Comment,
				Protocol: proto,
				Ports:    ports,
				To:       egress.To,
			})
		}

		var ingressrules []*apiv2.FirewallIngressRule
		for _, ingress := range alloc.FirewallRules.Ingress {
			var proto apiv2.IPProtocol
			if ingress.Protocol == "tcp" {
				proto = apiv2.IPProtocol_IP_PROTOCOL_TCP
			}
			if ingress.Protocol == "udp" {
				proto = apiv2.IPProtocol_IP_PROTOCOL_UDP
			}
			var ports []uint32
			for _, port := range ingress.Ports {
				ports = append(ports, uint32(port))
			}

			ingressrules = append(ingressrules, &apiv2.FirewallIngressRule{
				Comment:  ingress.Comment,
				Protocol: proto,
				Ports:    ports,
				To:       ingress.To,
				From:     ingress.From,
			})
		}

		firewallRules = &apiv2.FirewallRules{
			Egress:  egressrules,
			Ingress: ingressrules,
		}
	}

	var networks []*apiv2.MachineNetwork
	for _, nw := range alloc.Networks {

		natType := apiv2.NATType_NAT_TYPE_NONE
		if nw.Nat != nil && *nw.Nat {
			natType = apiv2.NATType_NAT_TYPE_IPV4_MASQUERADE
		}

		var networkType apiv2.NetworkType
		switch pointer.SafeDeref(nw.Networktypev2) {
		case "external":
			networkType = apiv2.NetworkType_NETWORK_TYPE_EXTERNAL
		case "underlay":
			networkType = apiv2.NetworkType_NETWORK_TYPE_UNDERLAY
		case "super":
			networkType = apiv2.NetworkType_NETWORK_TYPE_SUPER
		case "super-namespaced":
			networkType = apiv2.NetworkType_NETWORK_TYPE_SUPER_NAMESPACED
		case "child":
			networkType = apiv2.NetworkType_NETWORK_TYPE_CHILD
		case "child-shared":
			networkType = apiv2.NetworkType_NETWORK_TYPE_CHILD_SHARED
		}

		networks = append(networks, &apiv2.MachineNetwork{
			Network:             pointer.SafeDeref(nw.Networkid),
			Prefixes:            nw.Prefixes,
			DestinationPrefixes: nw.Destinationprefixes,
			Ips:                 nw.Ips,
			Vrf:                 uint64(pointer.SafeDeref(nw.Vrf)),
			Asn:                 uint32(pointer.SafeDeref(nw.Asn)),
			Project:             nw.Projectid,
			NatType:             natType,
			NetworkType:         networkType,
		})
	}

	machineAllocation := &apiv2.MachineAllocation{
		Uuid:        pointer.SafeDeref(alloc.Allocationuuid),
		Name:        pointer.SafeDeref(alloc.Name),
		Description: alloc.Description,
		CreatedBy:   pointer.SafeDeref(alloc.Creator),
		Project:     pointer.SafeDeref(alloc.Project),
		Image: &apiv2.Image{
			Id:  pointer.SafeDeref(alloc.Image.ID),
			Url: alloc.Image.URL,
		},
		Hostname:       pointer.SafeDeref(alloc.Hostname),
		SshPublicKeys:  alloc.SSHPubKeys,
		Userdata:       alloc.UserData,
		AllocationType: allocationType,
		FirewallRules:  firewallRules,
		Networks:       networks,
		DnsServers:     dnsservers,
		NtpServers:     ntpservers,
		Vpn:            vpn,
	}

	h.log.Info("generated apiv2 allocation", "allocation", machineAllocation)

	return machineDetails, machineAllocation, nil
}

func (h *hammer) onlyNicsWithNeighbors(nics []*models.V1MachineNic) []*apiv2.MachineNic {
	noNeighbors := func(neighbors []*models.V1MachineNic) bool {
		if len(neighbors) == 0 {
			return true
		}
		for _, n := range neighbors {
			if n.Mac == nil || *n.Mac == "" {
				return true
			}
		}
		return false
	}

	result := []*apiv2.MachineNic{}
	for i := range nics {
		nic := nics[i]
		if noNeighbors(nic.Neighbors) {
			continue
		}
		n := &apiv2.MachineNic{
			Mac:        pointer.SafeDeref(nic.Mac),
			Name:       pointer.SafeDeref(nic.Name),
			Identifier: pointer.SafeDeref(nic.Identifier),
			Neighbors: []*apiv2.MachineNic{
				{Mac: pointer.SafeDeref(nic.Neighbors[0].Mac), Name: pointer.SafeDeref(nic.Neighbors[0].Name)},
			},
		}
		result = append(result, n)
	}
	h.log.Info("onlyNicWithNeighbors add", "result", result)
	return result
}

func (h *hammer) onlyNicsWithNeighborsLegacy(nics []*models.V1MachineNic) []*installerv1.V1MachineNic {
	noNeighbors := func(neighbors []*models.V1MachineNic) bool {
		if len(neighbors) == 0 {
			return true
		}
		for _, n := range neighbors {
			if n.Mac == nil || *n.Mac == "" {
				return true
			}
		}
		return false
	}

	result := []*installerv1.V1MachineNic{}
	for i := range nics {
		nic := nics[i]
		if noNeighbors(nic.Neighbors) {
			continue
		}
		n := &installerv1.V1MachineNic{
			Mac:        nic.Mac,
			Name:       nic.Name,
			Identifier: nic.Identifier,
			Neighbors: []*installerv1.V1MachineNic{
				{Mac: nic.Neighbors[0].Mac, Name: nic.Neighbors[0].Name},
			},
		}
		result = append(result, n)
	}
	return result
}
