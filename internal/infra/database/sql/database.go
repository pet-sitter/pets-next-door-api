package sql

import (
	"context"
	"database/sql"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"log"
)

type DB struct {
	conn        *sql.DB
	databaseURL string
}

func OpenSqlDB(databaseURL string) (database.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	return &DB{conn: db, databaseURL: databaseURL}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) Flush() error {
	for _, tableName := range database.TableNames {
		_, err := db.conn.Exec("DELETE FROM " + tableName)
		if err != nil {
			log.Fatal(err)
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

func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.conn.ExecContext(ctx, query, args...)
}

func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.conn.QueryContext(ctx, query, args...)
}

func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.conn.QueryRowContext(ctx, query, args...)
}

func (db *DB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return db.conn.PrepareContext(ctx, query)
}
