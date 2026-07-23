package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/metal-stack/metal-hammer/cmd/event"
	"github.com/metal-stack/metal-hammer/cmd/firmware"
	"github.com/metal-stack/metal-hammer/pkg/kernel"
)

// powerCycleTimeout is the time we wait for the bmc to actually cut the power,
// PowerCycle returns before the machine is down.
const powerCycleTimeout = 2 * time.Minute

// UpdateFirmware updates the firmware of components which are behind the desired firmware
// version, currently only the nvm of intel e810 based network cards. It then restarts the
// machine in the way the flashed firmware needs to become active: some firmware is live
// after a warm reboot, the nvm of the intel cards is written into an inactive bank which
// only gets activated once the card lost power, so it needs a power cycle instead.
// A firmware update must never keep this machine from being provisioned: if the restart
// can not be performed, the machine continues with the firmware it is running, the update
// is retried on the next boot. A component which is behind the desired version but does not
// take the update is retried forever, so it is reported to metal-api instead of only being
// logged on the console of this run.
func (h *hammer) UpdateFirmware() {
	if h.hal.Board().VM {
		return
	}

	result := firmware.New(h.log).Update(context.Background())

	for _, name := range result.NotFlashed {
		// this machine attempts the same update on every boot without ever getting there,
		// which is only visible in the console log of this run, so report it to metal-api.
		// Preparing is the phase this machine is in, the message carries the detail.
		h.eventEmitter.Emit(event.ProvisioningEventPreparing,
			fmt.Sprintf("the firmware of %s is behind the desired version but did not take the update", name))
	}

	switch result.Activation {
	case firmware.ActivationNone:
		// nothing was flashed, the machine keeps running the firmware it booted with.
	case firmware.ActivationWarmReboot:
		if err := h.plannedRestart("firmware", "firmware update requires a reboot", kernel.Reboot); err != nil {
			h.log.Error("firmware", "message", "reboot failed, continuing with the firmware which is currently active", "error", err)
		}
	case firmware.ActivationPowerCycle:
		if err := h.plannedRestart("firmware", "firmware update requires a power cycle", h.hal.PowerCycle); err != nil {
			// a warm reboot can not activate the new firmware, the card keeps reporting the
			// old version and would be flashed again on every boot, so rebooting here would
			// only turn an unreachable bmc into an endless reboot loop.
			h.log.Error("firmware", "message", "power cycle failed, continuing with the firmware which is currently active", "error", err)
			return
		}

		// the bmc powers the machine off asynchronously, do not continue with the installation in the meantime
		time.Sleep(powerCycleTimeout)

		h.log.Error("firmware", "message", "machine did not power cycle after a firmware update, continuing", "waited", powerCycleTimeout)
	default:
		// a firmware was flashed which asks for an activation this hammer does not know how
		// to perform, e.g. a new Activation added to the firmware package. Do not silently
		// skip the restart, the firmware stays inactive and is retried on the next boot.
		h.log.Error("firmware", "message", "unknown firmware activation, not restarting the machine", "activation", result.Activation)
	}
}
