package cmd

import (
	"fmt"
	"time"

	"github.com/metal-stack/go-hal"
	v1 "github.com/metal-stack/metal-api/pkg/api/v1"
	"github.com/metal-stack/metal-go/api/models"
	"github.com/metal-stack/metal-hammer/cmd/event"
	"github.com/metal-stack/metal-hammer/cmd/network"
	"github.com/metal-stack/metal-hammer/cmd/register"
	"github.com/metal-stack/metal-hammer/cmd/report"
	"github.com/metal-stack/metal-hammer/cmd/storage"
	"github.com/metal-stack/metal-hammer/pkg/kernel"
	"github.com/metal-stack/metal-hammer/pkg/os/command"
	"github.com/metal-stack/metal-hammer/pkg/password"
	"github.com/metal-stack/v"
	"go.uber.org/zap"
)

// Hammer is the machine which forms a bare metal to a working server
type Hammer struct {
	Spec             *Specification
	log              *zap.SugaredLogger
	Hal              hal.InBand
	MetalAPIClient   *MetalAPIClient
	EventEmitter     *event.EventEmitter
	LLDPClient       *network.LLDPClient
	FilesystemLayout *models.V1FilesystemLayoutResponse
	// IPAddress is the ip of the eth0 interface during installation
	IPAddress          string
	Started            time.Time
	ChrootPrefix       string
	OsImageDestination string
}

// Run orchestrates the whole register/wipe/format/burn and reboot process
func Run(log *zap.SugaredLogger, spec *Specification, hal hal.InBand) (*event.EventEmitter, error) {
	log.Infow("metal-hammer run", "firmware", kernel.Firmware(), "bios", hal.Board().BIOS.String())
	metalAPIClient, err := NewMetalAPIClient(log, spec.PixieAPIUrl)
	if err != nil {
		log.Errorw("failed to fetch GRPC certificates", "error", err)
		return nil, err
	}

	bootService := metalAPIClient.BootService()

	eventEmitter := event.NewEventEmitter(log, metalAPIClient.Event(), spec.MachineUUID)

	eventEmitter.Emit(event.ProvisioningEventPreparing, fmt.Sprintf("starting metal-hammer version:%q", v.V))

	err = command.CommandsExist()
	if err != nil {
		return eventEmitter, err
	}

	hammer := &Hammer{
		Hal:                hal,
		Spec:               spec,
		log:                log,
		IPAddress:          spec.IP,
		EventEmitter:       eventEmitter,
		ChrootPrefix:       "/rootfs",
		OsImageDestination: "/tmp/os.tgz",
		MetalAPIClient:     metalAPIClient,
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

	reg := register.New(log, spec.MachineUUID, bootService, eventEmitter, n, hal)

	err = reg.RegisterMachine()
	if err != nil {
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
			err = hammer.installImage(eventEmitter, bootService, m)
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

	err = hammer.MetalAPIClient.WaitForAllocation(eventEmitter, spec.MachineUUID)
	if err != nil {
		return eventEmitter, fmt.Errorf("wait for installation %w", err)
	}
	m, err = hammer.fetchMachine(spec.MachineUUID)
	if err != nil {
		return eventEmitter, fmt.Errorf("wait for installation %w", err)
	}

	log.Infow("perform install", "machineID", m.ID, "imageID", *m.Allocation.Image.ID)
	hammer.FilesystemLayout = m.Allocation.Filesystemlayout
	err = hammer.installImage(eventEmitter, bootService, m)
	return eventEmitter, err
}

func (h *Hammer) installImage(eventEmitter *event.EventEmitter, bootService v1.BootServiceClient, m *models.V1MachineResponse) error {
	eventEmitter.Emit(event.ProvisioningEventInstalling, "start installation")
	installationStart := time.Now()
	info, err := h.Install(m)

	// FIXME, must not return here.
	if err != nil {
		return fmt.Errorf("install %w ", err)
	}

	// FIXME OSPartition and PrimaryDisk are not used anymore, remove from model in metal-api
	rep := &report.Report{
		MachineUUID:     h.Spec.MachineUUID,
		Client:          bootService,
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
