package cmd

import (
	"fmt"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/ipmi"

	log "github.com/inconshreveable/log15"
)

// EnsureUEFI check if the boot firmware is set to uefi when booting via pxe permanent.
// If not already set, make required modifications and reboot the machine.
func (h *Hammer) EnsureUEFI() error {
	i := ipmi.New()

	if !i.DevicePresent() {
		log.Info("ensureUEFI: we are virtual with no real ipmi device present, ignoring.")
		return nil
	}

	if i.UEFIEnabled() && i.BootOptionsPersistent() {
		log.Info("ensureUEFI: all requirements are met.")
		return nil
	}

	err := i.EnableUEFI(ipmi.PXE, true)
	if err != nil {
		return fmt.Errorf("unable to ensureUEFI %v", err)
	}

	log.Info("ensureUEFI: set uefi persistent, reboot now.")
	if h.Spec.DevMode {
		log.Warn("ensureUEFI required reboot skipped in devmode")
		return nil
	}
	err = pkg.Reboot()
	if err != nil {
		log.Error("reboot", "error", err)
	}
	return nil
}
