package main

import (
	"context"
	"database/sql"
	"olimotracker/config"
	"olimotracker/internal/api"
	"olimotracker/internal/auth"
	"olimotracker/internal/categories"
	"olimotracker/internal/galaxy"
	"olimotracker/internal/sessions"
	"olimotracker/internal/stats"
	"olimotracker/internal/user"
	cr "olimotracker/pkg/cron"
	"olimotracker/pkg/db"
	"olimotracker/pkg/jwttoken"
	"olimotracker/pkg/logger"
	"olimotracker/pkg/middleware"
	"olimotracker/pkg/migrator"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
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

	sqlDB, err := sql.Open("pgx", cfg.Postgres.DSN())
	if err != nil {
		log.Error("failed to open db for migrations", "err", err)
		panic(err)
	}
	mig := migrator.NewMigrator(sqlDB)
	if err := mig.Up(); err != nil {
		log.Error("migration failed", "err", err)
		panic(err)
	}
	defer sqlDB.Close()

	err = mig.Up()
	if err != nil {
		log.Error("failed UP migrations", "err", err)
		return
	}

	log.Info("postgres ping success")

	jwtTTL := time.Duration(cfg.JWT.JWT_TTL) * time.Hour
	jwt := jwttoken.NewTokenGenerator(cfg.JWT.Secret, jwtTTL)

	apiRepo := api.NewRepository(database, log)
	middleware := middleware.NewMiddlware(jwt, log, apiRepo)

	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.APIKeyMiddleware())

	authRepo := auth.NewRepository(database, log)
	authService := auth.NewService(authRepo, log, jwt)
	authHandler := auth.NewHandler(authService)

	authHandler.RegisterRoutes(r)

	categoriesRepo := categories.NewRepository(database, log)
	categoriesService := categories.NewService(categoriesRepo, log)
	categoriesHandler := categories.NewHandler(categoriesService, log, middleware)

	categoriesHandler.RegisterRoutes(r)
	categoriesHandler.RegisterAPIRoutes(apiV1)

	statsRepo := stats.NewRepo(database, log)
	statsService := stats.NewService(statsRepo, log)
	statsHandler := stats.NewHandler(statsService, middleware)

	statsHandler.RegisterRoutes(r)
	statsHandler.RegisterAPIRoutes(apiV1)

	sessionsRepo := sessions.NewRepository(database, log)
	sessionsService := sessions.NewService(sessionsRepo, log, statsService)
	sessionsHandler := sessions.NewHandler(sessionsService, log, middleware)

	sessionsHandler.RegisterRoutes(r)
	sessionsHandler.RegisterAPIRoutes(apiV1)

	userRepo := user.NewRepository(database, log)
	userService := user.NewService(userRepo, log)
	userHandler := user.NewHandler(userService, middleware)

	userHandler.RegisterRoutes(r)

	galaxyRepo := galaxy.NewRepository(database, log)
	galaxyService := galaxy.NewService(statsService, sessionsRepo, log, galaxyRepo)
	galaxyHandler := galaxy.NewHandler(galaxyService, middleware)
	galaxyHandler.RegisterRoutes(r)

	c := cron.New()

	cron := cr.NewCron(c, log, galaxyService, statsService)
	cron.AddsCronJobs()

	c.Start()
	defer c.Stop()

	r.Run(cfg.Http.Addr())
}
