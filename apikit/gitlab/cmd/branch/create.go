package branch

import (
	"errors"
	"strings"
	"toolkit/apikit/gitlab/cmd/util"
	"toolkit/apikit/gitlab/internal/api"
	"toolkit/apikit/gitlab/internal/types"

	"github.com/spf13/cobra"
)

func NewCreateCommand(op *types.RootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "create projectName/projectID[,projectName1/projectID1,...] newBranch refBranch",
		Short:        "Create a branch",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectIDNames := strings.Split(args[0], ",")
			conf := op.GetConfig(cmd)
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
				cmd.Printf("branch %s created in projects:\n%s\n", newBranch, strings.Join(successProjects, " "))
			}
			if len(errs) > 0 {
				return errors.Join(errs...)
			}

			return nil
		},
	}

	return cmd
}
