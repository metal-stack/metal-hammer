package cmd

import (
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/cmd/event"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/ipmi"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/kernel"
	"time"

	"github.com/pkg/errors"

	log "github.com/inconshreveable/log15"
)

// EnsureUEFI check if the boot firmware is set to uefi when booting via pxe permanent.
// If not already set, make required modifications and reboot the machine.
func (h *Hammer) EnsureUEFI() error {
	firmware := kernel.Firmware()
	if firmware == "efi" {
		log.Info("uefi", "message", "machine booted with efi, no action")
		return nil
	}

	i := ipmi.New()

	if !i.DevicePresent() {
		log.Info("uefi", "message", "no ipmi device present, no action")
		return nil
	}

	if i.BootOptionsPersistent() {
		log.Info("uefi", "message", "all requirements are met, no action")
		return nil
	}

	err := i.EnableUEFI(ipmi.PXE, true)
	if err != nil {
		return errors.Wrap(err, "unable to ensureUEFI")
	}

	log.Warn("uefi", "message", "set persistent, reboot in 10 sec.")
	if h.Spec.DevMode {
		log.Warn("required reboot skipped", "devmode", h.Spec.DevMode)
		return nil
	}

	h.EventEmitter.Emit(event.ProvisioningEventPlannedReboot, "need to reboot to get uefi set")
	time.Sleep(10 * time.Second)

	err = kernel.Reboot()
	if err != nil {
		log.Error("reboot", "error", err)
	}
	return nil
}
