package main

import (
	"context"
	dbSQL "database/sql"
	"errors"
	"flag"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database/sql"
	"log"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/cmd/import_breeds/breeds_importer_service"
	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"
)

func main() {
	flags := parseFlags()

	log.Printf("Starting to import pet types: %s to database\n", flags.petTypeToImport)

	db, err := sql.OpenSqlDB(configs.DatabaseURL)

	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}

	ctx := context.Background()
	client, err := breeds_importer_service.NewBreedsImporterService(ctx, configs.GoogleSheetsAPIKey)
	if err != nil {
		log.Fatalf("error initializing google sheets client: %v\n", err)
	}

	spreadsheet, err := client.GetSpreadsheet(configs.BreedsGoogleSheetsID)
	if err != nil {
		log.Fatalf("error getting spreadsheet: %v\n", err)
	}

	switch flags.petTypeToImport {
	case Cat:
		var catRows = client.GetCatNames(spreadsheet)
		importBreeds(ctx, &db, pet.PetTypeCat, &catRows)
	case Dog:
		var dogRows = client.GetDogNames(spreadsheet)
		importBreeds(ctx, &db, pet.PetTypeDog, &dogRows)
	case All:
		var catRows = client.GetCatNames(spreadsheet)
		var dogRows = client.GetDogNames(spreadsheet)

		importBreeds(ctx, &db, pet.PetTypeCat, &catRows)
		importBreeds(ctx, &db, pet.PetTypeDog, &dogRows)
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

func importBreed(ctx context.Context, conn *database.DB, petType pet.PetType, row breeds_importer_service.Row) (*pet.Breed, *pnd.AppError) {
	log.Printf("Importing breed with pet_type: %s, name: %s to database", petType, row.Breed)

	var breed *pet.Breed
	err := database.WithTransaction(ctx, conn, func(tx *database.Tx) *pnd.AppError {
		existing, err := postgres.FindBreedByPetTypeAndName(ctx, *tx, petType, row.Breed)
		if err != nil && !errors.Is(err.Err, dbSQL.ErrNoRows) {
			return err
		}

		if existing != nil {
			log.Printf("Breed with id: %d, pet_type: %s, name: %s already exists in database", existing.ID, existing.PetType, existing.Name)
			breed = existing
			return nil
		}

		breed, err = postgres.CreateBreed(ctx, *tx, &pet.Breed{PetType: petType, Name: row.Breed})
		if err != nil {
			return err
		}

		log.Printf("Succeeded to import breed with id: %d, pet_type: %s, name: %s to database", breed.ID, breed.PetType, breed.Name)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return breed, nil
}

func importBreeds(ctx context.Context, conn *database.DB, petType pet.PetType, rows *[]breeds_importer_service.Row) {
	for _, row := range *rows {
		breed, err := importBreed(ctx, conn, petType, row)
		if err != nil {
			log.Printf("Failed to import breed with pet_type: %s, name: %s to database", breed.PetType, breed.Name)
		}
	}
}
