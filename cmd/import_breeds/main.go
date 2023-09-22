package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"github.com/pet-sitter/pets-next-door-api/cmd/import_breeds/breeds_importer_service"
	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"
	"log"
)

func main() {
	flags := parseFlags()

	log.Printf("Starting to import pet types: %s to database\n", flags.petTypeToImport)

	db, err := database.Open(configs.DatabaseURL)

	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}

	breedStore := postgres.NewBreedPostgresStore(db)

	ctx := context.Background()
	client, err := breeds_importer_service.NewBreedsImporterService(ctx, configs.GoogleSheetsAPIKey)

	if err != nil {
		log.Fatalf("error initializing google sheets client: %v\n", err)
	}

	spreadsheet, err := client.GetSpreadsheet(configs.BreedsSheetsID)

	switch flags.petTypeToImport {
	case Cat:
		var catRows = client.GetCatNames(spreadsheet)
		importBreeds(breedStore, pet.PetTypeCat, &catRows)
		break
	case Dog:
		var dogRows = client.GetDogNames(spreadsheet)
		importBreeds(breedStore, pet.PetTypeDog, &dogRows)
		break
	case All:
		var catRows = client.GetCatNames(spreadsheet)
		var dogRows = client.GetDogNames(spreadsheet)

		importBreeds(breedStore, pet.PetTypeCat, &catRows)
		importBreeds(breedStore, pet.PetTypeDog, &dogRows)
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

func importBreed(breedStore pet.BreedStore, petType pet.PetType, row breeds_importer_service.Row) (*pet.Breed, error) {
	log.Printf("Importing breed with pet_type: %s, name: %s to database", petType, row.Breed)

	existing, err := breedStore.FindBreedByPetTypeAndName(petType, row.Breed)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if existing != nil {
		log.Printf("Breed with id: %d, pet_type: %s, name: %s already exists in database", existing.ID, existing.PetType, existing.Name)
		return existing, nil
	}

	breed, err := breedStore.CreateBreed(&pet.Breed{
		PetType: petType,
		Name:    row.Breed,
	})

	if err != nil {
		return breed, err
	}

	log.Printf("Succeeded to import breed with id: %d, pet_type: %s, name: %s to database", breed.ID, breed.PetType, breed.Name)
	return breed, nil
}

func importBreeds(breedStore pet.BreedStore, petType pet.PetType, rows *[]breeds_importer_service.Row) {
	for _, row := range *rows {
		breed, err := importBreed(breedStore, petType, row)
		if err != nil {
			log.Printf("Failed to import breed with pet_type: %s, name: %s to database", breed.PetType, breed.Name)
		}
	}
}
