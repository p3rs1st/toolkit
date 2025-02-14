package project

import (
	"toolkit/apikit/gitlab/cmd/util"
	"toolkit/apikit/gitlab/internal/api"
	"toolkit/apikit/gitlab/internal/types"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func NewListCommand(op *types.RootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "ls",
		Aliases:      []string{"list"},
		Short:        "List all projects",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			util.RequireNoArguments(cmd, args)

			projects, err := api.ListProjects(op.GetConfig(cmd), api.ListProjectsOption{})
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
