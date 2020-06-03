package main

import (
	"fmt"
	"github.com/metal-stack/v"
	"os"
	"time"

	log "github.com/inconshreveable/log15"
	"github.com/metal-stack/go-hal/detect"
	"github.com/metal-stack/metal-hammer/cmd"
	"github.com/metal-stack/metal-hammer/cmd/event"
	"github.com/metal-stack/metal-hammer/cmd/network"
	"github.com/metal-stack/metal-hammer/pkg/kernel"
	"github.com/pkg/errors"
)

func main() {
	fmt.Print(cmd.HammerBanner)

	if len(os.Args) > 1 {
		log.Error("cmd args are not supported")
		os.Exit(1)
	}

	err := updateResolvConf()
	if err != nil {
		log.Error("error updating resolv.conf", "error", err)
		os.Exit(1)
	}

	hal, err := detect.ConnectInBand()
	if err != nil {
		log.Error("unable to detect hardware", "error", err)
		os.Exit(1)
	}

	uuid, err := hal.UUID()
	if err != nil {
		log.Error("unable to get uuid hardware", "error", err)
		os.Exit(1)
	}

	ip := network.InternalIP()
	err = cmd.StartSSHD(ip)
	if err != nil {
		log.Error("sshd error", "error", err)
		os.Exit(1)
	}

	// Reboot if metal-hammer crashes after 60sec.
	go kernel.Watchdog()

	log.Info("metal-hammer", "version", v.V, "hal", hal.Describe())

	spec := cmd.NewSpec()
	spec.MachineUUID = uuid.String()
	spec.IP = ip

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

	emitter, err := cmd.Run(spec, hal)
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

func updateResolvConf() error {
	// when starting the metal-hammer u-root sets a static resolv.conf file containing 8.8.8.8
	// this can only be overriden by running dhclient
	// however, we can use the dhcp information that the kernel used during startup
	// this information is contained in /proc/net/pnp
	//
	// https://www.kernel.org/doc/Documentation/filesystems/nfs/nfsroot.txt
	symlink := "/etc/resolv.conf"
	target := "/proc/net/pnp"

	err := os.Remove(symlink)
	if err != nil {
		return err
	}

	return os.Symlink(target, symlink)
}
