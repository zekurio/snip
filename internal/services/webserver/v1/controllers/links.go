package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/snip/internal/models"
	"github.com/zekurio/snip/internal/services/database"
	"github.com/zekurio/snip/internal/util/static"
	"github.com/zekurio/snip/pkg/randutils"
)

type LinksController struct {
	db database.IDatabase
}

func (c *LinksController) Setup(ctn di.Container, router fiber.Router) {
	c.db = ctn.Get(static.DiDatabase).(database.IDatabase)

	router.Get("/", c.getLinks)
	router.Post("/create", c.postCreateLink)
	router.Get("/delete/:id", c.getDeleteLink)
}

func (c *LinksController) getLinks(ctx *fiber.Ctx) error {
	uuid := ctx.Locals("uuid").(string)

	links, err := c.db.GetLinksByUser(uuid)
	if err != nil {
		return ctx.JSON(models.Status{
			Code:    fiber.StatusInternalServerError,
			Message: "failed to get links",
		})
	}

	return ctx.JSON(models.NewListResponse(links))
}

func (c *LinksController) postCreateLink(ctx *fiber.Ctx) error {
	uuid := ctx.Locals("uuid").(string)
	link := new(models.Link)

	var body struct {
		URL string `json:"url"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.JSON(models.Status{
			Code:    fiber.StatusBadRequest,
			Message: "cannot parse json",
		})
	}

	link.ID = randutils.ForceRandBase64Str(8)

	for {
		_, err := c.db.GetLinkByID(link.ID)
		if err != nil {
			break
		}
		link.ID = randutils.ForceRandBase64Str(8)
	}

	link.URL = body.URL
	link.OwnerID = uuid
	link.CreatedAt = time.Now()

	err := c.db.AddUpdateLink(link)
	if err != nil {
		return ctx.JSON(models.Status{
			Code:    fiber.StatusInternalServerError,
			Message: "failed to create link",
		})
	}

	return ctx.JSON(models.Ok)
}

func (c *LinksController) getDeleteLink(ctx *fiber.Ctx) error {
	uuid := ctx.Locals("uuid").(string)
	id := ctx.Params("id")

	link, err := c.db.GetLinkByID(id)
	if err != nil {
		return ctx.JSON(models.Status{
			Code:    fiber.StatusInternalServerError,
			Message: "failed to get link",
		})
	}

	if link.OwnerID != uuid {
		return ctx.JSON(models.Status{
			Code:    fiber.StatusUnauthorized,
			Message: "you are not the owner of this link",
		})
	}

	err = c.db.DeleteLink(id)
	if err != nil {
		return ctx.JSON(models.Status{
			Code:    fiber.StatusInternalServerError,
			Message: "failed to delete link",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(models.Ok)
}
