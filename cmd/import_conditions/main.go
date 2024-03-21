package main

import (
	"context"
	"log"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sos_post"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"
)

func main() {
	log.Println("Starting to import condition")

	db, err := database.Open(configs.DatabaseURL)
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}

	var result string
	var err2 *pnd.AppError

	ctx := context.Background()
	err2 = database.WithTransaction(ctx, db, func(tx *database.Tx) *pnd.AppError {
		result, err2 = postgres.InitConditions(ctx, tx, sos_post.ConditionName)
		if err2 != nil {
			return err2
		}

		return nil
	})

	if err2 != nil {
		log.Fatalf("error initializing condition: %v\n", err2)
	}

	log.Println(result)
}
