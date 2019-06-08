package ghdb

import (
	"context"
	"testing"

	"github.com/joho/godotenv"
	"github.com/src-d/go-mysql-server/sql"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"go.bobheadxi.dev/ghdb/github"
)

func TestEngine_integration(t *testing.T) {
	godotenv.Load()
	auth := github.NewEnvAuth()
	t.Logf("%+v", auth)
	engine, err := New(EngineOpts{
		Logger: zaptest.NewLogger(t).Sugar(),
		Auth:   auth,
	})
	require.NoError(t, err)

	_, _, err = engine.Query(&sql.Context{
		Context: context.Background(),
	}, `
		SELECT * FROM this_is_dumb
	`)
	require.NoError(t, err)
}
