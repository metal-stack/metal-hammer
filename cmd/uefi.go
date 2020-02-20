package cmd

import (
	"time"

	"github.com/metal-stack/metal-hammer/pkg/sum"

	"github.com/metal-stack/metal-hammer/cmd/event"
	"github.com/metal-stack/metal-hammer/pkg/kernel"

	log "github.com/inconshreveable/log15"
)

// EnsureUEFI ensures that UEFI boot is enabled and reboots the machine
// if required and not in dev mode.
func (h *Hammer) EnsureUEFI() error {
	s, err := sum.New()
	if err != nil {
		return err
	}

	err = s.EnsureUEFIBoot(false)
	if err != nil {
		return err
	}

	h.EventEmitter.Emit(event.ProvisioningEventPlannedReboot, "update BIOS configuration, need to reboot to get uefi set")
	log.Info("uefi", "message", "set UEFI boot, reboot in 3 sec")
	time.Sleep(3 * time.Second)

	err = kernel.Reboot()
	if err != nil {
		log.Error("reboot", "error", err)
	}
	return nil
}

// EnsureBootOrder ensures that the BIOS boot order is properly set,
// i.e. first boot from OS image and then PXE boot.
// Will be skipped in dev mode.
func (h *Hammer) EnsureBootOrder(bootloaderID string) error {
	if h.Spec.DevMode {
		return nil
	}

	s, err := sum.New()
	if err != nil {
		return err
	}

	err = s.EnsureBootOrder(bootloaderID, false)
	if err != nil {
		return err
	}
	log.Info("uefi", "message", "boot order ensured")

	return nil
}
