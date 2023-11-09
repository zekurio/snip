package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/snip/internal/services/database"
	"github.com/zekurio/snip/internal/services/util/static"
)

type LinksController struct {
	db database.IDatabase
}

func (c *LinksController) Setup(ctn di.Container, router fiber.Router) {
	c.db = ctn.Get(static.DiDatabase).(database.IDatabase)

	router.Get("/", c.getLinks)
}

func (c *LinksController) getLinks(ctx *fiber.Ctx) error {
	// TODO get current user, and get all links for that user
	return ctx.SendString("TODO")
}
