package service

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/panuwat39/my-api/internal/user/model"
// )

// type fakeUserRepository struct {
// 	createFunc      func(context.Context, model.User) (model.User, error)
// 	findByIDFunc    func(context.Context, string) (model.User, error)
// 	findByEmailFunc func(context.Context, string) (model.User, error)
// 	listFunc        func(context.Context) ([]model.User, error)
// 	updateFunc      func(context.Context, model.User) (model.User, error)
// 	deleteFunc      func(context.Context, string) error
// }

// func (f *fakeUserRepository) Create(
// 	ctx context.Context,
// 	user model.User,
// ) (model.User, error) {
// 	return f.createFunc(ctx, user)
// }

// func (f *fakeUserRepository) FindByID(
// 	ctx context.Context,
// 	id string,
// ) (model.User, error) {
// 	return f.findByIDFunc(ctx, id)
// }

// func (f *fakeUserRepository) FindByEmail(
// 	ctx context.Context,
// 	email string,
// ) (model.User, error) {
// 	return f.findByEmailFunc(ctx, email)
// }

// func (f *fakeUserRepository) List(
// 	ctx context.Context,
// ) ([]model.User, error) {
// 	return f.listFunc(ctx)
// }

// func (f *fakeUserRepository) Update(
// 	ctx context.Context,
// 	user model.User,
// ) (model.User, error) {
// 	return f.updateFunc(ctx, user)
// }

// func (f *fakeUserRepository) Delete(
// 	ctx context.Context,
// 	id string,
// ) error {
// 	return f.deleteFunc(ctx, id)
// }

// func TestServiceCreateSuccess(t *testing.T) {
// 	repository := &fakeUserRepository{
// 		findByEmailFunc: func(
// 			ctx context.Context,
// 			email string,
// 		) (model.User, error) {
// 			return model.User{}, ErrUserNotFound
// 		},
// 		createFunc: func(
// 			ctx context.Context,
// 			user model.User,
// 		) (model.User, error) {
// 			user.ID = "user-001"
// 			return user, nil
// 		},
// 	}

// 	service := New(repository)

// 	user, err := service.Create(
// 		context.Background(),
// 		model.CreateUserRequest{
// 			Name:  "  Panuwat  ",
// 			Email: "PANUWAT@EXAMPLE.COM",
// 		},
// 	)
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}

// 	if user.ID != "user-001" {
// 		t.Errorf("expected ID user-001, got %s", user.ID)
// 	}

// 	if user.Name != "Panuwat" {
// 		t.Errorf("expected name Panuwat, got %s", user.Name)
// 	}

// 	if user.Email != "panuwat@example.com" {
// 		t.Errorf(
// 			"expected email panuwat@example.com, got %s",
// 			user.Email,
// 		)
// 	}

// 	if user.CreatedAt.IsZero() {
// 		t.Error("expected CreatedAt to be set")
// 	}

// 	if user.UpdatedAt.IsZero() {
// 		t.Error("expected UpdatedAt to be set")
// 	}
// }

// func TestServiceCreateValidation(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		request     model.CreateUserRequest
// 		expectedErr error
// 	}{
// 		{
// 			name: "empty name",
// 			request: model.CreateUserRequest{
// 				Name:  " ",
// 				Email: "panuwat@example.com",
// 			},
// 			expectedErr: ErrInvalidUserName,
// 		},
// 		{
// 			name: "invalid email",
// 			request: model.CreateUserRequest{
// 				Name:  "Panuwat",
// 				Email: "invalid-email",
// 			},
// 			expectedErr: ErrInvalidUserEmail,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			repository := &fakeUserRepository{}
// 			service := New(repository)

// 			_, err := service.Create(
// 				context.Background(),
// 				test.request,
// 			)

// 			if !errors.Is(err, test.expectedErr) {
// 				t.Errorf(
// 					"expected error %v, got %v",
// 					test.expectedErr,
// 					err,
// 				)
// 			}
// 		})
// 	}
// }

// func TestServiceCreateDuplicateEmail(t *testing.T) {
// 	repository := &fakeUserRepository{
// 		findByEmailFunc: func(
// 			ctx context.Context,
// 			email string,
// 		) (model.User, error) {
// 			return model.User{
// 				ID:    "existing-user",
// 				Email: email,
// 			}, nil
// 		},
// 	}

// 	service := New(repository)

// 	_, err := service.Create(
// 		context.Background(),
// 		model.CreateUserRequest{
// 			Name:  "Panuwat",
// 			Email: "panuwat@example.com",
// 		},
// 	)

// 	if !errors.Is(err, ErrEmailAlreadyExists) {
// 		t.Errorf(
// 			"expected ErrEmailAlreadyExists, got %v",
// 			err,
// 		)
// 	}
// }

// func TestServiceGetByIDSuccess(t *testing.T) {
// 	repository := &fakeUserRepository{
// 		findByIDFunc: func(
// 			ctx context.Context,
// 			id string,
// 		) (model.User, error) {
// 			return model.User{
// 				ID:    id,
// 				Name:  "Panuwat",
// 				Email: "panuwat@example.com",
// 			}, nil
// 		},
// 	}

// 	service := New(repository)

// 	user, err := service.GetByID(
// 		context.Background(),
// 		"user-001",
// 	)
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}

// 	if user.ID != "user-001" {
// 		t.Errorf("expected ID user-001, got %s", user.ID)
// 	}
// }

// func TestServiceGetByIDEmptyID(t *testing.T) {
// 	repository := &fakeUserRepository{}
// 	service := New(repository)

// 	_, err := service.GetByID(
// 		context.Background(),
// 		" ",
// 	)

// 	if !errors.Is(err, ErrUserNotFound) {
// 		t.Errorf(
// 			"expected ErrUserNotFound, got %v",
// 			err,
// 		)
// 	}
// }
