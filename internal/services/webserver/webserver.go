package webserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
)

type WebServer struct {
	app       *fiber.App
	cfg       WebserverConfig
	container di.Container
}
