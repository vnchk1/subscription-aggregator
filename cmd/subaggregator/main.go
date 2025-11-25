package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/vnchk1/subscription-aggregator/internal/config"
	"github.com/vnchk1/subscription-aggregator/internal/db"
	"github.com/vnchk1/subscription-aggregator/internal/handler"
	logging "github.com/vnchk1/subscription-aggregator/internal/logger"
	"github.com/vnchk1/subscription-aggregator/internal/migration"
	"github.com/vnchk1/subscription-aggregator/internal/repository"
	"github.com/vnchk1/subscription-aggregator/internal/server"
	"github.com/vnchk1/subscription-aggregator/internal/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	logger := logging.NewLogger(cfg.Logger.LogLevel)

	ctx := context.Background()

	migrator, err := migration.NewMigrator(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}

	defer migrator.Close()

	if err = migrator.Up(ctx); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	logger.Debug("Migrations completed")

	pool, err := db.NewPool(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("Failed to create database connection pool: %v", err)
	}
	defer db.ClosePool(pool)

	logger.Debug("Connected to database")

	subscriptionRepo := repository.NewSubscriptionRepository(pool)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)

	srv := server.New(cfg.Server, logger)

	server.SetupRouter(srv.GetEchoInstance(), subscriptionHandler, logger)

	go func() {
		if err = srv.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	logger.Debug("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Debug("Server exited")
}
