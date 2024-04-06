package main

import (
	"context"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database/pgx"
	"log"

	"github.com/pet-sitter/pets-next-door-api/internal/configs"
)

func main() {
	ctx := context.Background()
	_, err := pgx.OpenPgxDB(ctx, configs.DatabaseURL)
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}

	err2 := database.Migrate(configs.DatabaseURL, configs.MigrationPath)
	if err2 != nil {
		log.Fatalf("error migrating database: %v\n", err2)
	}
}
