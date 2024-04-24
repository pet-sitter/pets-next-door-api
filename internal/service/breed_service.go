package service

import (
	"context"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/breed"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/commonvo"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
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

func (s *BreedService) FindBreeds(
	ctx context.Context, page, size int, petType *string,
) (*breed.BreedListView, *pnd.AppError) {
	tx, err := s.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	breeds, err := postgres.FindBreeds(ctx, tx, page, size, petType)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return breeds.ToBreedListView(), nil
}

func (s *BreedService) FindBreedByPetTypeAndName(
	ctx context.Context, petType commonvo.PetType, name string,
) (*breed.BreedView, *pnd.AppError) {
	tx, err := s.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	breedData, err := postgres.FindBreedByPetTypeAndName(ctx, tx, petType, name)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &breed.BreedView{
		ID:      breedData.ID,
		PetType: breedData.PetType,
		Name:    breedData.Name,
	}, nil
}
