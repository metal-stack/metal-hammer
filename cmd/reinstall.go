package cmd

import (
	"context"
	"time"

	v1 "github.com/metal-stack/metal-api/pkg/api/v1"
	kernelapi "github.com/metal-stack/metal-hammer/pkg/api"
	"github.com/metal-stack/metal-hammer/pkg/kernel"
)

func (h *hammer) abortReinstall(reason error, machineID string, primaryDiskWiped bool) error {
	h.log.Error("reinstall cancelled => boot into existing OS...", "reason", reason)

	var bootInfo *kernelapi.Bootinfo

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := h.metalAPIClient.BootService().AbortReinstall(ctx, &v1.BootServiceAbortReinstallRequest{Uuid: machineID, PrimaryDiskWiped: primaryDiskWiped})
	if err != nil {
		h.log.Error("failed to abort reinstall", "error", err)
		time.Sleep(5 * time.Second)
	}

	if resp != nil && resp.BootInfo != nil {
		bootInfo = &kernelapi.Bootinfo{
			Initrd:       resp.BootInfo.Initrd,
			Cmdline:      resp.BootInfo.Cmdline,
			Kernel:       resp.BootInfo.Kernel,
			BootloaderID: resp.BootInfo.BootloaderId,
		}
	}

	return kernel.RunKexec(bootInfo)
}
