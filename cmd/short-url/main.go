package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"log/slog"

	"github.com/Alex322322/short-url/internal/config"
	"github.com/Alex322322/short-url/internal/http/server/handlers/url/save"
	mwLogger "github.com/Alex322322/short-url/internal/http/server/middleware/logger"
	"github.com/Alex322322/short-url/internal/storage/sqlite"
	"github.com/go-chi/chi"

	//"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5/middleware"
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
	router := chi.NewRouter()

	// setup middleware
	router.Use(middleware.RequestID)                 // assign request ID to every request
	router.Use(middleware.Logger)                    // log requests
	router.Use(mwLogger.NewLoggerMiddleware(logger)) // custom logger middleware
	router.Use(middleware.Recoverer)                 // recover from panics
	router.Use(middleware.URLFormat)                 // parse extensions from URL
	router.Use(middleware.Timeout(60 * time.Second))
	//

	router.Post("/url", save.New(logger, storage))

	// TODO run server - net/http
	logger.Info("Starting HTTP server", slog.String("address", cfg.HTTPServer.Address))
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("HTTP server failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Error("HTTP server stopped")
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
