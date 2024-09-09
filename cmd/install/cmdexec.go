package install

import (
	"context"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"
)

type cmdexec struct {
	log *slog.Logger
	c   func(ctx context.Context, name string, arg ...string) *exec.Cmd
}

type cmdParams struct {
	name     string
	args     []string
	dir      string
	timeout  time.Duration
	combined bool
	stdin    string
	env      []string
}

func (i *cmdexec) command(p *cmdParams) (out string, err error) {
	var (
		start  = time.Now()
		output []byte
	)
	i.log.Info("running command", "command", strings.Join(append([]string{p.name}, p.args...), " "), "start", start.String())

	ctx := context.Background()
	if p.timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, p.timeout)
		defer cancel()
	}

	cmd := i.c(ctx, p.name, p.args...)
	if p.dir != "" {
		cmd.Dir = "/etc/metal"
	}

	cmd.Env = append(cmd.Env, p.env...)

	// show stderr
	cmd.Stderr = os.Stderr

	if p.stdin != "" {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return "", err
		}

		go func() {
			defer stdin.Close()
			_, err = io.WriteString(stdin, p.stdin)
			if err != nil {
				i.log.Error("error when writing to command's stdin", "error", err)
			}
		}()
	}

	if p.combined {
		output, err = cmd.CombinedOutput()
	} else {
		output, err = cmd.Output()
	}

	out = string(output)
	took := time.Since(start)

	if err != nil {
		i.log.Error("executed command with error", "output", out, "duration", took.String(), "error", err)
		return "", err
	}

	i.log.Info("executed command", "output", out, "duration", took.String())

	return
}
