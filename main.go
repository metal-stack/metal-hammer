package main

import (
	"io"
	"net/http"
	"os"

	"git.f-i-ts.de/maas/discover/cmd"
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
	err := envconfig.Process("discover", &spec)
	if err != nil {
		log.Error("configuration error", "error", err)
		os.Exit(1)
	}
	if len(os.Args) > 1 {
		envconfig.Usage("discover", &spec)
		os.Exit(0)
	}

	log.Info("discover", "version", getVersionString())
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

	err = cmd.RegisterDevice(&spec)
	if err != nil {
		log.Error("register device", "error", err)
	}

	err = cmd.Install("ubuntu")
	if err != nil {
		log.Error("install", "error", err)
	}

	rootHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "discover\n")
	}

	http.HandleFunc("/", rootHandler)
	log.Info("waiting for a image to burn")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Error("http server not stared", "error", err)
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
