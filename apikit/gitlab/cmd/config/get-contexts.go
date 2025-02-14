package config

import (
	"toolkit/apikit/gitlab/cmd/util"
	"toolkit/apikit/gitlab/internal/types"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func NewGetContextsCommand(op *types.RootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-contexts",
		Short: "Use a specific context from the config file",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			config := op.GetRawConfig(cmd)

			util.RenderTable(
				cmd,
				table.Row{"Current", "Name", "Base URL"},
				config.Contexts,
				func(c types.ConfigContext) []interface{} {
					current := ""
					if c.Name == config.CurrentContext {
						current = "*"
					}
					return []interface{}{current, c.Name, c.BaseURL}
				},
			)
		},
	}

	return cmd
}
