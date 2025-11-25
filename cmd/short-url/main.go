package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log/slog"

	"github.com/Alex322322/short-url/internal/config"
	"github.com/Alex322322/short-url/internal/http/server/handlers/url/delete"
	"github.com/Alex322322/short-url/internal/http/server/handlers/url/redirect"
	"github.com/Alex322322/short-url/internal/http/server/handlers/url/save"
	mwLogger "github.com/Alex322322/short-url/internal/http/server/middleware/logger"
	"github.com/Alex322322/short-url/internal/storage/postgres"
	"github.com/go-chi/chi/v5"

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

	// init logger
	logger := setupLogger(cfg.Env)
	logger.Info("Logger short-url initialized", slog.String("env", cfg.Env))
	logger.Debug("Debug is turned on")

	// init storage
	//storage, err := sqlite.NewStorage(cfg.StoragePath)

	storage, err := postgres.NewStorage(config.BuildConnString(cfg.ConfigPostgres))
	if err != nil {
		logger.Error(
			"Failed to initialize storage",
			slog.String("port: ", cfg.ConfigPostgres.Port),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	// init router - chi
	router := chi.NewRouter()

	// setup middleware
	router.Use(middleware.RequestID)                 // assign request ID to every request
	router.Use(middleware.Logger)                    // log requests
	router.Use(mwLogger.NewLoggerMiddleware(logger)) // custom logger middleware
	router.Use(middleware.Recoverer)                 // recover from panics
	router.Use(middleware.URLFormat)                 // parse extensions from URL
	router.Use(middleware.Timeout(60 * time.Second))

	// setup authorization middleware
	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("short-url", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", save.New(logger, storage))
		r.Delete("/{alias}", delete.New(logger, storage))
	})

	// setup routes
	//router.Post("/url", save.New(logger, storage))
	router.Get("/{alias}", redirect.New(logger, storage))
	//router.Delete("/url/{alias}", delete.New(logger, storage))

	// setup server
	logger.Info("Starting HTTP server", slog.String("address", cfg.HTTPServer.Address))
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}


	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Gracefull shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("failed to start server")
		}
	}()

	logger.Info("server started")

	<-done
	logger.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("failed to stop server", slog.String("error", err.Error()))

		return
	}

	// TODO: close storage
	logger.Info("HTTP server stopped")
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
