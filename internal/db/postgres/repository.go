package postgres

import (
	"github.com/jmoiron/sqlx"
)

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}

type repository struct {
	db *sqlx.DB
}
