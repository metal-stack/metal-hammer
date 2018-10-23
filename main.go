package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"git.f-i-ts.de/cloud-native/maas/metal-hammer/cmd"
	"git.f-i-ts.de/cloud-native/maas/metal-hammer/pkg"
	log "github.com/inconshreveable/log15"
	"github.com/kelseyhightower/envconfig"
)

var (
	version   = "devel"
	revision  string
	gitsha1   string
	builddate string
)

func main() {
	var spec cmd.Specification
	err := envconfig.Process("metal-hammer", &spec)
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

	fmt.Print(cmd.Hammer)
	log.Info("metal-hammer", "version", getVersionString())
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
	if i, ok := envmap["METAL_CORE_URL"]; ok {
		spec.InstallURL = i + "/device/install"
		spec.RegisterURL = i + "/device/register"
		spec.ReportURL = i + "/device/report"
	}

	if i, ok := envmap["IMAGE_URL"]; ok {
		spec.ImageURL = i
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

func getVersionString() string {
	var versionString = version
	if gitsha1 != "" {
		versionString += " (" + gitsha1 + ")"
	}
	if revision != "" {
		versionString += ", " + revision
	}
	if builddate != "" {
		versionString += ", " + builddate
	}
	return versionString
}
