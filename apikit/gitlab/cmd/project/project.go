package project

import (
	"toolkit/apikit/gitlab/cmd/util"
	"toolkit/apikit/gitlab/internal/types"

	"github.com/spf13/cobra"
)

func NewProjectCommand(option *types.RootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "GitLab projects",
		Run:   util.NoArguemntsCommandRun(),
	}

	cmd.AddCommand(NewListCommand(option))

	return cmd
}
