package controller

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"github.com/panuwat39/my-api/internal/shared/fiberx"
	userservice "github.com/panuwat39/my-api/internal/user/service"
)

func (c *Controller) writeServiceError(
	ctx fiber.Ctx,
	err error,
) error {
	switch {
	case errors.Is(
		err,
		userservice.ErrInvalidUserName,
	):
		return fiberx.Error(
			ctx,
			fiber.StatusBadRequest,
			"INVALID_USER_NAME",
			"user name is required",
		)

	case errors.Is(
		err,
		userservice.ErrInvalidUserEmail,
	):
		return fiberx.Error(
			ctx,
			fiber.StatusBadRequest,
			"INVALID_USER_EMAIL",
			"user email is invalid",
		)

	case errors.Is(
		err,
		userservice.ErrEmailAlreadyExists,
	):
		return fiberx.Error(
			ctx,
			fiber.StatusConflict,
			"EMAIL_ALREADY_EXISTS",
			"email already exists",
		)

	case errors.Is(
		err,
		userservice.ErrUserNotFound,
	):
		return fiberx.Error(
			ctx,
			fiber.StatusNotFound,
			"USER_NOT_FOUND",
			"user not found",
		)

	case errors.Is(err, userservice.ErrInvalidUserPassword):
		return fiberx.Error(
			ctx,
			fiber.StatusBadRequest,
			"INVALID_USER_PASSWORD",
			"password must be between 8 and 72 bytes",
		)

	default:
		c.logger.Error(
			"user request failed",
			"method", ctx.Method(),
			"path", ctx.Path(),
			"error", err,
		)

		return fiberx.Error(
			ctx,
			fiber.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			"internal server error",
		)
	}
}
