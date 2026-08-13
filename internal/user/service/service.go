package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/panuwat39/my-api/internal/shared/pagination"
	"github.com/panuwat39/my-api/internal/user/model"
	"github.com/panuwat39/my-api/internal/user/port"
)

type Service struct {
	repository     port.UserRepository
	passwordHasher port.PasswordHasher
}

func New(
	repository port.UserRepository,
	passwordHasher port.PasswordHasher,
) *Service {
	return &Service{
		repository:     repository,
		passwordHasher: passwordHasher,
	}
}

func (s *Service) Create(
	ctx context.Context,
	request model.CreateUserRequest,
) (model.User, error) {
	name := strings.TrimSpace(request.Name)
	email := strings.ToLower(
		strings.TrimSpace(request.Email),
	)

	if name == "" {
		return model.User{}, ErrInvalidUserName
	}

	if !isValidEmail(email) {
		return model.User{}, ErrInvalidUserEmail
	}

	passwordBytes := []byte(request.Password)

	if len(passwordBytes) < 8 ||
		len(passwordBytes) > 72 {
		return model.User{}, ErrInvalidUserPassword
	}

	existingUser, err := s.repository.FindByEmail(
		ctx,
		email,
	)

	switch {
	case err == nil && existingUser.ID != "":
		return model.User{}, ErrEmailAlreadyExists

	case err != nil &&
		!errors.Is(err, ErrUserNotFound):
		return model.User{}, err
	}

	passwordHash, err := s.passwordHasher.Hash(
		request.Password,
	)
	if err != nil {
		return model.User{}, fmt.Errorf(
			"hash user password: %w",
			err,
		)
	}

	now := time.Now().UTC()

	user := model.User{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         model.RoleUser,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return s.repository.Create(ctx, user)
}

func (s *Service) GetByID(
	ctx context.Context,
	id string,
) (model.User, error) {
	if strings.TrimSpace(id) == "" {
		return model.User{}, ErrUserNotFound
	}

	return s.repository.FindByID(ctx, id)
}

func (s *Service) List(
	ctx context.Context,
	query pagination.Query,
) ([]model.User, int64, error) {
	return s.repository.List(ctx, query)
}

func (s *Service) Update(
	ctx context.Context,
	id string,
	request model.UpdateUserRequest,
) (model.User, error) {
	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	if request.Name != nil {
		name := strings.TrimSpace(*request.Name)

		if name == "" {
			return model.User{}, ErrInvalidUserName
		}

		user.Name = name
	}

	if request.Email != nil {
		email := strings.ToLower(
			strings.TrimSpace(*request.Email),
		)

		if !isValidEmail(email) {
			return model.User{}, ErrInvalidUserEmail
		}

		existingUser, err := s.repository.FindByEmail(
			ctx,
			email,
		)

		switch {
		case err == nil &&
			existingUser.ID != user.ID:
			return model.User{}, ErrEmailAlreadyExists

		case err != nil &&
			!errors.Is(err, ErrUserNotFound):
			return model.User{}, err
		}

		user.Email = email
	}

	user.UpdatedAt = time.Now().UTC()

	return s.repository.Update(ctx, user)
}

func (s *Service) Delete(
	ctx context.Context,
	id string,
) error {
	if strings.TrimSpace(id) == "" {
		return ErrUserNotFound
	}

	return s.repository.Delete(ctx, id)
}

func isValidEmail(email string) bool {
	if email == "" {
		return false
	}

	address, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	return address.Address == email
}

var _ port.UserService = (*Service)(nil)
