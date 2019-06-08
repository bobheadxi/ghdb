package github

import (
	"github.com/shurcooL/githubv4"
	"github.com/src-d/go-mysql-server/sql"
	"golang.org/x/oauth2"

	"go.bobheadxi.dev/ghdb/log"
)

type Database struct {
	l log.Logger
}

func NewDatabase(l log.Logger, auth oauth2.TokenSource) (sql.Database, error) {
	githubv4.NewClient(nil)
	return nil, nil
}
