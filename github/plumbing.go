package github

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

func newHTTPClient(auth TokenSource) *http.Client {
	return oauth2.NewClient(context.Background(), auth)
}
