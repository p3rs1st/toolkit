package branch

import (
	"errors"
	"strings"

	"toolkit/apikit/gitlab/cmd/util"
	"toolkit/apikit/gitlab/internal/api"
	"toolkit/apikit/gitlab/internal/types"
	"toolkit/apikit/gitlab/pkg"

	"github.com/spf13/cobra"
)

func NewCreateCommand(option *types.RootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create projectName/projectID[,projectName1/projectID1,...] newBranch refBranch",
		Short:             "Create a branch",
		SilenceUsage:      true,
		Args:              cobra.ExactArgs(3),
		ValidArgsFunction: createCommandValidArgsFunction(option),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectIDNames := strings.Split(args[0], ",")
			conf := option.GetConfig(cmd)
			newBranch := args[1]
			refBranch := args[2]
			projects := []api.Project{}

			for _, projectIDName := range projectIDNames {
				project, err := util.GetProjectByIDName(conf, projectIDName)
				if err != nil {
					return err
				}
				projects = append(projects, project)
			}

			successProjects := []string{}
			errs := []error{}
			for i, project := range projects {
				err := api.CreateProjectBranch(conf, project.ID, newBranch, refBranch)
				if err != nil {
					errs = append(errs, err)
				} else {
					successProjects = append(successProjects, projectIDNames[i])
				}
			}

			if len(successProjects) > 0 {
				cmd.Printf(
					"branch %s created in projects:\n%s\n",
					newBranch,
					strings.Join(successProjects, " "),
				)
			}
			if len(errs) > 0 {
				return errors.Join(errs...)
			}

			return nil
		},
	}

	return cmd
}

func createCommandValidArgsFunction(
	option *types.RootOptions,
) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(
		cmd *cobra.Command, args []string, toComplete string,
	) ([]string, cobra.ShellCompDirective) {
		conf := option.GetConfig(cmd)

		if len(args) == 0 {
			projects, err := api.ListProjects(
				conf,
				api.ListProjectsOption{SearchName: toComplete},
			)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			return pkg.MapFunc(
				projects,
				func(p api.Project) string { return p.Name },
			), cobra.ShellCompDirectiveNoFileComp
		} else if len(args) == 1 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		keys, err := util.ListUnionProjectBranches(
			conf,
			strings.Split(args[0], ","),
			toComplete,
		)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return keys, cobra.ShellCompDirectiveNoFileComp
	}
}
