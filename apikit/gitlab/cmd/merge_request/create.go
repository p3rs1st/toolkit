package mergerequest

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"toolkit/apikit/gitlab/internal/api"
	"toolkit/apikit/gitlab/internal/types"
	"toolkit/apikit/gitlab/pkg"

	"github.com/spf13/cobra"
)

func NewCreateCommand(op *types.RootOptions) *cobra.Command {
	removeSourceBranch := false

	cmd := &cobra.Command{
		Use:   "create projectName/projectID[,projectName1/projectID1] sourceBranch [targetBranch]",
		Short: "Create a new merge request",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.MinimumNArgs(2)(cmd, args); err != nil {
				return err
			}
			if err := cobra.MaximumNArgs(3)(cmd, args); err != nil {
				return err
			}
			return nil
		},
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectIDNames := strings.Split(args[0], ",")
			conf := op.GetConfig(cmd)

			projectIDs := []int{}
			projectNames := []string{}
			sourceBranches := []string{}
			targetBranches := []string{}

			for _, projectIDName := range projectIDNames {
				projectID, err := strconv.Atoi(projectIDName)
				project := api.Project{}
				if err != nil {
					projectName := projectIDName
					projects, err := api.ListProjects(
						conf, api.ListProjectsOption{SearchName: projectName},
					)
					if err != nil {
						return err
					}
					projects = pkg.FilterFunc(projects, func(p api.Project) bool { return p.Name == projectName })
					if len(projects) == 0 {
						return fmt.Errorf("project %q not found", projectName)
					}
					if len(projects) > 1 {
						return fmt.Errorf(
							"multiple projects found: %q",
							strings.Join(
								pkg.MapFunc(projects, func(p api.Project) string { return p.PathWithNamespace }),
								", ",
							),
						)
					}
					project = projects[0]
					projectID = projects[0].ID
				} else {
					project, err = api.GetProject(conf, projectID)
					if err != nil {
						return err
					}
				}

				sourceBranch := args[1]
				targetBranch := project.DefaultBranch
				if len(args) > 2 {
					targetBranch = args[2]
				}
				branches, err := api.ListProjectBranches(conf, projectID, api.ListProjectBranchesOption{})
				if err != nil {
					return err
				}
				if !slices.ContainsFunc(branches, func(b api.Branch) bool { return b.Name == sourceBranch }) {
					return fmt.Errorf("source branch %q not found", sourceBranch)
				}
				if !slices.ContainsFunc(branches, func(b api.Branch) bool { return b.Name == targetBranch }) {
					return fmt.Errorf("target branch %q not found", targetBranch)
				}
				projectIDs = append(projectIDs, projectID)
				projectNames = append(projectNames, project.Name)
				sourceBranches = append(sourceBranches, sourceBranch)
				targetBranches = append(targetBranches, targetBranch)
			}

			webURLs := []string{}
			errs := []error{}
			for i, projectID := range projectIDs {
				sourceBranch, targetBranch := sourceBranches[i], targetBranches[i]

				mr, err := api.CreateMergeRequest(
					conf,
					strconv.Itoa(projectID),
					sourceBranch,
					targetBranch,
					fmt.Sprintf("Merge %s into %s", sourceBranch, targetBranch),
					api.CreateMergeRequestOption{RemoveSourceBranch: removeSourceBranch},
				)
				if err != nil {
					errs = append(errs, fmt.Errorf("create merge request for %q failed: %w", projectNames[i], err))
					continue
				}
				webURLs = append(webURLs, mr.WebURL)
			}
			cmd.Printf("merge requests created:\n%s\n", strings.Join(webURLs, "\n"))

			if len(errs) > 0 {
				return errors.Join(errs...)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&removeSourceBranch, "rm", "r", false, "remove source branch after merge")

	return cmd
}
