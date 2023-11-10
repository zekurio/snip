package auth

import (
	"time"

	"github.com/zekurio/snip/internal/models"
	"github.com/zekurio/snip/internal/util/static"
	"github.com/zekurio/snip/pkg/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
)

type LoginHandler struct {
	accessTokenHandler  AccessTokenHandler
	refreshTokenHandler RefreshTokenHandler
}

func NewLoginHandler(container di.Container) *LoginHandler {
	return &LoginHandler{
		accessTokenHandler:  container.Get(static.DiAuthAccessTokenHandler).(AccessTokenHandler),
		refreshTokenHandler: container.Get(static.DiAuthRefreshTokenHandler).(RefreshTokenHandler),
	}
}

func (h *LoginHandler) LoginFailedHandler(ctx *fiber.Ctx, status int, msg string) error {
	return fiber.NewError(status, msg)
}

func (h *LoginHandler) BindRefreshToken(ctx *fiber.Ctx, uuid string) error {
	ctx.Locals("uuid", uuid)

	refreshToken, err := h.refreshTokenHandler.GetRefreshToken(uuid)
	if err != nil {
		return err
	}

	expires := time.Now().Add(static.AuthSessionExpiration)
	ctx.Cookie(&fiber.Cookie{
		Name:     static.RefreshTokenCookieName,
		Value:    refreshToken,
		Path:     "/",
		Expires:  expires,
		HTTPOnly: true,
		Secure:   !debug.Enabled(),
	})

	return nil
}

func (h *LoginHandler) LoginSuccessHandler(ctx *fiber.Ctx, res *models.User) error {
	if err := h.BindRefreshToken(ctx, res.ID); err != nil {
		return err
	}

	return nil
}

func (h *LoginHandler) LogoutHandler(ctx *fiber.Ctx) error {
	if uid, ok := ctx.Locals("uuid").(string); ok && uid != "" {
		if err := h.refreshTokenHandler.RevokeToken(uid); err != nil {
			return err
		}
	}

	ctx.ClearCookie(static.RefreshTokenCookieName)

	return ctx.JSON(models.Ok)
}
