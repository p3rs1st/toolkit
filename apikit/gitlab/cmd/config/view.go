package config

import (
	"toolkit/apikit/gitlab/cmd/util"
	"toolkit/apikit/gitlab/internal/types"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewViewCommand(op *types.RootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "view",
		Short:        "View the current configuration",
		GroupID:      "config sub",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			util.RequireNoArguments(cmd, args)

			settings := op.AllSettings(cmd)
			settingsYAML, err := yaml.Marshal(settings)
			if err != nil {
				return err
			}
			cmd.Print(string(settingsYAML))

			return nil
		},
	}

	return cmd
}
