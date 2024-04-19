package pet

import (
	"context"

	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

type Breed struct {
	ID        int     `field:"id"`
	Name      string  `field:"name"`
	PetType   PetType `field:"pet_type"`
	CreatedAt string  `field:"created_at"`
	UpdatedAt string  `field:"updated_at"`
	DeletedAt string  `field:"deleted_at"`
}

type BreedList struct {
	*pnd.PaginatedView[Breed]
}

func NewBreedList(page, size int) *BreedList {
	return &BreedList{PaginatedView: pnd.NewPaginatedView(
		page, size, false, make([]Breed, 0),
	)}
}

type BreedStore interface {
	FindBreeds(ctx context.Context, tx *database.Tx, page, size int, petType *string) (*BreedList, *pnd.AppError)
	FindBreedByPetTypeAndName(ctx context.Context, tx *database.Tx, petType PetType, name string) (*Breed, *pnd.AppError)
	CreateBreed(ctx context.Context, tx *database.Tx, breed *Breed) (*Breed, *pnd.AppError)
}
