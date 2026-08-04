package port

import (
	"context"

	"github.com/panuwat39/my-api/internal/user/model"
)

type UserRepository interface {
	Create(
		ctx context.Context,
		user model.User,
	) (model.User, error)

	FindByID(
		ctx context.Context,
		id string,
	) (model.User, error)

	FindByEmail(
		ctx context.Context,
		email string,
	) (model.User, error)

	List(
		ctx context.Context,
	) ([]model.User, error)

	Update(
		ctx context.Context,
		user model.User,
	) (model.User, error)

	Delete(
		ctx context.Context,
		id string,
	) error
}
