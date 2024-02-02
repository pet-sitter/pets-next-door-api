package pet

import pnd "github.com/pet-sitter/pets-next-door-api/api"

type Breed struct {
	ID        int     `field:"id"`
	Name      string  `field:"name"`
	PetType   PetType `field:"pet_type"`
	CreatedAt string  `field:"created_at"`
	UpdatedAt string  `field:"updated_at"`
	DeletedAt string  `field:"deleted_at"`
}

type BreedStore interface {
	FindBreeds(page int, size int, petType *string) ([]*Breed, *pnd.AppError)
	FindBreedByPetTypeAndName(petType PetType, name string) (*Breed, *pnd.AppError)
	CreateBreed(breed *Breed) (*Breed, *pnd.AppError)
}
