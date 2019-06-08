package github

import "github.com/shurcooL/githubv4"

// pageInfo declares metadata for paging
type pageInfo struct {
	EndCursor   githubv4.String
	HasNextPage bool
}
