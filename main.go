package main

import (
	"fmt"

	"github.com/metal-pod/v"

	"os"
	"time"

	"git.f-i-ts.de/cloud-native/metal/metal-hammer/cmd"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/cmd/event"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/cmd/network"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/kernel"
	log "github.com/inconshreveable/log15"
	"github.com/pkg/errors"
)

func main() {
	fmt.Print(cmd.HammerBanner)
	// Reboot if metal-hammer crashes after 60sec.
	go kernel.Watchdog()
	ip := network.InternalIP()
	err := cmd.StartSSHD(ip)
	if err != nil {
		log.Error("sshd error", "error", err)
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		log.Error("cmd args are not supported")
		os.Exit(1)
	}

	log.Info("metal-hammer", "version", v.V)

	spec := cmd.NewSpec(ip)
	spec.Log()

	var level log.Lvl
	if spec.Debug {
		level = log.LvlDebug
	} else {
		level = log.LvlInfo
	}

	h := log.CallerFileHandler(log.StdoutHandler)
	h = log.LvlFilterHandler(level, h)
	log.Root().SetHandler(h)

	emitter, err := cmd.Run(spec)
	if err != nil {
		wait := 5 * time.Second
		st := errors.WithStack(err)
		fmt.Printf("%+v", st)
		log.Error("metal-hammer failed", "rebooting in", wait, "error", err)
		emitter.Emit(event.ProvisioningEventCrashed, fmt.Sprintf("%s", err))
		time.Sleep(wait)
		err := kernel.Reboot()
		if err != nil {
			log.Error("metal-hammer reboot failed", "error", err)
			emitter.Emit(event.ProvisioningEventCrashed, fmt.Sprintf("%s", err))
		}
	}
}
