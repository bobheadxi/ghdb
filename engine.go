package ghdb

import (
	sqle "github.com/src-d/go-mysql-server"
	"github.com/src-d/go-mysql-server/sql"
	"go.bobheadxi.dev/ghdb/github"
	"go.bobheadxi.dev/ghdb/log"
)

// EngineOpts configures a ghdb engine instance.
type EngineOpts struct {
	Auth     github.TokenSource
	Database github.DatabaseOpts
	Logger   log.Logger
}

// Engine provides a MySQL-compatible interface to your GitHub issues
type Engine struct {
	e  *sqle.Engine
	db github.Database
	//
}

// New instantiates a new ghdb engine
func New(opts EngineOpts) (*Engine, error) {
	driver := sqle.NewDefault()

	gh, err := github.NewDatabase(opts.Logger, opts.Auth, opts.Database)
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
