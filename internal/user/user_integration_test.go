//go:build integration

package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/panuwat39/my-api/internal/user/controller"
	"github.com/panuwat39/my-api/internal/user/repository"
	"github.com/panuwat39/my-api/internal/user/route"
	"github.com/panuwat39/my-api/internal/user/service"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	sharedpassword "github.com/panuwat39/my-api/internal/shared/password"
)

func TestUserHTTPIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		15*time.Second,
	)
	defer cancel()

	client, database := setupMongoDB(t, ctx)

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cleanupCancel()

		if err := database.Drop(cleanupCtx); err != nil {
			t.Errorf("drop test database: %v", err)
		}

		if err := client.Disconnect(cleanupCtx); err != nil {
			t.Errorf("disconnect MongoDB: %v", err)
		}
	})

	app := setupApp(database)

	t.Run("create user", func(t *testing.T) {
		body := []byte(`{
			"name": "Panuwat",
			"email": "panuwat.integration@example.com",
			"password": "Password123"
		}`)

		request := httptest.NewRequest(
			http.MethodPost,
			"/api/v1/users",
			bytes.NewReader(body),
		)

		request.Header.Set(
			"Content-Type",
			"application/json",
		)

		response, err := app.Test(request)
		if err != nil {
			t.Fatalf("send request: %v", err)
		}
		defer response.Body.Close()

		if response.StatusCode != fiber.StatusCreated {
			responseBody, _ := io.ReadAll(response.Body)

			t.Fatalf(
				"expected status %d, got %d: %s",
				fiber.StatusCreated,
				response.StatusCode,
				string(responseBody),
			)
		}

		var result struct {
			Data struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"data"`
		}

		if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if result.Data.ID == "" {
			t.Fatal("expected user ID")
		}

		if result.Data.Name != "Panuwat" {
			t.Errorf(
				"expected name Panuwat, got %s",
				result.Data.Name,
			)
		}

		if result.Data.Email != "panuwat.integration@example.com" {
			t.Errorf(
				"unexpected email: %s",
				result.Data.Email,
			)
		}
	})

	t.Run("list users", func(t *testing.T) {
		request := httptest.NewRequest(
			http.MethodGet,
			"/api/v1/users?page=1&limit=10",
			nil,
		)

		response, err := app.Test(request)
		if err != nil {
			t.Fatalf("send request: %v", err)
		}
		defer response.Body.Close()

		if response.StatusCode != fiber.StatusOK {
			responseBody, _ := io.ReadAll(response.Body)

			t.Fatalf(
				"expected status %d, got %d: %s",
				fiber.StatusOK,
				response.StatusCode,
				string(responseBody),
			)
		}

		var result struct {
			Data struct {
				Items []struct {
					ID    string `json:"id"`
					Name  string `json:"name"`
					Email string `json:"email"`
				} `json:"items"`

				Meta struct {
					Page       int   `json:"page"`
					Limit      int   `json:"limit"`
					Total      int64 `json:"total"`
					TotalPages int64 `json:"total_pages"`
				} `json:"meta"`
			} `json:"data"`
		}

		if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if len(result.Data.Items) != 1 {
			t.Fatalf(
				"expected 1 user, got %d",
				len(result.Data.Items),
			)
		}

		if result.Data.Meta.Total != 1 {
			t.Errorf(
				"expected total 1, got %d",
				result.Data.Meta.Total,
			)
		}

		if result.Data.Meta.Page != 1 {
			t.Errorf(
				"expected page 1, got %d",
				result.Data.Meta.Page,
			)
		}

		if result.Data.Meta.Limit != 10 {
			t.Errorf(
				"expected limit 10, got %d",
				result.Data.Meta.Limit,
			)
		}
	})
}

func setupMongoDB(
	t *testing.T,
	ctx context.Context,
) (*mongo.Client, *mongo.Database) {
	t.Helper()

	uri := os.Getenv("MONGO_TEST_URI")

	baseDatabaseName := os.Getenv("MONGO_TEST_DATABASE")

	if baseDatabaseName == "" {
		t.Fatal("MONGO_TEST_DATABASE is required")
	}

	databaseName := baseDatabaseName + "_http"

	if uri == "" {
		t.Fatal("MONGO_TEST_URI is required")
	}

	if databaseName == "" {
		t.Fatal("MONGO_TEST_DATABASE is required")
	}

	client, err := mongo.Connect(
		options.Client().ApplyURI(uri),
	)
	if err != nil {
		t.Fatalf("connect MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(context.Background())

		t.Fatalf("ping MongoDB: %v", err)
	}

	database := client.Database(databaseName)

	if err := database.Collection("users").Drop(ctx); err != nil {
		_ = client.Disconnect(context.Background())

		t.Fatalf("drop users collection before test: %v", err)
	}

	return client, database
}

func setupApp(
	database *mongo.Database,
) *fiber.App {
	userRepository := repository.NewMongoDBRepository(
		database,
	)

	passwordHasher := sharedpassword.NewHasher()

	userService := service.New(
		userRepository,
		passwordHasher,
	)

	logger := slog.New(
		slog.NewTextHandler(
			io.Discard,
			nil,
		),
	)

	userController := controller.New(
		userService,
		logger,
	)

	app := fiber.New()

	route.RegisterV1(
		app,
		userController,
	)

	return app
}
