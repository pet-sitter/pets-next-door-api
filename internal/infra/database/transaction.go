package database

import (
	"context"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

type Tx interface {
	Rollback(ctx context.Context) *pnd.AppError
	Commit(ctx context.Context) *pnd.AppError
	EndTx(ctx context.Context, f func() *pnd.AppError) *pnd.AppError
}
