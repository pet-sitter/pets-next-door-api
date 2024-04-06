package sql

import (
	"context"
	"database/sql"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

type DB struct {
	conn        *sql.DB
	databaseURL string
}

func OpenSqlDB(databaseURL string) (database.DB, *pnd.AppError) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	return &DB{conn: db, databaseURL: databaseURL}, nil
}

func (db *DB) Close() *pnd.AppError {
	err := db.conn.Close()
	if err != nil {
		return pnd.FromPostgresError(err)
	}

	return nil
}

func (db *DB) Flush() *pnd.AppError {
	for _, tableName := range database.TableNames {
		_, err := db.conn.Exec("DELETE FROM " + tableName)
		if err != nil {
			return pnd.FromPostgresError(err)
		}
	}

	return nil
}

func (db *DB) BeginTx(ctx context.Context) (database.Tx, *pnd.AppError) {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return &SqlTx{tx}, nil
}

func (db *DB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, *pnd.AppError) {
	result, err := db.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return result, nil
}

func (db *DB) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, *pnd.AppError) {
	rows, err := db.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return rows, nil
}

func (db *DB) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.conn.QueryRowContext(ctx, query, args...)
}

func (db *DB) Prepare(ctx context.Context, query string) (*sql.Stmt, *pnd.AppError) {
	stmt, err := db.conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return stmt, nil
}
