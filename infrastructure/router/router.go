package router

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	sharedmiddleware "github.com/panuwat39/my-api/internal/shared/middleware"
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

	app.Use(recoverer.New(recoverer.Config{
		EnableStackTrace: true,
		PanicHandler: func(c fiber.Ctx, recovered any) error {
			dependencies.Logger.Error(
				"panic recovered",
				"request_id", requestid.FromContext(c),
				"method", c.Method(),
				"path", c.Path(),
				"panic", recovered,
			)

			return fiber.ErrInternalServerError
		},
	}))

	app.Use(requestid.New())
	app.Use(sharedmiddleware.Logging(dependencies.Logger))

	registerV1Routes(app)

	if dependencies.UserController != nil {
		userroute.RegisterV1(
			app,
			dependencies.UserController,
		)
	}

	return app
}

func errorHandler(logger *slog.Logger) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		statusCode := fiber.StatusInternalServerError
		errorCode := "INTERNAL_SERVER_ERROR"
		message := "internal server error"

		if fiberError, ok := err.(*fiber.Error); ok {
			statusCode = fiberError.Code

			switch statusCode {
			case fiber.StatusNotFound:
				errorCode = "ROUTE_NOT_FOUND"
				message = "route not found"

			case fiber.StatusMethodNotAllowed:
				errorCode = "METHOD_NOT_ALLOWED"
				message = "method not allowed"

			case fiber.StatusBadRequest:
				errorCode = "BAD_REQUEST"
				message = "bad request"
			}
		}

		logger.Error(
			"http error handled",
			"request_id", requestid.FromContext(c),
			"method", c.Method(),
			"path", c.Path(),
			"status", statusCode,
			"error", err,
		)

		return c.Status(statusCode).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    errorCode,
				"message": message,
			},
		})
	}
}
