package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/snip/internal/models"
	"github.com/zekurio/snip/internal/services/database"
	"github.com/zekurio/snip/internal/services/util/static"
	"github.com/zekurio/snip/internal/services/webserver/auth"
)

type UsersController struct {
	db          database.IDatabase
	userHandler *auth.UserHandler
}

func (c *UsersController) Setup(ctn di.Container, router fiber.Router) {
	c.db = ctn.Get(static.DiDatabase).(database.IDatabase)
	c.userHandler = ctn.Get(static.DiUserHandler).(*auth.UserHandler)

	router.Post("/", c.createUser)
	router.Get("/", c.getUser)
	router.Delete("/", c.deleteUser)
}

func (c *UsersController) createUser(ctx *fiber.Ctx) error {
	user := new(models.User)

	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	err := c.userHandler.RegisterUser(user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot create user, database error",
		})
	}

	return err
}

func (c *UsersController) getUser(ctx *fiber.Ctx) error {
	panic("not implemented")
}

func (c *UsersController) deleteUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.db.DeleteUser(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot delete user",
		})
	}

	// Return success response
	return ctx.SendString("User deleted successfully")
}
