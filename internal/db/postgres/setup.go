package postgres

import "github.com/jmoiron/sqlx"

func MustConnect(dsn string) *sqlx.DB {
	conn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic("failed to connect to database: " + dsn)
	}

	if conn.Ping() != nil {
		panic("failed to ping database: " + dsn)
	}

	mustMigrate()

	return conn
}

func mustMigrate() {

}
