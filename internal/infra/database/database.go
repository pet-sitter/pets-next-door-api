package database

import (
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate/v4"

	// To use Postgres driver
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// To use file source
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type DB struct {
	DB          *sql.DB
	databaseURL string
}

type Tx struct {
	*sql.Tx
}

func Open(databaseURL string) (*DB, error) {
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
	tableNames := []string{
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
		_, err := db.DB.Exec("DELETE FROM " + tableName) //nolint:gosec
		if err != nil {
			return err
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

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
