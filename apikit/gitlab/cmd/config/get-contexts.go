package config

import (
	"toolkit/apikit/gitlab/cmd/util"
	"toolkit/apikit/gitlab/internal/types"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func NewGetContextsCommand(option *types.RootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get-contexts",
		Short:   "Use a specific context from the config file",
		GroupID: "context",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			config := option.GetRawConfig(cmd)

			util.RenderTable(
				cmd,
				table.Row{"Current", "Name", "Base URL"},
				config.Contexts,
				func(conf types.ConfigContext) []interface{} {
					current := ""
					if conf.Name == config.CurrentContext {
						current = "*"
					}

					return []interface{}{current, conf.Name, conf.BaseURL}
				},
			)
		},
	}

	return cmd
}
