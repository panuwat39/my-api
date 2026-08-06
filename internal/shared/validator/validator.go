package validator

import playgroundvalidator "github.com/go-playground/validator/v10"

type Validator struct {
	validate *playgroundvalidator.Validate
}

func New() *Validator {
	return &Validator{
		validate: playgroundvalidator.New(),
	}
}

func (v *Validator) Validate(target any) error {
	return v.validate.Struct(target)
}
