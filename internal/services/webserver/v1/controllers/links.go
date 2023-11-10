package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/snip/internal/services/database"
	"github.com/zekurio/snip/internal/util/static"
)

type LinksController struct {
	db database.IDatabase
}

func (c *LinksController) Setup(ctn di.Container, router fiber.Router) {
	c.db = ctn.Get(static.DiDatabase).(database.IDatabase)

	router.Get("/", c.getLinks)
}

func (c *LinksController) getLinks(ctx *fiber.Ctx) error {
	uuid := ctx.Locals("uuid").(string)

	// query all links
	links, err := c.db.GetLinksByUser(uuid)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get links",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    links,
	})
}
