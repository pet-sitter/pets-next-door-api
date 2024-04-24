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

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type WithProfileImage struct {
	ID              int64
	OwnerID         int64
	Name            string
	PetType         commonvo.PetType
	Sex             Gender
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
		Sex:             Gender(row.Sex),
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
		Sex:             Gender(row.Sex),
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
		Sex:             Gender(row.Sex),
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
