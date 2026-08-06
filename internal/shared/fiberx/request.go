package fiberx

import (
	"errors"

	"github.com/gofiber/fiber/v3"
)

var ErrInvalidRequestBody = errors.New("invalid request body")

func BindBody(
	c fiber.Ctx,
	target any,
) error {
	if err := c.Bind().Body(target); err != nil {
		return errors.Join(
			ErrInvalidRequestBody,
			err,
		)
	}

	return nil
}
