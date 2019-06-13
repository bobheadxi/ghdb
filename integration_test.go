package ghdb

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/src-d/go-mysql-server/sql"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"golang.org/x/oauth2"

	"go.bobheadxi.dev/ghdb/github"
)

func TestEngine_integration(t *testing.T) {
	godotenv.Load()
	engine, err := New(EngineOpts{
		Logger: zaptest.NewLogger(t, zaptest.WrapOptions(zap.AddCaller())).Sugar(),
		Auth: github.NewStaticTokenSource("test", oauth2.Token{
			AccessToken: os.Getenv("GITHUB_TOKEN"),
		}),
		Database: github.DatabaseOpts{PoolSize: 5},
	})
	require.NoError(t, err)

	ctx := context.Background()
	schema, rows, err := engine.Query(sql.NewContext(ctx), `SHOW TABLES;`)
	require.NoError(t, err)
	t.Log("schema length: ", len(schema))
	defer rows.Close()
	for r := rows.Next(); r != nil; r = rows.Next() {

	}
}
