package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/snip/internal/services/database"
	"github.com/zekurio/snip/internal/services/util/static"
)

type RedirectController struct {
	db database.IDatabase
}

func (c *RedirectController) Setup(ctn di.Container, router fiber.Router) {
	c.db = ctn.Get(static.DiDatabase).(database.IDatabase)

	router.Get("/:id", c.redirect)
}

func (c *RedirectController) redirect(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	link, err := c.db.GetLinkByID(id)
	if err != nil {
		// TODO return an error page, for now just return a 404
		return fiber.ErrNotFound
	}

	// TODO update link to set the new last access time

	return ctx.Redirect(link.URL)
}
