package config

import (
	"toolkit/apikit/gitlab/cmd/util"
	"toolkit/apikit/gitlab/internal/types"

	"github.com/spf13/cobra"
)

func NewConfigCommand(op *types.RootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Edit local configuration file",
		Run:   util.NoArguemntsCommandRun(),
	}

	cmd.AddCommand(NewSetCommand(op))
	cmd.AddCommand(NewViewCommand(op))

	cmd.AddGroup(
		&cobra.Group{
			ID:    "config sub",
			Title: "Operate configuration",
		},
		&cobra.Group{
			ID:    "context",
			Title: "Context",
		},
	)
	cmd.AddCommand(NewUseContextCommand(op))
	cmd.AddCommand(NewGetContextsCommand(op))

	return cmd
}
