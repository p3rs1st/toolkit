package project

import (
	"toolkit/apikit/gitlab/cmd/util"
	"toolkit/apikit/gitlab/internal/api"
	"toolkit/apikit/gitlab/internal/types"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func NewListCommand(option *types.RootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "ls",
		Aliases:      []string{"list"},
		Short:        "List all projects",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			projects, err := api.ListProjects(option.GetConfig(cmd), api.ListProjectsOption{})
			if err != nil {
				return err
			}

			util.RenderTable(
				cmd,
				table.Row{"Name", "Path With Namespace", "ID"},
				projects,
				func(p api.Project) []interface{} {
					return []interface{}{p.Name, p.PathWithNamespace, p.ID}
				},
			)

			return nil
		},
	}

	return cmd
}
