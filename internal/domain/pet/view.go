package pet

import pnd "github.com/pet-sitter/pets-next-door-api/api"

type AddPetsToOwnerRequest struct {
	Pets []AddPetRequest `json:"pets" validate:"required"`
}

type AddPetRequest struct {
	Name           string  `json:"name" validate:"required"`
	PetType        PetType `json:"petType" validate:"required,oneof=dog cat"`
	Sex            PetSex  `json:"sex" validate:"required,oneof=male female"`
	Neutered       bool    `json:"neutered" validate:"required"`
	Breed          string  `json:"breed" validate:"required"`
	BirthDate      string  `json:"birthDate" validate:"required"`
	WeightInKg     float64 `json:"weightInKg" validate:"required"`
	ProfileImageID *int    `json:"profileImageId"`
}

type FindMyPetsView struct {
	Pets []PetView `json:"pets"`
}

func NewFindMyPetsView(pets []PetWithProfileImage) *FindMyPetsView {
	petViews := make([]PetView, len(pets))
	for i, pet := range pets {
		petViews[i] = *pet.ToPetView()
	}
	return &FindMyPetsView{Pets: petViews}
}

type PetView struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	PetType         PetType `json:"petType"`
	Sex             PetSex  `json:"sex"`
	Neutered        bool    `json:"neutered"`
	Breed           string  `json:"breed"`
	BirthDate       string  `json:"birthDate"`
	WeightInKg      float64 `json:"weightInKg"`
	ProfileImageURL *string `json:"profileImageUrl"`
}

func (pet *PetWithProfileImage) ToPetView() *PetView {
	return &PetView{
		ID:              pet.ID,
		Name:            pet.Name,
		PetType:         pet.PetType,
		Sex:             pet.Sex,
		Neutered:        pet.Neutered,
		Breed:           pet.Breed,
		BirthDate:       pet.BirthDate,
		WeightInKg:      pet.WeightInKg,
		ProfileImageURL: pet.ProfileImageURL,
	}
}

func NewPetViewList(pets []PetWithProfileImage) []PetView {
	petViews := make([]PetView, len(pets))
	for i, pet := range pets {
		petViews[i] = *pet.ToPetView()
	}
	return petViews
}

type BreedView struct {
	ID      int     `json:"id"`
	PetType PetType `json:"petType"`
	Name    string  `json:"name"`
}

type BreedListView struct {
	*pnd.PaginatedView[*BreedView]
}

func FromBreedList(breeds *BreedList) *BreedListView {
	breedViews := make([]*BreedView, 0)
	for _, breed := range breeds.Items {
		breedViews = append(breedViews, &BreedView{
			ID:      breed.ID,
			PetType: breed.PetType,
			Name:    breed.Name,
		})
	}

	return &BreedListView{
		PaginatedView: pnd.NewPaginatedView(
			breeds.Page, breeds.Size, breeds.IsLastPage, breedViews,
		),
	}
}
