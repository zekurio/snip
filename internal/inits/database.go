package inits

import (
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/snip/internal/models"
	"github.com/zekurio/snip/internal/services/database"
	"github.com/zekurio/snip/internal/services/database/postgres"
	"github.com/zekurio/snip/internal/util/static"
)

func InitDatabase(ctn di.Container) (db database.IDatabase, err error) {
	cfg := ctn.Get(static.DiConfig).(models.Config)

	db, err = postgres.NewPostgres(cfg.Postgres)
	if err != nil {
		return nil, err
	}

	db, err = database.WrapCache(db, err)
	if err != nil {
		return nil, err
	}

	return db, nil
}
