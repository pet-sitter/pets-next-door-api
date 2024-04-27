package tests

import (
	"testing"

	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

func SetUp(t *testing.T) (*database.DB, func(t *testing.T)) {
	t.Helper()

	db, _ := database.Open(TestDatabaseURL)
	db.Flush()

	return db, func(t *testing.T) {
		t.Helper()

		db.Close()
	}
}

func CreateForEach(setUp, tearDown func()) func(func()) {
	return func(test func()) {
		setUp()
		test()
		tearDown()
	}
}
