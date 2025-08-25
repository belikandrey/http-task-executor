package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	root "task-service"
)

func MigratePostgresql(db *sqlx.DB) error {

	goose.SetBaseFS(root.MigrationFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.Up(db.DB, "migrations")
}
