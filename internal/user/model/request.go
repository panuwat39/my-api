package model

type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`

	Email string `json:"email" validate:"required,email,max=254"`

	Password string `json:"password" validate:"required,min=8,max=72"`
}

type UpdateUserRequest struct {
	Name *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`

	Email *string `json:"email,omitempty" validate:"omitempty,email,max=254"`
}
