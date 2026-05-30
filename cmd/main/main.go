package main

import (
	"context"
	"olimotracker/config"
	"olimotracker/internal/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	dbconn, err := pgxpool.New(ctx, cfg.Postgres.DSN())
	if err != nil {
		panic(err)
	}
	defer dbconn.Close()

	database := db.NewDB(dbconn)
	err = database.Ping(ctx)
	if err != nil {
		panic(err)
	}

}
