package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/metal-stack/go-hal"
	v1 "github.com/metal-stack/metal-api/pkg/api/v1"
	apigrpc "github.com/metal-stack/metal-api/pkg/grpc"
	"github.com/metal-stack/metal-go/api/client/machine"
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
)

const defaultWaitTimeOut = 3 * time.Second

// hammer is the machine which forms a bare metal to a working server
type hammer struct {
	log              *slog.Logger
	spec             *Specification
	hal              hal.InBand
	metalAPIClient   *MetalAPIClient
	eventEmitter     *event.EventEmitter
	filesystemLayout *models.V1FilesystemLayoutResponse
	// IPAddress is the ip of the eth0 interface during installation
	chrootPrefix       string
	osImageDestination string
}

// Run orchestrates the whole register/wipe/format/burn and reboot process
func Run(log *slog.Logger, spec *Specification, hal hal.InBand) (*event.EventEmitter, error) {
	log.Info("metal-hammer run", "firmware", kernel.Firmware(), "bios", hal.Board().BIOS.String())
	metalAPIClient, err := NewMetalAPIClient(log, spec)
	if err != nil {
		log.Error("failed to fetch GRPC certificates", "error", err)
		return nil, err
	}

	bootService := metalAPIClient.BootService()

	eventEmitter := event.NewEventEmitter(log, metalAPIClient.Event(), spec.MachineUUID)

	eventEmitter.Emit(event.ProvisioningEventPreparing, fmt.Sprintf("starting metal-hammer version:%q", v.V))

	err = command.CommandsExist()
	if err != nil {
		return eventEmitter, err
	}

	hammer := &hammer{
		hal:                hal,
		spec:               spec,
		log:                log,
		eventEmitter:       eventEmitter,
		chrootPrefix:       "/rootfs",
		osImageDestination: "/tmp/os.tgz",
		metalAPIClient:     metalAPIClient,
	}

	// Reboot after 24Hours if no allocation was requested.
	go kernel.AutoReboot(log, 1*24*time.Hour, 24*time.Hour, func() {
		eventEmitter.Emit(event.ProvisioningEventPlannedReboot, "autoreboot after 24h")
	})

	hammer.spec.ConsolePassword = password.Generate(16)

	err = hammer.createBmcSuperuser()
	if err != nil {
		log.Error("failed to update bmc superuser password", "error", err)
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
	network.NtpDate(log, spec.MetalConfig.NTPServers)

	reg := register.New(log, spec.MachineUUID, spec.MetalConfig.Partition, bootService, eventEmitter, n, hal)

	err = reg.RegisterMachine()
	if err != nil {
		return eventEmitter, fmt.Errorf("register %w", err)
	}

	resp, err := metalAPIClient.Machine().FindMachine(machine.NewFindMachineParams().WithID(spec.MachineUUID), nil)
	if err != nil {
		return eventEmitter, fmt.Errorf("fetch %w", err)
	}
	m := resp.Payload
	if m != nil && m.Allocation != nil && m.Allocation.Reinstall != nil && *m.Allocation.Reinstall {
		hammer.filesystemLayout = m.Allocation.Filesystemlayout
		primaryDiskWiped := false
		if m.Allocation.Image == nil || m.Allocation.Image.ID == nil {
			err = fmt.Errorf("no image specified")
		} else {
			log.Info("perform reinstall", "machineID", *m.ID, "imageID", *m.Allocation.Image.ID)
			err = hammer.installImage(eventEmitter, bootService, m)
			primaryDiskWiped = true
		}
		if err != nil {
			log.Error("reinstall failed", "error", err)
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
		log.Error("failed to configure BIOS", "error", err)
		return eventEmitter, err
	}

	eventEmitter.Emit(event.ProvisioningEventWaiting, "waiting for allocation")

	err = apigrpc.WaitForAllocation(context.Background(), log, metalAPIClient.BootService(), spec.MachineUUID, defaultWaitTimeOut)
	if err != nil {
		return eventEmitter, fmt.Errorf("wait for installation %w", err)
	}

	resp, err = metalAPIClient.Machine().FindMachine(machine.NewFindMachineParams().WithID(spec.MachineUUID), nil)
	if err != nil {
		return eventEmitter, fmt.Errorf("wait for installation %w", err)
	}
	m = resp.Payload

	log.Info("perform install", "machineID", m.ID, "imageID", *m.Allocation.Image.ID)
	hammer.filesystemLayout = m.Allocation.Filesystemlayout
	err = hammer.installImage(eventEmitter, bootService, m)
	return eventEmitter, err
}

func (h *hammer) installImage(eventEmitter *event.EventEmitter, bootService v1.BootServiceClient, m *models.V1MachineResponse) error {
	eventEmitter.Emit(event.ProvisioningEventInstalling, "start installation")
	installationStart := time.Now()
	info, err := h.Install(m)

	// FIXME, must not return here.
	if err != nil {
		return fmt.Errorf("install %w ", err)
	}

	rep := &report.Report{
		MachineUUID:     h.spec.MachineUUID,
		Client:          bootService,
		ConsolePassword: h.spec.ConsolePassword,
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

	h.log.Info("installation", "took", time.Since(installationStart))

	// this can be useful for metal-images os debugging
	// h.log.Info("waiting 10 sec to enable os debugging")
	// time.Sleep(10 * time.Second)

	eventEmitter.Emit(event.ProvisioningEventBootingNewKernel, "booting into distro kernel")
	return kernel.RunKexec(info)
}
