package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"path/filepath"
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
	"gopkg.in/yaml.v3"
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

// install will execute /install.sh in the pulled docker image which was extracted onto disk
// to finish installation e.g. install mbr, grub, write network and filesystem config
func (h *hammer) install(prefix string, machine *models.V1MachineResponse, rootUUID string) (*installerv1.Bootinfo, error) {
	h.log.Info("install", "image", machine.Allocation.Image.URL)

	machineDetails, machineAllocation, err := h.writeInstallerConfig(machine, rootUUID)
	if err != nil {
		return nil, fmt.Errorf("writing configuration install.yaml failed %w", err)
	}

	err = h.writeUserData(machine)
	if err != nil {
		return nil, fmt.Errorf("writing userdata failed %w", err)
	}

	err = h.writeLVMLocalConf()
	if err != nil {
		return nil, err
	}

	if err := chroot.RunInChroot(h.log, prefix, func() error {
		return installer.Install(context.TODO(), h.log, machineDetails, machineAllocation)
	}); err != nil {
		return nil, fmt.Errorf("unable to run the installer %w", err)
	}

	h.log.Info("finish running the installer")

	info, err := kernel.ReadBootinfo(path.Join(prefix, "etc", "metal", "boot-info.yaml"))
	if err != nil {
		return info, fmt.Errorf("unable to read boot-info.yaml %w", err)
	}

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

func (h *hammer) writeInstallerConfig(machine *models.V1MachineResponse, rootUUiD string) (*installerv1.MachineDetails, *apiv2.MachineAllocation, error) {
	h.log.Info("write installation configuration")
	configdir := path.Join(h.chrootPrefix, "etc", "metal")
	err := os.MkdirAll(configdir, 0755)
	if err != nil {
		return nil, nil, fmt.Errorf("mkdir of %s target os failed %w", configdir, err)
	}

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

	lldpdConfig := &installerv1.LLDPDConfig{
		MachineUUID: h.spec.MachineUUID,
		Timestamp:   time.Now().Format(time.RFC3339),
	}
	yamlContent, err := yaml.Marshal(lldpdConfig)
	if err != nil {
		return nil, nil, err
	}

	err = os.WriteFile(installerv1.LLDPDConfigPath, yamlContent, os.ModePerm)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to write lldpd config %w", err)
	}

	machineDetails := &installerv1.MachineDetails{
		ID:          h.spec.MachineUUID,
		Password:    h.spec.ConsolePassword,
		Console:     console,
		RaidEnabled: raidEnabled,
		RootUUID:    rootUUiD,
		Nics:        h.onlyNicsWithNeighbors(machine.Hardware.Nics),
	}

	yamlContent, err = yaml.Marshal(machineDetails)
	if err != nil {
		return nil, nil, err
	}

	err = os.WriteFile(installerv1.MachineDetailsPath, yamlContent, os.ModePerm)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to write machine details %w", err)
	}

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

	var dnsserver []*apiv2.DNSServer
	for _, dns := range alloc.DNSServers {
		dnsserver = append(dnsserver, &apiv2.DNSServer{
			Ip: pointer.SafeDeref(dns.IP),
		})
	}
	var ntpserver []*apiv2.NTPServer
	for _, ntp := range alloc.NtpServers {
		ntpserver = append(ntpserver, &apiv2.NTPServer{
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
		switch pointer.SafeDeref(nw.Nattypev2) {
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
		Uuid:        *alloc.Allocationuuid,
		Name:        *alloc.Name,
		Description: alloc.Description,
		CreatedBy:   *alloc.Creator,
		Project:     *alloc.Project,
		Image: &apiv2.Image{
			Id:  *alloc.Image.ID,
			Url: alloc.Image.URL,
		},
		Hostname:       *alloc.Hostname,
		SshPublicKeys:  alloc.SSHPubKeys,
		Userdata:       alloc.UserData,
		AllocationType: allocationType,
		FirewallRules:  firewallRules,
		Networks:       networks,
		DnsServer:      dnsserver,
		NtpServer:      ntpserver,
		Vpn:            vpn,
	}

	yamlContent, err = yaml.Marshal(machineAllocation)
	if err != nil {
		return nil, nil, err
	}

	err = os.WriteFile(installerv1.MachineAllocationPath, yamlContent, os.ModePerm)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to write machine allocation %w", err)
	}

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
		}
		result = append(result, n)
	}
	h.log.Info("onlyNicWithNeighbors add", "result", result)
	return result
}
