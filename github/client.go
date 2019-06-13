package github

import (
	"context"

	"github.com/shurcooL/githubv4"
)

type client struct {
	github *githubv4.Client

	pool *connPool
	auth TokenSource
}

func newClient(auth TokenSource, pool *connPool) *client {
	return &client{
		github: githubv4.NewClient(newHTTPClient(auth)),
		pool:   pool,
		auth:   auth,
	}
}

func (c *client) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	return c.pool.Exec(func() error { return c.github.Query(ctx, q, variables) })
}

func (c *client) DirectQuery(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	return c.github.Query(ctx, q, variables)
}

func (c *client) Close() error { return c.pool.Close() }
