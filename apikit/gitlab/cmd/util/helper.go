package util

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func RequireNoArguments(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		cmd.PrintErrf(
			"Error: unexpected arguments: %q\nSee '%s -h' for help\n",
			strings.Join(args, " "),
			cmd.CommandPath(),
		)
		os.Exit(1)
	}
}

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
