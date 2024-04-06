package database

import (
	"context"
	"database/sql"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

type Tx interface {
	Rollback() *pnd.AppError
	Commit() *pnd.AppError
	EndTx(f func() *pnd.AppError) *pnd.AppError

	ExecContext(context context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(context context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(context context.Context, query string, args ...any) *sql.Row
	PrepareContext(context context.Context, query string) (*sql.Stmt, error)
}
