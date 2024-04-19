package tests

import "os"

var TestDatabaseURL = os.Getenv("TEST_DATABASE_URL")

//nolint:gochecknoinits
func init() {
	if TestDatabaseURL == "" {
		TestDatabaseURL = "postgresql://postgres:postgres@localhost:5455/pets_next_door_api_test?sslmode=disable"
	}
}
