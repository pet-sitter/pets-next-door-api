package pet

import (
	"github.com/pet-sitter/pets-next-door-api/internal/datatype"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/commonvo"
	"github.com/shopspring/decimal"
)

type BasePet struct {
	ID         int              `field:"id" json:"id"`
	OwnerID    int              `field:"owner_id" json:"owner_id"`
	Name       string           `field:"name" json:"name"`
	PetType    commonvo.PetType `field:"pet_type" json:"pet_type"`
	Sex        Gender           `field:"sex" json:"sex"`
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

type ViewForSOSPost struct {
	BasePet
	ProfileImageURL *string `field:"profile_image_url" json:"profile_image_url"`
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
