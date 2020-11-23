package cmd

import (
	"time"

	"github.com/metal-stack/metal-hammer/cmd/event"
	"github.com/metal-stack/metal-hammer/pkg/kernel"

	log "github.com/inconshreveable/log15"
)

// ConfigureBIOS ensures that UEFI boot is enabled and CSM-support is disabled.
// It then reboots the machine.
func (h *Hammer) ConfigureBIOS() error {
	if h.Spec.DevMode || h.Hal.Board().VM {
		return nil
	}

	reboot, err := h.Hal.ConfigureBIOS()
	if err != nil {
		return err
	}
	log.Info("bios", "message", "successfully configured BIOS")

	if reboot {
		msg := "BIOS configuration requires a reboot"
		h.EventEmitter.Emit(event.ProvisioningEventPlannedReboot, msg)
		log.Info("bios", msg, "reboot in 1 sec")
		time.Sleep(1 * time.Second)
		err = kernel.Reboot()
		if err != nil {
			return err
		}
	}

	return nil
}

// EnsureBootOrder ensures that the BIOS boot order is properly set,
// i.e. first boot from OS image and then PXE boot
func (h *Hammer) EnsureBootOrder(bootloaderID string) error {
	if h.Spec.DevMode || h.Hal.Board().VM {
		return nil
	}

	err := h.Hal.EnsureBootOrder(bootloaderID)
	if err != nil {
		return err
	}
	log.Info("bios", "message", "successfully ensured boot order")

	return nil
}
