package cmd

import (
	"fmt"
	log "github.com/inconshreveable/log15"
	"github.com/jaypipes/ghw"
	"github.com/metal-stack/metal-hammer/cmd/event"
	"github.com/metal-stack/metal-hammer/cmd/storage"
	"github.com/metal-stack/metal-hammer/metal-core/client/machine"
	"github.com/metal-stack/metal-hammer/metal-core/models"
	"github.com/metal-stack/metal-hammer/pkg/kernel"
	"github.com/pkg/errors"
	"strings"
	"time"
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

// wipe only the disk that has the OS installed on one of its partitions, keep all other disks untouched
func (h *Hammer) reinstall(m *models.ModelsV1MachineResponse, hw *models.DomainMetalHammerRegisterMachineRequest, eventEmitter *event.EventEmitter) (bool, error) {
	if m.Allocation.BootInfo == nil {
		return false, errors.New("machine is not yet ready for reinstallations")
	}
	h.Disk = storage.GetDisk(*m.Allocation.BootInfo.ImageID, m.Size, hw.Disks)
	currentPrimaryDiskName := h.Disk.Device
	h.Disk = storage.GetDisk(*m.Allocation.Image.ID, m.Size, hw.Disks)
	primaryDiskName := h.Disk.Device
	if currentPrimaryDiskName != primaryDiskName {
		return false, fmt.Errorf("current primary disk %s differs from the one that  %s", currentPrimaryDiskName, primaryDiskName)
	}
	if strings.HasPrefix(primaryDiskName, "/dev/") {
		primaryDiskName = primaryDiskName[5:]
	}

	block, err := ghw.Block()
	if err != nil {
		log.Error("ghw.Block() failed", "error", err)
		return false, errors.Wrap(err, "unable to gather disks")
	}
	var primaryDisk *ghw.Disk
	for _, d := range block.Disks {
		if primaryDiskName == d.Name {
			primaryDisk = d
			break
		}
	}
	if primaryDisk == nil {
		log.Warn("Unable to find primary disk", "primary disk", primaryDiskName)
		return false, errors.Wrapf(err, "unable to find primary disk %s", primaryDiskName)
	}

	log.Info("Wipe primary disk", "primary disk", primaryDisk.Name)
	err = storage.WipeDisk(primaryDisk)
	if err != nil {
		log.Error("failed to wipe primary disk", "error", err)
		return false, errors.Wrap(err, "wipe")
	}

	return true, h.installImage(eventEmitter, m, hw.Nics)
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
