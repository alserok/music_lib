package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"

	_ "github.com/lib/pq"
)

func MustConnect(dsn string, dir ...string) *sqlx.DB {
	conn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic("failed to connect to database: " + dsn)
	}

	if conn.Ping() != nil {
		panic("failed to ping database: " + dsn)
	}

	mustMigrate(conn, dir...)

	return conn
}

const (
	migrationsDir = "./internal/db/migrations"
)

func mustMigrate(conn *sqlx.DB, dir ...string) {
	if err := goose.SetDialect("postgres"); err != nil {
		panic("failed to set dialect: " + err.Error())
	}

	if len(dir) == 0 {
		if err := goose.Up(conn.DB, migrationsDir); err != nil {
			panic("failed to migrate: " + err.Error())
		}
	} else {
		if err := goose.Up(conn.DB, dir[0]); err != nil {
			panic("failed to migrate: " + err.Error())
		}
	}
}
