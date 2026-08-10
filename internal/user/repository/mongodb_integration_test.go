//go:build integration

package repository

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/panuwat39/my-api/internal/shared/pagination"
	"github.com/panuwat39/my-api/internal/user/model"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestMongoDBRepositoryListPagination(t *testing.T) {
	uri := os.Getenv("MONGO_TEST_URI")
	databaseName := os.Getenv("MONGO_TEST_DATABASE")

	if uri == "" {
		t.Fatal("MONGO_TEST_URI is required")
	}

	if databaseName == "" {
		t.Fatal("MONGO_TEST_DATABASE is required")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	client, err := mongo.Connect(
		options.Client().ApplyURI(uri),
	)
	if err != nil {
		t.Fatalf("connect MongoDB: %v", err)
	}

	t.Cleanup(func() {
		disconnectCtx, disconnectCancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer disconnectCancel()

		if err := client.Disconnect(disconnectCtx); err != nil {
			t.Errorf("disconnect MongoDB: %v", err)
		}
	})

	if err := client.Ping(ctx, nil); err != nil {
		t.Fatalf("ping MongoDB: %v", err)
	}

	database := client.Database(databaseName)

	collection := database.Collection("users")

	if err := collection.Drop(ctx); err != nil {
		t.Fatalf("drop users collection before test: %v", err)
	}

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cleanupCancel()

		if err := collection.Drop(cleanupCtx); err != nil {
			t.Errorf("drop users collection after test: %v", err)
		}
	})

	repository := NewMongoDBRepository(database)

	baseTime := time.Date(
		2026,
		time.January,
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	for i := 1; i <= 5; i++ {
		createdAt := baseTime.Add(
			time.Duration(i) * time.Hour,
		)

		_, err := repository.Create(
			ctx,
			model.User{
				Name:      "Integration User",
				Email:     integrationTestEmail(i),
				CreatedAt: createdAt,
				UpdatedAt: createdAt,
			},
		)
		if err != nil {
			t.Fatalf(
				"create test user %d: %v",
				i,
				err,
			)
		}
	}

	query := pagination.NewQuery(2, 2)

	users, total, err := repository.List(
		ctx,
		query,
	)
	if err != nil {
		t.Fatalf("list users: %v", err)
	}

	if total != 5 {
		t.Fatalf(
			"expected total 5, got %d",
			total,
		)
	}

	if len(users) != 2 {
		t.Fatalf(
			"expected 2 users, got %d",
			len(users),
		)
	}

	if users[0].Email != "integration-3@example.com" {
		t.Errorf(
			"expected first email integration-3@example.com, got %s",
			users[0].Email,
		)
	}

	if users[1].Email != "integration-2@example.com" {
		t.Errorf(
			"expected second email integration-2@example.com, got %s",
			users[1].Email,
		)
	}
}

func integrationTestEmail(number int) string {
	return "integration-" +
		strconv.Itoa(number) +
		"@example.com"
}
