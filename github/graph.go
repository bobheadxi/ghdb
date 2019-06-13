package github

import "github.com/shurcooL/githubv4"

// pageInfo declares metadata for paging
type pageInfo struct {
	EndCursor   githubv4.String
	HasNextPage bool
}

type repoTableInfo struct {
	NameWithOwner string
	Labels        struct {
		Nodes []labelInfo
	} `graphql:"labels(first:100)"` // no more than 100 columns plz
}

type labelInfo struct {
	// ID   string
	Name string
}
