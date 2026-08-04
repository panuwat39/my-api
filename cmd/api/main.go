package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/panuwat39/my-api/infrastructure/config"
	"github.com/panuwat39/my-api/infrastructure/database"
	"github.com/panuwat39/my-api/infrastructure/logger"
	"github.com/panuwat39/my-api/infrastructure/router"
	"github.com/panuwat39/my-api/infrastructure/server"
)

func main() {
	cfg := config.Load()
	appLogger := logger.New(cfg.AppEnv)

	mongoClient, err := database.ConnectMongoDB(cfg.MongoURI)
	if err != nil {
		appLogger.Error(
			"mongodb connection failed",
			"error", err,
		)
		os.Exit(1)
	}

	mongoDatabase := mongoClient.Database(cfg.MongoDatabase)
	_ = mongoDatabase

	handler := router.New(appLogger)
	httpServer := server.NewHTTPServer(cfg.HTTPPort, handler)

	serverErrors := make(chan error, 1)

	go func() {
		appLogger.Info(
			"server starting",
			"environment", cfg.AppEnv,
			"address", "http://localhost:"+cfg.HTTPPort,
			"database", cfg.MongoDatabase,
		)

		serverErrors <- httpServer.Start()
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(
		shutdownSignal,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer signal.Stop(shutdownSignal)

	select {
	case signalReceived := <-shutdownSignal:
		appLogger.Warn(
			"shutdown signal received",
			"signal", signalReceived.String(),
		)

	case err := <-serverErrors:
		if err != nil {
			appLogger.Error(
				"server stopped unexpectedly",
				"error", err,
			)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		appLogger.Error(
			"http server shutdown failed",
			"error", err,
		)
	}

	if err := database.DisconnectMongoDB(mongoClient); err != nil {
		appLogger.Error(
			"mongodb disconnect failed",
			"error", err,
		)
	}

	appLogger.Info("server stopped")
}
