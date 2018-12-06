package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	log "github.com/inconshreveable/log15"
)

// StartDHClient will start dhclient to enforce an ip on all interfaces.
func StartDHClient() error {
	sshd, err := exec.LookPath("dhclient")
	if err != nil {
		return fmt.Errorf("unable to locate dhclient info:%v", err)
	}
	cmd := exec.Command(sshd, "-ipv4")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Env = os.Environ()
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("unable to start dhclient info:%v", err)
	}
	log.Info("dhclient started")
	return nil
}
