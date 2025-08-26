package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	root "http-task-executor/task-service"
)

// MigratePostgresql runs migrations in database.
func MigratePostgresql(db *sqlx.DB) error {

	goose.SetBaseFS(root.MigrationFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.Up(db.DB, "migrations")
}
