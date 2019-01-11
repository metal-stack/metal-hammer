package os

import (
	"github.com/pkg/errors"
	"os"
	"os/exec"
)

// ExecuteCommand small helper to execute a command, redirect stdout/stderr.
func ExecuteCommand(name string, arg ...string) error {
	path, err := exec.LookPath(name)
	if err != nil {
		return errors.Wrapf(err, "unable to locate program:%s in path", name)
	}
	cmd := exec.Command(path, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}
