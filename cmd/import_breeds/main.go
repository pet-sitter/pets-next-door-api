package main

import (
	"context"
	"flag"
	"log"

	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/breed"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/commonvo"

	"github.com/pet-sitter/pets-next-door-api/cmd/import_breeds/breedsimporterservice"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

func main() {
	flags := parseFlags()

	log.Printf("Starting to import pet types: %s to database\n", flags.petTypeToImport)

	db, err := database.Open(configs.DatabaseURL)
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}

	ctx := context.Background()
	client, err := breedsimporterservice.NewBreedsImporterService(ctx, configs.GoogleSheetsAPIKey)
	if err != nil {
		log.Fatalf("error initializing google sheets client: %v\n", err)
	}

	spreadsheet, err := client.GetSpreadsheet(configs.BreedsGoogleSheetsID)
	if err != nil {
		log.Fatalf("error getting spreadsheet: %v\n", err)
	}

	switch flags.petTypeToImport {
	case Cat:
		catRows := client.GetCatNames(spreadsheet)
		importBreeds(ctx, db, commonvo.PetTypeCat, &catRows)
	case Dog:
		dogRows := client.GetDogNames(spreadsheet)
		importBreeds(ctx, db, commonvo.PetTypeDog, &dogRows)
	case All:
		catRows := client.GetCatNames(spreadsheet)
		dogRows := client.GetDogNames(spreadsheet)

		importBreeds(ctx, db, commonvo.PetTypeCat, &catRows)
		importBreeds(ctx, db, commonvo.PetTypeDog, &dogRows)
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

func importBreed(
	ctx context.Context, conn *database.DB, petType commonvo.PetType, row breedsimporterservice.Row,
) (*breed.DetailView, *pnd.AppError) {
	log.Printf("Importing breed with pet_type: %s, name: %s to database", petType, row.Breed)

	existingList, err := databasegen.New(conn).FindBreeds(ctx, databasegen.FindBreedsParams{
		PetType: utils.StrToNullStr(petType.String()),
		Name:    utils.StrToNullStr(row.Breed),
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	if len(existingList) > 1 {
		existing := existingList[0]
		log.Printf(
			"Breed with id: %d, pet_type: %s, name: %s already exists in database",
			existing.ID,
			existing.PetType,
			existing.Name,
		)
		return breed.ToDetailViewFromRows(existing), nil
	}

	breedData, err := databasegen.New(conn).CreateBreed(ctx, databasegen.CreateBreedParams{
		Name:    row.Breed,
		PetType: petType.String(),
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	log.Printf(
		"Succeeded to import breed with id: %d, pet_type: %s, name: %s to database",
		breedData.ID,
		breedData.PetType,
		breedData.Name,
	)

	return &breed.DetailView{
		ID:      int64(breedData.ID),
		PetType: commonvo.PetType(breedData.PetType),
		Name:    breedData.Name,
	}, nil
}

func importBreeds(ctx context.Context, conn *database.DB, petType commonvo.PetType, rows *[]breedsimporterservice.Row) {
	for _, row := range *rows {
		breedData, err := importBreed(ctx, conn, petType, row)
		if err != nil {
			log.Printf("Failed to import breed with pet_type: %s, name: %s to database", breedData.PetType, breedData.Name)
		}
	}
}
