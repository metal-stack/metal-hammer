package cmd

import (
	"time"

	"git.f-i-ts.de/cloud-native/metal/metal-hammer/cmd/event"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/cmd/firmware"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/cmd/network"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/cmd/register"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/cmd/report"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/cmd/storage"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/client/machine"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/metal-core/models"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/kernel"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/password"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	log "github.com/inconshreveable/log15"
	"github.com/pkg/errors"
)

// Hammer is the machine which forms a bare metal to a working server
type Hammer struct {
	Client     *machine.Client
	Spec       *Specification
	Disk       storage.Disk
	LLDPClient *network.LLDPClient
	// IPAddress is the ip of the eth0 interface during installation
	IPAddress    string
	Started      time.Time
	EventEmitter *event.EventEmitter
}

// Run orchestrates the whole register/wipe/format/burn and reboot process
func Run(spec *Specification) (*event.EventEmitter, error) {
	log.Info("metal-hammer run", "firmware", kernel.Firmware())

	transport := httptransport.New(spec.MetalCoreURL, "", nil)
	client := machine.New(transport, strfmt.Default)
	eventEmitter := event.NewEventEmitter(client, spec.MachineUUID)

	eventEmitter.Emit(event.ProvisioningEventPreparing, "starting metal-hammer")

	hammer := &Hammer{
		Client:       client,
		Spec:         spec,
		IPAddress:    spec.IP,
		EventEmitter: eventEmitter,
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

	firmware := firmware.New()
	firmware.Update()

	lsi := storage.NewStorcli()
	err := lsi.EnableJBOD()
	if err != nil {
		log.Warn("root", "unable to format raid controller", err)
	}

	err = n.UpAllInterfaces()
	if err != nil {
		return eventEmitter, errors.Wrap(err, "interfaces")
	}

	// Set Time from ntp
	network.NtpDate()

	err = storage.WipeDisks()
	if err != nil {
		return eventEmitter, errors.Wrap(err, "wipe")
	}

	reg := &register.Register{
		MachineUUID: spec.MachineUUID,
		Client:      client,
		Network:     n,
	}

	eventEmitter.Emit(event.ProvisioningEventRegistering, "start registering")
	// Remove uuid return use MachineUUID() above.
	hw, uuid, err := reg.RegisterMachine()
	if !spec.DevMode && err != nil {
		return eventEmitter, errors.Wrap(err, "register")
	}

	err = hammer.EnsureUEFI()
	if err != nil {
		return eventEmitter, errors.Wrap(err, "uefi")
	}

	eventEmitter.Emit(event.ProvisioningEventWaiting, "waiting for installation")

	// Ensure we can run without metal-core, given IMAGE_URL is configured as kernel cmdline
	var machine *models.ModelsV1MachineResponse
	if spec.DevMode {
		cidr := "10.0.1.2/24"
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
		machine = &models.ModelsV1MachineResponse{
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
		machine, err = hammer.Wait(uuid)
		if err != nil {
			return eventEmitter, errors.Wrap(err, "wait for installation")
		}
	}

	hammer.Disk = storage.GetDisk(machine.Allocation.Image, machine.Size, hw.Disks)

	eventEmitter.Emit(event.ProvisioningEventInstalling, "start installation")
	installationStart := time.Now()
	info, err := hammer.Install(machine, hw)

	// FIXME, must not return here.
	if err != nil {
		return eventEmitter, errors.Wrap(err, "install")
	}

	rep := &report.Report{
		MachineUUID:     spec.MachineUUID,
		Client:          client,
		ConsolePassword: spec.ConsolePassword,
		InstallError:    err,
	}

	err = rep.ReportInstallation()
	if err != nil {
		wait := 10 * time.Second
		log.Error("report installation failed", "reboot in", wait, "error", err)
		time.Sleep(wait)
		if !spec.DevMode {
			err = kernel.Reboot()
			if err != nil {
				log.Error("reboot", "error", err)
			}
		}
	}

	log.Info("installation", "took", time.Since(installationStart))
	eventEmitter.Emit(event.ProvisioningEventBootingNewKernel, "booting into distro kernel")
	return eventEmitter, kernel.RunKexec(info)
}
