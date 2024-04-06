package database

import (
	"context"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

type DB interface {
	Close() error
	Flush() error
	Migrate(migrationPath string) error
	BeginTx(ctx context.Context) (Tx, *pnd.AppError)
}
