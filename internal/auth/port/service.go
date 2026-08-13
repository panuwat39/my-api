package port

import (
	"context"

	"github.com/panuwat39/my-api/internal/auth/model"
)

type AuthService interface {
	Login(
		ctx context.Context,
		request model.LoginRequest,
	) (model.LoginResponse, error)
}
