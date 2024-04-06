package main

import (
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database/sql"
	"log"

	"github.com/pet-sitter/pets-next-door-api/internal/configs"
)

func main() {
	_, err := sql.OpenSqlDB(configs.DatabaseURL)
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}

	err = database.Migrate(configs.DatabaseURL, configs.MigrationPath)
	if err != nil {
		log.Fatalf("error migrating database: %v\n", err)
	}
}
