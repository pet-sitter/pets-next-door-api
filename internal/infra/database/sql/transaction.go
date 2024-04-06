package sql

import (
	"context"
	"database/sql"
	"errors"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

type SqlTx struct {
	*sql.Tx
}

func (sct *SqlTx) EndTx(_ context.Context, f func() *pnd.AppError) *pnd.AppError {
	var err *pnd.AppError
	tx := sct.Tx

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err := f(); err != nil {
			tx.Rollback()
		} else if err := tx.Commit(); err != nil {
			tx.Rollback()
		}
	}()

	err = f()
	return err
}

func (sct *SqlTx) Rollback(_ context.Context) *pnd.AppError {
	if err := sct.Tx.Rollback(); err != nil {
		if errors.Is(err, sql.ErrTxDone) {
			return nil
		}
		return pnd.FromPostgresError(err)
	}

	return nil
}

func (sct *SqlTx) Commit(_ context.Context) *pnd.AppError {
	if err := sct.Tx.Commit(); err != nil {
		return pnd.FromPostgresError(err)
	}

	return nil
}

// SqlTransactioner 인터페이스는 database/sql의 Tx를 구현하는 인터페이스이다.
type SqlTransactioner interface {
	Exec(context context.Context, query string, args ...any) (sql.Result, *pnd.AppError)
	Query(context context.Context, query string, args ...any) (*sql.Rows, *pnd.AppError)
	QueryRow(context context.Context, query string, args ...any) *sql.Row
	Prepare(context context.Context, query string) (*sql.Stmt, *pnd.AppError)
}
