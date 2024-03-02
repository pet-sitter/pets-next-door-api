package pet

import pnd "github.com/pet-sitter/pets-next-door-api/api"

type AddPetsToOwnerRequest struct {
	Pets []AddPetRequest `json:"pets" validate:"required"`
}

type AddPetRequest struct {
	Name       string  `json:"name" validate:"required"`
	PetType    PetType `json:"petType" validate:"required,oneof=dog cat"`
	Sex        PetSex  `json:"sex" validate:"required,oneof=male female"`
	Neutered   bool    `json:"neutered" validate:"required"`
	Breed      string  `json:"breed" validate:"required"`
	BirthDate  string  `json:"birthDate" validate:"required"`
	WeightInKg float64 `json:"weightInKg" validate:"required"`
}

type FindMyPetsView struct {
	Pets []PetView `json:"pets"`
}

type PetView struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	PetType    PetType `json:"petType"`
	Sex        PetSex  `json:"sex"`
	Neutered   bool    `json:"neutered"`
	Breed      string  `json:"breed"`
	BirthDate  string  `json:"birthDate"`
	WeightInKg float64 `json:"weightInKg"`
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
