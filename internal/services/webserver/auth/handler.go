package auth

import (
	"github.com/gofiber/fiber/v2"
)

// AuthenticationHandler provides fiber endpoints and handlers
// to authenticate users via a username and password.
type AuthenticationHandler interface {

	// LoginFailedHandler is called when either the login
	// or the refresh token validation failed.
	LoginFailedHandler(ctx *fiber.Ctx, status int, msg string) error

	// BindRefreshToken generates a refresh token for the
	// given user ident and binds it to the given context.
	BindRefreshToken(ctx *fiber.Ctx, uid string) error

	// LoginSuccessHandler is called when the user
	// successfully logged in.
	LoginSuccessHandler(ctx *fiber.Ctx, username string, password string) error

	// LogoutHandler is called when the user wants
	// to log out.
	LogoutHandler(ctx *fiber.Ctx) error
}
