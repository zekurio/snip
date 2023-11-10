package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekurio/snip/internal/models"
	"github.com/zekurio/snip/internal/services/database"
	"github.com/zekurio/snip/internal/services/webserver/auth"
	"github.com/zekurio/snip/internal/util/static"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	db     database.IDatabase
	lgh    *auth.LoginHandler
	rth    auth.RefreshTokenHandler
	ath    auth.AccessTokenHandler
	authMw auth.Middleware
}

func (c *AuthController) Setup(ctn di.Container, router fiber.Router) {
	c.db = ctn.Get(static.DiDatabase).(database.IDatabase)
	c.lgh = ctn.Get(static.DiAuthLoginHandler).(*auth.LoginHandler)
	c.rth = ctn.Get(static.DiAuthRefreshTokenHandler).(auth.RefreshTokenHandler)
	c.ath = ctn.Get(static.DiAuthAccessTokenHandler).(auth.AccessTokenHandler)
	c.authMw = ctn.Get(static.DiAuthMiddleware).(auth.Middleware)

	router.Post("/login", c.postLogin)
	router.Post("/signup", c.postSignup)
	router.Post("/accesstoken", c.postAccessToken)
	router.Post("/logout", c.authMw.Handle, c.postLogout)
}

func (c *AuthController) postLogin(ctx *fiber.Ctx) error {
	var redirect string

	login := new(models.UserLoginDetails)

	if err := ctx.BodyParser(login); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if redirect = ctx.Query("redirect"); redirect == "" {
		redirect = "/"
	}

	user, err := c.db.GetUserByUsername(login.Username)
	if err != nil || user == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Credentials don't match the records",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Credentials don't match the records",
		})
	}

	if err := c.lgh.LoginSuccessHandler(ctx, user); err != nil {
		return err
	}

	return ctx.Redirect(redirect)
}

func (c *AuthController) postSignup(ctx *fiber.Ctx) error {
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

	return ctx.Status(fiber.StatusCreated).JSON(models.Ok)
}

func (c *AuthController) postAccessToken(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies(static.RefreshTokenCookieName)
	if refreshToken == "" {
		return fiber.ErrUnauthorized
	}

	ident, err := c.rth.ValidateRefreshToken(refreshToken)
	if err != nil {
		logrus.Error("Failed validating refresh token", err)
	}
	if ident == "" {
		return fiber.ErrUnauthorized
	}

	token, expires, err := c.ath.GetAccessToken(ident)
	if err != nil {
		return err
	}

	return ctx.JSON(&models.AccessTokenResponse{
		Token:   token,
		Expires: expires,
	})
}

func (c *AuthController) postLogout(ctx *fiber.Ctx) error {
	id := ctx.Locals("uuid").(string)

	if err := c.rth.RevokeToken(id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot revoke token",
		})
	}

	ctx.ClearCookie(static.RefreshTokenCookieName)

	return ctx.JSON(models.Ok)
}
