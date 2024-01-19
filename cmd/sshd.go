package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"syscall"

	"github.com/metal-stack/metal-hammer/pkg/os/command"
)

const sshdCommand = command.SSHD

// StartSSHD will start sshd to be able to diagnose problems on the pxe bootet machine.
func StartSSHD(log *slog.Logger, ip string) error {
	sshd, err := exec.LookPath(sshdCommand)
	if err != nil {
		return fmt.Errorf("unable to locate sshd %w", err)
	}
	cmd := exec.Command(sshd, "-port", "22")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Env = os.Environ()
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("unable to start sshd %w", err)
	}
	log.Info(fmt.Sprintf("sshd started, connect via ssh -i metal.key root@%s", ip))
	return nil
}
