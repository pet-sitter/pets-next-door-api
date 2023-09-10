package main

import (
	"log"

	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

func main() {
	db, err := database.Open(configs.DatabaseURL)
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}

	err = db.Migrate(configs.MigrationPath)
	if err != nil {
		log.Fatalf("error migrating database: %v\n", err)
	}
}
