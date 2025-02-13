package api

import (
	"strconv"
	"toolkit/apikit/gitlab/internal/types"

	"resty.dev/v3"
)

var (
	perPage    int    = 100
	perPageStr string = strconv.Itoa(perPage)
)

func AuthRequest(conf types.Config) (*resty.Client, *resty.Request) {
	if conf.Token == "" {
		// mark token as invalid
		conf.Token = "*"
	}
	client := resty.New()
	return client, client.R().SetHeader("Private-Token", conf.Token)
}
