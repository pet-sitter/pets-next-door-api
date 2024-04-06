package database

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
)

func Migrate(databaseURL, migrationPath string) error {
	m, err := migrate.New(
		"file://"+migrationPath,
		databaseURL,
	)

	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
