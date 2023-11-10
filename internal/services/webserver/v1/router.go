package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/snip/internal/services/webserver/auth"
	"github.com/zekurio/snip/internal/services/webserver/v1/controllers"
	"github.com/zekurio/snip/internal/util/static"
)

type Router struct {
	ctn di.Container
}

func (r *Router) SetContainer(ctn di.Container) {
	r.ctn = ctn
}

func (r *Router) Route(router fiber.Router) {
	authMw := r.ctn.Get(static.DiAuthMiddleware).(auth.Middleware)

	new(controllers.AuthController).Setup(r.ctn, router.Group("/auth"))

	router.Use(authMw.Handle)

	// Authentification required
	new(controllers.LinksController).Setup(r.ctn, router.Group("/links"))
}
