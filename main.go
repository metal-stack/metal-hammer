package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"git.f-i-ts.de/cloud-native/maas/metal-hammer/cmd"
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
	cmdline, err := ioutil.ReadFile("/proc/cmdline")
	if err != nil {
		log.Error("unable to read /proc/cmdline", "error", err)
	}

	cmdLineValues := strings.Split(string(cmdline), " ")
	envmap := make(map[string]string)
	for _, v := range cmdLineValues {
		keyValue := strings.Split(v, "=")
		if len(keyValue) == 2 {
			key := strings.TrimSpace(keyValue[0])
			value := strings.TrimSpace(keyValue[1])
			envmap[key] = value
		}
	}

	// METAL_CORE_URL must be in the form http://metal-core:4242
	if i, ok := envmap["METAL_CORE_URL"]; ok {
		spec.InstallURL = i + "/device/install"
		spec.RegisterURL = i + "/device/register"
		spec.ReportURL = i + "/device/report"
	}

	if i, ok := envmap["IMAGE_URL"]; ok {
		spec.ImageURL = i
	}

	fmt.Print(cmd.Hammer)
	log.Info("metal-hammer", "version", getVersionString())
	var level log.Lvl
	if spec.Debug {
		level = log.LvlDebug
	} else {
		level = log.LvlInfo
	}
	spec.Log()

	h := log.CallerFileHandler(log.StdoutHandler)
	h = log.LvlFilterHandler(level, h)
	log.Root().SetHandler(h)

	err = cmd.Run(&spec)
	if err != nil {
		log.Error("metal-hammer run", "error", err)
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
