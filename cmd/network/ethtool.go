package network

import (
	"bufio"
	log "github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"os/exec"
	"strings"
)

// Ethtool to query/set ethernet interfaces
type Ethtool struct {
	command string
}

// NewEthtool create a new Ethtool with the default command
func NewEthtool() *Ethtool {
	return &Ethtool{command: "ethtool"}
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
			log.Error("ethtool", "interface", ifi, "error disabling fw-lldp", err)
			return
		}
		log.Info("ethtool", "interface", ifi, "fw-lldp", "disabled")
	}
}
