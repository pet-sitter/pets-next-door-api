package pet

import (
	"database/sql"
	"time"

	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/datatype"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/commonvo"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
	"github.com/shopspring/decimal"
)

type BasePet struct {
	ID         int              `field:"id" json:"id"`
	OwnerID    int              `field:"owner_id" json:"owner_id"`
	Name       string           `field:"name" json:"name"`
	PetType    commonvo.PetType `field:"pet_type" json:"pet_type"`
	Sex        PetSex           `field:"sex" json:"sex"`
	Neutered   bool             `field:"neutered" json:"neutered"`
	Breed      string           `field:"breed" json:"breed"`
	BirthDate  datatype.Date    `field:"birth_date" json:"birth_date"`
	WeightInKg decimal.Decimal  `field:"weight_in_kg" json:"weight_in_kg"`
	Remarks    string           `field:"remarks" json:"remarks"`
	CreatedAt  string           `field:"created_at" json:"created_at"`
	UpdatedAt  string           `field:"updated_at" json:"updated_at"`
	DeletedAt  string           `field:"deleted_at" json:"deleted_at"`
}

type Pet struct {
	BasePet
	ProfileImageID *int `field:"profile_image_id"`
}

type PetList []*Pet

type PetWithProfileImage struct {
	BasePet
	ProfileImageURL *string `field:"profile_image_url" json:"profile_image_url"`
}

type PetWithProfileList []*PetWithProfileImage

type PetSex string

const (
	PetSexMale   PetSex = "male"
	PetSexFemale PetSex = "female"
)

type WithProfileImage struct {
	ID              int64
	OwnerID         int64
	Name            string
	PetType         commonvo.PetType
	Sex             PetSex
	Neutered        bool
	Breed           string
	BirthDate       datatype.Date
	WeightInKg      decimal.Decimal
	Remarks         string
	ProfileImageURL *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       sql.NullTime
}

func ToWithProfileImage(row databasegen.FindPetRow) *WithProfileImage {
	weightInKg, _ := decimal.NewFromString(row.WeightInKg)
	birthDate := datatype.DateOf(row.BirthDate)

	return &WithProfileImage{
		ID:              int64(row.ID),
		OwnerID:         row.OwnerID,
		Name:            row.Name,
		PetType:         commonvo.PetType(row.PetType),
		Sex:             PetSex(row.Sex),
		Neutered:        row.Neutered,
		Breed:           row.Breed,
		BirthDate:       birthDate,
		WeightInKg:      weightInKg,
		Remarks:         row.Remarks,
		ProfileImageURL: utils.NullStrToStrPtr(row.ProfileImageUrl),
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
		DeletedAt:       row.DeletedAt,
	}
}

func ToWithProfileImageFromRows(row databasegen.FindPetsRow) *WithProfileImage {
	weightInKg, _ := decimal.NewFromString(row.WeightInKg)
	birthDate := datatype.DateOf(row.BirthDate)

	return &WithProfileImage{
		ID:              int64(row.ID),
		OwnerID:         row.OwnerID,
		Name:            row.Name,
		PetType:         commonvo.PetType(row.PetType),
		Sex:             PetSex(row.Sex),
		Neutered:        row.Neutered,
		Breed:           row.Breed,
		BirthDate:       birthDate,
		WeightInKg:      weightInKg,
		Remarks:         row.Remarks,
		ProfileImageURL: utils.NullStrToStrPtr(row.ProfileImageUrl),
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
		DeletedAt:       row.DeletedAt,
	}
}

func ToWithProfileImageFromIDsRows(row databasegen.FindPetsByIDsRow) *WithProfileImage {
	weightInKg, _ := decimal.NewFromString(row.WeightInKg)
	birthDate := datatype.DateOf(row.BirthDate)

	return &WithProfileImage{
		ID:              int64(row.ID),
		OwnerID:         row.OwnerID,
		Name:            row.Name,
		PetType:         commonvo.PetType(row.PetType),
		Sex:             PetSex(row.Sex),
		Neutered:        row.Neutered,
		Breed:           row.Breed,
		BirthDate:       birthDate,
		WeightInKg:      weightInKg,
		Remarks:         row.Remarks,
		ProfileImageURL: utils.NullStrToStrPtr(row.ProfileImageUrl),
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
		DeletedAt:       row.DeletedAt,
	}
}
