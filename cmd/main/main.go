package main

import (
	"context"
	"olimotracker/config"
	"olimotracker/internal/auth"
	"olimotracker/pkg/db"
	"olimotracker/pkg/jwttoken"
	"olimotracker/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(gin.Recovery())

	log := logger.New(cfg.LoggerLevel)
	log.Info("logger init")

	dbconn, err := pgxpool.New(ctx, cfg.Postgres.DSN())
	if err != nil {
		log.Error("postgres connect", "error", err)
		panic(err)
	}
	defer dbconn.Close()

	log.Info("postgres connect", "dsn", cfg.Postgres.DSN())

	database := db.NewDB(dbconn, log)
	err = database.Ping(ctx)
	if err != nil {
		log.Error("postgres ping", "error", err)
		panic(err)
	}

	log.Info("postgres ping success")

	jwt := jwttoken.NewTokenGenerator(cfg.JWT.Secret, time.Duration(cfg.JWT.JWT_TTL))

	authRepo := auth.NewRepository(database, log)
	authService := auth.NewService(authRepo, log, jwt)
	authHandler := auth.NewHandler(authService)

	authHandler.RegisterRoutes(r)

	r.Run(cfg.Http.Addr())
}
