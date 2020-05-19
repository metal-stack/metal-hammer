package cmd

import (
	"github.com/metal-stack/metal-hammer/cmd/event"
	"github.com/metal-stack/metal-hammer/pkg/kernel"
	"github.com/metal-stack/metal-hammer/pkg/sum"
	"time"

	log "github.com/inconshreveable/log15"
)

// UpdateBIOS ensures that UEFI boot is enabled and CSM-support is disabled.
// It then reboots the machine.
func (h *Hammer) UpdateBIOS() error {
	s, err := sum.New()
	if err != nil {
		return err
	}

	reboot, err := s.UpdateBIOS()
	if err != nil {
		log.Warn("BIOS updates for this machine type are intentionally not supported, skipping UpdateBIOS", "error", err)
		return nil
	}
	if reboot {
		h.EventEmitter.Emit(event.ProvisioningEventPlannedReboot, "update BIOS configuration, need to reboot")

		log.Info("bios", "message", "updated BIOS configuration, reboot in 1 sec")
		time.Sleep(1 * time.Second)
		err = kernel.Reboot()
		if err != nil {
			log.Error("reboot", "error", err)
		}
	}

	return nil
}

// EnsureBootOrder ensures that the BIOS boot order is properly set,
// i.e. first boot from OS image and then PXE boot
func (h *Hammer) EnsureBootOrder(bootloaderID string) error {
	if h.Spec.DevMode {
		return nil
	}

	s, err := sum.New()
	if err != nil {
		return err
	}

	err = s.EnsureBootOrder(bootloaderID)
	if err != nil {
		return err
	}
	log.Info("bios", "message", "boot order ensured")

	return nil
}
