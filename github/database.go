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

// Database collects the interfaces that the GitHub database exposes
type Database interface {
	sql.Database
	io.Closer

	Ping() error
}

// DatabaseOpts defines options for this database
type DatabaseOpts struct {
	PoolSize           int
	ConnectTimeout     time.Duration
	TransactionTimeout time.Duration
}

func (d *DatabaseOpts) setDefaults() {
	if d.PoolSize == 0 {
		d.PoolSize = runtime.NumCPU()
	}
	if d.ConnectTimeout == 0 {
		d.ConnectTimeout = 30 * time.Second
	}
	if d.TransactionTimeout == 0 {
		d.TransactionTimeout = 1 * time.Minute
	}
}

type database struct {
	l    log.Logger
	c    *client
	pool *connPool

	opts     DatabaseOpts
	ctx      context.Context
	cancelFn context.CancelFunc
}

// NewDatabase instantiates a new GitHub-backed sql.Database
// TODO caching?
func NewDatabase(l log.Logger, auth TokenSource, opts DatabaseOpts) (Database, error) {
	opts.setDefaults()
	ctx := context.Background()
	cancellable, cancel := context.WithCancel(ctx)
	d := &database{
		c:    newClient(auth),
		l:    l,
		pool: newPool(opts),

		opts:     opts,
		ctx:      cancellable,
		cancelFn: cancel,
	}
	return d, d.Ping()
}

func (d *database) Name() string { return d.c.auth.Name() }

func (d *database) Tables() map[string]sql.Table {
	tables := make(map[string]sql.Table)

	var q struct {
		Viewer struct {
			Repository struct {
				Nodes []struct {
					NameWithOwner string
				}
				PageInfo pageInfo
			} `graphql:"repository(affiliations: $affs, ownerAffiliations: $oAffs, first: 100, after: $cursor)"`
		}
	}
	if err := d.pool.Exec(func() error {
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
		if err := d.c.Query(d.ctx, &q, vars); err != nil {
			return err
		}

		for {
			if err := d.c.Query(d.ctx, &q, vars); err != nil {
				return err
			}
			for _, n := range q.Viewer.Repository.Nodes {
				tables[n.NameWithOwner] = nil // TODO
			}
			if !q.Viewer.Repository.PageInfo.HasNextPage {
				break
			}
			vars["cursor"] = githubv4.NewString(q.Viewer.Repository.PageInfo.EndCursor)
		}

		return nil
	}); err != nil {
		d.l.Error("failed to fetch tables: %v", err)
		return nil
	}

	return tables
}

func (d *database) Ping() error {
	// check the client
	// don't use pool here
	var q struct{ Viewer struct{ Name string } }
	timeout, cancel := context.WithTimeout(d.ctx, 30*time.Second)
	defer cancel()
	if err := d.c.Query(timeout, q, nil); err != nil {
		d.l.Error("ping unsuccessful", "error", err)
		return err
	}
	d.l.Info("ping successful", "viewer.name", q.Viewer.Name)
	return nil
}

func (d *database) Close() error {
	d.l.Debug("closing database")
	// TODO: maybe try to collect errors if there are any?
	d.cancelFn()
	return d.pool.Close()
}
