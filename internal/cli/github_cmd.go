package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func newGitHubCmd(runner GitHubRunner, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "github <file[:line]>",
		Short: "Copy GitHub URL for a repository file",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			url, err := runner.Run(args[0])
			if err != nil {
				return err
			}

			fmt.Fprintln(out, url)
			return nil
		},
	}

	cmd.Aliases = []string{"ytg"}
	return cmd
}
