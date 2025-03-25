package api

import (
	"errors"
	"strconv"

	"toolkit/apikit/gitlab/internal/types"

	"resty.dev/v3"
)

const (
	perPage int = 100
)

var perPageStr = strconv.Itoa(perPage)

var errHTTPCode = errors.New("error http code")

func AuthRequest(conf types.ConfigContext) (*resty.Client, *resty.Request) {
	if conf.Token == "" {
		// mark token as invalid
		conf.Token = "*"
	}

	client := resty.New()

	return client, client.R().SetHeader("Private-Token", conf.Token)
}
