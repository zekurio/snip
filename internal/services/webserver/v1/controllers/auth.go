package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/snip/internal/models"
	"github.com/zekurio/snip/internal/services/database"
	"github.com/zekurio/snip/internal/services/webserver/auth"
	"github.com/zekurio/snip/internal/util/static"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	db     database.IDatabase
	rth    auth.RefreshTokenHandler
	ath    auth.AccessTokenHandler
	authMw auth.Middleware
}

func (c *AuthController) Setup(ctn di.Container, router fiber.Router) {
	c.db = ctn.Get(static.DiDatabase).(database.IDatabase)
	c.rth = ctn.Get(static.DiAuthRefreshTokenHandler).(auth.RefreshTokenHandler)
	c.ath = ctn.Get(static.DiAuthAccessTokenHandler).(auth.AccessTokenHandler)
	c.authMw = ctn.Get(static.DiAuthMiddleware).(auth.Middleware)

	router.Post("/login", c.login)
	router.Post("/signup", c.signup)
	router.Get("/logout", c.logout)
}

func (c *AuthController) login(ctx *fiber.Ctx) error {
	user := new(models.User)

	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	existingUser, err := c.db.GetUserByUsername(user.Username)
	if err != nil || existingUser == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Credentials don't match the records",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Credentials don't match the records",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Login successful",
	})
}

func (c *AuthController) signup(ctx *fiber.Ctx) error {
	user := new(models.User)

	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	user.ID = uuid.New().String()

	// check if user already exists by UUID
	existingUser, err := c.db.GetUserByID(user.ID)
	if err == nil && existingUser != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User UUID already exists",
		})
	}

	// check if user already exists by username
	existingUser, err = c.db.GetUserByUsername(user.Username)
	if err == nil && existingUser != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User already exists",
		})
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot hash password",
		})
	}

	user.Password = string(hashBytes)

	if err := c.db.AddUpdateUser(user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot add user",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created",
	})
}

func (c *AuthController) logout(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uuid").(string)

	if err := c.rth.RevokeToken(uid); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot revoke token",
		})
	}

	ctx.ClearCookie(static.RefreshTokenCookieName)

	return ctx.JSON(models.Ok)
}
