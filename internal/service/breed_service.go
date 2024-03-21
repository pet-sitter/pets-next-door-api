package service

import (
	"context"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"
)

type BreedService struct {
	conn *database.DB
}

func NewBreedService(conn *database.DB) *BreedService {
	return &BreedService{
		conn: conn,
	}
}

func (s *BreedService) FindBreeds(ctx context.Context, page int, size int, petType *string) (*pet.BreedListView, *pnd.AppError) {
	var breeds *pet.BreedList
	var err *pnd.AppError

	err = database.WithTransaction(ctx, s.conn, func(tx *database.Tx) *pnd.AppError {
		breeds, err = postgres.FindBreeds(ctx, tx, page, size, petType)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return pet.FromBreedList(breeds), nil
}

func (s *BreedService) FindBreedByPetTypeAndName(ctx context.Context, petType pet.PetType, name string) (*pet.BreedView, *pnd.AppError) {
	var breedView *pet.BreedView

	err := database.WithTransaction(ctx, s.conn, func(tx *database.Tx) *pnd.AppError {
		breed, err := postgres.FindBreedByPetTypeAndName(ctx, tx, petType, name)
		if err != nil {
			return err
		}

		breedView = &pet.BreedView{
			ID:      breed.ID,
			PetType: breed.PetType,
			Name:    breed.Name,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return breedView, nil
}
