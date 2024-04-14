package tests

import (
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
)

func GenerateDummyRegisterUserRequest(profileImageID *int) *user.RegisterUserRequest {
	return &user.RegisterUserRequest{
		Email:                "test@example.com",
		Nickname:             "nickname",
		Fullname:             "fullname",
		ProfileImageID:       profileImageID,
		FirebaseProviderType: user.FirebaseProviderTypeKakao,
		FirebaseUID:          "uid",
	}
}

func GenerateDummyAddPetRequest(profileImageID *int) *pet.AddPetRequest {
	return &pet.AddPetRequest{
		Name:           "name",
		PetType:        "dog",
		Sex:            "male",
		Neutered:       true,
		Breed:          "poodle",
		BirthDate:      "2020-01-01",
		WeightInKg:     10.0,
		ProfileImageID: profileImageID,
	}
}
