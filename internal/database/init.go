package database

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

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
	var tableNames []string = []string{
		"users",
		"media",
		"pets",
	}

	for _, tableName := range tableNames {
		_, err := db.DB.Exec("DELETE FROM " + tableName)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func (db *DB) Migrate(migrationPath string) {

	m, err := migrate.New(
		"file://"+migrationPath,
		db.databaseURL,
	)

	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}

func (db *DB) Begin() (*Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}
