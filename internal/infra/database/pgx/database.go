package pgx

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

type DB struct {
	Pool        *pgxpool.Pool
	databaseURL string
}

func OpenPgxDB(ctx context.Context, databaseURL string) (*DB, *pnd.AppError) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, pnd.FromPgxError(err)
	}
	return &DB{Pool: pool, databaseURL: databaseURL}, nil
}

func (db *DB) Close() *pnd.AppError {
	db.Pool.Close()
	return nil
}

func (db *DB) Flush(ctx context.Context) *pnd.AppError {
	for _, tableName := range database.TableNames {
		_, err := db.Pool.Exec(ctx, "DELETE FROM "+tableName)
		if err != nil {
			return pnd.FromPgxError(err)
		}
	}

	return nil
}

// BeginTx 는 database.Tx 를 반환한다.
// 사용하지는 않지만 database.DB 인터페이스 구현을 위해 정의되었다.
func (db *DB) BeginTx(ctx context.Context) (database.Tx, *pnd.AppError) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, pnd.FromPgxError(err)
	}
	return &PgxTx{tx: tx}, nil
}

// BeginPgxTx 는 pgx.Tx 를 반환한다.
// database.DB 인터페이스와 호환이 되지 않는다.
func (db *DB) BeginPgxTx(ctx context.Context) (*PgxTx, *pnd.AppError) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, pnd.FromPgxError(err)
	}
	return &PgxTx{tx: tx}, nil
}
