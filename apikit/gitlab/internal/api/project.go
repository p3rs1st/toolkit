package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"toolkit/apikit/gitlab/internal/types"
)

type ListProjectsOption struct {
	SearchName string
}

type Branch struct {
	Name string `json:"name"`
}

type ListProjectBranchesOption struct {
	SearchName string
}

type Project struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	PathWithNamespace string `json:"path_with_namespace"`
	DefaultBranch     string `json:"default_branch"`
}

func CreateProjectBranch(conf types.Config, projectID int, branch, ref string) error {
	client, req := AuthRequest(conf)
	defer client.Close()

	res, err := req.
		SetBody(map[string]string{
			"branch": branch,
			"ref":    ref,
		}).
		Post(conf.BaseURL + "/api/v4/projects/" + strconv.Itoa(projectID) + "/repository/branches")
	if err != nil {
		return err
	}

	if res.StatusCode() == 201 {
		return nil
	} else if res.StatusCode() == 401 {
		return ErrNoAuthorization
	}

	return fmt.Errorf("code %d: %s", res.StatusCode(), res.String())
}

func GetProject(conf types.Config, projectID int) (Project, error) {
	client, req := AuthRequest(conf)
	defer client.Close()

	res, err := req.
		Get(conf.BaseURL + "/api/v4/project/" + strconv.Itoa(projectID))
	if err != nil {
		return Project{}, err
	}

	if res.StatusCode() == 200 {
		project := Project{}
		if err := json.Unmarshal(res.Bytes(), &project); err != nil {
			return Project{}, err
		}
		return project, nil
	} else if res.StatusCode() == 401 {
		return Project{}, ErrNoAuthorization
	}

	return Project{}, fmt.Errorf("code %d: %s", res.StatusCode(), res.String())
}

func ListProjects(conf types.Config, op ListProjectsOption) ([]Project, error) {
	projects := []Project{}
	lastProjectID := 0
	for {
		params := map[string]string{
			"membership": "true",
			"order_by":   "id",
			"sort":       "asc",
			"id_after":   strconv.Itoa(lastProjectID),
			"per_page":   perPageStr,
		}
		if op.SearchName != "" {
			params["search"] = op.SearchName
		}
		curProjects, err := listProjects(conf, params)
		if err != nil {
			return nil, err
		}
		if len(curProjects) == 0 {
			return projects, nil
		}
		lastProjectID = curProjects[len(curProjects)-1].ID
		projects = append(projects, curProjects...)
		if len(curProjects) < perPage {
			return projects, nil
		}
	}
}

func ListProjectBranches(conf types.Config, projectID int, op ListProjectBranchesOption) ([]Branch, error) {
	branches := []Branch{}
	page := 1
	for {
		params := map[string]string{
			"page":     strconv.Itoa(page),
			"per_page": perPageStr,
		}
		if op.SearchName != "" {
			params["search"] = op.SearchName
		}
		curBranches, err := listProjectBranches(conf, projectID, params)
		if err != nil {
			return nil, err
		}
		if len(curBranches) == 0 {
			return branches, nil
		}
		branches = append(branches, curBranches...)
		page += 1
		if len(curBranches) < perPage {
			return branches, nil
		}
	}
}

func listProjects(conf types.Config, params map[string]string) ([]Project, error) {
	client, req := AuthRequest(conf)
	defer client.Close()

	res, err := req.
		SetQueryParams(params).
		Get(conf.BaseURL + "/api/v4/projects")
	if err != nil {
		return nil, err
	}

	if res.StatusCode() == 200 {
		projects := []Project{}
		if err := json.Unmarshal(res.Bytes(), &projects); err != nil {
			return nil, err
		}
		return projects, nil
	} else if res.StatusCode() == 401 {
		return nil, ErrNoAuthorization
	}

	return nil, fmt.Errorf("code %d: %s", res.StatusCode(), res.String())
}

func listProjectBranches(conf types.Config, projectID int, params map[string]string) ([]Branch, error) {
	client, req := AuthRequest(conf)
	defer client.Close()

	res, err := req.
		SetQueryParams(params).
		Get(conf.BaseURL + "/api/v4/projects/" + strconv.Itoa(projectID) + "/repository/branches")
	if err != nil {
		return nil, err
	}

	if res.StatusCode() == 200 {
		branches := []Branch{}
		if err := json.Unmarshal(res.Bytes(), &branches); err != nil {
			return nil, err
		}
		return branches, nil
	} else if res.StatusCode() == 401 {
		return nil, ErrNoAuthorization
	}

	return nil, fmt.Errorf("code %d: %s", res.StatusCode(), res.String())
}
