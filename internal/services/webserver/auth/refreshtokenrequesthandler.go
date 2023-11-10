package auth

import (
	"strings"
	"time"

	"github.com/zekurio/snip/internal/models"
	"github.com/zekurio/snip/internal/util/static"
	"github.com/zekurio/snip/pkg/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
)

type RefreshTokenRequestHandler struct {
	accessTokenHandler  AccessTokenHandler
	refreshTokenHandler RefreshTokenHandler
}

func NewRefreshTokenRequestHandler(container di.Container) *RefreshTokenRequestHandler {
	return &RefreshTokenRequestHandler{
		accessTokenHandler:  container.Get(static.DiAuthAccessTokenHandler).(AccessTokenHandler),
		refreshTokenHandler: container.Get(static.DiAuthRefreshTokenHandler).(RefreshTokenHandler),
	}
}

func (h *RefreshTokenRequestHandler) LoginFailedHandler(ctx *fiber.Ctx, status int, msg string) error {
	return fiber.NewError(status, msg)
}

func (h *RefreshTokenRequestHandler) BindRefreshToken(ctx *fiber.Ctx, uuid string) error {
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

func (h *RefreshTokenRequestHandler) LoginSuccessHandler(ctx *fiber.Ctx, res models.LoginSuccessResponse) error {
	if err := h.BindRefreshToken(ctx, res.User.ID); err != nil {
		return err
	}

	location := "/"
	if res.Redirect != "" {
		location += strings.TrimLeft(res.Redirect, "/")
	}

	return ctx.Redirect(location, fiber.StatusTemporaryRedirect)
}

func (h *RefreshTokenRequestHandler) LogoutHandler(ctx *fiber.Ctx) error {
	if uid, ok := ctx.Locals("uuid").(string); ok && uid != "" {
		if err := h.refreshTokenHandler.RevokeToken(uid); err != nil {
			return err
		}
	}

	ctx.ClearCookie(static.RefreshTokenCookieName)

	return ctx.JSON(models.Ok)
}
