package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"
)

type Transactioner interface {
	Rollback() error
	Commit() error
	EndTx(f func() error) error
	BeginTx() (*DB, error)
}

func (db *DB) BeginTx(ctx context.Context) (*Tx, error) {
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &Tx{tx}, nil
}

func (tx *Tx) EndTx(f func() error) error {
	var err error
	defer func() {
		if p := recover(); p != nil {
			if err2 := tx.Tx.Rollback(); err2 != nil {
				log.Error().Err(err2).Msg("error rolling back transaction")
				panic(p)
			}
		} else if err = f(); err != nil {
			if err2 := tx.Tx.Rollback(); err2 != nil {
				log.Error().Err(err2).Msg("error rolling back transaction")
			}
		} else if err := tx.Tx.Commit(); err != nil {
			if err2 := tx.Tx.Rollback(); err2 != nil {
				log.Error().Err(err2).Msg("error rolling back transaction")
			}
		}
	}()

	err = f()
	return err
}

func (tx *Tx) Rollback() error {
	if err := tx.Tx.Rollback(); err != nil {
		if errors.Is(err, sql.ErrTxDone) {
			return nil
		}
		return err
	}

	return nil
}

func (tx *Tx) Commit() error {
	return tx.Tx.Commit()
}

func WithTransaction(ctx context.Context, conn *DB, f func(tx *Tx) error) error {
	tx, err := conn.BeginTx(ctx)
	if err != nil {
		return err
	}

	if err := f(tx); err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	return tx.Commit()
}

func (tx *Tx) ExecContext(
	ctx context.Context,
	query string,
	args ...interface{},
) (sql.Result, error) {
	return tx.Tx.ExecContext(ctx, query, args...)
}

func (tx *Tx) QueryContext(
	ctx context.Context,
	query string,
	args ...interface{},
) (*sql.Rows, error) {
	return tx.Tx.QueryContext(ctx, query, args...)
}

func (tx *Tx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return tx.Tx.QueryRowContext(ctx, query, args...)
}

func (tx *Tx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return tx.Tx.PrepareContext(ctx, query)
}
