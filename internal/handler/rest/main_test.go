package rest_test

import (
	"os"
	"testing"
	"time"

	"github.com/BruceCompiler/bank/internal/repository/postgres"
	"github.com/BruceCompiler/bank/internal/server"
	"github.com/BruceCompiler/bank/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store postgres.Store) *server.HTTPServer {
	config := utils.Config{
		TokenSynmmetricKey:  utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := server.NewHTTPServer(config, store)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
