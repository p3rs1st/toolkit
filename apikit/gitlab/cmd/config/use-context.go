package config

import (
	"toolkit/apikit/gitlab/internal/types"

	"github.com/spf13/cobra"
)

func NewUseContextCommand(op *types.RootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "use-context context",
		Short:   "Use a specific context from the config file",
		GroupID: "context",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			op.SaveCurrentContext(cmd, args[0])
		},
	}

	return cmd
}
