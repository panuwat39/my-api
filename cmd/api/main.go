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

	sharedpassword "github.com/panuwat39/my-api/internal/shared/password"
	usercontroller "github.com/panuwat39/my-api/internal/user/controller"
	userrepository "github.com/panuwat39/my-api/internal/user/repository"
	userservice "github.com/panuwat39/my-api/internal/user/service"

	authcontroller "github.com/panuwat39/my-api/internal/auth/controller"
	authservice "github.com/panuwat39/my-api/internal/auth/service"

	sharedtoken "github.com/panuwat39/my-api/internal/shared/token"

	mongodbmigration "github.com/panuwat39/my-api/migrations/mongodb"

	authmiddleware "github.com/panuwat39/my-api/internal/auth/middleware"

	usermodel "github.com/panuwat39/my-api/internal/user/model"
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

	migrationCtx, migrationCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	if err := mongodbmigration.CreateUserIndexes(
		migrationCtx,
		mongoDatabase,
	); err != nil {
		migrationCancel()

		appLogger.Error(
			"mongodb migration failed",
			"error", err,
		)

		if disconnectErr := database.DisconnectMongoDB(mongoClient); disconnectErr != nil {
			appLogger.Error(
				"mongodb disconnect failed",
				"error", disconnectErr,
			)
		}

		os.Exit(1)
	}

	migrationCancel()

	appLogger.Info(
		"mongodb indexes ready",
		"database", cfg.MongoDatabase,
	)

	userRepository := userrepository.NewMongoDBRepository(
		mongoDatabase,
	)

	passwordHasher := sharedpassword.NewHasher()

	tokenIssuer, err := sharedtoken.NewJWTIssuer(
		cfg.JWTSecret,
		cfg.JWTIssuer,
		cfg.JWTAccessTTL,
	)
	if err != nil {
		appLogger.Error(
			"jwt issuer initialization failed",
			"error", err,
		)

		_ = database.DisconnectMongoDB(
			mongoClient,
		)

		os.Exit(1)
	}

	authenticate := authmiddleware.Authenticate(tokenIssuer)

	requireAdmin := authmiddleware.RequireRole(
		string(usermodel.RoleAdmin),
	)

	requireSelfOrAdmin := authmiddleware.RequireSelfOrRole(
		"id",
		string(usermodel.RoleAdmin),
	)

	userService := userservice.New(
		userRepository,
		passwordHasher,
	)

	userController := usercontroller.New(
		userService,
		appLogger,
	)

	authService := authservice.New(
		userRepository,
		passwordHasher,
		tokenIssuer,
	)

	authController := authcontroller.New(
		authService,
		appLogger,
	)

	app := router.New(
		router.Dependencies{
			Logger:             appLogger,
			UserController:     userController,
			AuthController:     authController,
			AuthMiddleware:     authenticate,
			RequireAdmin:       requireAdmin,
			RequireSelfOrAdmin: requireSelfOrAdmin,
		},
	)

	httpServer := server.NewHTTPServer(
		cfg.HTTPPort,
		app,
	)

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

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer shutdownCancel()

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
