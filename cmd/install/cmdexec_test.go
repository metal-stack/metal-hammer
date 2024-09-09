package install

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// tests were inspired by this blog article: https://npf.io/2015/06/testing-exec-command/

type fakeexec struct {
	t         *testing.T
	mockCount int
	mocks     []fakeexecparams
}

// nolint:musttag
type fakeexecparams struct {
	WantCmd  []string `json:"want_cmd"`
	Output   string   `json:"output"`
	ExitCode int      `json:"exit_code"`
}

func fakeCmd(t *testing.T, params ...fakeexecparams) func(ctx context.Context, command string, args ...string) *exec.Cmd {
	f := fakeexec{
		t:     t,
		mocks: params,
	}
	return f.command
}

func (f *fakeexec) command(ctx context.Context, command string, args ...string) *exec.Cmd {
	if f.mockCount >= len(f.mocks) {
		require.Fail(f.t, "more commands called than mocks are available")
	}

	params := f.mocks[f.mockCount]
	f.mockCount++

	assert.Equal(f.t, params.WantCmd, append([]string{command}, args...))

	j, err := json.Marshal(params)
	require.NoError(f.t, err)

	cs := []string{"-test.run=TestHelperProcess", "--", string(j)}
	cmd := exec.CommandContext(ctx, os.Args[0], cs...) //nolint
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	var f fakeexecparams
	err := json.Unmarshal([]byte(os.Args[3]), &f)
	require.NoError(t, err)

	fmt.Fprint(os.Stdout, f.Output)

	os.Exit(f.ExitCode)
}
