package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"git.f-i-ts.de/cloud-native/maas/metal-hammer/cmd"
	"git.f-i-ts.de/cloud-native/maas/metal-hammer/pkg"
	"git.f-i-ts.de/cloud-native/metallib/version"
	log "github.com/inconshreveable/log15"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	err := startSSHD()
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

	fmt.Print(cmd.Hammer)
	log.Info("metal-hammer", "version", version.V)
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
		spec.InstallURL = url
		spec.RegisterURL = url
		spec.ReportURL = url
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

	os.Setenv("DEBUG", "1")
	err = cmd.Run(&spec)
	if err != nil {
		wait := 5 * time.Second
		log.Error("metal-hammer failed", "rebooting in", wait, "error", err)
		time.Sleep(wait)
		pkg.Reboot()
	}
}

func startSSHD() error {
	cmd := exec.Command("/bbin/sshd", "-port", "22")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Env = os.Environ()
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("unable to start sshd info:%v", err)
	}
	log.Info(fmt.Sprintf("sshd started, connect via ssh -i metal.key root@%s", getInternalIP()))
	return nil
}

func getInternalIP() string {
	itf, _ := net.InterfaceByName("eth0")
	item, _ := itf.Addrs()
	var ip net.IP
	for _, addr := range item {
		switch v := addr.(type) {
		case *net.IPNet:
			if !v.IP.IsLoopback() {
				if v.IP.To4() != nil { //Verify if IP is IPV4
					ip = v.IP
				}
			}
		}
	}
	if ip != nil {
		return ip.String()
	}
	return ""
}
