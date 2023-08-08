package types

type PetType string

const (
	PetTypeDog PetType = "dog"
	PetTypeCat PetType = "cat"
)

type PetSex string

const (
	PetSexMale   PetSex = "male"
	PetSexFemale PetSex = "female"
)

type DogBreed string

// TODO: add more dog breeds
const (
	DogBreedPoodle DogBreed = "poodle"
)

type CatBreed string

// TODO: add more cat breeds
const (
	CatBreedPersian CatBreed = "persian"
)
