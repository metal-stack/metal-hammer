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

	"github.com/metal-stack/api/go/enum"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-hammer/cmd/utils"
	"github.com/metal-stack/metal-lib/pkg/net"

	installerv1 "github.com/metal-stack/os-installer/api/v1"
	"github.com/metal-stack/os-installer/pkg/installer"

	img "github.com/metal-stack/metal-hammer/cmd/image"
	"github.com/metal-stack/metal-hammer/cmd/storage"
	"github.com/metal-stack/metal-hammer/pkg/chroot"
	"github.com/metal-stack/metal-hammer/pkg/kernel"
)

// Install a given image to the disk by using genuinetools/img
func (h *hammer) Install(alloc *apiv2.MachineAllocation, machineHardware *apiv2.MachineHardware) (*installerv1.Bootinfo, error) {
	s := storage.New(h.log, h.chrootPrefix, h.filesystemLayout)
	err := s.Run()
	if err != nil {
		return nil, err
	}

	image := alloc.Image.Url

	err = img.NewImage(h.log).Pull(image, h.osImageDestination)
	if err != nil {
		return nil, err
	}

	err = img.NewImage(h.log).Burn(h.chrootPrefix, image, h.osImageDestination)
	if err != nil {
		return nil, err
	}

	info, err := h.install(h.chrootPrefix, alloc, s.RootUUID, machineHardware)
	if err != nil {
		return nil, err
	}

	// This is executed after installation to be compatible with images which create fstab by their own
	// TODO can be removed and be done in s.Run() once all images do not create fstab anymore
	err = s.CreateFSTab()
	if err != nil {
		return nil, err
	}

	err = s.Umount()
	if err != nil {
		return nil, err
	}

	return info, nil
}

// install will execute Install from os-installer in the chroot where the os-image was extracted
// to finish installation e.g. install mbr, grub, write network and filesystem config
func (h *hammer) install(prefix string, alloc *apiv2.MachineAllocation, rootUUID string, machineHardware *apiv2.MachineHardware) (*installerv1.Bootinfo, error) {
	h.log.Info("install", "image", alloc.Image.Url)

	machineDetails, err := h.convertConfigs(alloc, rootUUID, machineHardware)
	if err != nil {
		return nil, fmt.Errorf("error converting configuration: %w", err)
	}

	legacyConfig, err := h.generateLegacyConfig(alloc, machineDetails, machineHardware)
	if err != nil {
		return nil, fmt.Errorf("error converting legacy configuration: %w", err)
	}

	err = h.writeUserData(alloc)
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
		err = h.writeConfigs(legacyConfig, machineDetails, alloc)
		if err != nil {
			return fmt.Errorf("error writing configuration: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		h.log.Debug("start install in chroot", "details", machineDetails, "allocation", alloc)
		i := installer.New(h.log, machineDetails, alloc)
		err = i.Install(ctx)
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

func (h *hammer) writeUserData(alloc *apiv2.MachineAllocation) error {
	configdir := path.Join(h.chrootPrefix, "etc", "metal")
	destination := path.Join(configdir, "userdata")

	base64UserData := alloc.Userdata
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

	err = i.PersistLegacyInstallYaml(legacyConfig)
	if err != nil {
		return fmt.Errorf("unable to persist configuration: %w", err)
	}

	return nil
}

func (h *hammer) generateLegacyConfig(alloc *apiv2.MachineAllocation, details *installerv1.MachineDetails, machineHardware *apiv2.MachineHardware) (*installerv1.InstallerConfig, error) {
	var vpn *installerv1.V1MachineVPN
	if alloc.Vpn != nil {
		vpn = &installerv1.V1MachineVPN{
			Address: &alloc.Vpn.ControlPlaneAddress,
			AuthKey: &alloc.Vpn.AuthKey,
		}
	}

	var (
		dnsServers []*installerv1.V1DNSServer
		ntpServers []*installerv1.V1NTPServer
	)
	for _, dns := range alloc.DnsServers {
		dnsServers = append(dnsServers, &installerv1.V1DNSServer{
			IP: &dns.Ip,
		})
	}
	for _, ntp := range alloc.NtpServers {
		ntpServers = append(ntpServers, &installerv1.V1NTPServer{
			Address: &ntp.Address,
		})
	}

	var firewallRules *installerv1.V1FirewallRules
	if alloc.FirewallRules != nil {
		var (
			egress  []*installerv1.V1FirewallEgressRule
			ingress []*installerv1.V1FirewallIngressRule
		)
		for _, e := range alloc.FirewallRules.Egress {
			protocol, err := enum.GetStringValue(e.Protocol)
			if err != nil {
				return nil, err
			}
			var ports []int32
			for _, port := range e.Ports {
				ports = append(ports, int32(port))
			}
			egress = append(egress, &installerv1.V1FirewallEgressRule{
				Comment:  e.Comment,
				Protocol: *protocol,
				Ports:    ports,
				To:       e.To,
			})
		}

		for _, i := range alloc.FirewallRules.Ingress {
			protocol, err := enum.GetStringValue(i.Protocol)
			if err != nil {
				return nil, err
			}
			var ports []int32
			for _, port := range i.Ports {
				ports = append(ports, int32(port))
			}
			ingress = append(ingress, &installerv1.V1FirewallIngressRule{
				Comment:  i.Comment,
				Protocol: *protocol,
				Ports:    ports,
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
	for _, nw := range alloc.Networks {
		var (
			nat         bool
			underlay    bool
			private     bool
			networkType string
		)
		switch nw.NetworkType {
		case apiv2.NetworkType_NETWORK_TYPE_CHILD:
			private = true
			networkType = net.PrivatePrimaryUnshared
		case apiv2.NetworkType_NETWORK_TYPE_CHILD_SHARED:
			private = true
			networkType = net.PrivatePrimaryShared
		case apiv2.NetworkType_NETWORK_TYPE_EXTERNAL:
			networkType = net.External
		case apiv2.NetworkType_NETWORK_TYPE_UNDERLAY:
			underlay = true
			networkType = net.Underlay
		}

		natTypeV2, err := enum.GetStringValue(nw.NatType)
		if err != nil {
			return nil, err
		}

		networkTypeV2, err := enum.GetStringValue(nw.NetworkType)
		if err != nil {
			return nil, err
		}

		networks = append(networks, &installerv1.V1MachineNetwork{
			Asn:                 new(int64(nw.Asn)),
			Destinationprefixes: nw.DestinationPrefixes,
			Ips:                 nw.Ips,
			Nat:                 &nat,
			Nattypev2:           natTypeV2,
			Networkid:           &nw.Network,
			Networktype:         &networkType,
			Networktypev2:       networkTypeV2,
			Prefixes:            nw.Prefixes,
			Private:             &private,
			Projectid:           nw.Project,
			Underlay:            &underlay,
			Vrf:                 new(int64(nw.Vrf)),
		})
	}

	role, err := enum.GetStringValue(alloc.AllocationType)
	if err != nil {
		return nil, err
	}

	legacyConfig := &installerv1.InstallerConfig{
		Hostname:      alloc.Hostname,
		Password:      details.Password,
		Console:       details.Console,
		RaidEnabled:   details.RaidEnabled,
		RootUUID:      details.RootUUID,
		SSHPublicKey:  strings.Join(alloc.SshPublicKeys, "\n"),
		MachineUUID:   h.spec.MachineUUID,
		Timestamp:     time.Now().Format(time.RFC3339),
		Networks:      networks,
		Nics:          h.onlyNicsWithNeighborsLegacy(machineHardware.Nics),
		VPN:           vpn,
		Role:          *role,
		FirewallRules: firewallRules,
		DNSServers:    dnsServers,
		NTPServers:    ntpServers,
	}

	return legacyConfig, nil
}

func (h *hammer) convertConfigs(alloc *apiv2.MachineAllocation, rootUUiD string, machineHardware *apiv2.MachineHardware) (*installerv1.MachineDetails, error) {
	cmdline, err := kernel.ParseCmdline()
	if err != nil {
		return nil, fmt.Errorf("unable to get kernel cmdline map %w", err)
	}

	console, ok := cmdline["console"]
	if !ok {
		console = "ttyS0"
	}

	var raidEnabled bool
	if alloc != nil && alloc.FilesystemLayout != nil && len(alloc.FilesystemLayout.Raid) > 0 {
		raidEnabled = true
	}

	machineDetails := &installerv1.MachineDetails{
		ID:          h.spec.MachineUUID,
		Password:    h.spec.ConsolePassword,
		Console:     console,
		RaidEnabled: raidEnabled,
		RootUUID:    rootUUiD,
		Nics:        machineHardware.Nics,
	}

	h.log.Info("generated apiv2 machinedetails", "details", machineDetails)

	return machineDetails, nil
}

func (h *hammer) onlyNicsWithNeighborsLegacy(nics []*apiv2.MachineNic) []*installerv1.V1MachineNic {
	noNeighbors := func(neighbors []*apiv2.MachineNic) bool {
		if len(neighbors) == 0 {
			return true
		}
		for _, n := range neighbors {
			if n.Mac == "" {
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
			Mac:        &nic.Mac,
			Name:       &nic.Name,
			Identifier: &nic.Identifier,
			Neighbors: []*installerv1.V1MachineNic{
				{
					Mac:  &nic.Neighbors[0].Mac,
					Name: &nic.Neighbors[0].Name,
				},
			},
		}
		result = append(result, n)
	}
	return result
}
