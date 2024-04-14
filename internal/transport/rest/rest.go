package rest

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	"log/slog"
	"open-api-games/internal/config"
	"open-api-games/internal/repository"
	"open-api-games/internal/service/game_processor"
	"open-api-games/internal/service/seed"
	"open-api-games/internal/transport/rest/game_processor_handler"
	"os"
	"os/signal"
	"syscall"
)

func Start() error {
	// Base context
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Create a slog logger
	logLevel := new(slog.LevelVar)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: logLevel}))

	// Config
	logger.Info("config initializing...")
	cfg := config.Get(logger)

	// Set log level
	logLevel.Set(cfg.GetSlogLevel())

	// Initialize repository
	logger.Info("repositories initializing...")
	repo, err := repository.NewRepo(ctx, logger)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		return err
	}

	// Ensure repo indexes
	err = repo.EnsureIndexes(ctx)
	if err != nil {
		logger.Error("failed to ensure indexes", "error", err)
		return err
	}

	// Initialize service
	logger.Info("services initializing...")
	gameProcessor := game_processor.New(repo, logger)

	// Initialize handler
	logger.Info("handlers initializing...")
	gameProcessorHandler := game_processor_handler.New(gameProcessor, logger)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(slogecho.New(logger))
	e.Use(middleware.Recover())

	// Routes with check sign middleware
	rootGroup := e.Group("/open-api-games", gameProcessorHandler.CheckSign)
	v1Group := rootGroup.Group("/v1")
	v1Group.POST("/games-processor", gameProcessorHandler.Process)

	// Healthcheck
	e.GET("/healthcheck", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	// Seeds for testing
	// TODO: remove this block after testing
	seedService := seed.New(repo, logger)
	//repo.DropTest(ctx, logger)
	seedService.Seed(ctx)
	// end block to remove

	// Start server
	logger.Info("transport server starting...")

	go func() {
		err := e.Start(cfg.HTTPAddr)
		if err != nil {
			logger.Error("transport error", "error", err)
		}
	}()

	// TODO: wait for shutdown of dependencies with tasks in progress
	<-ctx.Done()

	logger.Info("closing repository...")
	err = repo.Close(ctx)
	if err != nil {
		logger.Error("failed to close repository", "error", err)
		return err
	}

	logger.Info("transport server stopping...")
	err = e.Shutdown(ctx)
	if err != nil {
		logger.Error("transport stop error", "error", err)
		return err
	}

	return nil
}
