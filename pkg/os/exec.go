package os

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// terminateGrace is the time a command gets to shut down on its own after its context
// is done, before it is killed. A tool which is in the middle of writing firmware must
// get the chance to finish that write instead of being killed right away.
const terminateGrace = 1 * time.Minute

// ExecuteCommand small helper to execute a command, redirect stdout/stderr.
func ExecuteCommand(name string, arg ...string) error {
	cmd, err := command(context.Background(), name, arg...)
	if err != nil {
		return err
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}

// ExecuteCommandWithOutput execute a command and return its stdout and stderr separately,
// use this instead of ExecuteCommand if the output must be inspected or logged.
// stdout and stderr are kept apart on purpose, a caller which parses the output of a tool
// must not have to deal with warnings the tool writes to stderr.
// Once the given context is done the command is asked to terminate and only killed if it
// does not exit within terminateGrace, see command.
func ExecuteCommandWithOutput(ctx context.Context, name string, arg ...string) (stdout, stderr string, err error) {
	cmd, err := command(ctx, name, arg...)
	if err != nil {
		return "", "", err
	}

	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err = cmd.Run()
	if err != nil && ctx.Err() != nil {
		err = fmt.Errorf("%s did not finish (%w): %w", name, ctx.Err(), err)
	}

	return out.String(), errOut.String(), err
}

func command(ctx context.Context, name string, arg ...string) (*exec.Cmd, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return nil, fmt.Errorf("unable to locate program:%s in path %w", name, err)
	}

	cmd := exec.CommandContext(ctx, path, arg...)

	cmd.Cancel = func() error {
		return cmd.Process.Signal(syscall.SIGTERM)
	}
	cmd.WaitDelay = terminateGrace

	return cmd, nil
}
