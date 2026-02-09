package cli

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type recordingRunner struct {
	received []string
	result   string
	err      error
}

func (r *recordingRunner) Run(input string) (string, error) {
	r.received = append(r.received, input)
	return r.result, r.err
}

type recordingAliasInstaller struct {
	received []string
	result   []string
	err      error
}

func (r *recordingAliasInstaller) Install(shell string) ([]string, error) {
	r.received = append(r.received, shell)
	return r.result, r.err
}

func TestPathCommand_DefaultsToCurrentDirInput(t *testing.T) {
	runner := &recordingRunner{result: "/tmp"}
	out := bytes.NewBuffer(nil)

	cmd := newPathCmd(runner, out)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.NoError(t, err)
	require.Len(t, runner.received, 1)
	assert.Equal(t, "", runner.received[0])
	assert.Equal(t, "/tmp\n", out.String())
}

func TestContentCommand_RequiresOneArg(t *testing.T) {
	runner := &recordingRunner{result: "/tmp/file"}
	out := bytes.NewBuffer(nil)

	cmd := newContentCmd(runner, out)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "accepts 1 arg(s)")
}

func TestPathCommand_PropagatesRunnerError(t *testing.T) {
	runner := &recordingRunner{err: errors.New("path failed")}
	out := bytes.NewBuffer(nil)

	cmd := newPathCmd(runner, out)
	cmd.SetArgs([]string{"README.md"})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "path failed")
}

func TestContentCommand_PropagatesRunnerError(t *testing.T) {
	runner := &recordingRunner{err: errors.New("content failed")}
	out := bytes.NewBuffer(nil)

	cmd := newContentCmd(runner, out)
	cmd.SetArgs([]string{"README.md"})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "content failed")
}

func TestGitHubCommand_CopiesAndPrintsURL(t *testing.T) {
	runner := &recordingRunner{result: "https://github.com/bnema/yank-that/blob/main/README.md#L1"}
	out := bytes.NewBuffer(nil)

	cmd := newGitHubCmd(runner, out)
	cmd.SetArgs([]string{"README.md:1"})

	err := cmd.Execute()
	require.NoError(t, err)
	require.Len(t, runner.received, 1)
	assert.Equal(t, "README.md:1", runner.received[0])
	assert.Equal(t, "https://github.com/bnema/yank-that/blob/main/README.md#L1\n", out.String())
}

func TestGitHubCommand_RequiresOneArg(t *testing.T) {
	runner := &recordingRunner{result: "ok"}
	out := bytes.NewBuffer(nil)

	cmd := newGitHubCmd(runner, out)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "accepts 1 arg(s)")
}

func TestGitHubCommand_PropagatesRunnerError(t *testing.T) {
	runner := &recordingRunner{err: errors.New("github failed")}
	out := bytes.NewBuffer(nil)

	cmd := newGitHubCmd(runner, out)
	cmd.SetArgs([]string{"README.md"})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "github failed")
}

func TestInstallAliasesCommand_DefaultsToAll(t *testing.T) {
	installer := &recordingAliasInstaller{result: []string{"installed"}}
	out := bytes.NewBuffer(nil)

	cmd := newInstallAliasesCmd(installer, out)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.NoError(t, err)
	require.Len(t, installer.received, 1)
	assert.Equal(t, "all", installer.received[0])
	assert.Equal(t, "installed\n", out.String())
}

func TestInstallAliasesCommand_PropagatesErrors(t *testing.T) {
	installer := &recordingAliasInstaller{err: errors.New("boom")}
	out := bytes.NewBuffer(nil)

	cmd := newInstallAliasesCmd(installer, out)
	cmd.SetArgs([]string{"bash"})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
}
