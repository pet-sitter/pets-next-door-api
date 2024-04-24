package breed

import (
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/commonvo"
)

type BreedView struct {
	ID      int              `json:"id"`
	PetType commonvo.PetType `json:"petType"`
	Name    string           `json:"name"`
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
