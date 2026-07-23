package cmd

import (
	"time"

	"github.com/metal-stack/metal-hammer/cmd/event"
	"github.com/metal-stack/metal-hammer/pkg/kernel"
)

// restartDelay gives the event emitter a moment to ship the planned reboot event
// before the machine goes down.
const restartDelay = 1 * time.Second

// ConfigureBIOS ensures that UEFI boot is enabled and CSM-support is disabled.
// It then reboots the machine.
func (h *hammer) ConfigureBIOS() error {
	if h.hal.Board().VM {
		return nil
	}

	reboot, err := h.hal.ConfigureBIOS()
	if err != nil {
		return err
	}
	h.log.Info("bios", "message", "successfully configured BIOS")

	if reboot {
		return h.plannedRestart("bios", "BIOS configuration requires a reboot", kernel.Reboot)
	}

	return nil
}

// plannedRestart announces a restart which is required to apply a configuration or
// firmware change and then performs it with the given restart function.
func (h *hammer) plannedRestart(component, msg string, restart func() error) error {
	h.eventEmitter.Emit(event.ProvisioningEventPlannedReboot, msg)
	h.log.Info(component, "message", msg, "restarting in", restartDelay)
	time.Sleep(restartDelay)

	return restart()
}

// EnsureBootOrder ensures that the BIOS boot order is properly set,
// i.e. first boot from OS image and then PXE boot
func (h *hammer) EnsureBootOrder(bootloaderID string) error {
	if h.hal.Board().VM {
		return nil
	}

	err := h.hal.EnsureBootOrder(bootloaderID)
	if err != nil {
		return err
	}
	h.log.Info("bios", "message", "successfully ensured boot order")

	return nil
}
