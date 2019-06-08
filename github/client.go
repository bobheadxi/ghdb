package github

import (
	"github.com/shurcooL/githubv4"
)

type client struct {
	*githubv4.Client
	auth TokenSource
}

func newClient(auth TokenSource) *client {
	return &client{
		Client: githubv4.NewClient(newHTTPClient(auth)),
		auth:   auth,
	}
}
