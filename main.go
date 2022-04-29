package main

import (
	"fmt"
	"os"
	"time"

	"github.com/metal-stack/v"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/metal-stack/go-hal/connect"
	halzap "github.com/metal-stack/go-hal/pkg/logger/zap"
	"github.com/metal-stack/metal-hammer/cmd"
	"github.com/metal-stack/metal-hammer/cmd/event"
	"github.com/metal-stack/metal-hammer/cmd/network"
	"github.com/metal-stack/metal-hammer/pkg/kernel"
)

func main() {
	fmt.Print(cmd.HammerBanner)
	if len(os.Args) > 1 {
		panic("cmd args are not supported")
	}

	err := updateResolvConf()
	if err != nil {
		fmt.Printf("error updating resolv.conf %s", err)
		os.Exit(1)
	}

	zcfg := zap.NewProductionConfig()
	lvl, err := zap.ParseAtomicLevel("debug")
	if err != nil {
		fmt.Printf("unable to set log level %s", err)
		os.Exit(1)
	}
	zcfg.Level = lvl
	zcfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zlog, err := zcfg.Build()
	if err != nil {
		panic(err)
	}
	log := zlog.Sugar()

	// Reboot if metal-hammer crashes after 60sec.
	go kernel.Watchdog(log)

	hal, err := connect.InBand(halzap.New(log))
	if err != nil {
		log.Errorw("unable to detect hardware", "error", err)
		os.Exit(1)
	}

	uuid, err := hal.UUID()
	if err != nil {
		log.Errorw("unable to get uuid hardware", "error", err)
		os.Exit(1)
	}

	ip := network.InternalIP()
	err = cmd.StartSSHD(log, ip)
	if err != nil {
		log.Errorw("sshd error", "error", err)
		os.Exit(1)
	}

	log.Infow("metal-hammer", "version", v.V, "hal", hal.Describe())

	spec := cmd.NewSpec(log)
	spec.MachineUUID = uuid.String()
	spec.IP = ip

	spec.Log()

	// FIXME set loglevel from spec.Debug

	emitter, err := cmd.Run(log, spec, hal)
	if err != nil {
		wait := 5 * time.Second
		log.Errorw("metal-hammer failed", "rebooting in", wait, "error", err)
		emitter.Emit(event.ProvisioningEventCrashed, fmt.Sprintf("%s", err))
		time.Sleep(wait)
		err := kernel.Reboot()
		if err != nil {
			log.Errorw("metal-hammer reboot failed", "error", err)
			emitter.Emit(event.ProvisioningEventCrashed, fmt.Sprintf("%s", err))
		}
	}
}

func updateResolvConf() error {
	// when starting the metal-hammer u-root sets a static resolv.conf file containing 8.8.8.8
	// this can only be overwritten by running dhclient
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
