package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"git.f-i-ts.de/cloud-native/maas/metal-hammer/cmd"
	"git.f-i-ts.de/cloud-native/maas/metal-hammer/pkg"
	"git.f-i-ts.de/cloud-native/metallib/version"
	log "github.com/inconshreveable/log15"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	err := cmd.StartSSHD()
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

	fmt.Print(cmd.HammerBanner)
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

	if port, ok := envmap["IPMI_PORT"]; ok {
		spec.IPMIPort = port
		spec.DevMode = true
	}

	if i, ok := envmap["IMAGE_URL"]; ok {
		spec.ImageURL = i
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

	spec.Log()

	err = cmd.Run(&spec)
	if err != nil {
		wait := 5 * time.Second
		log.Error("metal-hammer failed", "rebooting in", wait, "error", err)
		time.Sleep(wait)
		pkg.Reboot()
	}
}
