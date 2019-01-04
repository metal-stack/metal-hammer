package os

import (
	"fmt"
	"os"
	"os/exec"
)

// ExecuteCommand small helper to execute a command, redirect stdout/stderr.
func ExecuteCommand(name string, arg ...string) error {
	path, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("unable to locate program:%s in path info:%v", name, err)
	}
	cmd := exec.Command(path, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}
