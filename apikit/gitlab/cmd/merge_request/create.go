package mergerequest

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"toolkit/apikit/gitlab/cmd/util"
	"toolkit/apikit/gitlab/internal/api"
	"toolkit/apikit/gitlab/internal/types"

	"github.com/spf13/cobra"
)

func NewCreateCommand(op *types.RootOptions) *cobra.Command {
	removeSourceBranch := false

	cmd := &cobra.Command{
		Use:          "create projectName/projectID[,projectName1/projectID1,...] sourceBranch [targetBranch]",
		Short:        "Create a new merge request",
		Args:         cobra.RangeArgs(2, 3),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectIDNames := strings.Split(args[0], ",")
			conf := op.GetConfig(cmd)

			projectIDs := []int{}
			projectNames := []string{}
			sourceBranches := []string{}
			targetBranches := []string{}

			for _, projectIDName := range projectIDNames {
				project, err := util.GetProjectByIDName(conf, projectIDName)
				if err != nil {
					return err
				}
				projectID := project.ID

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
			if len(webURLs) > 0 {
				cmd.Printf("merge requests created:\n%s\n", strings.Join(webURLs, "\n"))
			}
			if len(errs) > 0 {
				return errors.Join(errs...)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&removeSourceBranch, "rm", "r", false, "remove source branch after merge")

	return cmd
}
