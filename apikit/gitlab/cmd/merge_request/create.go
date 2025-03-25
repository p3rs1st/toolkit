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
	"toolkit/apikit/gitlab/pkg"

	"github.com/spf13/cobra"
)

var errBranchNotFound = errors.New("branch not found")

type createOptions struct {
	*types.RootOptions
	removeSourceBranch bool
	title              string
}

func NewCreateCommand(option *types.RootOptions) *cobra.Command {
	createOption := createOptions{
		RootOptions: option,
	}

	cmd := &cobra.Command{
		Use:          "create projectName/projectID[,projectName1/projectID1,...] sourceBranch [targetBranch]",
		Short:        "Create a new merge request",
		SilenceUsage: true,
		Args:         cobra.RangeArgs(2, 3),
		ValidArgsFunction: func(
			cmd *cobra.Command, args []string, toComplete string,
		) ([]string, cobra.ShellCompDirective) {
			conf := option.GetConfig(cmd)
			if len(args) == 0 {
				projects, err := api.ListProjects(conf, api.ListProjectsOption{SearchName: toComplete})
				if err != nil {
					return nil, cobra.ShellCompDirectiveError
				}

				return pkg.MapFunc(projects, func(p api.Project) string { return p.Name }),
					cobra.ShellCompDirectiveNoFileComp
			}

			keys, err := util.ListUnionProjectBranches(conf, strings.Split(args[0], ","), toComplete)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			return keys, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: createCommandRunE(createOption),
	}

	cmd.Flags().BoolVarP(&createOption.removeSourceBranch, "rm", "r", false, "remove source branch after merge")
	cmd.Flags().StringVar(&createOption.title, "title", "", "title of the merge request")

	return cmd
}

func createCommandRunE(
	createOption createOptions,
) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		projectIDNames := strings.Split(args[0], ",")
		conf := createOption.RootOptions.GetConfig(cmd)

		targetBranch := ""
		if len(args) > 2 {
			targetBranch = args[2]
		}

		projectIDs,
			projectNames,
			sourceBranches,
			targetBranches,
			err := prepareCreateMergeRequestInfo(conf, projectIDNames, args[1], targetBranch)
		if err != nil {
			return err
		}

		webURLs := []string{}
		errs := []error{}

		for idx, projectID := range projectIDs {
			sourceBranch, targetBranch := sourceBranches[idx], targetBranches[idx]

			curTitle := createOption.title
			if curTitle == "" {
				curTitle = fmt.Sprintf("Merge %s into %s", sourceBranch, targetBranch)
			}

			mergeRequest, err := api.CreateMergeRequest(
				conf,
				strconv.Itoa(projectID),
				sourceBranch,
				targetBranch,
				curTitle,
				api.CreateMergeRequestOption{RemoveSourceBranch: createOption.removeSourceBranch},
			)
			if err != nil {
				errs = append(errs, fmt.Errorf("create merge request for %q failed: %w", projectNames[idx], err))

				continue
			}

			webURLs = append(webURLs, mergeRequest.WebURL)
		}

		if len(webURLs) > 0 {
			cmd.Printf("merge requests created:\n%s\n", strings.Join(webURLs, "\n"))
		}

		if len(errs) > 0 {
			return errors.Join(errs...)
		}

		return nil
	}
}

func prepareCreateMergeRequestInfo(
	conf types.ConfigContext,
	projectIDNames []string,
	sourceBranch,
	targetBranch string,
) ([]int, []string, []string, []string, error) {
	projectIDs := []int{}
	projectNames := []string{}
	sourceBranches := []string{}
	targetBranches := []string{}

	for _, projectIDName := range projectIDNames {
		project, err := util.GetProjectByIDName(conf, projectIDName)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		projectID := project.ID

		sourceBranch := sourceBranch
		curTargetBranch := project.DefaultBranch

		if targetBranch != "" {
			curTargetBranch = targetBranch
		}

		branches, err := api.ListProjectBranches(conf, projectID, api.ListProjectBranchesOption{})
		if err != nil {
			return nil, nil, nil, nil, err
		}

		if !slices.ContainsFunc(branches, func(b api.Branch) bool { return b.Name == sourceBranch }) {
			return nil, nil, nil, nil, fmt.Errorf("%w: %s", errBranchNotFound, sourceBranch)
		}

		if !slices.ContainsFunc(branches, func(b api.Branch) bool { return b.Name == curTargetBranch }) {
			return nil, nil, nil, nil, fmt.Errorf("%w: %s", errBranchNotFound, curTargetBranch)
		}

		projectIDs = append(projectIDs, projectID)
		projectNames = append(projectNames, project.Name)
		sourceBranches = append(sourceBranches, sourceBranch)
		targetBranches = append(targetBranches, curTargetBranch)
	}

	return projectIDs, projectNames, sourceBranches, targetBranches, nil
}
