package network

import (
	"bufio"
	"os/exec"
	"strings"
	"syscall"

	"os"
	"path"
	"path/filepath"

	"io/ioutil"

	log "github.com/inconshreveable/log15"
	"github.com/metal-stack/metal-hammer/pkg/os/command"

	"github.com/pkg/errors"
)

// EthtoolCommand to gather ethernet informations
const ethtoolCommand = command.Ethtool

// Ethtool to query/set ethernet interfaces
type Ethtool struct {
	command string
}

// NewEthtool create a new Ethtool with the default command
func NewEthtool() *Ethtool {
	err := syscall.Mount("debugfs", "/sys/kernel/debug", "debugfs", 0, "")
	if err != nil {
		log.Warn("ethtool", "mounting debugfs failed", err)
	}
	return &Ethtool{command: ethtoolCommand}
}

// Run execute ethtool
func (e *Ethtool) Run(args ...string) (string, error) {
	path, err := exec.LookPath(e.command)
	if err != nil {
		return "", errors.Wrapf(err, "unable to locate program:%s in path", e.command)
	}
	cmd := exec.Command(path, args...)
	output, err := cmd.Output()

	log.Debug("run", "command", e.command, "args", args, "output", string(output), "error", err)
	return string(output), err
}

// disableFirmwareLLDP Intel i40e based 10G+ network cards (e.g. XXV710)
// have network card based firmware lldp sending enabled.
// this prevents receiving lldp pdu`s from our switches, so turn it off.
// Another approach which is not reboot safe an will only turn off lldp is:
func (e *Ethtool) disableFirmwareLLDP(ifi string) {

	output, err := e.Run("--show-priv-flags", ifi)
	if err != nil {
		log.Info("ethtool", "interface", ifi, "msg", "no priv-flags or disable-fw-lldp not present")
		return
	}

	log.Debug("ethtool", "show-priv-flags", output)
	scanner := bufio.NewScanner(strings.NewReader(output))
	fwLLDP := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "disable-fw-lldp") && strings.Contains(line, "off") {
			fwLLDP = "off"
		}
		if strings.Contains(line, "disable-fw-lldp") && strings.Contains(line, "on") {
			fwLLDP = "on"
		}
	}

	if fwLLDP != "" {
		log.Info("ethtool", "interface", ifi, "disable-fw-lldp is set to", fwLLDP)
	}

	if fwLLDP == "off" {
		_, err := e.Run("--set-priv-flags", ifi, "disable-fw-lldp", "on")
		if err != nil {
			log.Error("ethtool", "interface", ifi, "error disabling fw-lldp try to stop it", err)
			e.stopFirmwareLLDP()
			return
		}
		log.Info("ethtool", "interface", ifi, "fw-lldp", "disabled")
	}
}

var buggyIntelNicDriverNames = []string{"i40e"}

// stopFirmwareLLDP stop Firmeware LLDP not persistent over reboots, only during runtime.
// mount -t debugfs none /sys/kernel/debug
// echo lldp stop > /sys/kernel/debug/i40e/0000:01:00.2/command
// where <0000:01:00.2> is the pci address of the ethernet nic, this can be inspected by lspci,
// or a loop over all directories in /sys/kernel/debug/i40e/*/command
func (e *Ethtool) stopFirmwareLLDP() {
	for _, driver := range buggyIntelNicDriverNames {
		debugFSPath := path.Join("/sys/kernel/debug", driver)
		err := filepath.Walk(debugFSPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Warn("ethtool", "stopFirmwareLLDP in path", path, "error", err)
				return err
			}
			if !info.IsDir() && info.Name() == "command" {
				log.Info("ethtool", "stopFirmwareLLDP found command", path)
				stopCommand := []byte("lldp stop")
				err := ioutil.WriteFile(path, stopCommand, os.ModePerm)
				if err != nil {
					log.Error("ethtool", "stopFirmwareLLDP stop lldp > command", path, "error", err)
				}
			}
			return nil
		})
		if err != nil {
			log.Error("ethtool", "stopFirmwareLLDP unable to walk through debugfs", debugFSPath, "error", err)
		}
	}
}
