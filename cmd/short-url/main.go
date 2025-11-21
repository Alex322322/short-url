package main

import (
	"log"
	"os"

	"log/slog"

	"github.com/Alex322322/short-url/internal/config"
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

	// TODO init config - cleanenv
	cfg := config.MustLoad()
	//fmt.Println(cfg)

	// TODO init logger - slog log/slog
	logger := setupLogger(cfg.Env)
	logger.Info("Logger short-url initialized", slog.String("env", cfg.Env))
	logger.Debug("Debug is turned on")

	// TODO init storage - sqlite/postgres
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
