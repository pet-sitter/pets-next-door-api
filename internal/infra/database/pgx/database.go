package pgx

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

type DB struct {
	pool        *pgxpool.Pool
	databaseURL string
}

func OpenPgxDB(ctx context.Context, databaseURL string) (*DB, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}
	return &DB{pool: pool, databaseURL: databaseURL}, nil
}

func (db *DB) Close() error {
	db.pool.Close()
	return nil
}

func (db *DB) Flush(ctx context.Context) error {
	for _, tableName := range database.TableNames {
		_, err := db.pool.Exec(ctx, "DELETE FROM "+tableName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) BeginTx(ctx context.Context) (database.Tx, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &PgxTx{tx: tx, ctx: ctx}, nil
}
