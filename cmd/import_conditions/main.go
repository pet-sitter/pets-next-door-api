package main

import (
	"context"
	"log"

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

	conditionStore := postgres.NewConditionPostgresStore(db)

	ctx := context.Background()
	result, err2 := conditionStore.InitConditions(ctx, sos_post.ConditionName)
	if err2 != nil {
		log.Fatalf("error initializing condition: %v\n", err2)
	}

	log.Println(result)
}
