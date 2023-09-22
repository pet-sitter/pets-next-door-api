package pet

type Breed struct {
	ID        int     `field:"id"`
	Name      string  `field:"name"`
	PetType   PetType `field:"pet_type"`
	CreatedAt string  `field:"created_at"`
	UpdatedAt string  `field:"updated_at"`
	DeletedAt string  `field:"deleted_at"`
}

type BreedStore interface {
	FindBreeds(page int, size int, petType *string) ([]*Breed, error)
	FindBreedByPetTypeAndName(petType PetType, name string) (*Breed, error)
	CreateBreed(breed *Breed) (*Breed, error)
}
