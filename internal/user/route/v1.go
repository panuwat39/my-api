package route

import (
	"github.com/gofiber/fiber/v3"
)

type UserController interface {
	Create(fiber.Ctx) error
	List(fiber.Ctx) error
	GetByID(fiber.Ctx) error
	Update(fiber.Ctx) error
	Delete(fiber.Ctx) error
}

func RegisterV1(
	app *fiber.App,
	controller UserController,
	authenticate fiber.Handler,
	requireAdmin fiber.Handler,
	requireSelfOrAdmin fiber.Handler,
) {
	users := app.Group("/api/v1/users")

	users.Post(
		"",
		controller.Create,
	)

	users.Get(
		"",
		authenticate,
		requireAdmin,
		controller.List,
	)

	users.Get(
		"/:id",
		authenticate,
		requireSelfOrAdmin,
		controller.GetByID,
	)

	users.Patch(
		"/:id",
		authenticate,
		requireSelfOrAdmin,
		controller.Update,
	)

	users.Delete(
		"/:id",
		authenticate,
		requireAdmin,
		controller.Delete,
	)
}
