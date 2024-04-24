package tests

import (
	"fmt"

	"github.com/pet-sitter/pets-next-door-api/internal/datatype"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sospost"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/shopspring/decimal"
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
	birthDate, _ := datatype.ParseDate("2020-01-01")
	return &pet.AddPetRequest{
		Name:           "name",
		PetType:        "dog",
		Sex:            "male",
		Neutered:       true,
		Breed:          "poodle",
		BirthDate:      birthDate,
		WeightInKg:     decimal.NewFromFloat(10.0),
		ProfileImageID: profileImageID,
	}
}

func GenerateDummyAddPetsRequest(profileImageID *int) []pet.AddPetRequest {
	birthDate1, _ := datatype.ParseDate("2020-01-01")
	birthDate2, _ := datatype.ParseDate("2020-02-01")
	birthDate3, _ := datatype.ParseDate("2020-03-01")

	return []pet.AddPetRequest{
		{
			Name:           "dog_1",
			PetType:        "dog",
			Sex:            "male",
			Neutered:       true,
			Breed:          "poodle",
			BirthDate:      birthDate1,
			WeightInKg:     decimal.NewFromFloat(10.0),
			Remarks:        "remarks",
			ProfileImageID: profileImageID,
		},
		{
			Name:           "dog_2",
			PetType:        "dog",
			Sex:            "male",
			Neutered:       true,
			Breed:          "poodle",
			BirthDate:      birthDate2,
			WeightInKg:     decimal.NewFromFloat(10.0),
			Remarks:        "remarks",
			ProfileImageID: profileImageID,
		},
		{
			Name:           "cat_1",
			PetType:        "cat",
			Sex:            "female",
			Neutered:       true,
			Breed:          "munchkin",
			BirthDate:      birthDate3,
			WeightInKg:     decimal.NewFromFloat(8.0),
			Remarks:        "remarks",
			ProfileImageID: profileImageID,
		},
	}
}

func GenerateDummyWriteSOSPostRequest(imageID, petIDs []int64, sosPostCnt int) *sospost.WriteSOSPostRequest {
	return &sospost.WriteSOSPostRequest{
		Title:    fmt.Sprintf("Title%d", sosPostCnt),
		Content:  fmt.Sprintf("Content%d", sosPostCnt),
		ImageIDs: imageID,
		Reward:   "Reward",
		Dates: []sospost.SOSDateView{
			{DateStartAt: fmt.Sprintf("2024-04-1%d", sosPostCnt), DateEndAt: fmt.Sprintf("2024-04-2%d", sosPostCnt)},
			{DateStartAt: fmt.Sprintf("2024-05-1%d", sosPostCnt), DateEndAt: fmt.Sprintf("2024-05-2%d", sosPostCnt)},
		},
		CareType:     sospost.CareTypeFoster,
		CarerGender:  sospost.CarerGenderMale,
		RewardType:   sospost.RewardTypeFee,
		ConditionIDs: []int{1, 2},
		PetIDs:       petIDs,
	}
}
