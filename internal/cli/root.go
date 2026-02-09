package cli

import (
	"io"
	"os"

	"github.com/bnema/yank-that/internal/aliases"
	"github.com/bnema/yank-that/internal/app"
	"github.com/bnema/yank-that/internal/content"
	"github.com/bnema/yank-that/internal/githuburl"
	"github.com/spf13/cobra"
)

type PathRunner interface {
	Run(input string) (string, error)
}

type ContentRunner interface {
	Run(input string) (string, error)
}

type GitHubRunner interface {
	Run(input string) (string, error)
}

type AliasInstaller interface {
	Install(shell string) ([]string, error)
}

type Dependencies struct {
	PathRunner     PathRunner
	ContentRunner  ContentRunner
	GitHubRunner   GitHubRunner
	AliasInstaller AliasInstaller
	Out            io.Writer
}

func NewRootCmd(deps *Dependencies) *cobra.Command {
	resolvedDeps := mergeDependencies(deps)

	cmd := &cobra.Command{
		Use:          "yt",
		Short:        "Yank that",
		SilenceUsage: true,
	}

	cmd.AddCommand(
		newPathCmd(resolvedDeps.PathRunner, resolvedDeps.Out),
		newContentCmd(resolvedDeps.ContentRunner, resolvedDeps.Out),
		newGitHubCmd(resolvedDeps.GitHubRunner, resolvedDeps.Out),
		newInstallAliasesCmd(resolvedDeps.AliasInstaller, resolvedDeps.Out),
	)

	return cmd
}

func Execute() error {
	return NewRootCmd(nil).Execute()
}

func mergeDependencies(overrides *Dependencies) *Dependencies {
	deps := defaultDependencies()
	if overrides == nil {
		return deps
	}

	if overrides.PathRunner != nil {
		deps.PathRunner = overrides.PathRunner
	}
	if overrides.ContentRunner != nil {
		deps.ContentRunner = overrides.ContentRunner
	}
	if overrides.GitHubRunner != nil {
		deps.GitHubRunner = overrides.GitHubRunner
	}
	if overrides.AliasInstaller != nil {
		deps.AliasInstaller = overrides.AliasInstaller
	}
	if overrides.Out != nil {
		deps.Out = overrides.Out
	}

	return deps
}

func defaultDependencies() *Dependencies {
	return &Dependencies{
		PathRunner:     app.NewDefault(),
		ContentRunner:  content.NewDefault(),
		GitHubRunner:   githuburl.NewDefault(),
		AliasInstaller: aliases.NewInstaller(),
		Out:            os.Stdout,
	}
}
