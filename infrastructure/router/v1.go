package router

import "github.com/gofiber/fiber/v3"

func registerV1Routes(app *fiber.App) {
	v1 := app.Group("/api/v1")

	v1.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data": fiber.Map{
				"status":  "ok",
				"version": "v1",
			},
		})
	})
}
