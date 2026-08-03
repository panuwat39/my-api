package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/panuwat39/my-api/infrastructure/config"
	"github.com/panuwat39/my-api/infrastructure/database"
	"github.com/panuwat39/my-api/infrastructure/router"
	"github.com/panuwat39/my-api/infrastructure/server"
)

func main() {
	cfg := config.Load()

	mongoClient, err := database.ConnectMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	mongoDatabase := mongoClient.Database(cfg.MongoDatabase)
	_ = mongoDatabase

	handler := router.New()
	httpServer := server.NewHTTPServer(cfg.HTTPPort, handler)

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf(
			"server running in %s mode on http://localhost:%s",
			cfg.AppEnv,
			cfg.HTTPPort,
		)

		serverErrors <- httpServer.Start()
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(
		shutdownSignal,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	select {
	case signalReceived := <-shutdownSignal:
		log.Printf("shutdown signal received: %s", signalReceived)

	case err := <-serverErrors:
		log.Printf("server stopped unexpectedly: %v", err)
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	if err := database.DisconnectMongoDB(mongoClient); err != nil {
		log.Printf("MongoDB disconnect error: %v", err)
	}

	log.Println("server stopped")
}
