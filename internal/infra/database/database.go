package database

import (
	"context"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

var TableNames = []string{
	"users",
	"breeds",
	"resource_media",
	"sos_posts_pets",
	"media",
	"pets",
	"sos_posts_conditions",
	"sos_conditions",
	"sos_posts_dates",
	"sos_dates",
	"sos_posts",
	"base_posts",
}

type DB interface {
	Close() *pnd.AppError
	Flush() *pnd.AppError
	BeginTx(ctx context.Context) (Tx, *pnd.AppError)
}

func WithTransaction(ctx context.Context, conn *DB, f func(tx *Tx) *pnd.AppError) *pnd.AppError {
	tx, err := (*conn).BeginTx(ctx)
	if err != nil {
		return err
	}

	if err := f(&tx); err != nil {
		if err := (tx).Rollback(ctx); err != nil {
			return err
		}

		return err
	}

	if err := (tx).Commit(ctx); err != nil {
		return err
	}

	return nil
}
