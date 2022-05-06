package cmd

import (
	"fmt"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/metal-stack/go-hal"
	"github.com/metal-stack/metal-hammer/cmd/event"
	"github.com/metal-stack/metal-hammer/cmd/network"
	"github.com/metal-stack/metal-hammer/cmd/register"
	"github.com/metal-stack/metal-hammer/cmd/report"
	"github.com/metal-stack/metal-hammer/cmd/storage"
	"github.com/metal-stack/metal-hammer/metal-core/client/certs"
	"github.com/metal-stack/metal-hammer/metal-core/client/machine"
	"github.com/metal-stack/metal-hammer/metal-core/models"
	"github.com/metal-stack/metal-hammer/pkg/kernel"
	"github.com/metal-stack/metal-hammer/pkg/os/command"
	"github.com/metal-stack/metal-hammer/pkg/password"
	mn "github.com/metal-stack/metal-lib/pkg/net"
	"github.com/metal-stack/v"
	"go.uber.org/zap"
)

// Hammer is the machine which forms a bare metal to a working server
type Hammer struct {
	Spec             *Specification
	log              *zap.SugaredLogger
	Hal              hal.InBand
	Client           machine.ClientService
	GrpcClient       *GrpcClient
	EventEmitter     *event.EventEmitter
	LLDPClient       *network.LLDPClient
	FilesystemLayout *models.ModelsV1FilesystemLayoutResponse
	// IPAddress is the ip of the eth0 interface during installation
	IPAddress          string
	Started            time.Time
	ChrootPrefix       string
	OsImageDestination string
}

// Run orchestrates the whole register/wipe/format/burn and reboot process
func Run(log *zap.SugaredLogger, spec *Specification, hal hal.InBand) (*event.EventEmitter, error) {
	log.Infow("metal-hammer run", "firmware", kernel.Firmware(), "bios", hal.Board().BIOS.String())

	transport := httptransport.New(spec.MetalCoreURL, "", nil)
	client := machine.New(transport, strfmt.Default)
	certsClient := certs.New(transport, strfmt.Default)

	grpcClient, err := NewGrpcClient(log, certsClient)
	if err != nil {
		log.Errorw("failed to fetch GRPC certificates", "error", err)
		return nil, err
	}

	eventEmitter := event.NewEventEmitter(log, grpcClient.Event(), spec.MachineUUID)

	eventEmitter.Emit(event.ProvisioningEventPreparing, fmt.Sprintf("starting metal-hammer version:%q", v.V))

	err = command.CommandsExist()
	if err != nil {
		return eventEmitter, err
	}

	hammer := &Hammer{
		Hal:                hal,
		Client:             client,
		Spec:               spec,
		log:                log,
		IPAddress:          spec.IP,
		EventEmitter:       eventEmitter,
		ChrootPrefix:       "/rootfs",
		OsImageDestination: "/tmp/os.tgz",
		GrpcClient:         grpcClient,
	}

	// Reboot after 24Hours if no allocation was requested.
	go kernel.AutoReboot(log, 3*24*time.Hour, 24*time.Hour, func() {
		eventEmitter.Emit(event.ProvisioningEventPlannedReboot, "autoreboot after 24h")
	})

	hammer.Spec.ConsolePassword = password.Generate(16)

	err = hammer.createBmcSuperuser()
	if err != nil {
		log.Errorw("failed to update bmc superuser password", "error", err)
		return eventEmitter, err
	}

	n := &network.Network{
		MachineUUID: spec.MachineUUID,
		IPAddress:   spec.IP,
		Started:     time.Now(),
		Log:         log,
	}

	// TODO: Does not work yet, needs to be done manually
	// firmware := firmware.New()
	// firmware.Update()

	err = n.UpAllInterfaces()
	if err != nil {
		return eventEmitter, fmt.Errorf("interfaces %w", err)
	}

	// Set Time from ntp
	network.NtpDate(log)

	reg := &register.Register{
		MachineUUID: spec.MachineUUID,
		Client:      client,
		Network:     n,
		Hal:         hal,
		Log:         log,
	}

	hw, err := reg.ReadHardwareDetails()
	if err != nil {
		return eventEmitter, fmt.Errorf("unable to read all hardware details %w", err)
	}

	eventEmitter.Emit(event.ProvisioningEventRegistering, "start registering")
	err = reg.RegisterMachine(hw)
	if !spec.DevMode && err != nil {
		return eventEmitter, fmt.Errorf("register %w", err)
	}

	m, err := hammer.fetchMachine(spec.MachineUUID)
	if err != nil {
		return eventEmitter, fmt.Errorf("fetch %w", err)
	}
	if m != nil && m.Allocation != nil && m.Allocation.Reinstall != nil && *m.Allocation.Reinstall {
		hammer.FilesystemLayout = m.Allocation.Filesystemlayout
		primaryDiskWiped := false
		if m.Allocation.Image == nil || m.Allocation.Image.ID == nil {
			err = fmt.Errorf("no image specified")
		} else {
			log.Infow("perform reinstall", "machineID", *m.ID, "imageID", *m.Allocation.Image.ID)
			err = hammer.installImage(eventEmitter, m, hw.Nics)
			primaryDiskWiped = true
		}
		if err != nil {
			log.Errorw("reinstall failed", "error", err)
			err = hammer.abortReinstall(err, *m.ID, primaryDiskWiped)
		}
		return eventEmitter, err
	}

	err = storage.NewDisks(log).Wipe()
	if err != nil {
		return eventEmitter, fmt.Errorf("wipe %w", err)
	}

	err = hammer.ConfigureBIOS()
	if err != nil {
		log.Errorw("failed to configure BIOS", "error", err)
		return eventEmitter, err
	}

	// Ensure we can run without metal-core, given IMAGE_URL is configured as kernel cmdline
	if spec.DevMode {
		eventEmitter.Emit(event.ProvisioningEventWaiting, "waiting for installation")

		cidr := "10.0.1.2"
		if spec.Cidr != "" {
			cidr = spec.Cidr
		}

		if !spec.BGPEnabled {
			cidr = "dhcp"
		}
		asn := int64(4200000001)
		underlay := mn.Underlay
		privatePrimaryUnshared := mn.PrivatePrimaryUnshared
		boolTrue := true
		boolFalse := false
		vrf0 := int64(0)
		vrf1 := int64(4200000001)
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
					{
						Ips:                 []string{cidr},
						Asn:                 &asn,
						Private:             &boolTrue,
						Underlay:            &boolFalse,
						Networktype:         &privatePrimaryUnshared,
						Destinationprefixes: []string{"0.0.0.0/0"},
						Vrf:                 &vrf1,
						Nat:                 &boolFalse,
					},
					{
						Ips:                 []string{"1.2.3.4"},
						Asn:                 &asn,
						Private:             &boolFalse,
						Underlay:            &boolTrue,
						Networktype:         &underlay,
						Destinationprefixes: []string{"2.3.4.5/24"},
						Vrf:                 &vrf0,
						Nat:                 &boolTrue,
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
		err := hammer.GrpcClient.WaitForAllocation(eventEmitter, spec.MachineUUID)
		if err != nil {
			return eventEmitter, fmt.Errorf("wait for installation %w", err)
		}
		m, err = hammer.fetchMachine(spec.MachineUUID)
		if err != nil {
			return eventEmitter, fmt.Errorf("wait for installation %w", err)
		}
	}

	log.Infow("perform install", "machineID", m.ID, "imageID", *m.Allocation.Image.ID)
	hammer.FilesystemLayout = m.Allocation.Filesystemlayout
	err = hammer.installImage(eventEmitter, m, hw.Nics)
	return eventEmitter, err
}

func (h *Hammer) installImage(eventEmitter *event.EventEmitter, m *models.ModelsV1MachineResponse, nics []*models.ModelsV1MachineNicExtended) error {
	eventEmitter.Emit(event.ProvisioningEventInstalling, "start installation")
	installationStart := time.Now()
	info, err := h.Install(m, nics)

	// FIXME, must not return here.
	if err != nil {
		return fmt.Errorf("install %w ", err)
	}

	// FIXME OSPartition and PrimaryDisk are not used anymore, remove from model in metal-api
	rep := &report.Report{
		MachineUUID:     h.Spec.MachineUUID,
		Client:          h.Client,
		ConsolePassword: h.Spec.ConsolePassword,
		Initrd:          info.Initrd,
		Cmdline:         info.Cmdline,
		Kernel:          info.Kernel,
		BootloaderID:    info.BootloaderID,
		InstallError:    err,
		Log:             h.log,
	}

	err = rep.ReportInstallation()
	if err != nil {
		return err
	}

	h.log.Infow("installation", "took", time.Since(installationStart))
	eventEmitter.Emit(event.ProvisioningEventBootingNewKernel, "booting into distro kernel")
	return kernel.RunKexec(info)
}
