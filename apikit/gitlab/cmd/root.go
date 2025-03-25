package cmd

import (
	"toolkit/apikit/gitlab/cmd/branch"
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

	option := &types.RootOptions{}

	cmd.PersistentFlags().StringVarP(
		&option.ConfigFilepath,
		"config",
		"c",
		types.DefaultConfigPath,
		"Path to the config file",
	)
	cmd.PersistentFlags().StringVarP(
		&option.CurrentContext,
		"context",
		"t",
		"",
		"Tempoary override for current context",
	)

	cmd.AddCommand(config.NewConfigCommand(option))
	cmd.AddCommand(project.NewProjectCommand(option))
	cmd.AddCommand(mergerequest.NewMergeRequestCommand(option))
	cmd.AddCommand(branch.NewBranchCommand(option))

	return cmd
}
