package api

import (
	"encoding/json"
	"fmt"
	"toolkit/apikit/gitlab/internal/types"
)

type CreateMergeRequestOption struct {
	RemoveSourceBranch bool
}

type MergeRequest struct {
	ID     int    `json:"id"`
	WebURL string `json:"web_url"`
}

func CreateMergeRequest(
	conf types.ConfigContext, project, sourceBranch, targetBranch, title string, op CreateMergeRequestOption,
) (MergeRequest, error) {

	client, req := AuthRequest(conf)
	defer client.Close()

	body := map[string]any{
		"source_branch": sourceBranch,
		"target_branch": targetBranch,
		"title":         title,
	}
	if op.RemoveSourceBranch {
		body["remove_source_branch"] = true
	}

	res, err := req.
		SetBody(body).
		Post(conf.BaseURL + "/api/v4/projects/" + project + "/merge_requests")
	if err != nil {
		return MergeRequest{}, err
	}

	if res.StatusCode() == 201 {
		var mr MergeRequest
		if err := json.Unmarshal(res.Bytes(), &mr); err != nil {
			return MergeRequest{}, err
		}
		return mr, nil
	} else if res.StatusCode() == 401 {
		return MergeRequest{}, ErrNoAuthorization
	}

	return MergeRequest{}, fmt.Errorf("code %d: %s", res.StatusCode(), res.String())
}
