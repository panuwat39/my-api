package service

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidUserName     = errors.New("user name is required")
	ErrInvalidUserEmail    = errors.New("user email is invalid")
	ErrInvalidUserPassword = errors.New("user password is invalid")
)
