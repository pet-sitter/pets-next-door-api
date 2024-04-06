package main

import (
	"context"
	dbSQL "database/sql"
	"errors"
	"flag"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database/pgx"
	"log"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/cmd/import_breeds/breeds_importer_service"
	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"
)

func main() {
	flags := parseFlags()

	log.Printf("Starting to import pet types: %s to database\n", flags.petTypeToImport)

	ctx := context.Background()
	db, err := pgx.OpenPgxDB(ctx, configs.DatabaseURL)
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}

	client, err2 := breeds_importer_service.NewBreedsImporterService(ctx, configs.GoogleSheetsAPIKey)
	if err2 != nil {
		log.Fatalf("error initializing google sheets client: %v\n", err2)
	}

	spreadsheet, err2 := client.GetSpreadsheet(configs.BreedsGoogleSheetsID)
	if err2 != nil {
		log.Fatalf("error getting spreadsheet: %v\n", err2)
	}

	switch flags.petTypeToImport {
	case Cat:
		var catRows = client.GetCatNames(spreadsheet)
		importBreeds(ctx, db, pet.PetTypeCat, &catRows)
	case Dog:
		var dogRows = client.GetDogNames(spreadsheet)
		importBreeds(ctx, db, pet.PetTypeDog, &dogRows)
	case All:
		var catRows = client.GetCatNames(spreadsheet)
		var dogRows = client.GetDogNames(spreadsheet)

		importBreeds(ctx, db, pet.PetTypeCat, &catRows)
		importBreeds(ctx, db, pet.PetTypeDog, &dogRows)
	}

	log.Println("Completed importing pet types to database")
}

type PetTypeToImport string

const (
	Cat PetTypeToImport = "cat"
	Dog PetTypeToImport = "dog"
	All PetTypeToImport = "all"
)

func fromString(breedToImport string) PetTypeToImport {
	switch breedToImport {
	case "cat":
		return Cat
	case "dog":
		return Dog
	default:
		return All
	}
}

type Flags struct {
	petTypeToImport PetTypeToImport
}

func parseFlags() Flags {
	flag.String("petType", "", "Pet type to import to database")
	flag.Parse()

	petTypeToImportArg := flag.Arg(0)

	petTypeToImport := fromString(petTypeToImportArg)

	return Flags{petTypeToImport: petTypeToImport}
}

func importBreed(ctx context.Context, conn *pgx.DB, petType pet.PetType, row breeds_importer_service.Row) (*pet.Breed, *pnd.AppError) {
	log.Printf("Importing breed with pet_type: %s, name: %s to database", petType, row.Breed)

	var breed *pet.Breed

	tx, err := conn.BeginPgxTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	existing, err := postgres.FindBreedByPetTypeAndName(ctx, tx, petType, row.Breed)
	if err != nil && !errors.Is(err.Err, dbSQL.ErrNoRows) {
		return nil, err
	}

	if existing != nil {
		log.Printf("Breed with id: %d, pet_type: %s, name: %s already exists in database", existing.ID, existing.PetType, existing.Name)
		breed = existing
	}

	breed, err = postgres.CreateBreed(ctx, tx, &pet.Breed{PetType: petType, Name: row.Breed})
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	log.Printf("Succeeded to import breed with id: %d, pet_type: %s, name: %s to database", breed.ID, breed.PetType, breed.Name)

	return breed, nil
}

func importBreeds(ctx context.Context, conn *pgx.DB, petType pet.PetType, rows *[]breeds_importer_service.Row) {
	for _, row := range *rows {
		breed, err := importBreed(ctx, conn, petType, row)
		if err != nil {
			log.Printf("Failed to import breed with pet_type: %s, name: %s to database", breed.PetType, breed.Name)
		}
	}
}
