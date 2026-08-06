package router

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"

	userroute "github.com/panuwat39/my-api/internal/user/route"
)

type Dependencies struct {
	Logger         *slog.Logger
	UserController userroute.UserController
}

func New(dependencies Dependencies) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "my-api",
		BodyLimit:    1 * 1024 * 1024,
		ErrorHandler: errorHandler(dependencies.Logger),
	})

	registerV1Routes(app)

	if dependencies.UserController != nil {
		userroute.RegisterV1(app, dependencies.UserController)
	}

	return app
}

func errorHandler(logger *slog.Logger) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError

		if fiberError, ok := err.(*fiber.Error); ok {
			code = fiberError.Code
		}

		logger.Error(
			"http request failed",
			"method", c.Method(),
			"path", c.Path(),
			"status", code,
			"error", err,
		)

		return c.Status(code).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "HTTP_ERROR",
				"message": "request failed",
			},
		})
	}
}
