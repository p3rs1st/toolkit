package config

import (
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
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
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
