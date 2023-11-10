package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/zekurio/snip/internal/util/static"
	static2 "github.com/zekurio/snip/internal/util/static"
	"github.com/zekurio/snip/pkg/debug"

	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekurio/snip/internal/inits"
	"github.com/zekurio/snip/internal/models"
	"github.com/zekurio/snip/internal/services/config"
	"github.com/zekurio/snip/internal/services/database"
	"github.com/zekurio/snip/internal/services/webserver/auth"
)

var (
	flagConfigPath = flag.String("c", "config.toml", "Path to config file")
	flagDebug      = flag.Bool("debug", false, "Enable debug mode")
)

func main() {
	// Parse command line flags
	flag.Parse()

	// Set debug mode
	debug.SetEnabled(*flagDebug)

	diBuilder, err := di.NewBuilder()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create DI builder")
	}

	// Config dependency
	diBuilder.Add(di.Def{
		Name: static2.DiConfig,
		Build: func(ctn di.Container) (interface{}, error) {
			return config.Parse(*flagConfigPath, "SNIP_", models.DefaultConfig)
		},
	})

	// Database and cache dependency
	diBuilder.Add(di.Def{
		Name: static2.DiDatabase,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitDatabase(ctn)
		},
		Close: func(obj interface{}) error {
			d := obj.(database.IDatabase)
			logrus.Info("Closing database connection")
			return d.Close()
		},
	})

	// Auth refresh token dependency
	diBuilder.Add(di.Def{
		Name: static.DiAuthRefreshTokenHandler,
		Build: func(ctn di.Container) (interface{}, error) {
			return auth.NewRefreshTokenHandlerImpl(ctn), nil
		},
	})

	// Auth access token dependency
	diBuilder.Add(di.Def{
		Name: static.DiAuthAccessTokenHandler,
		Build: func(ctn di.Container) (interface{}, error) {
			return auth.NewAccessTokenHandlerImpl(ctn), nil
		},
	})

	// Login Handler dependency
	diBuilder.Add(di.Def{
		Name: static.DiAuthLoginHandler,
		Build: func(ctn di.Container) (interface{}, error) {
			return auth.NewLoginHandler(ctn), nil
		},
	})

	// Auth middleware dependency
	diBuilder.Add(di.Def{
		Name: static.DiAuthMiddleware,
		Build: func(ctn di.Container) (interface{}, error) {
			return auth.NewAccessTokenMiddleware(ctn), nil
		},
	})

	// Webserver dependency
	diBuilder.Add(di.Def{
		Name: static2.DiWebserver,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitWebserver(ctn), nil
		},
	})

	// Build dependency injection container
	ctn := diBuilder.Build()

	// Tear down dependency instances
	defer func(ctn di.Container) {
		err := ctn.DeleteWithSubContainers()
		if err != nil {
			logrus.WithError(err).Fatal("Failed to tear down dependency instances")
		}
	}(ctn)

	ctn.Get(static2.DiWebserver)

	// Block main go routine until one of the following
	// specified exit sys calls occure.
	logrus.Info("Started event loop. Stop with CTRL-C...")

	logrus.Info("Initialization finished")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
