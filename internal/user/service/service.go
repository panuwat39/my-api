package service

import (
	"context"
	"net/mail"
	"strings"
	"time"

	"github.com/panuwat39/my-api/internal/user/model"
	"github.com/panuwat39/my-api/internal/user/port"
)

type Service struct {
	repository port.UserRepository
}

func New(repository port.UserRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(ctx context.Context, request model.CreateUserRequest) (model.User, error) {
	name := strings.TrimSpace(request.Name)
	email := strings.ToLower(strings.TrimSpace(request.Email))

	if name == "" {
		return model.User{}, ErrInvalidUserName
	}

	if !isValidEmail(email) {
		return model.User{}, ErrInvalidUserEmail
	}

	existingUser, err := s.repository.FindByEmail(ctx, email)
	if err == nil && existingUser.ID != "" {
		return model.User{}, ErrEmailAlreadyExists
	}

	now := time.Now().UTC()

	user := model.User{
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
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
) ([]model.User, error) {
	return s.repository.List(ctx)
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
		email := strings.ToLower(strings.TrimSpace(*request.Email))
		if !isValidEmail(email) {
			return model.User{}, ErrInvalidUserEmail
		}

		existingUser, err := s.repository.FindByEmail(ctx, email)
		if err == nil && existingUser.ID != user.ID {
			return model.User{}, ErrEmailAlreadyExists
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
