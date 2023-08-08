package models

import "github.com/pet-sitter/pets-next-door-api/internal/types"

type Pet struct {
	ID         int           `field:"id"`
	OwnerID    int           `field:"owner_id"`
	Name       string        `field:"name"`
	PetType    types.PetType `field:"pet_type"`
	Sex        types.PetSex  `field:"sex"`
	Neutered   bool          `field:"neutered"`
	Breed      string        `field:"breed"`
	BirthDate  string        `field:"birth_date"`
	WeightInKg float64       `field:"weight_in_kg"`
	CreatedAt  string        `field:"created_at"`
	UpdatedAt  string        `field:"updated_at"`
	DeletedAt  string        `field:"deleted_at"`
}
