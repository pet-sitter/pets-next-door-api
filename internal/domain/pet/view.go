package pet

import (
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
)

type FindMyPetsView struct {
	Pets []PetView `json:"pets"`
}

func (pets *PetWithProfileList) ToFindMyPetsView() *FindMyPetsView {
	petViews := make([]PetView, len(*pets))
	for i, pet := range *pets {
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
	Remarks         string  `json:"remarks"`
	ProfileImageURL *string `json:"profileImageUrl"`
}

func (pet *Pet) ToPetView() *PetView {
	return &PetView{
		ID:         pet.ID,
		Name:       pet.Name,
		PetType:    pet.PetType,
		Sex:        pet.Sex,
		Neutered:   pet.Neutered,
		Breed:      pet.Breed,
		BirthDate:  utils.FormatDate(pet.BirthDate),
		WeightInKg: pet.WeightInKg,
		Remarks:    pet.Remarks,
	}
}

func (pets *PetList) ToPetViewList() []PetView {
	petViews := make([]PetView, len(*pets))
	for i, pet := range *pets {
		petViews[i] = *pet.ToPetView()
	}
	return petViews
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
		Remarks:         pet.Remarks,
		ProfileImageURL: pet.ProfileImageURL,
	}
}

func (pets *PetWithProfileList) ToPetViewList() []PetView {
	petViews := make([]PetView, len(*pets))
	for i, pet := range *pets {
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

func (breeds *BreedList) ToBreedListView() *BreedListView {
	breedViews := make([]*BreedView, len(breeds.Items))
	for i, breed := range breeds.Items {
		breedViews[i] = &BreedView{
			ID:      breed.ID,
			PetType: breed.PetType,
			Name:    breed.Name,
		}
	}

	return &BreedListView{
		PaginatedView: pnd.NewPaginatedView(
			breeds.Page, breeds.Size, breeds.IsLastPage, breedViews,
		),
	}
}
