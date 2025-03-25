package uurl

import (
	"errors"
	"fmt"
	"net/url"
)

var errInvalidURL = errors.New("invalid url")

func CheckURLValid(raw string) error {
	uri, err := url.ParseRequestURI(raw)
	if err != nil {
		return err
	}

	if (uri.Scheme != "http" && uri.Scheme != "https") || uri.Host == "" {
		return fmt.Errorf("%w: %s", errInvalidURL, raw)
	}

	return nil
}
