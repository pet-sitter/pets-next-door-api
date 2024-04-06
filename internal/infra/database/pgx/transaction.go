package pgx

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

type PgxTx struct {
	tx pgx.Tx
}

func (ptx *PgxTx) EndTx(ctx context.Context, f func() *pnd.AppError) *pnd.AppError {
	var err *pnd.AppError
	tx := ptx.tx

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		} else if err := f(); err != nil {
			tx.Rollback(ctx)
		} else if err := tx.Commit(ctx); err != nil {
			tx.Rollback(ctx)
		}
	}()

	err = f()
	return err
}

func (ptx *PgxTx) Rollback(ctx context.Context) *pnd.AppError {
	if err := ptx.tx.Rollback(ctx); err != nil {
		if errors.Is(err, pgx.ErrTxClosed) {
			return nil
		}
		return pnd.FromPgxError(err)
	}

	return nil
}

func (ptx *PgxTx) Commit(ctx context.Context) *pnd.AppError {
	if err := ptx.tx.Commit(ctx); err != nil {
		return pnd.FromPgxError(err)
	}

	return nil
}

func (ptx *PgxTx) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, *pnd.AppError) {
	ct, err := ptx.tx.Exec(ctx, query, args...)
	if err != nil {
		return pgconn.CommandTag{}, pnd.FromPgxError(err)
	}

	return ct, nil
}

func (ptx *PgxTx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, *pnd.AppError) {
	rows, err := ptx.tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, pnd.FromPgxError(err)
	}

	return rows, nil
}

func (ptx *PgxTx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return ptx.tx.QueryRow(ctx, sql, args...)
}

func (ptx *PgxTx) Prepare(ctx context.Context, name string, sql string) (*pgconn.StatementDescription, *pnd.AppError) {
	ps, err := ptx.tx.Prepare(ctx, name, sql)
	if err != nil {
		return nil, pnd.FromPgxError(err)
	}

	return ps, nil
}
