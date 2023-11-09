package controllers

import (
	"time"

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
		return fiber.ErrNotFound // TODO return a not found page later on
	}

	link.LastAccess = time.Now()

	err = c.db.AddUpdateLink(link)
	if err != nil {
		// TODO this shouldn't happen
		return fiber.ErrInternalServerError
	}

	return ctx.Redirect(link.URL)
}
