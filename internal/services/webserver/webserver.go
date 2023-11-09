package webserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekurio/snip/internal/models"
	"github.com/zekurio/snip/internal/services/util/static"
	v1 "github.com/zekurio/snip/internal/services/webserver/v1"
)

type WebServer struct {
	app       *fiber.App
	cfg       models.WebserverConfig
	container di.Container
}

func New(ctn di.Container) (ws *WebServer, err error) {
	ws = new(WebServer)
	ws.container = ctn

	ws.cfg = ctn.Get(static.DiConfig).(models.Config).Webserver

	ws.app = fiber.New(fiber.Config{
		AppName:               "snip",
		DisableStartupMessage: true,
		ProxyHeader:           "X-Forwarded-For",
	})

	ws.registerRouter(new(v1.Router), []string{"/api/v1", "/api"})

	return
}

func (ws *WebServer) registerRouter(router *v1.Router, routes []string, middlewares ...fiber.Handler) {
	router.SetContainer(ws.container)
	for _, r := range routes {
		router.Route(ws.app.Group(r, middlewares...))
	}
}

func (ws *WebServer) ListenAndServeBlocking() error {
	tls := ws.cfg.TLS

	if tls.Enabled {
		logrus.Infof("Starting webserver on %s with TLS", ws.cfg.Addr)
		return ws.app.ListenTLS(ws.cfg.Addr, tls.Cert, tls.Key)
	}

	logrus.Infof("Starting webserver on %s", ws.cfg.Addr)
	return ws.app.Listen(ws.cfg.Addr)
}
