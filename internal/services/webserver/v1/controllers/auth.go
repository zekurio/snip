package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/snip/internal/models"
	"github.com/zekurio/snip/internal/services/database"
)

type AuthController struct {
	db database.IDatabase
}

func (c *AuthController) Setup(ctn di.Container, router fiber.Router) {
	router.Post("/login", c.login)
	router.Post("/signup", c.signup)
	router.Get("/logout", c.logout)
}

func (c *AuthController) login(ctx *fiber.Ctx) error {
	return nil
}

func (c *AuthController) signup(ctx *fiber.Ctx) error {
	user := new(models.User)

	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Check if the user already exists
	return nil
}

func (c *AuthController) logout(ctx *fiber.Ctx) error {
	return nil
}
