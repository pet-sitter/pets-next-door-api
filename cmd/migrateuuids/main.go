package main

import (
	"context"
	"flag"
	"log"

	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

func main() {
	readOnlyPtr := flag.Bool("readonly", false, "Run migration in read-only mode")
	forcePtr := flag.Bool("force", false, "Run migration in force mode")
	flag.Parse()

	var readOnly bool
	if *readOnlyPtr {
		readOnly = true
	}

	var force bool
	if *forcePtr {
		force = true
	}

	log.Printf("Running UUID migration script with readonly: %v, force: %v\n", readOnly, force)

	log.Println("Running UUID migration script")

	db, err := database.Open(configs.DatabaseURL)
	if err != nil {
		panic(err)
	}

	for _, target := range Targets {
		log.Printf("Processing table %s\n", target.Table)
		for _, fk := range target.FKs {
			log.Printf(" - FK %s\n", fk.Column)
		}
	}

	log.Println("Starting migration")

	ctx := context.Background()

	pndErr := migrate(ctx, db, MigrateOptions{ReadOnly: *readOnlyPtr, Force: *forcePtr})
	if pndErr != nil {
		panic(pndErr)
	}

	log.Println("Completed migration")
}
