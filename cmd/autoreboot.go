package cmd

import (
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/kernel"
	log "github.com/inconshreveable/log15"
	"time"
)

// AutoReboot will start a timer and reboot after given duration
func AutoReboot(after time.Duration) {
	log.Info("autoreboot", "after", after)
	rebootTimer := time.NewTimer(after)
	<-rebootTimer.C
	log.Info("autoreboot", "timeout reached", "rebooting in 10sec")
	time.Sleep(10 * time.Second)
	err := kernel.Reboot()
	if err != nil {
		log.Error("autoreboot", "unable to reboot, error", err)
	}
}
