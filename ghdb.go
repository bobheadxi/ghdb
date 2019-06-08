package ghdb

import (
	sqle "github.com/src-d/go-mysql-server"
	"github.com/src-d/go-mysql-server/sql"
	"go.bobheadxi.dev/ghdb/github"
	"go.bobheadxi.dev/ghdb/log"
	"golang.org/x/oauth2"
)

type EngineOpts struct {
	Auth   oauth2.TokenSource
	Logger log.Logger
}

type Engine struct {
	d *sqle.Engine
	//
}

func New(opts EngineOpts) (*Engine, error) {
	driver := sqle.NewDefault()

	gh, err := github.NewDatabase(opts.Logger, opts.Auth)
	if err != nil {
		return nil, err
	}

	driver.AddDatabase(gh)

	return &Engine{driver}, driver.Init()
}

func (e *Engine) Query(ctx *sql.Context, query string) (sql.Schema, sql.RowIter, error) {
	return e.d.Query(ctx, query)
}
