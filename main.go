package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"syscall"
	"time"

	"github.com/metal-stack/v"

	"github.com/metal-stack/go-hal/connect"
	"github.com/metal-stack/go-hal/pkg/logger"
	"github.com/metal-stack/metal-hammer/cmd"
	"github.com/metal-stack/metal-hammer/cmd/event"
	"github.com/metal-stack/metal-hammer/cmd/network"
	"github.com/metal-stack/metal-hammer/pkg/kernel"
	"github.com/moby/sys/mountinfo"
	"github.com/pkg/term"
)

func main() {
	br := bufio.NewWriter(os.Stdout)

	t, err := term.Open("/dev/ttyS0")
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for range ticker.C {
			_ = br.Flush()
			err = t.Flush()

		}
	}()

	jsonHandler := slog.NewJSONHandler(br, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	log := slog.New(jsonHandler)

	fmt.Print(cmd.HammerBanner)
	if len(os.Args) > 1 {
		panic("cmd args are not supported")
	}

	mounted, err := mountinfo.Mounted("/etc")
	if err != nil {
		log.Error("unable to check if /etc is a mountpoint", "error", err)
		os.Exit(1)
	}
	if mounted {
		err := syscall.Unmount("/etc", syscall.MNT_FORCE)
		if err != nil {
			log.Error("unable to umount /etc, which is overmounted with tmpfs", "error", err)
			os.Exit(1)
		}
	}

	err = updateResolvConf()
	if err != nil {
		log.Error("error updating resolv.conf", "error", err)
		os.Exit(1)
	}

	// Reboot if metal-hammer crashes after 60sec.
	// go kernel.Watchdog(log)

	hal, err := connect.InBand(logger.NewSlog(log))
	if err != nil {
		log.Error("unable to detect hardware", "error", err)
		// os.Exit(1)
	}

	uuid, err := hal.UUID()
	if err != nil {
		log.Error("unable to get uuid hardware", "error", err)
		os.Exit(1)
	}
	log = log.With("machineID", uuid.String())

	ip := network.InternalIP()
	err = cmd.StartSSHD(log, ip)
	if err != nil {
		log.Error("sshd error", "error", err)
		os.Exit(1)
	}

	log.Info("starting", "version", v.V.String(), "hal", hal.Describe())

	spec := cmd.NewSpec(log)

	// Synchronize time using NTP
	network.NtpDate(log, spec.MetalConfig.NTPServers)

	spec.MachineUUID = uuid.String()
	spec.IP = ip

	spec.Log()

	withRemoteHandler, err := cmd.AddRemoteHandler(spec, jsonHandler)
	if err != nil {
		log.Error("unable to add remote logging", "error", err)
	} else {
		log = slog.New(withRemoteHandler).With("machineID", uuid.String())
		log.Info("remote logging enabled")
	}

	// FIXME set loglevel from spec.Debug

	emitter, err := cmd.Run(log, spec, hal)
	if err != nil {
		wait := 5 * time.Second
		log.Error("metal-hammer failed", "rebooting in", wait, "error", err)
		if emitter != nil {
			emitter.Emit(event.ProvisioningEventCrashed, fmt.Sprintf("%s", err))
		}
		time.Sleep(wait)
		err := kernel.Reboot()
		if err != nil {
			log.Error("metal-hammer reboot failed", "error", err)
			if emitter != nil {
				emitter.Emit(event.ProvisioningEventCrashed, fmt.Sprintf("%s", err))
			}
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

	if _, err := os.Stat(symlink); !os.IsNotExist(err) {
		err := os.Remove(symlink)
		if err != nil {
			return err
		}
	}

	return os.Symlink(target, symlink)
}
