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
func (h *Hammer) reinstall(m *models.ModelsV1MachineResponse, hw *models.DomainMetalHammerRegisterMachineRequest, eventEmitter *event.EventEmitter) (*kernel.Bootinfo, error) {
	h.Disk = storage.GetDisk(*m.Allocation.Image.ID, m.Size, hw.Disks)
	primaryDiskName := h.Disk.Device
	if strings.HasPrefix(primaryDiskName, "/dev/") {
		primaryDiskName = primaryDiskName[5:]
	}

	info := &kernel.Bootinfo{
		Initrd:       *m.Allocation.Reinstall.Initrd,
		Cmdline:      *m.Allocation.Reinstall.Cmdline,
		Kernel:       *m.Allocation.Reinstall.Kernel,
		BootloaderID: *m.Allocation.Reinstall.Bootloaderid,
	}

	block, err := ghw.Block()
	if err != nil {
		log.Error("ghw.Block() failed", "error", err)
		return nil, errors.Wrap(err, "unable to gather disks")
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
		return info, errors.Wrapf(err, "unable to find primary disk %s", primaryDiskName)
	}

	log.Info("Wipe primary disk", "primary disk", primaryDisk.Name)
	err = storage.WipeDisk(primaryDisk)
	if err != nil {
		log.Error("failed to wipe primary disk", "error", err)
		return info, errors.Wrap(err, "wipe")
	}

	newInfo, err := h.installImage(eventEmitter, m, hw.Nics)
	if err != nil {
		return info, err
	}
	return newInfo, nil
}

func (h *Hammer) abortReinstall(reason error, bootInfo *kernel.Bootinfo) error {
	log.Error("reinstall cancelled => boot into existing OS...", "reason", reason)

	h.EventEmitter.Emit(event.ProvisioningEventReinstallAborted, fmt.Sprintf("Reinstall aborted: %s", reason.Error()))

	time.Sleep(5 * time.Second)

	return kernel.RunKexec(bootInfo)
}
