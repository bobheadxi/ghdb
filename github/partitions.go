package github

import (
	"strconv"

	"github.com/src-d/go-mysql-server/sql"
)

type partitionIter struct {
	start  int
	offset int
}

func newPartitionIter(start int) (sql.PartitionIter, error) {
	// return nil, io.EOF

	return &partitionIter{
		start:  start,
		offset: 0,
	}, nil
}

func (p *partitionIter) Next() (sql.Partition, error) {
	return &partition{key: []byte(strconv.Itoa(p.start + p.offset))}, nil
}

func (p *partitionIter) Close() error { return nil }

type partition struct {
	key []byte
}

func (p *partition) Key() []byte { return p.key }
