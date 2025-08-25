package postgres

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"task-executor/internal/config"
	"time"
)

const (
	maxOpenConnections    = 30
	connectionMaxLifetime = 100 * time.Second
	maxIdleConnections    = 10
	connectionMaxIdleTime = 10 * time.Second
)

func NewPostgresqlDatabase(c *config.Config) (*sqlx.DB, error) {
	dbUrl := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.Name,
		c.Postgres.SslMode,
		c.Postgres.Password)

	db, err := sqlx.Open(c.Postgres.Driver, dbUrl)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConnections)
	db.SetConnMaxLifetime(connectionMaxLifetime)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetConnMaxIdleTime(connectionMaxIdleTime)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
