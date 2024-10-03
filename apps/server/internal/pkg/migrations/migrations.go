package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
	"github.com/taskemapp/server/apps/server/internal/config"
	"go.uber.org/zap"
)

func Invoke(c config.Config, log *zap.Logger) error {
	if err := goose.SetDialect("pgx"); err != nil {
		log.Sugar().Error("Failed to set dialect: ", err)
		return err
	}
	db, err := sql.Open("pgx", c.PostgresUrl)
	if err != nil {
		log.Sugar().Error("Failed to open db conn: ", err)
		return err
	}
	defer db.Close()

	log.Sugar().Info("Run migrations")
	err = goose.Up(db, "migrations")
	if err != nil {
		log.Sugar().Error("Migration failed: ", err)
		return err
	}

	return nil
}
