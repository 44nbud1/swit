// This file contains the repository implementation layer.
package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/lib/pq"
)

type DbTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

type Repository struct {
	Db DbTX
}

func NewRepository(db DbTX) *Repository {
	return &Repository{
		Db: db,
	}
}

func (r *Repository) WithTx(tx pgx.Tx) *Repository {
	return &Repository{
		Db: tx,
	}
}
