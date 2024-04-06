package sql

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"log"
)

type DB struct {
	DB          *sql.DB
	databaseURL string
}

func OpenSqlDB(databaseURL string) (database.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db, databaseURL: databaseURL}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

func (db *DB) Flush() error {
	var tableNames = []string{
		"users",
		"breeds",
		"resource_media",
		"sos_posts_pets",
		"media",
		"pets",
		"sos_posts_conditions",
		"sos_conditions",
		"sos_posts_dates",
		"sos_dates",
		"sos_posts",
		"base_posts",
	}

	for _, tableName := range tableNames {
		_, err := db.DB.Exec("DELETE FROM " + tableName)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func (db *DB) Migrate(migrationPath string) error {
	m, err := migrate.New(
		"file://"+migrationPath,
		db.databaseURL,
	)

	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (db *DB) BeginTx(ctx context.Context) (database.Tx, *pnd.AppError) {
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return &SqlTx{tx}, nil
}
