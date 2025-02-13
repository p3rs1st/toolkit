package cmd

import (
	"toolkit/apikit/gitlab/cmd/config"
	mergerequest "toolkit/apikit/gitlab/cmd/merge_request"
	"toolkit/apikit/gitlab/cmd/project"
	"toolkit/apikit/gitlab/cmd/util"
	"toolkit/apikit/gitlab/internal/types"

	"github.com/spf13/cobra"
)

func NewGitlabCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gitlab",
		Short: "Gitlab API commands",
		Run:   util.NoArguemntsCommandRun(),
	}

	var op types.RootOptions

	cmd.PersistentFlags().StringVarP(
		&op.ConfigFilepath,
		"config",
		"c",
		types.DefaultConfigPath,
		"Path to the config file",
	)

	cmd.AddCommand(config.NewConfigCommand(&op))
	cmd.AddCommand(project.NewProjectCommand(&op))
	cmd.AddCommand(mergerequest.NewMergeRequestCommand(&op))

	return cmd
}
