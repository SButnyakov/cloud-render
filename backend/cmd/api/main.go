package main

import (
	"cloud-render/internal/db/postgres"
	"cloud-render/internal/db/redis"
	"cloud-render/internal/http/api"
	"cloud-render/internal/http/middleware/auth"
	"cloud-render/internal/http/middleware/cors"
	mwLogger "cloud-render/internal/http/middleware/logger"
	"cloud-render/internal/lib/config"
	"cloud-render/internal/lib/sl"
	"cloud-render/internal/lib/tokenManager"
	"cloud-render/internal/repository"
	"cloud-render/internal/service"
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
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Envs
	cfgPath := os.Getenv("API_CONFIG_PATH")
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	inputPath := os.Getenv("FILES_INPUT_PATH")

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

	// Redis
	client, err := redis.New(cfg)
	if err != nil {
		log.Error("failed to initialize redis", sl.Err(err))
		os.Exit(-1)
	}
	defer client.Close()

	// Migrating
	if cfg.Env == "dev" || cfg.Env == "prod" {
		postgres.MigrateTop(pg, cfg.DB.MigrationsPath)
	}

	// JWT manager
	jwtManager, err := tokenManager.New(jwtSecretKey)
	if err != nil {
		log.Error("failed to initialize jwt token manager", sl.Err(err))
		os.Exit(-1)
	}

	// Static repos
	subscriptionTypesRepository := repository.NewSubscriptionTypeRepository(pg)
	paymentTypesRepository := repository.NewPaymentTypeRepository(pg)
	orderStatusRepository := repository.NewOrderStatusRepository(pg)

	// Static maps
	subscriptionTypes, err := subscriptionTypesRepository.GetTypesMap()
	if err != nil {
		log.Error("failed to get subscription types", sl.Err(err))
		os.Exit(-1)
	}
	paymentTypes, err := paymentTypesRepository.GetTypesMap()
	if err != nil {
		log.Error("failed to get paymeny types", sl.Err(err))
		os.Exit(-1)
	}
	orderStatusesStrToInt, err := orderStatusRepository.GetStatusesMapStringToInt()
	if err != nil {
		log.Error("failed to get order statuses str to int", sl.Err(err))
		os.Exit(-1)
	}
	orderStatusesIntToStr, err := orderStatusRepository.GetStatusesMapIntToString()
	if err != nil {
		log.Error("failed to get order statuses int to str", sl.Err(err))
		os.Exit(-1)
	}

	// Dynamic repos
	subscriptionRepository := repository.NewSubscriptionRepository(pg)
	orderRepository := repository.NewOrderRepository(pg)

	// Services
	subscriptionService := service.NewSubscriptionService(subscriptionRepository, cfg, subscriptionTypes, paymentTypes)
	orderService := service.NewOrderService(orderRepository, orderStatusesStrToInt, orderStatusesIntToStr,
		inputPath, cfg, client)

	// Router
	router := chi.NewRouter()

	// Router middleware
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(cors.New())
	router.Use(auth.New(log, jwtManager))

	// Router handlers
	router.Post(cfg.Paths.Subscribe, api.Subscribe(log, cfg, subscriptionService))
	router.Get(cfg.Paths.User, api.User(log, subscriptionService))
	router.Post(cfg.Paths.Send, api.Send(log, orderService))

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
