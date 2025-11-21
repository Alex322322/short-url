package main

import (
	"log"
	"os"
	"time"

	"log/slog"

	"github.com/Alex322322/short-url/internal/config"
	"github.com/Alex322322/short-url/internal/storage/sqlite"
	"github.com/joho/godotenv"
)

const (
	envLocal       = "local"
	envProduction  = "prod"
	envDevelopment = "dev"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// init config
	cfg := config.MustLoad()
	//fmt.Println(cfg)

	// init logger
	logger := setupLogger(cfg.Env)
	logger.Info("Logger short-url initialized", slog.String("env", cfg.Env))
	logger.Debug("Debug is turned on")

	// init storage
	storage, err := sqlite.NewStorage(cfg.StoragePath)
	if err != nil {
		logger.Error(
			"Failed to initialize storage",
			slog.String("storage_path", cfg.StoragePath), 
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}


	_ = storage // TODO remove after storage used

	// TODO init router - chi


	// TODO run server - net/http
}

// setup logger based on environment
func setupLogger(env string) *slog.Logger {

	var logger *slog.Logger
	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDevelopment:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProduction:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logger
}
