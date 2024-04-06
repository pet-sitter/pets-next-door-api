package pet

import (
	"context"
	"database/sql"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"time"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
)

type BasePet struct {
	ID         int          `field:"id"`
	OwnerID    int          `field:"owner_id"`
	Name       string       `field:"name"`
	PetType    PetType      `field:"pet_type"`
	Sex        PetSex       `field:"sex"`
	Neutered   bool         `field:"neutered"`
	Breed      string       `field:"breed"`
	BirthDate  time.Time    `field:"birth_date"`
	WeightInKg float64      `field:"weight_in_kg"`
	Remarks    string       `field:"remarks"`
	CreatedAt  time.Time    `field:"created_at"`
	UpdatedAt  time.Time    `field:"updated_at"`
	DeletedAt  sql.NullTime `field:"deleted_at"`
}

type Pet struct {
	BasePet
	ProfileImageID *int `field:"profile_image_id"`
}

type PetList []*Pet

type PetWithProfileImage struct {
	BasePet
	ProfileImageURL *string `field:"profile_image_url"`
}

type PetWithProfileList []*PetWithProfileImage

type PetStore interface {
	CreatePet(ctx context.Context, tx database.Tx, pet *Pet) (*PetWithProfileImage, *pnd.AppError)
	FindPetByID(ctx context.Context, tx database.Tx, petID int) (*PetWithProfileImage, *pnd.AppError)
	FindPetsByOwnerID(ctx context.Context, tx database.Tx, ownerID int) (*PetWithProfileList, *pnd.AppError)
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
