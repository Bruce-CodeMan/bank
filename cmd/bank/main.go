package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/BruceCompiler/bank/internal/repository/postgres"
	"github.com/BruceCompiler/bank/internal/server"
)

const (
	dbSource      = "postgresql://root:secret@localhost:5432/bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	pool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect db: ", err)
	}

	store := postgres.NewStore(pool)

	server := server.NewHTTPServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

}
