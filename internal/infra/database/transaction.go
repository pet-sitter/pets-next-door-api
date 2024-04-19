package database

import (
	"context"
	"database/sql"
	"errors"

	pnd "github.com/pet-sitter/pets-next-door-api/api"

	"github.com/rs/zerolog/log"
)

type DBTx interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

type Transactioner interface {
	Rollback() *pnd.AppError
	Commit() *pnd.AppError
	EndTx(f func() *pnd.AppError) *pnd.AppError
	BeginTx() (*DB, *pnd.AppError)
}

func (db *DB) BeginTx(ctx context.Context) (*Tx, *pnd.AppError) {
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return &Tx{tx}, nil
}

func (sct *Tx) EndTx(f func() *pnd.AppError) *pnd.AppError {
	var err *pnd.AppError
	tx := sct.Tx

	defer func() {
		if p := recover(); p != nil {
			if err2 := tx.Rollback(); err2 != nil {
				log.Error().Err(err2).Msg("error rolling back transaction")
				panic(p)
			}
		} else if err = f(); err != nil {
			if err2 := tx.Rollback(); err2 != nil {
				log.Error().Err(err2).Msg("error rolling back transaction")
			}
		} else if err := tx.Commit(); err != nil {
			if err2 := tx.Rollback(); err2 != nil {
				log.Error().Err(err2).Msg("error rolling back transaction")
			}
		}
	}()

	err = f()
	return err
}

func (sct *Tx) Rollback() *pnd.AppError {
	if err := sct.Tx.Rollback(); err != nil {
		if errors.Is(err, sql.ErrTxDone) {
			return nil
		}
		return pnd.FromPostgresError(err)
	}

	return nil
}

func (sct *Tx) Commit() *pnd.AppError {
	if err := sct.Tx.Commit(); err != nil {
		return pnd.FromPostgresError(err)
	}

	return nil
}

func WithTransaction(ctx context.Context, conn *DB, f func(tx *Tx) *pnd.AppError) *pnd.AppError {
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
