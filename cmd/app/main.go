package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alextotalk/feline-intelligence/internal/config"
	"github.com/alextotalk/feline-intelligence/internal/delivery/handlers"
	"github.com/alextotalk/feline-intelligence/internal/infrastructure/catapi"
	"github.com/alextotalk/feline-intelligence/internal/infrastructure/repository"
	"github.com/alextotalk/feline-intelligence/internal/lib/logger/handlers/slogpretty"
	"github.com/alextotalk/feline-intelligence/internal/lib/logger/sl"
	"github.com/alextotalk/feline-intelligence/internal/storage/pg"
	"github.com/alextotalk/feline-intelligence/internal/usecase"

	_ "github.com/alextotalk/feline-intelligence/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @title Feline Intelligence API
// @version 1.0
// @description APIs to manage spy cats, missions and goals.
// @host localhost:8080
// @BasePath /
func main() {
	cfg, err := config.LoadConfig("config/local.yaml")
	if err != nil {
		log.Fatalf("cannot read config: %v", err)
	}

	logger := InitLogger(cfg.App.Env)
	logger.Info("Starting application", "app", cfg.App.Name, "env", cfg.App.Env)

	db, err := pg.NewPostgres(cfg)
	if err != nil {
		logger.Error("Failed to initialize Postgres", sl.Err(err))
		os.Exit(1)
	}
	logger.Info("Successfully connected to Postgres", "host", cfg.Database.Host)

	catAPI := catapi.NewCatAPI("https://api.thecatapi.com", "") // наприклад, cfg.App.TheCatAPIKey

	catRepo := repository.NewCatPgRepository(db)
	missionRepo := repository.NewMissionPgRepository(db)
	targetRepo := repository.NewTargetPgRepository(db)

	catUC := usecase.NewCatUsecase(catRepo, catAPI)
	missionUC := usecase.NewMissionUsecase(missionRepo, targetRepo, catRepo)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	handlers.NewCatHandler(e, catUC)
	handlers.NewMissionHandler(e, missionUC)

	go func() {
		address := fmt.Sprintf(":%d", cfg.Server.Port)
		logger.Info("Starting server", "address", address)
		if err := e.Start(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Error starting server", sl.Err(err))
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	logger.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		logger.Error("Error shutting down server", sl.Err(err))
	}

	if err := db.Close(); err != nil {
		logger.Error("Error closing database connection", sl.Err(err))
	}

	logger.Info("Application shut down gracefully.")
}

func InitLogger(env string) *slog.Logger {
	switch env {
	case envLocal:
		return setupPrettySlog()
	case envDev:
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
