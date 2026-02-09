package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRootCmd_HasExpectedSubcommands(t *testing.T) {
	cmd := NewRootCmd(&Dependencies{
		PathRunner:     stubRunner{returnValue: "/tmp"},
		ContentRunner:  stubRunner{returnValue: "/tmp/file"},
		GitHubRunner:   stubRunner{returnValue: "https://github.com/bnema/yank-that/blob/main/README.md"},
		AliasInstaller: stubAliasInstaller{},
	})
	names := make([]string, 0)
	for _, c := range cmd.Commands() {
		names = append(names, c.Name())
	}

	assert.Contains(t, names, "path")
	assert.Contains(t, names, "content")
	assert.Contains(t, names, "github")
	assert.Contains(t, names, "install-aliases")
}

func TestRootCmd_UsageListsPathContentGithub(t *testing.T) {
	cmd := NewRootCmd(&Dependencies{
		PathRunner:     stubRunner{returnValue: "/tmp"},
		ContentRunner:  stubRunner{returnValue: "/tmp/file"},
		GitHubRunner:   stubRunner{returnValue: "https://github.com/bnema/yank-that/blob/main/README.md"},
		AliasInstaller: stubAliasInstaller{},
	})

	out := bytes.NewBuffer(nil)
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, out.String(), "path")
	assert.Contains(t, out.String(), "content")
	assert.Contains(t, out.String(), "github")
	assert.Contains(t, out.String(), "install-aliases")
}

type stubRunner struct {
	returnValue string
	err         error
}

func (s stubRunner) Run(input string) (string, error) {
	return s.returnValue, s.err
}

type stubAliasInstaller struct{}

func (s stubAliasInstaller) Install(shell string) ([]string, error) {
	return []string{"ok"}, nil
}
