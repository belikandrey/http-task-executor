package postgres

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"http-task-executor/task-executor/internal/task-executor/config"
)

// NewPostgresqlDatabase creates new database instance.
func NewPostgresqlDatabase(c *config.Config) (*sqlx.DB, error) {
	dbURL := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.Name,
		c.Postgres.SslMode,
		c.Postgres.Password)

	db, err := sqlx.Open(c.Postgres.Driver, dbURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(int(c.Postgres.MaxOpenConnections))
	db.SetConnMaxLifetime(c.Postgres.ConnectionMaxLifetime)
	db.SetMaxIdleConns(int(c.Postgres.MaxIdleConnections))
	db.SetConnMaxIdleTime(c.Postgres.ConnectionMaxIdleTime)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
