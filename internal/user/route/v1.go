package route

import "github.com/gofiber/fiber/v3"

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
) {
	users := app.Group("/api/v1/users")

	users.Post("", controller.Create)
	users.Get("", controller.List)
	users.Get("/:id", controller.GetByID)
	users.Patch("/:id", controller.Update)
	users.Delete("/:id", controller.Delete)
}
