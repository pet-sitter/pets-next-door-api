package main

import (
	"context"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database/pgx"
	"log"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sos_post"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"
)

func main() {
	log.Println("Starting to import condition")

	ctx := context.Background()
	db, err := pgx.OpenPgxDB(ctx, configs.DatabaseURL)
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}

	var result string
	var err2 *pnd.AppError

	tx, err := db.BeginPgxTx(ctx)
	if err != nil {
		log.Fatalf("error beginning transaction: %v\n", err)
	}

	result, err2 = postgres.InitConditions(ctx, tx, sos_post.ConditionName)
	if err2 != nil {
		log.Fatalf("error initializing condition: %v\n", err2)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("error committing transaction: %v\n", err)
	}

	log.Println(result)
}
