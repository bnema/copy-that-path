package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func newContentCmd(runner ContentRunner, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "content <file>",
		Short: "Copy file content",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			resolved, err := runner.Run(args[0])
			if err != nil {
				return err
			}

			fmt.Fprintln(out, resolved)
			return nil
		},
	}

	cmd.Aliases = []string{"ytc"}
	return cmd
}
