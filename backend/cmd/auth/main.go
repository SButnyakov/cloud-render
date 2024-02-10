package main

import (
	"cloud-render/internal/db/postgres"
	"cloud-render/internal/lib/config"
	"cloud-render/internal/lib/sl"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	// Envs
	cfgPath := os.Getenv("AUTH_CONFIG_PATH")

	// Config
	cfg := config.MustLoad(cfgPath)

	// Logger
	log := sl.SetupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))
	log.Info("initializing server", slog.String("host", cfg.HTTPServer.Host), slog.Int("port", cfg.HTTPServer.Port))
	log.Debug("logger debug mode enabled")

	// DB
	pg, err := postgres.New(cfg.DB)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
		os.Exit(-1)
	}
	defer pg.Close()

	// Migrating
	if cfg.Env == "dev" || cfg.Env == "prod" {
		postgres.MigrateTop(pg, "file://../../../migrations/auth/postgres")
	}

	// Router
	router := chi.NewRouter()

	// Server
	httpServer := http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port),
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// Startup
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("starting server", slog.String("address", httpServer.Addr))
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error.", sl.Err(err))
			os.Exit(1)
		}
	}()

	<-interrupt
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error("server shutdown failed.", sl.Err(err))
		os.Exit(1)
	}
}
