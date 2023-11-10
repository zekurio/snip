package auth

import (
	"strings"

	"github.com/zekurio/snip/internal/util/static"

	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
)

var (
	errInvalidAccessToken = fiber.NewError(fiber.StatusUnauthorized, "invalid access token")
)

type AccessTokenMiddleware struct {
	ath AccessTokenHandler
}

func NewAccessTokenMiddleware(container di.Container) *AccessTokenMiddleware {
	return &AccessTokenMiddleware{
		ath: container.Get(static.DiAuthAccessTokenHandler).(AccessTokenHandler),
	}
}

func (m *AccessTokenMiddleware) Handle(ctx *fiber.Ctx) (err error) {
	var ident string

	authHeader := ctx.Get("authorization")
	if authHeader == "" {
		return errInvalidAccessToken
	}

	split := strings.Split(authHeader, " ")
	if len(split) < 2 {
		return errInvalidAccessToken
	}

	switch strings.ToLower(split[0]) {

	case "accesstoken":
		if ident, err = m.ath.ValidateAccessToken(split[1]); err != nil || ident == "" {
			return errInvalidAccessToken
		}
	default:
		return fiber.ErrUnauthorized
	}

	return next(ctx, ident)
}

func next(ctx *fiber.Ctx, ident string) error {
	ctx.Locals("uuid", ident)
	return ctx.Next()
}
