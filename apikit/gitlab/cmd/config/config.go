package config

import (
	"toolkit/apikit/gitlab/cmd/util"
	"toolkit/apikit/gitlab/internal/types"

	"github.com/spf13/cobra"
)

func NewConfigCommand(option *types.RootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Edit local configuration file",
		Run:   util.NoArguemntsCommandRun(),
	}

	cmd.AddCommand(NewSetCommand(option))
	cmd.AddCommand(NewViewCommand(option))

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
	cmd.AddCommand(NewUseContextCommand(option))
	cmd.AddCommand(NewGetContextsCommand(option))

	return cmd
}
