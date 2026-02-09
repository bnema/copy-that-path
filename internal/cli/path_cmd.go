package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func newPathCmd(runner PathRunner, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "path [target]",
		Short: "Copy resolved absolute path",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			input := ""
			if len(args) == 1 {
				input = args[0]
			}

			resolved, err := runner.Run(input)
			if err != nil {
				return err
			}

			fmt.Fprintln(out, resolved)
			return nil
		},
	}

	cmd.Aliases = []string{"ytp"}
	return cmd
}
