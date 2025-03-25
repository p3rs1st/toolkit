package mergerequest

import (
	"toolkit/apikit/gitlab/cmd/util"
	"toolkit/apikit/gitlab/internal/types"

	"github.com/spf13/cobra"
)

func NewMergeRequestCommand(option *types.RootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mr",
		Short: "GitLab merge request commands",
		Run:   util.NoArguemntsCommandRun(),
	}

	cmd.AddCommand(NewCreateCommand(option))

	return cmd
}
