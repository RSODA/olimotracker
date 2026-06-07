package main

import (
	"context"
	"olimotracker/config"
	"olimotracker/internal/auth"
	"olimotracker/internal/categories"
	"olimotracker/internal/sessions"
	"olimotracker/internal/stats"
	"olimotracker/internal/user"
	"olimotracker/pkg/db"
	"olimotracker/pkg/jwttoken"
	"olimotracker/pkg/logger"
	"olimotracker/pkg/middleware"
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

	jwtTTL := time.Duration(cfg.JWT.JWT_TTL) * time.Hour
	jwt := jwttoken.NewTokenGenerator(cfg.JWT.Secret, jwtTTL)

	middleware := middleware.NewMiddlware(jwt, log)

	authRepo := auth.NewRepository(database, log)
	authService := auth.NewService(authRepo, log, jwt)
	authHandler := auth.NewHandler(authService)

	authHandler.RegisterRoutes(r)

	categoriesRepo := categories.NewRepository(database, log)
	categoriesService := categories.NewService(categoriesRepo, log)
	categoriesHandler := categories.NewHandler(categoriesService, log, middleware)

	categoriesHandler.RegisterRoutes(r)

	statsRepo := stats.NewRepo(database, log)
	statsService := stats.NewService(statsRepo, log)
	statsHandler := stats.NewHandler(statsService, middleware)

	statsHandler.RegisterRoutes(r)

	sessionsRepo := sessions.NewRepository(database, log)
	sessionsService := sessions.NewService(sessionsRepo, log, statsService)
	sessionsHandler := sessions.NewHandler(sessionsService, log, middleware)

	sessionsHandler.RegisterRoutes(r)

	userRepo := user.NewRepository(database, log)
	userService := user.NewService(userRepo, log)
	userHandler := user.NewHandler(userService, middleware)

	userHandler.RegisterRoutes(r)

	r.Run(cfg.Http.Addr())
}
