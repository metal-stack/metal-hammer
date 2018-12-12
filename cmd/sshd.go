package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	log "github.com/inconshreveable/log15"
)

// StartSSHD will start sshd to be able to diagnose problems on the pxe bootet machine.
func StartSSHD(ip string) error {
	sshd, err := exec.LookPath("sshd")
	if err != nil {
		return fmt.Errorf("unable to locate sshd info:%v", err)
	}
	cmd := exec.Command(sshd, "-port", "22")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Env = os.Environ()
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("unable to start sshd info:%v", err)
	}
	log.Info(fmt.Sprintf("sshd started, connect via ssh -i metal.key root@%s", ip))
	return nil
}
