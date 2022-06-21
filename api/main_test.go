package api

import (
	"os"
	db "simple_bank/db/sqlc"
	"simple_bank/utils"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	// create a new config object
	config := utils.Config{
		TokenSymmeticKey:    utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server

}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	// Run the test
	os.Exit(m.Run())
}
