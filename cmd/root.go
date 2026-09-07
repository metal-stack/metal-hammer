package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/api/go/metalstack/infra/v2/infrav2connect"
	"github.com/metal-stack/go-hal"
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
	filesystemLayout *apiv2.FilesystemLayout
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

	eventEmitter.Emit(apiv2.MachineProvisioningEventType_MACHINE_PROVISIONING_EVENT_TYPE_PREPARING, fmt.Sprintf("starting metal-hammer version:%q", v.V))

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
		eventEmitter.Emit(apiv2.MachineProvisioningEventType_MACHINE_PROVISIONING_EVENT_TYPE_PLANNED_REBOOT, "autoreboot after 24h")
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

	reg := register.New(log, spec.MachineUUID, spec.MetalConfig.Partition, bootService, eventEmitter, n, hal)

	machineHardware, err := reg.RegisterMachine()
	if err != nil {
		return eventEmitter, fmt.Errorf("register %w", err)
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

	eventEmitter.Emit(apiv2.MachineProvisioningEventType_MACHINE_PROVISIONING_EVENT_TYPE_WAITING, "waiting for allocation")

	alloc, err := WaitForAllocation(context.Background(), log, metalAPIClient.BootService(), spec.MachineUUID, defaultWaitTimeOut)
	if err != nil {
		return eventEmitter, fmt.Errorf("wait for installation %w", err)
	}

	log.Info("perform install", "machineID", spec.MachineUUID, "imageID", alloc.Image.Id)
	hammer.filesystemLayout = alloc.FilesystemLayout
	err = hammer.installImage(eventEmitter, metalAPIClient.BootService(), alloc, machineHardware)
	return eventEmitter, err
}

func (h *hammer) installImage(eventEmitter *event.EventEmitter, bootService infrav2connect.BootServiceClient, alloc *apiv2.MachineAllocation, machineHardware *apiv2.MachineHardware) error {
	eventEmitter.Emit(apiv2.MachineProvisioningEventType_MACHINE_PROVISIONING_EVENT_TYPE_INSTALLING, "start installation")
	installationStart := time.Now()
	info, installErr := h.Install(alloc, machineHardware)

	rep := &report.Report{
		MachineUUID:     h.spec.MachineUUID,
		Client:          bootService,
		ConsolePassword: h.spec.ConsolePassword,
		InstallError:    installErr,
		Log:             h.log,
	}

	// info is nil when the installation failed
	if info != nil {
		rep.Initrd = info.Initrd
		rep.Cmdline = info.Cmdline
		rep.Kernel = info.Kernel
		rep.BootloaderID = info.BootloaderID
	}

	reportErr := rep.ReportInstallation()

	err := errors.Join(installErr, reportErr)
	if err != nil {
		return err
	}

	h.log.Info("installation", "took", time.Since(installationStart))

	// this can be useful for metal-images os debugging
	// h.log.Info("waiting 10 sec to enable os debugging")
	// time.Sleep(10 * time.Second)

	eventEmitter.Emit(apiv2.MachineProvisioningEventType_MACHINE_PROVISIONING_EVENT_TYPE_BOOTING_NEW_KERNEL, "booting into distro kernel")
	return kernel.RunKexec(info)
}
