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

func (s *BreedService) FindBreeds(
	ctx context.Context, page int, size int, petType *string,
) (*pet.BreedListView, *pnd.AppError) {
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
	ctx context.Context, petType pet.PetType, name string,
) (*pet.BreedView, *pnd.AppError) {
	tx, err := s.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	breed, err := postgres.FindBreedByPetTypeAndName(ctx, tx, petType, name)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &pet.BreedView{
		ID:      breed.ID,
		PetType: breed.PetType,
		Name:    breed.Name,
	}, nil
}
