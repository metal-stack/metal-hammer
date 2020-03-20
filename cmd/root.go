package cmd

import (
	"strings"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	log "github.com/inconshreveable/log15"
	"github.com/metal-stack/metal-hammer/cmd/event"
	"github.com/metal-stack/metal-hammer/cmd/network"
	"github.com/metal-stack/metal-hammer/cmd/register"
	"github.com/metal-stack/metal-hammer/cmd/report"
	"github.com/metal-stack/metal-hammer/cmd/storage"
	"github.com/metal-stack/metal-hammer/metal-core/client/machine"
	"github.com/metal-stack/metal-hammer/metal-core/models"
	"github.com/metal-stack/metal-hammer/pkg/bios"
	"github.com/metal-stack/metal-hammer/pkg/kernel"
	"github.com/metal-stack/metal-hammer/pkg/os/command"
	"github.com/metal-stack/metal-hammer/pkg/password"
	"github.com/pkg/errors"
)

// Hammer is the machine which forms a bare metal to a working server
type Hammer struct {
	Client     *machine.Client
	Spec       *Specification
	Disk       storage.Disk
	LLDPClient *network.LLDPClient
	// IPAddress is the ip of the eth0 interface during installation
	IPAddress          string
	Started            time.Time
	EventEmitter       *event.EventEmitter
	ChrootPrefix       string
	OsImageDestination string
}

// Run orchestrates the whole register/wipe/format/burn and reboot process
func Run(spec *Specification) (*event.EventEmitter, error) {
	log.Info("metal-hammer run", "firmware", kernel.Firmware(), "bios", bios.Bios().String())

	transport := httptransport.New(spec.MetalCoreURL, "", nil)
	client := machine.New(transport, strfmt.Default)
	eventEmitter := event.NewEventEmitter(client, spec.MachineUUID)

	eventEmitter.Emit(event.ProvisioningEventPreparing, "starting metal-hammer")

	err := command.CommandsExist()
	if err != nil {
		return eventEmitter, err
	}

	hammer := &Hammer{
		Client:             client,
		Spec:               spec,
		IPAddress:          spec.IP,
		EventEmitter:       eventEmitter,
		ChrootPrefix:       "/rootfs",
		OsImageDestination: "/tmp/os.tgz",
	}

	// Reboot after 24Hours if no allocation was requested.
	go kernel.AutoReboot(24*time.Hour, func() {
		eventEmitter.Emit(event.ProvisioningEventPlannedReboot, "autoreboot after 24h")
	})

	hammer.Spec.ConsolePassword = password.Generate(16)

	n := &network.Network{
		MachineUUID: spec.MachineUUID,
		IPAddress:   spec.IP,
		Started:     time.Now(),
	}

	// TODO: Does not work yet, needs to be done manually
	// firmware := firmware.New()
	// firmware.Update()

	err = n.UpAllInterfaces()
	if err != nil {
		return eventEmitter, errors.Wrap(err, "interfaces")
	}

	// Set Time from ntp
	network.NtpDate()

	reg := &register.Register{
		MachineUUID: spec.MachineUUID,
		Client:      client,
		Network:     n,
	}

	hw, err := reg.ReadHardwareDetails()
	if err != nil {
		return eventEmitter, errors.Wrap(err, "unable to read all hardware details")
	}

	m, err := hammer.fetchMachine(spec.MachineUUID)
	if err == nil && m != nil && m.Allocation != nil && m.Allocation.Reinstall != nil {
		info, err := hammer.reinstall(m, hw, eventEmitter)
		if err != nil {
			log.Error("reinstall failed -> abort", "error", err)
			err = hammer.abortReinstall(err, info)
			if err != nil {
				return eventEmitter, err
			}
		}
		return eventEmitter, nil
	}

	err = storage.WipeDisks()
	if err != nil {
		return eventEmitter, errors.Wrap(err, "wipe")
	}

	eventEmitter.Emit(event.ProvisioningEventRegistering, "start registering")
	err = reg.RegisterMachine(hw)
	if !spec.DevMode && err != nil {
		return eventEmitter, errors.Wrap(err, "register")
	}

	firmware := kernel.Firmware()
	log.Info("firmware", "is", firmware)

	if firmware != "efi" && !spec.DevMode {
		log.Info("firmware is not efi, enforce efi boot mode using preparation image")
		err = hammer.EnsureUEFI()
		if err != nil {
			log.Warn("BIOS updates for this machine type are intentionally not supported, skipping EnsureUEFI", "error", err)
		}
	}

	eventEmitter.Emit(event.ProvisioningEventWaiting, "waiting for installation")

	// Ensure we can run without metal-core, given IMAGE_URL is configured as kernel cmdline
	if spec.DevMode {
		cidr := "10.0.1.2"
		if spec.Cidr != "" {
			cidr = spec.Cidr
		}

		if !spec.BGPEnabled {
			cidr = "dhcp"
		}
		asn := int64(4200000001)
		private := true
		private2 := false
		underlay := false
		underlay2 := true
		nat := false
		nat2 := true
		vrf := int64(0)
		vrf2 := int64(4200000001)
		hostname := "devmode"
		sshkeys := []string{"not a valid ssh public key, can be specified during machine create.", "second public key"}
		m = &models.ModelsV1MachineResponse{
			Allocation: &models.ModelsV1MachineAllocation{
				Image: &models.ModelsV1ImageResponse{
					URL: spec.ImageURL,
					ID:  &spec.ImageID,
				},
				Hostname:   &hostname,
				SSHPubKeys: sshkeys,
				Networks: []*models.ModelsV1MachineNetwork{
					&models.ModelsV1MachineNetwork{
						Ips:                 []string{cidr},
						Asn:                 &asn,
						Private:             &private,
						Underlay:            &underlay,
						Destinationprefixes: []string{"0.0.0.0/0"},
						Vrf:                 &vrf,
						Nat:                 &nat,
					},
					&models.ModelsV1MachineNetwork{
						Ips:                 []string{"1.2.3.4"},
						Asn:                 &asn,
						Private:             &private2,
						Underlay:            &underlay2,
						Destinationprefixes: []string{"2.3.4.5/24"},
						Vrf:                 &vrf2,
						Nat:                 &nat2,
					},
				},
			},
			Size: &models.ModelsV1SizeResponse{
				ID: &spec.SizeID,
			},
		}
		mac1 := "00:00:00:00:01:01"
		mac2 := "00:00:00:00:01:02"
		mac3 := "00:00:00:00:01:03"
		name1 := "eth0"
		name2 := "eth1"
		hw = &models.DomainMetalHammerRegisterMachineRequest{
			Nics: []*models.ModelsV1MachineNicExtended{
				{
					Mac:  &mac1,
					Name: &name1,
					Neighbors: []*models.ModelsV1MachineNicExtended{
						{
							Mac: &mac3,
						},
					},
				},
				{
					Mac:  &mac2,
					Name: &name2,
				},
			},
		}
	} else {
		m, err = hammer.Wait(spec.MachineUUID)
		if err != nil {
			return eventEmitter, errors.Wrap(err, "wait for installation")
		}
	}

	hammer.Disk = storage.GetDisk(*m.Allocation.Image.ID, m.Size, hw.Disks)
	_, err = hammer.installImage(eventEmitter, m, hw.Nics)
	return eventEmitter, err
}

func (h *Hammer) installImage(eventEmitter *event.EventEmitter, m *models.ModelsV1MachineResponse, nics []*models.ModelsV1MachineNicExtended) (*kernel.Bootinfo, error) {
	eventEmitter.Emit(event.ProvisioningEventInstalling, "start installation")
	installationStart := time.Now()
	info, err := h.Install(m, nics)

	// FIXME, must not return here.
	if err != nil {
		return info, errors.Wrap(err, "install")
	}

	var osPartition string
	for _, p := range h.Disk.Partitions {
		if p.MountPoint == "/" {
			osPartition = p.Device
			break
		}
	}
	primaryDisk := h.Disk.Device
	if strings.HasPrefix(primaryDisk, "/dev/") {
		primaryDisk = primaryDisk[5:]
	}
	rep := &report.Report{
		MachineUUID:     h.Spec.MachineUUID,
		Client:          h.Client,
		ConsolePassword: h.Spec.ConsolePassword,
		PrimaryDisk:     primaryDisk,
		OSPartition:     osPartition,
		Initrd:          info.Initrd,
		Cmdline:         info.Cmdline,
		Kernel:          info.Kernel,
		BootloaderID:    info.BootloaderID,
		InstallError:    err,
	}

	err = rep.ReportInstallation()
	if err != nil {
		wait := 10 * time.Second
		log.Error("report installation failed", "reboot in", wait, "error", err)
		time.Sleep(wait)
		if !h.Spec.DevMode {
			err = kernel.Reboot()
			if err != nil {
				log.Error("reboot", "error", err)
			}
		}
	}

	log.Info("installation", "took", time.Since(installationStart))
	eventEmitter.Emit(event.ProvisioningEventBootingNewKernel, "booting into distro kernel")
	return info, kernel.RunKexec(info)
}
