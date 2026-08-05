package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func CreateUserIndexes(
	ctx context.Context,
	database *mongo.Database,
) error {
	collection := database.Collection("users")

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "email", Value: 1},
			},
			Options: options.Index().
				SetName("users_email_unique").
				SetUnique(true),
		},
	}

	if _, err := collection.Indexes().CreateMany(ctx, indexes); err != nil {
		return fmt.Errorf("create user indexes: %w", err)
	}

	return nil
}
