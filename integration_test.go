package ghdb

import (
	"context"
	"testing"

	"github.com/src-d/go-mysql-server/sql"
	"github.com/stretchr/testify/require"

	"go.uber.org/zap/zaptest"

	"go.bobheadxi.dev/ghdb/github"
)

func TestEngine_integration(t *testing.T) {
	engine, err := New(EngineOpts{
		Logger: zaptest.NewLogger(t).Sugar(),
		Auth:   github.NewEnvAuth(),
	})
	require.NoError(t, err)

	_, _, err = engine.Query(&sql.Context{
		Context: context.Background(),
	}, `
		SELECT * FROM this_is_dumb
	`)
	require.NoError(t, err)
}
