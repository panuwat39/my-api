package service

import (
	"context"
	"errors"
	"strings"

	authmodel "github.com/panuwat39/my-api/internal/auth/model"
	authport "github.com/panuwat39/my-api/internal/auth/port"
	userport "github.com/panuwat39/my-api/internal/user/port"
	userservice "github.com/panuwat39/my-api/internal/user/service"
)

type Service struct {
	userRepository userport.UserRepository
	passwordHasher userport.PasswordHasher
	tokenIssuer    authport.TokenIssuer
}

func New(
	userRepository userport.UserRepository,
	passwordHasher userport.PasswordHasher,
	tokenIssuer authport.TokenIssuer,
) *Service {
	return &Service{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
		tokenIssuer:    tokenIssuer,
	}
}

func (s *Service) Login(
	ctx context.Context,
	request authmodel.LoginRequest,
) (authmodel.LoginResponse, error) {
	email := strings.ToLower(
		strings.TrimSpace(request.Email),
	)

	user, err := s.userRepository.FindByEmail(
		ctx,
		email,
	)
	if err != nil {
		if errors.Is(err, userservice.ErrUserNotFound) {
			return authmodel.LoginResponse{},
				ErrInvalidCredentials
		}

		return authmodel.LoginResponse{}, err
	}

	if err := s.passwordHasher.Compare(
		user.PasswordHash,
		request.Password,
	); err != nil {
		return authmodel.LoginResponse{},
			ErrInvalidCredentials
	}

	accessToken, err := s.tokenIssuer.IssueAccessToken(ctx, user.ID, user.Email, string(user.Role))
	if err != nil {
		return authmodel.LoginResponse{}, err
	}

	return authmodel.LoginResponse{
		UserID:      user.ID,
		Email:       user.Email,
		AccessToken: accessToken,
		TokenType:   "Bearer",
	}, nil
}

var _ authport.AuthService = (*Service)(nil)
