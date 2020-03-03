package os

import (
	"github.com/pkg/errors"
	"os"
	"os/exec"
)

// ExecuteCommand small helper to execute a command, redirect stdout/stderr.
func ExecuteCommand(name string, args ...string) error {
	cmd, err := CreateCommand(name, args...)
	if err != nil {
		return err
	}
	return cmd.Run()
}

// ExecuteCommandCombinedOutput small helper to execute a command and return the output, redirect stdout/stderr.
func ExecuteCommandCombinedOutput(name string, args ...string) (string, error) {
	cmd, err := CreateCommand(name, args...)
	if err != nil {
		return "", err
	}
	bb, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	if len(bb) == 0 {
		return "", nil
	}
	return string(bb[:len(bb)-1]), nil
}

// CreateCommand small helper to create a command, redirect stdout/stderr.
func CreateCommand(name string, args ...string) (*exec.Cmd, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to locate program:%s in path", name)
	}
	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd, nil
}
