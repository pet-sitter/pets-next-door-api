package database

import "database/sql"

type DB struct {
	DB          *sql.DB
	databaseURL string
}

type Tx struct {
	*sql.Tx
}
