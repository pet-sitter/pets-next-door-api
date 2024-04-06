package pgx

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

type PgxTx struct {
	tx  pgx.Tx
	ctx context.Context
}

func (ptx *PgxTx) EndTx(f func() *pnd.AppError) *pnd.AppError {
	var err *pnd.AppError
	tx := ptx.tx

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ptx.ctx)
			panic(p)
		} else if err := f(); err != nil {
			tx.Rollback(ptx.ctx)
		} else if err := tx.Commit(ptx.ctx); err != nil {
			tx.Rollback(ptx.ctx)
		}
	}()

	err = f()
	return err
}

func (ptx *PgxTx) Rollback() *pnd.AppError {
	if err := ptx.tx.Rollback(ptx.ctx); err != nil {
		if errors.Is(err, pgx.ErrTxClosed) {
			return nil
		}
		return pnd.FromPgxError(err)
	}

	return nil
}

func (ptx *PgxTx) Commit() *pnd.AppError {
	if err := ptx.tx.Commit(ptx.ctx); err != nil {
		return pnd.FromPgxError(err)
	}

	return nil
}

// TODO: return the actual result
func (ptx *PgxTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	_, err := ptx.tx.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// TODO: return the actual result
func (ptx *PgxTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	_, err := ptx.tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// TODO: return the actual result
func (ptx *PgxTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	ptx.tx.QueryRow(ctx, query, args...)
	return nil
}

// TODO: return the actual result
func (ptx *PgxTx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	_, err := ptx.tx.Prepare(ctx, "query", query)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
