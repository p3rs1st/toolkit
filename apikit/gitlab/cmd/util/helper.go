package util

import (
	"os"

	"github.com/spf13/cobra"
)

func UnknownCommand(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		cmd.PrintErrf(
			"Error: unknown command: %q\nSee '%s -h' for help\n",
			args[0],
			cmd.CommandPath(),
		)
		os.Exit(1)
	}
}

func NoArguemntsCommandRun() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		UnknownCommand(cmd, args)
	}
}
