package port

import (
	"context"

	"github.com/panuwat39/my-api/internal/shared/pagination"
	"github.com/panuwat39/my-api/internal/user/model"
)

type UserService interface {
	Create(
		ctx context.Context,
		request model.CreateUserRequest,
	) (model.User, error)

	GetByID(
		ctx context.Context,
		id string,
	) (model.User, error)

	List(
		ctx context.Context,
		query pagination.Query,
	) ([]model.User, int64, error)

	Update(
		ctx context.Context,
		id string,
		request model.UpdateUserRequest,
	) (model.User, error)

	Delete(
		ctx context.Context,
		id string,
	) error
}
