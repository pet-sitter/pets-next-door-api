package database

type DBManager interface {
	Migrate() error
}
