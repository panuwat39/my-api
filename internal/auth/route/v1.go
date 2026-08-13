package route

import (
	"github.com/gofiber/fiber/v3"
)

type AuthController interface {
	Login(fiber.Ctx) error
	Me(fiber.Ctx) error
}

func RegisterV1(
	app *fiber.App,
	controller AuthController,
	authenticate fiber.Handler,
) {
	auth := app.Group("/api/v1/auth")

	auth.Post(
		"/login",
		controller.Login,
	)

	auth.Get(
		"/me",
		authenticate,
		controller.Me,
	)
}
