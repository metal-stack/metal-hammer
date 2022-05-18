package cmd

import (
	"context"
	"time"

	v1 "github.com/metal-stack/metal-api/pkg/api/v1"
	"github.com/metal-stack/metal-go/api/models"
	"github.com/metal-stack/metal-hammer/pkg/kernel"
)

// fetchMachine requests the machine data of given machine ID
func (h *Hammer) fetchMachine(machineID string) (*models.V1MachineResponse, error) {

	resp, err := h.MetalAPIClient.Driver.MachineGet(machineID)
	if err != nil {
		return nil, err
	}

	return resp.Machine, nil
}

func (h *Hammer) abortReinstall(reason error, machineID string, primaryDiskWiped bool) error {
	h.log.Errorw("reinstall cancelled => boot into existing OS...", "reason", reason)

	var bootInfo *kernel.Bootinfo

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := h.MetalAPIClient.BootService().AbortReinstall(ctx, &v1.BootServiceAbortReinstallRequest{Uuid: machineID, PrimaryDiskWiped: primaryDiskWiped})
	if err != nil {
		h.log.Errorw("failed to abort reinstall", "error", err)
		time.Sleep(5 * time.Second)
	}

	if resp != nil && resp.BootInfo != nil {
		bootInfo = &kernel.Bootinfo{
			Initrd:       resp.BootInfo.Initrd,
			Cmdline:      resp.BootInfo.Cmdline,
			Kernel:       resp.BootInfo.Kernel,
			BootloaderID: resp.BootInfo.BootloaderId,
		}
	}

	return kernel.RunKexec(bootInfo)
}
