package service

import (
	"context"

	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/breed"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
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
	ctx context.Context, params *breed.FindBreedsParams,
) (*breed.ListView, error) {
	rows, err := databasegen.New(s.conn).FindBreeds(ctx, params.ToDBParams())
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return breed.ToListViewFromRows(params.Page, params.Size, rows), nil
}
