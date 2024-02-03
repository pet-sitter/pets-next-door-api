package pet

import pnd "github.com/pet-sitter/pets-next-door-api/api"

type BreedService struct {
	breedStore BreedStore
}

func NewBreedService(breedStore BreedStore) *BreedService {
	return &BreedService{
		breedStore: breedStore,
	}
}

func (service *BreedService) FindBreeds(page int, size int, petType *string) (*BreedListView, *pnd.AppError) {
	breeds, err := service.breedStore.FindBreeds(page, size, petType)
	if err != nil {
		return nil, err
	}

	return FromBreedList(breeds), nil
}

func (service *BreedService) FindBreedByPetTypeAndName(petType PetType, name string) (*BreedView, *pnd.AppError) {
	breed, err := service.breedStore.FindBreedByPetTypeAndName(petType, name)
	if err != nil {
		return nil, err
	}

	return &BreedView{
		ID:      breed.ID,
		PetType: breed.PetType,
		Name:    breed.Name,
	}, nil
}
