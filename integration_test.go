package ghdb

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/src-d/go-mysql-server/sql"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"golang.org/x/oauth2"

	"go.bobheadxi.dev/ghdb/github"
)

func TestEngine_integration(t *testing.T) {
	godotenv.Load()
	engine, err := New(EngineOpts{
		Logger: zaptest.NewLogger(t).Sugar(),
		Auth: github.NewStaticTokenSource("test", oauth2.Token{
			AccessToken: os.Getenv("GITHUB_TOKEN"),
		}),
	})
	require.NoError(t, err)

	_, _, err = engine.Query(sql.NewContext(context.Background()), `
		SELECT * FROM ubclaunchpad/inertia
	`)
	require.NoError(t, err)
}
