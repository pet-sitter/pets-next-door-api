package pet

import (
	"github.com/google/uuid"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/commonvo"
	"github.com/shopspring/decimal"
)

type AddPetsToOwnerRequest struct {
	Pets []AddPetRequest `json:"pets" validate:"required"`
}

type AddPetRequest struct {
	Name           string           `json:"name"           validate:"required"`
	PetType        commonvo.PetType `json:"petType"        validate:"required,oneof=dog cat"`
	Sex            Gender           `json:"sex"            validate:"required,oneof=male female"`
	Neutered       bool             `json:"neutered"       validate:"required"`
	Breed          string           `json:"breed"          validate:"required"`
	BirthDate      string           `json:"birthDate"      validate:"required"`
	WeightInKg     decimal.Decimal  `json:"weightInKg"     validate:"required"`
	Remarks        string           `json:"remarks"`
	ProfileImageID uuid.NullUUID    `json:"profileImageId"`
}

type UpdatePetRequest struct {
	Name           string          `json:"name"           validate:"required"`
	Neutered       bool            `json:"neutered"       validate:"required"`
	Breed          string          `json:"breed"          validate:"required"`
	BirthDate      string          `json:"birthDate"      validate:"required"`
	WeightInKg     decimal.Decimal `json:"weightInKg"     validate:"required"`
	Remarks        string          `json:"remarks"`
	ProfileImageID uuid.NullUUID   `json:"profileImageId"`
}
