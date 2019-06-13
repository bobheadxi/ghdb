package github

import (
	"github.com/src-d/go-mysql-server/sql"
)

// queryable collects the interfaces that the internal table implements
type queryable interface {
	sql.Table
	// sql.Inserter
}

type table struct {
	repo string
	cols []labelInfo
	c    *client
}

func newTable(c *client, repo string, cols []labelInfo) queryable {
	return &table{
		repo: repo,
		cols: cols,
		c:    c,
	}
}

func (t *table) Name() string { return t.repo }

func (t *table) String() string { return t.Name() }

func (t *table) Schema() sql.Schema {
	// set up default columns
	cols := sql.Schema{
		{
			Name:     "id",
			Type:     sql.Int32,
			Nullable: false,
			Source:   t.repo,
		},
		{
			Name:     "title",
			Type:     sql.Text,
			Nullable: false,
			Source:   t.repo,
		},
		{
			Name:     "body",
			Type:     sql.Text,
			Nullable: false,
			Source:   t.repo,
		},
		{
			Name:     "open",
			Type:     sql.Boolean,
			Nullable: false,
			Source:   t.repo,
		},
	}

	// add labels as columns
	for _, c := range t.cols {
		name, tp := toType(c.Name)
		cols = append(cols, &sql.Column{
			Name:     name,
			Type:     tp,
			Nullable: true,
			Source:   t.repo,
		})
	}

	return sql.Schema(cols)
}

func (t *table) Partitions(*sql.Context) (sql.PartitionIter, error) {
	return nil, nil
}

func (t *table) PartitionRows(*sql.Context, sql.Partition) (sql.RowIter, error) {
	return &tableIter{}, nil
}

type tableIter struct {
}

func (i *tableIter) Next() (sql.Row, error) {
	return []interface{}{}, nil
}

func (i *tableIter) Close() error {
	return nil
}
