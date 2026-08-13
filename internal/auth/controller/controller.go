package controller

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v3"

	authmodel "github.com/panuwat39/my-api/internal/auth/model"
	authport "github.com/panuwat39/my-api/internal/auth/port"
	authservice "github.com/panuwat39/my-api/internal/auth/service"
	"github.com/panuwat39/my-api/internal/shared/fiberx"

	authmiddleware "github.com/panuwat39/my-api/internal/auth/middleware"
)

type Controller struct {
	service authport.AuthService
	logger  *slog.Logger
}

func New(
	service authport.AuthService,
	logger *slog.Logger,
) *Controller {
	return &Controller{
		service: service,
		logger:  logger,
	}
}

func (c *Controller) Login(ctx fiber.Ctx) error {
	var request authmodel.LoginRequest

	if err := fiberx.BindBody(ctx, &request); err != nil {
		return fiberx.Error(
			ctx,
			fiber.StatusBadRequest,
			"INVALID_REQUEST_BODY",
			"request body is invalid",
		)
	}

	result, err := c.service.Login(ctx, request)
	if err != nil {
		switch {
		case errors.Is(
			err,
			authservice.ErrInvalidCredentials,
		):
			return fiberx.Error(
				ctx,
				fiber.StatusUnauthorized,
				"INVALID_CREDENTIALS",
				"email or password is invalid",
			)

		default:
			c.logger.Error(
				"login failed",
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

	return fiberx.Success(
		ctx,
		fiber.StatusOK,
		result,
	)
}

func (c *Controller) Me(
	ctx fiber.Ctx,
) error {
	user, ok := authmiddleware.CurrentUser(ctx)
	if !ok {
		return fiberx.Error(
			ctx,
			fiber.StatusUnauthorized,
			"UNAUTHORIZED",
			"authentication is required",
		)
	}

	return fiberx.Success(
		ctx,
		fiber.StatusOK,
		fiber.Map{
			"user_id": user.UserID,
			"email":   user.Email,
		},
	)
}
