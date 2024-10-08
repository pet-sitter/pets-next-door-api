package pet

import (
	"github.com/google/uuid"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/commonvo"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
	"github.com/shopspring/decimal"
)

type DetailView struct {
	ID              uuid.UUID        `json:"id"`
	Name            string           `json:"name"`
	PetType         commonvo.PetType `json:"petType"`
	Sex             Gender           `json:"sex"`
	Neutered        bool             `json:"neutered"`
	Breed           string           `json:"breed"`
	BirthDate       string           `json:"birthDate"`
	WeightInKg      decimal.Decimal  `json:"weightInKg"`
	Remarks         string           `json:"remarks"`
	ProfileImageURL *string          `json:"profileImageUrl"`
}

func (pet *WithProfileImage) ToDetailView() *DetailView {
	return &DetailView{
		ID:              pet.ID,
		Name:            pet.Name,
		PetType:         pet.PetType,
		Sex:             pet.Sex,
		Neutered:        pet.Neutered,
		Breed:           pet.Breed,
		BirthDate:       pet.BirthDate.String(),
		WeightInKg:      pet.WeightInKg,
		Remarks:         pet.Remarks,
		ProfileImageURL: pet.ProfileImageURL,
	}
}

func ToDetailViewList(rows []databasegen.FindPetsBySOSPostIDRow) []DetailView {
	pl := make([]DetailView, len(rows))
	for i, row := range rows {
		pl[i] = *ToWithProfileImageFromSOSPostIDRow(row).ToDetailView()
	}
	return pl
}

type ListView struct {
	Pets []DetailView `json:"pets"`
}

func ToListView(rows []databasegen.FindPetsRow) *ListView {
	pl := &ListView{Pets: make([]DetailView, len(rows))}
	for i, row := range rows {
		pl.Pets[i] = *ToWithProfileImageFromRows(row).ToDetailView()
	}
	return pl
}

func ToListViewFromIDsRows(rows []databasegen.FindPetsByIDsRow) *ListView {
	pl := &ListView{Pets: make([]DetailView, len(rows))}
	for i, row := range rows {
		pl.Pets[i] = *ToWithProfileImageFromIDsRows(row).ToDetailView()
	}
	return pl
}
