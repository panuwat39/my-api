package controller

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"

	"github.com/panuwat39/my-api/internal/shared/fiberx"
	"github.com/panuwat39/my-api/internal/user/model"
	"github.com/panuwat39/my-api/internal/user/port"
)

type Controller struct {
	service port.UserService
	logger  *slog.Logger
}

func New(
	service port.UserService,
	logger *slog.Logger,
) *Controller {
	return &Controller{
		service: service,
		logger:  logger,
	}
}

func (c *Controller) Create(ctx fiber.Ctx) error {
	var request model.CreateUserRequest

	if err := fiberx.BindBody(ctx, &request); err != nil {
		c.logger.Warn(
			"invalid create user request body",
			"method", ctx.Method(),
			"path", ctx.Path(),
			"error", err,
		)

		return fiberx.Error(
			ctx,
			fiber.StatusBadRequest,
			"INVALID_REQUEST_BODY",
			"request body is invalid",
		)
	}

	user, err := c.service.Create(
		ctx,
		request,
	)
	if err != nil {
		return c.writeServiceError(ctx, err)
	}

	return fiberx.Success(
		ctx,
		fiber.StatusCreated,
		model.ToUserResponse(user),
	)
}

func (c *Controller) List(ctx fiber.Ctx) error {
	users, err := c.service.List(ctx)
	if err != nil {
		return c.writeServiceError(ctx, err)
	}

	result := make(
		[]model.UserResponse,
		0,
		len(users),
	)

	for _, user := range users {
		result = append(
			result,
			model.ToUserResponse(user),
		)
	}

	return fiberx.Success(
		ctx,
		fiber.StatusOK,
		result,
	)
}

func (c *Controller) GetByID(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	user, err := c.service.GetByID(
		ctx,
		id,
	)
	if err != nil {
		return c.writeServiceError(ctx, err)
	}

	return fiberx.Success(
		ctx,
		fiber.StatusOK,
		model.ToUserResponse(user),
	)
}

func (c *Controller) Update(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	var request model.UpdateUserRequest

	if err := fiberx.BindBody(ctx, &request); err != nil {
		c.logger.Warn(
			"invalid update user request body",
			"method", ctx.Method(),
			"path", ctx.Path(),
			"error", err,
		)

		return fiberx.Error(
			ctx,
			fiber.StatusBadRequest,
			"INVALID_REQUEST_BODY",
			"request body is invalid",
		)
	}

	user, err := c.service.Update(
		ctx,
		id,
		request,
	)
	if err != nil {
		return c.writeServiceError(ctx, err)
	}

	return fiberx.Success(
		ctx,
		fiber.StatusOK,
		model.ToUserResponse(user),
	)
}

func (c *Controller) Delete(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.service.Delete(ctx, id); err != nil {
		return c.writeServiceError(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
