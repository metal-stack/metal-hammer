package cmd

import (
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-hammer/pkg/kernel"
)

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
		msg := "BIOS configuration requires a reboot"
		h.eventEmitter.Emit(apiv2.MachineProvisioningEventType_MACHINE_PROVISIONING_EVENT_TYPE_PLANNED_REBOOT, msg)
		h.log.Info("bios", msg, "reboot in 1 sec")
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
