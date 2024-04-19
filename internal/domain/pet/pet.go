package pet

import (
	"context"

	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

type BasePet struct {
	ID         int     `field:"id" json:"id"`
	OwnerID    int     `field:"owner_id" json:"ownerId"`
	Name       string  `field:"name" json:"name"`
	PetType    PetType `field:"pet_type" json:"petType"`
	Sex        PetSex  `field:"sex" json:"sex"`
	Neutered   bool    `field:"neutered" json:"neutered"`
	Breed      string  `field:"breed" json:"breed"`
	BirthDate  string  `field:"birth_date" json:"birthDate"`
	WeightInKg float64 `field:"weight_in_kg" json:"weightInKg"`
	Remarks    string  `field:"remarks" json:"remarks"`
	CreatedAt  string  `field:"created_at" json:"createdAt"`
	UpdatedAt  string  `field:"updated_at" json:"updatedAt"`
	DeletedAt  string  `field:"deleted_at" json:"deletedAt"`
}

type Pet struct {
	BasePet
	ProfileImageID *int `field:"profile_image_id"`
}

type PetList []*Pet

type PetWithProfileImage struct {
	BasePet
	ProfileImageURL *string `field:"profile_image_url" json:"profileImageUrl"`
}

type PetWithProfileList []*PetWithProfileImage

type PetStore interface {
	CreatePet(ctx context.Context, tx *database.Tx, pet *Pet) (*PetWithProfileImage, *pnd.AppError)
	FindPetByID(ctx context.Context, tx *database.Tx, petID int) (*PetWithProfileImage, *pnd.AppError)
	FindPetsByOwnerID(ctx context.Context, tx *database.Tx, ownerID int) (*PetWithProfileList, *pnd.AppError)
	UpdatePet(ctx context.Context, tx *database.Tx, updatePetRequest *UpdatePetRequest) *pnd.AppError
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
