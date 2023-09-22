package pet

type AddPetsToOwnerRequest struct {
	Pets []AddPetRequest `json:"pets" validate:"required"`
}

type AddPetRequest struct {
	Name       string  `json:"name" validate:"required"`
	PetType    PetType `json:"pet_type" validate:"required,oneof=dog cat"`
	Sex        PetSex  `json:"sex" validate:"required,oneof=male female"`
	Neutered   bool    `json:"neutered" validate:"required"`
	Breed      string  `json:"breed" validate:"required"`
	BirthDate  string  `json:"birth_date" validate:"required"`
	WeightInKg float64 `json:"weight_in_kg" validate:"required"`
}

type FindMyPetsView struct {
	Pets []PetView `json:"pets"`
}

type PetView struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	PetType    PetType `json:"pet_type"`
	Sex        PetSex  `json:"sex"`
	Neutered   bool    `json:"neutered"`
	Breed      string  `json:"breed"`
	BirthDate  string  `json:"birth_date"`
	WeightInKg float64 `json:"weight_in_kg"`
}

type BreedView struct {
	ID      int     `json:"id"`
	PetType PetType `json:"pet_type"`
	Name    string  `json:"name"`
}
