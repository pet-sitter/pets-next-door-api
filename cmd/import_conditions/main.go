package main

import (
	"context"
	"log"

	"github.com/pet-sitter/pets-next-door-api/internal/service"

	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

func main() {
	log.Println("Starting to import condition")

	db, err := database.Open(configs.DatabaseURL)
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}

	ctx := context.Background()

	conditionService := service.NewSOSConditionService(db)
	conditionList, err := conditionService.InitConditions(ctx)
	if err != nil {
		log.Fatalf("error initializing conditions: %v\n", err)
	}

	log.Println("Total conditions imported: ", len(conditionList))
	for _, condition := range conditionList {
		log.Println("Condition ID: ", condition.ID, "Condition Name: ", condition.Name)
	}

	log.Println("Finished importing condition")
}
