package pet

import (
	"context"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

type Pet struct {
	ID         int     `field:"id"`
	OwnerID    int     `field:"owner_id"`
	Name       string  `field:"name"`
	PetType    PetType `field:"pet_type"`
	Sex        PetSex  `field:"sex"`
	Neutered   bool    `field:"neutered"`
	Breed      string  `field:"breed"`
	BirthDate  string  `field:"birth_date"`
	WeightInKg float64 `field:"weight_in_kg"`
	CreatedAt  string  `field:"created_at"`
	UpdatedAt  string  `field:"updated_at"`
	DeletedAt  string  `field:"deleted_at"`
}

type PetStore interface {
	CreatePet(ctx context.Context, pet *Pet) (*Pet, *pnd.AppError)
	FindPetsByOwnerID(ctx context.Context, ownerID int) ([]Pet, *pnd.AppError)
}

type PetType string

const (
	PetTypeDog PetType = "dog"
	PetTypeCat PetType = "cat"
)

type PetSex string

const (
	PetSexMale   PetSex = "male"
	PetSexFemale PetSex = "female"
)
