package tests

import (
	"fmt"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sos_post"
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

func GenerateDummyAddPetsRequest(profileImageID *int) []pet.AddPetRequest {
	return []pet.AddPetRequest{
		{
			Name:           "dog_1",
			PetType:        "dog",
			Sex:            "male",
			Neutered:       true,
			Breed:          "poodle",
			BirthDate:      "2020-01-01",
			WeightInKg:     10.0,
			Remarks:        "remarks",
			ProfileImageID: profileImageID,
		},
		{
			Name:           "dog_2",
			PetType:        "dog",
			Sex:            "male",
			Neutered:       true,
			Breed:          "poodle",
			BirthDate:      "2020-02-01",
			WeightInKg:     10.0,
			Remarks:        "remarks",
			ProfileImageID: profileImageID,
		},
		{
			Name:           "cat_1",
			PetType:        "cat",
			Sex:            "female",
			Neutered:       true,
			Breed:          "munchkin",
			BirthDate:      "2020-03-01",
			WeightInKg:     8.0,
			Remarks:        "remarks",
			ProfileImageID: profileImageID,
		},
	}
}

func GenerateDummyWriteSosPostRequest(imageID []int, petIDs []int, sosPostCnt int) *sos_post.WriteSosPostRequest {
	return &sos_post.WriteSosPostRequest{
		Title:    fmt.Sprintf("Title%d", sosPostCnt),
		Content:  fmt.Sprintf("Content%d", sosPostCnt),
		ImageIDs: imageID,
		Reward:   "Reward",
		Dates: []sos_post.SosDateView{
			{fmt.Sprintf("2024-04-1%d", sosPostCnt), fmt.Sprintf("2024-04-2%d", sosPostCnt)},
			{fmt.Sprintf("2024-05-1%d", sosPostCnt), fmt.Sprintf("2024-05-2%d", sosPostCnt)},
		},
		CareType:     sos_post.CareTypeFoster,
		CarerGender:  sos_post.CarerGenderMale,
		RewardType:   sos_post.RewardTypeFee,
		ConditionIDs: []int{1, 2},
		PetIDs:       petIDs,
	}
}
