package pet

type BreedService struct {
	breedStore BreedStore
}

func NewBreedService(breedStore BreedStore) *BreedService {
	return &BreedService{
		breedStore: breedStore,
	}
}

func (service *BreedService) FindBreeds(page int, size int, petType *string) ([]*BreedView, error) {
	breeds, err := service.breedStore.FindBreeds(page, size, petType)
	if err != nil {
		return nil, err
	}

	breedViews := make([]*BreedView, 0)
	for _, breed := range breeds {
		breedViews = append(breedViews, &BreedView{
			ID:      breed.ID,
			PetType: breed.PetType,
			Name:    breed.Name,
		})
	}

	return breedViews, nil
}

func (service *BreedService) FindBreedByPetTypeAndName(petType PetType, name string) (*BreedView, error) {
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
