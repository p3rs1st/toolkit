package util

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"toolkit/apikit/gitlab/internal/api"
	"toolkit/apikit/gitlab/internal/types"
	"toolkit/apikit/gitlab/pkg"
)

var (
	errProjectNotFound      = errors.New("project not found")
	errProjectMultipleFound = errors.New("multiple projects found")
)

func GetProjectByIDName(conf types.ConfigContext, projectIDName string) (api.Project, error) {
	projectID, err := strconv.Atoi(projectIDName)

	if err == nil {
		return GetProjectByName(conf, projectIDName)
	}

	return api.GetProject(conf, projectID)
}

func GetProjectByName(conf types.ConfigContext, projectName string) (api.Project, error) {
	projects, err := api.ListProjects(
		conf, api.ListProjectsOption{SearchName: projectName},
	)
	if err != nil {
		return api.Project{}, err
	}

	projects = pkg.FilterFunc(projects, func(p api.Project) bool { return p.Name == projectName })
	if len(projects) != 1 {
		return api.Project{}, fmt.Errorf("%w: %q", errProjectNotFound, projectName)
	}

	if len(projects) == 0 {
		return api.Project{}, fmt.Errorf("%w: %s", errProjectNotFound, projectName)
	}

	if len(projects) > 1 {
		return api.Project{}, fmt.Errorf(
			"%w: %q",
			errProjectMultipleFound,
			strings.Join(
				pkg.MapFunc(projects, func(p api.Project) string { return p.PathWithNamespace }),
				", ",
			),
		)
	}

	return projects[0], nil
}

func ListUnionProjectBranches(
	conf types.ConfigContext, projectIDNames []string, searchBranch string,
) ([]string, error) {
	var branchSet map[string]struct{}

	for _, projectIDName := range projectIDNames {
		project, err := GetProjectByIDName(conf, projectIDName)
		if err != nil {
			return nil, err
		}

		branches, err := api.ListProjectBranches(
			conf, project.ID, api.ListProjectBranchesOption{SearchName: searchBranch},
		)
		if err != nil {
			return nil, err
		}

		set := make(map[string]struct{})
		for _, branch := range branches {
			set[branch.Name] = struct{}{}
		}

		if branchSet == nil {
			branchSet = set
		} else {
			for branch := range branchSet {
				if _, ok := set[branch]; !ok {
					delete(branchSet, branch)
				}
			}
		}
	}

	keys := make([]string, 0, len(branchSet))
	for k := range branchSet {
		keys = append(keys, k)
	}

	return keys, nil
}
