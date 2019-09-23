package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/pkg/errors"

	"git.f-i-ts.de/cloud-native/metal/metal-hammer/pkg/os/command"
	log "github.com/inconshreveable/log15"
)

const sshdCommand = command.SSHD

// StartSSHD will start sshd to be able to diagnose problems on the pxe bootet machine.
func StartSSHD(ip string) error {
	sshd, err := exec.LookPath(sshdCommand)
	if err != nil {
		return errors.Wrap(err, "unable to locate sshd")
	}
	cmd := exec.Command(sshd, "-port", "22")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Env = os.Environ()
	err = cmd.Start()
	if err != nil {
		return errors.Wrap(err, "unable to start sshd")
	}
	log.Info(fmt.Sprintf("sshd started, connect via ssh -i metal.key root@%s", ip))
	return nil
}
