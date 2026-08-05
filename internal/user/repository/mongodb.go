package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/panuwat39/my-api/internal/user/model"
	"github.com/panuwat39/my-api/internal/user/port"
	userservice "github.com/panuwat39/my-api/internal/user/service"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const userCollectionName = "users"

type userDocument struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Name      string        `bson:"name"`
	Email     string        `bson:"email"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

type MongoDBRepository struct {
	collection *mongo.Collection
}

func NewMongoDBRepository(database *mongo.Database) *MongoDBRepository {
	return &MongoDBRepository{
		collection: database.Collection(userCollectionName),
	}
}

func (r *MongoDBRepository) Create(
	ctx context.Context,
	user model.User,
) (model.User, error) {
	document := userDocument{
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	result, err := r.collection.InsertOne(ctx, document)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return model.User{}, userservice.ErrEmailAlreadyExists
		}

		return model.User{}, fmt.Errorf("insert user: %w", err)
	}

	objectID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return model.User{}, errors.New("insert user: invalid inserted ID")
	}

	user.ID = objectID.Hex()

	return user, nil
}

func (r *MongoDBRepository) FindByID(
	ctx context.Context,
	id string,
) (model.User, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return model.User{}, userservice.ErrUserNotFound
	}

	var document userDocument

	err = r.collection.FindOne(
		ctx,
		bson.M{"_id": objectID},
	).Decode(&document)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.User{}, userservice.ErrUserNotFound
	}

	if err != nil {
		return model.User{}, fmt.Errorf("find user by ID: %w", err)
	}

	return document.toModel(), nil
}

func (r *MongoDBRepository) FindByEmail(
	ctx context.Context,
	email string,
) (model.User, error) {
	var document userDocument

	err := r.collection.FindOne(
		ctx,
		bson.M{"email": email},
	).Decode(&document)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.User{}, userservice.ErrUserNotFound
	}

	if err != nil {
		return model.User{}, fmt.Errorf("find user by email: %w", err)
	}

	return document.toModel(), nil
}

func (r *MongoDBRepository) List(
	ctx context.Context,
) ([]model.User, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{},
		options.Find().SetSort(bson.D{
			{Key: "created_at", Value: -1},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("find users: %w", err)
	}
	defer cursor.Close(ctx)

	users := make([]model.User, 0)

	for cursor.Next(ctx) {
		var document userDocument

		if err := cursor.Decode(&document); err != nil {
			return nil, fmt.Errorf("decode user: %w", err)
		}

		users = append(users, document.toModel())
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}

func (r *MongoDBRepository) Update(
	ctx context.Context,
	user model.User,
) (model.User, error) {
	objectID, err := bson.ObjectIDFromHex(user.ID)
	if err != nil {
		return model.User{}, userservice.ErrUserNotFound
	}

	update := bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"email":      user.Email,
			"updated_at": user.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		update,
	)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return model.User{}, userservice.ErrEmailAlreadyExists
		}

		return model.User{}, fmt.Errorf("update user: %w", err)
	}

	if result.MatchedCount == 0 {
		return model.User{}, userservice.ErrUserNotFound
	}

	return user, nil
}

func (r *MongoDBRepository) Delete(
	ctx context.Context,
	id string,
) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return userservice.ErrUserNotFound
	}

	result, err := r.collection.DeleteOne(
		ctx,
		bson.M{"_id": objectID},
	)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	if result.DeletedCount == 0 {
		return userservice.ErrUserNotFound
	}

	return nil
}

func (d userDocument) toModel() model.User {
	return model.User{
		ID:        d.ID.Hex(),
		Name:      d.Name,
		Email:     d.Email,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

var _ port.UserRepository = (*MongoDBRepository)(nil)
