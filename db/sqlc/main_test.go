package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/BruceCompiler/bank/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig()
	pool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect db: ", err)
	}
	testPool = pool

	// Test ping
	if err := pingPGXPool(context.Background(), testPool); err != nil {
		testPool.Close()
		log.Fatal("cannot ping db: ", err)
	}

	testQueries = New(testPool)

	code := m.Run()
	testPool.Close()
	os.Exit(code)
}

func pingPGXPool(ctx context.Context, pool *pgxpool.Pool) error {
	c, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer c.Release()

	_, err = c.Exec(ctx, "SELECT 1")
	return err
}
