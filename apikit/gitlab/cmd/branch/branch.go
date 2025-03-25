package branch

import (
	"toolkit/apikit/gitlab/cmd/util"
	"toolkit/apikit/gitlab/internal/types"

	"github.com/spf13/cobra"
)

func NewBranchCommand(option *types.RootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "branch",
		Short: "GitLab branches",
		Run:   util.NoArguemntsCommandRun(),
	}

	cmd.AddCommand(NewCreateCommand(option))

	return cmd
}
