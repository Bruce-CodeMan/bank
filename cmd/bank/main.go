package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/BruceCompiler/bank/internal/repository/postgres"
	"github.com/BruceCompiler/bank/internal/server"
	"github.com/BruceCompiler/bank/utils"
)

func main() {

	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatal("cannot load configuration: ", err)
	}
	pool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect db: ", err)
	}

	store := postgres.NewStore(pool)

	server, err := server.NewHTTPServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

}
