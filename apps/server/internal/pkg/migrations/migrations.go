package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
	"github.com/taskemapp/server/apps/server/internal/config"
	"github.com/taskemapp/server/apps/server/internal/pkg/logger"
	"go.uber.org/zap"
)

func Invoke(c config.Config, log logger.Logger) error {
	if err := goose.SetDialect("pgx"); err != nil {
		log.Error("Failed to set dialect: ", zap.Error(err))
		return err
	}
	db, err := sql.Open("pgx", c.PostgresUrl)
	if err != nil {
		log.Error("Failed to open db conn: ", zap.Error(err))
		return err
	}
	defer db.Close()

	log.Info("Run migrations")
	err = goose.Up(db, "migrations")
	if err != nil {
		log.Error("Migration failed: ", zap.Error(err))
		return err
	}

	return nil
}
