package util

import (
	"fmt"
	"strconv"
	"strings"
	"toolkit/apikit/gitlab/internal/api"
	"toolkit/apikit/gitlab/internal/types"
	"toolkit/apikit/gitlab/pkg"
)

func GetProjectByIDName(conf types.Config, projectIDName string) (api.Project, error) {
	projectID, err := strconv.Atoi(projectIDName)
	var project api.Project
	if err != nil {
		projectName := projectIDName
		projects, err := api.ListProjects(
			conf, api.ListProjectsOption{SearchName: projectName},
		)
		if err != nil {
			return api.Project{}, err
		}
		projects = pkg.FilterFunc(projects, func(p api.Project) bool { return p.Name == projectName })
		if len(projects) == 0 {
			return api.Project{}, fmt.Errorf("project %q not found", projectName)
		}
		if len(projects) > 1 {
			return api.Project{}, fmt.Errorf(
				"multiple projects found: %q",
				strings.Join(
					pkg.MapFunc(projects, func(p api.Project) string { return p.PathWithNamespace }),
					", ",
				),
			)
		}
		project = projects[0]
	} else {
		project, err = api.GetProject(conf, projectID)
		if err != nil {
			return api.Project{}, err
		}
	}

	return project, nil
}
