package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func Logging(logger *slog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		startedAt := time.Now()

		err := c.Next()

		statusCode := c.Response().StatusCode()

		if err != nil {
			statusCode = statusCodeFromError(err)
		}

		attributes := []any{
			"request_id", requestid.FromContext(c),
			"method", c.Method(),
			"path", c.Path(),
			"status", statusCode,
			"duration_ms", time.Since(startedAt).Milliseconds(),
			"ip", c.IP(),
		}

		if err != nil {
			attributes = append(attributes, "error", err)

			logger.Error(
				"http request failed",
				attributes...,
			)

			return err
		}

		logger.Info(
			"http request completed",
			attributes...,
		)

		return nil
	}
}

func statusCodeFromError(err error) int {
	if fiberError, ok := err.(*fiber.Error); ok {
		return fiberError.Code
	}

	return fiber.StatusInternalServerError
}
