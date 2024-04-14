package pet

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
	Remarks        string  `json:"remarks"`
	ProfileImageID *int    `json:"profileImageId"`
}

func (r *AddPetRequest) ToBasePet(ownerID int) *BasePet {
	return &BasePet{
		OwnerID:    ownerID,
		Name:       r.Name,
		PetType:    r.PetType,
		Sex:        r.Sex,
		Neutered:   r.Neutered,
		Breed:      r.Breed,
		BirthDate:  r.BirthDate,
		WeightInKg: r.WeightInKg,
		Remarks:    r.Remarks,
	}
}

func (r *AddPetRequest) ToPet(ownerID int) *Pet {
	return &Pet{
		BasePet:        *r.ToBasePet(ownerID),
		ProfileImageID: r.ProfileImageID,
	}
}
