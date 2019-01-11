package cmd

import (
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/ipmi"
	"github.com/pkg/errors"

	log "github.com/inconshreveable/log15"
)

// EnsureUEFI check if the boot firmware is set to uefi when booting via pxe permanent.
// If not already set, make required modifications and reboot the machine.
func (h *Hammer) EnsureUEFI() error {
	i := ipmi.New()

	if !i.DevicePresent() {
		log.Info("uefi no ipmi device present, no action")
		return nil
	}

	if i.UEFIEnabled() && i.BootOptionsPersistent() {
		log.Info("uefi all requirements are met, no action")
		return nil
	}

	err := i.EnableUEFI(ipmi.PXE, true)
	if err != nil {
		return errors.Wrap(err, "unable to ensureUEFI")
	}

	log.Info("uefi set persistent, reboot now.")
	if h.Spec.DevMode {
		log.Warn("required reboot skipped", "devmode", h.Spec.DevMode)
		return nil
	}
	err = pkg.Reboot()
	if err != nil {
		log.Error("reboot", "error", err)
	}
	return nil
}
