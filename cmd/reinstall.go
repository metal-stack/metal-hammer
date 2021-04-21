package cmd

import (
	"strings"
	"time"

	log "github.com/inconshreveable/log15"
	"github.com/metal-stack/metal-hammer/metal-core/client/machine"
	"github.com/metal-stack/metal-hammer/metal-core/models"
	"github.com/metal-stack/metal-hammer/pkg/kernel"
)

// fetchMachine requests the machine data of given machine ID
func (h *Hammer) fetchMachine(machineID string) (*models.ModelsV1MachineResponse, error) {
	params := machine.NewFindMachineParams()
	params.SetID(machineID)
	resp, err := h.Client.FindMachine(params)
	if err != nil {
		return nil, err
	}

	return resp.Payload, nil
}

func (h *Hammer) abortReinstall(reason error, machineID string, primaryDiskWiped bool) error {
	log.Error("reinstall cancelled => boot into existing OS...", "reason", reason)

	params := machine.NewAbortReinstallParams()
	params.ID = machineID
	params.Body = &models.DomainMetalHammerAbortReinstallRequest{
		PrimaryDiskWiped: &primaryDiskWiped,
	}

	var bootInfo *kernel.Bootinfo

	resp, err := h.Client.AbortReinstall(params)
	if err != nil {
		log.Error("failed to abort reinstall", "error", err)
		time.Sleep(5 * time.Second)
	}

	if resp != nil && resp.Payload != nil {
		bootInfo = &kernel.Bootinfo{
			Initrd:       *resp.Payload.Initrd,
			Cmdline:      *resp.Payload.Cmdline,
			Kernel:       *resp.Payload.Kernel,
			BootloaderID: *resp.Payload.Bootloaderid,
		}
	}

	return kernel.RunKexec(bootInfo)
}

func isValidBootInfo(m *models.ModelsV1MachineResponse) bool {
	if m.Allocation == nil || m.Allocation.BootInfo == nil {
		return false
	}

	pd := m.Allocation.BootInfo.PrimaryDisk
	if pd != nil && *pd != "" {
		return true
	}
	id := m.Allocation.BootInfo.ImageID
	if id != nil && *id != "" {
		return true
	}

	return false
}

func sanitizeDisk(disk string) string {
	if strings.HasPrefix(disk, "/dev/") {
		return disk[5:]
	}
	return disk
}
