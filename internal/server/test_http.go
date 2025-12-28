package server

import (
	"testing"
	"time"

	"github.com/BruceCompiler/bank/internal/repository/postgres"
	"github.com/BruceCompiler/bank/utils"
	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T, store postgres.Store) *HTTPServer {
	config := utils.Config{
		TokenSynmmetricKey:  utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewHTTPServer(config, store)
	require.NoError(t, err)
	return server
}
