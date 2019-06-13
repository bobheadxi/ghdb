package ghdb

import (
	sqle "github.com/src-d/go-mysql-server"
	"github.com/src-d/go-mysql-server/sql"
	"go.bobheadxi.dev/ghdb/github"
	"go.bobheadxi.dev/ghdb/log"
)

// DatabaseOpts configures a ghdb database instance.
type DatabaseOpts struct {
	Auth   github.TokenSource
	Engine github.EngineOpts
	Logger log.Logger
}

// Engine provides a MySQL-compatible interface to your GitHub issues
type Engine struct {
	e  *sqle.Engine
	db github.Engine
	//
}

// New instantiates a new ghdb engine
func New(opts DatabaseOpts) (*Engine, error) {
	driver := sqle.NewDefault()

	gh, err := github.NewEngine(opts.Logger, opts.Auth, opts.Engine)
	if err != nil {
		return nil, err
	}

	driver.AddDatabase(gh)

	return &Engine{
		e:  driver,
		db: gh,
	}, driver.Init()
}

// Query executes a MySQL query against the configured GitHub database.
func (e *Engine) Query(ctx *sql.Context, query string) (sql.Schema, sql.RowIter, error) {
	return e.e.Query(ctx, query)
}

// Close shuts down this ghdb engine, attempts to cancel ongoing transactions,
// and releases its resources
func (e *Engine) Close() error { return e.db.Close() }
