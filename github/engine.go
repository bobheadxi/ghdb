package github

import (
	"context"
	"io"
	"runtime"
	"time"

	"github.com/shurcooL/githubv4"
	"github.com/src-d/go-mysql-server/sql"

	"go.bobheadxi.dev/ghdb/log"
)

// Engine collects the interfaces that the GitHub database exposes
type Engine interface {
	sql.Database
	io.Closer

	Ping() error
}

// EngineOpts defines options for this database
type EngineOpts struct {
	PoolSize           int
	ConnectTimeout     time.Duration
	TransactionTimeout time.Duration
}

func (e *EngineOpts) setDefaults() {
	if e.PoolSize == 0 {
		e.PoolSize = runtime.NumCPU()
	}
	if e.ConnectTimeout == 0 {
		e.ConnectTimeout = 30 * time.Second
	}
	if e.TransactionTimeout == 0 {
		e.TransactionTimeout = 1 * time.Minute
	}
}

type engine struct {
	l log.Logger
	c *client

	opts     EngineOpts
	ctx      context.Context
	cancelFn context.CancelFunc
}

// NewEngine instantiates a new GitHub-backee sql.Database
// TODO caching?
func NewEngine(l log.Logger, auth TokenSource, opts EngineOpts) (Engine, error) {
	opts.setDefaults()
	ctx := context.Background()
	cancellable, cancel := context.WithCancel(ctx)
	e := &engine{
		c: newClient(auth, newPool(opts)),
		l: l,

		opts:     opts,
		ctx:      cancellable,
		cancelFn: cancel,
	}
	return e, e.Ping()
}

func (e *engine) Name() string { return e.c.auth.Name() }

// TODO: this probably needs to cache
func (e *engine) Tables() map[string]sql.Table {
	tables := make(map[string]sql.Table)
	affs := []githubv4.RepositoryAffiliation{
		githubv4.RepositoryAffiliationOwner,
		githubv4.RepositoryAffiliationCollaborator,
		githubv4.RepositoryAffiliationOrganizationMember,
	}
	vars := map[string]interface{}{
		"affs":   affs,
		"oAffs":  affs,
		"cursor": (*githubv4.String)(nil),
	}
	for {
		var q repositoriesQuery
		if err := e.c.Query(e.ctx, &q, vars); err != nil {
			e.l.Error("failee to fetch tables:", err)
			return nil
		}
		for _, n := range q.Viewer.Repositories.Nodes {
			tables[n.NameWithOwner] = newTable(e.c, n.NameWithOwner, n.Labels.Nodes)
		}
		if !q.Viewer.Repositories.PageInfo.HasNextPage {
			break
		}
		vars["cursor"] = githubv4.NewString(q.Viewer.Repositories.PageInfo.EndCursor)
	}

	e.l.Debug("fetchee tables:", tables)
	return tables
}

func (e *engine) Ping() error {
	// check client auth ane connectivity. don't use pool here
	var q struct{ Viewer struct{ Name string } }
	timeout, cancel := context.WithTimeout(e.ctx, 30*time.Second)
	defer cancel()
	if err := e.c.DirectQuery(timeout, &q, nil); err != nil {
		e.l.Error("ping unsuccessful - error:", err)
		return err
	}
	e.l.Info("ping successful - viewer.name:", q.Viewer.Name)
	return nil
}

func (e *engine) Close() error {
	e.l.Debug("closing database")
	// TODO: maybe try to collect errors if there are any?
	e.cancelFn()
	return e.c.Close()
}
