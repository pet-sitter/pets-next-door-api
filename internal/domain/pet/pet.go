package pet

import (
	"context"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

type BasePet struct {
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

func NewBasePet(id int, ownerID int, name string, petType PetType, sex PetSex, neutered bool, breed string, birthDate string, weightInKg float64, createdAt string, updatedAt string, deletedAt string) *BasePet {
	return &BasePet{
		ID:         id,
		OwnerID:    ownerID,
		Name:       name,
		PetType:    petType,
		Sex:        sex,
		Neutered:   neutered,
		Breed:      breed,
		BirthDate:  birthDate,
		WeightInKg: weightInKg,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
		DeletedAt:  deletedAt,
	}
}

type Pet struct {
	BasePet
	ProfileImageID *int `field:"profile_image_id"`
}

func NewPet(id int, ownerID int, name string, petType PetType, sex PetSex, neutered bool, breed string, birthDate string, weightInKg float64, createAt string, updatedAt string, deletedAt string, profileImageID *int) *Pet {
	return &Pet{
		BasePet:        *NewBasePet(id, ownerID, name, petType, sex, neutered, breed, birthDate, weightInKg, createAt, updatedAt, deletedAt),
		ProfileImageID: profileImageID,
	}
}

type PetWithProfileImage struct {
	BasePet
	ProfileImageURL *string `field:"profile_image_url"`
}

func NewPetWithProfileImage(id int, ownerID int, name string, petType PetType, sex PetSex, neutered bool, breed string, birthDate string, weightInKg float64, createAt string, updatedAt string, deletedAt string, profileImageURL *string) *PetWithProfileImage {
	return &PetWithProfileImage{
		BasePet:         *NewBasePet(id, ownerID, name, petType, sex, neutered, breed, birthDate, weightInKg, createAt, updatedAt, deletedAt),
		ProfileImageURL: profileImageURL,
	}
}

type PetStore interface {
	CreatePet(ctx context.Context, pet *Pet) (*PetWithProfileImage, *pnd.AppError)
	FindPetByID(ctx context.Context, petID int) (*PetWithProfileImage, *pnd.AppError)
	FindPetsByOwnerID(ctx context.Context, ownerID int) ([]PetWithProfileImage, *pnd.AppError)
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
