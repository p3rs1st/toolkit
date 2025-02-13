package api

import (
	"fmt"
	"toolkit/apikit/gitlab/internal/types"
)

var ErrNoAuthorization = fmt.Errorf("access token is invalid or expired")

func CheckAccessToken(conf types.Config) (bool, error) {
	client, req := AuthRequest(conf)
	defer client.Close()
	res, err := req.Get(conf.BaseURL + "/api/v4/personal_access_tokens")
	if err != nil {
		return false, err
	}

	if res.StatusCode() == 200 {
		return true, nil
	} else if res.StatusCode() == 401 {
		return false, nil
	}

	return false, fmt.Errorf("unexpected status code: %d", res.StatusCode())
}
