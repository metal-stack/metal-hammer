package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"git.f-i-ts.de/cloud-native/metal/metal-hammer/cmd"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/cmd/network"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg"
	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/uuid"
	"git.f-i-ts.de/cloud-native/metallib/version"
	log "github.com/inconshreveable/log15"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func main() {
	fmt.Print(cmd.HammerBanner)
	go handleSignals()
	// in case of an error/panic which let the whole app die, recover.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	ip := network.InternalIP()
	err := cmd.StartSSHD(ip)
	if err != nil {
		log.Error("sshd error", "error", err)
		os.Exit(1)
	}
	var spec cmd.Specification
	err = envconfig.Process("metal-hammer", &spec)
	if err != nil {
		log.Error("configuration error", "error", err)
		os.Exit(1)
	}
	if len(os.Args) > 1 {
		envconfig.Usage("metal-hammer", &spec)
		os.Exit(0)
	}

	// Grab metal-hammer configuration from kernel commandline
	envmap, err := pkg.ParseCmdline()
	if err != nil {
		log.Error("parse cmdline", "error", err)
		os.Exit(1)
	}

	log.Info("metal-hammer", "version", version.V)

	if d, ok := envmap["DEBUG"]; ok && (d == "1" || strings.ToLower(d) == "true") {
		spec.Debug = true
		os.Setenv("DEBUG", "1")
	}

	var level log.Lvl
	if spec.Debug {
		level = log.LvlDebug
	} else {
		level = log.LvlInfo
	}

	h := log.CallerFileHandler(log.StdoutHandler)
	h = log.LvlFilterHandler(level, h)
	log.Root().SetHandler(h)

	// METAL_CORE_URL must be in the form http://metal-core:4242
	if url, ok := envmap["METAL_CORE_ADDRESS"]; ok {
		spec.MetalCoreURL = url
	}

	if i, ok := envmap["IMAGE_URL"]; ok {
		spec.ImageURL = i
		spec.DevMode = true
	}

	if i, ok := envmap["IMAGE_ID"]; ok {
		spec.ImageID = i
		spec.DevMode = true
	}

	if c, ok := envmap["CIDR"]; ok {
		spec.Cidr = c
		spec.DevMode = true
	}

	if bgp, ok := envmap["BGP"]; ok {
		enabled, err := strconv.ParseBool(bgp)
		if err == nil {
			spec.BGPEnabled = enabled
		}
	}

	spec.MachineUUID = uuid.MachineUUID()
	spec.Ip = ip

	spec.Log()

	err = cmd.Run(&spec)
	if err != nil {
		wait := 5 * time.Second
		st := errors.WithStack(err)
		fmt.Printf("%+v", st)
		log.Error("metal-hammer failed", "rebooting in", wait, "error", err)
		time.Sleep(wait)
		pkg.Reboot()
	}
}

// handleSignals is just used to show what signals where sent to metal-hammer.
// this is useful to nail down production problems.
// FIXME think if we need to have this forever.
func handleSignals() {
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)

	// Passing no signals to Notify means that
	// all signals will be sent to the channel.
	signal.Notify(c)

	// Block until any signal is received.
	s := <-c
	log.Info("main", "got signal", s)
}
