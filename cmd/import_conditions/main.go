package main

import (
	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sos_post"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"
	"log"
)

func main() {
	log.Println("Starting to import condition")

	db, err := database.Open(configs.DatabaseURL)

	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}

	conditionStore := postgres.NewConditionPostgresStore(db)

	result, err := conditionStore.InitConditions(sos_post.ConditionName)

	if err != nil {
		log.Fatalf("error initializing condition: %v\n", err)
	}

	log.Println(result)
}
