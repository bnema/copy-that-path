package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func newInstallAliasesCmd(installer AliasInstaller, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "install-aliases [shell]",
		Short: "Install ytp/ytc/ytg aliases for bash, zsh, or fish",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			shell := "all"
			if len(args) == 1 {
				shell = args[0]
			}

			updates, err := installer.Install(shell)
			if err != nil {
				return err
			}

			for _, update := range updates {
				fmt.Fprintln(out, update)
			}

			return nil
		},
	}
}
