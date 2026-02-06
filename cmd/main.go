package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"auth-go-skd/config"
	httpApp "auth-go-skd/internal/http"
	"auth-go-skd/internal/providers/google"
	"auth-go-skd/internal/service"
	"auth-go-skd/internal/storage/postgres"
	"auth-go-skd/internal/storage/redis"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger, _ := zap.NewProduction()
	if cfg.Log.Level == "debug" {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Infof("Starting %s version %s", cfg.App.Name, cfg.App.Version)

	pg, err := postgres.New(cfg.Postgres)
	if err != nil {
		sugar.Fatalf("failed to connect to postgres: %v", err)
	}
	defer pg.Close()
	sugar.Info("Connected to PostgreSQL")

	rds, err := redis.New(cfg.Redis)
	if err != nil {
		sugar.Fatalf("failed to connect to redis: %v", err)
	}
	defer rds.Close()
	sugar.Info("Connected to Redis")

	googleProv := google.NewGoogleProvider(cfg.OAuth.Google)

	authService := service.NewAuthService(
		pg,
		pg,
		pg,
		googleProv,
		cfg,
	)

	handler := httpApp.NewHandler(authService, sugar, cfg)
	router := handler.InitRoutes()

	srv := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	go func() {
		sugar.Infof("Server is listening on port %s", cfg.HTTP.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	sugar.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		sugar.Fatal("Server forced to shutdown:", err)
	}

	sugar.Info("Server exiting")
}
