package postgres

import (
	"Rsbot_only/internal/config"
	"Rsbot_only/pkg/clientDB/postgresLocal"
	"context"
	"github.com/mentalisit/logger"
)

type Db struct {
	db    postgresLocal.Client
	log   *logger.Logger
	debug bool
}

func NewDb(log *logger.Logger, cfg *config.ConfigBot) *Db {
	db, err := postgresLocal.NewClient(context.Background(), log, 5, cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &Db{
		db:    db,
		log:   log,
		debug: cfg.IsDebug,
	}
}
