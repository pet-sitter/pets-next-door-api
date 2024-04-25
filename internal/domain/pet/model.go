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

func ToWithProfileImageFromSOSPostIDRow(row databasegen.FindPetsBySOSPostIDRow) *WithProfileImage {
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

// ViewForSOSPost
// v_pets_for_sos_posts 뷰를 위한 구조체
type ViewForSOSPost struct {
	ID              int              `field:"id" json:"id"`
	OwnerID         int              `field:"owner_id" json:"owner_id"`
	Name            string           `field:"name" json:"name"`
	PetType         commonvo.PetType `field:"pet_type" json:"pet_type"`
	Sex             Gender           `field:"sex" json:"sex"`
	Neutered        bool             `field:"neutered" json:"neutered"`
	Breed           string           `field:"breed" json:"breed"`
	BirthDate       datatype.Date    `field:"birth_date" json:"birth_date"`
	WeightInKg      decimal.Decimal  `field:"weight_in_kg" json:"weight_in_kg"`
	Remarks         string           `field:"remarks" json:"remarks"`
	CreatedAt       string           `field:"created_at" json:"created_at"`
	UpdatedAt       string           `field:"updated_at" json:"updated_at"`
	DeletedAt       string           `field:"deleted_at" json:"deleted_at"`
	ProfileImageURL *string          `field:"profile_image_url" json:"profile_image_url"`
}

func (v *ViewForSOSPost) ToDetailView() *DetailView {
	return &DetailView{
		ID:              int64(v.ID),
		Name:            v.Name,
		PetType:         v.PetType,
		Sex:             v.Sex,
		Neutered:        v.Neutered,
		Breed:           v.Breed,
		BirthDate:       v.BirthDate,
		WeightInKg:      v.WeightInKg,
		Remarks:         v.Remarks,
		ProfileImageURL: v.ProfileImageURL,
	}
}

type ViewListForSOSPost []*ViewForSOSPost

func (vl *ViewListForSOSPost) ToDetailViewList() []DetailView {
	pl := make([]DetailView, len(*vl))
	for i, v := range *vl {
		pl[i] = *v.ToDetailView()
	}
	return pl
}
