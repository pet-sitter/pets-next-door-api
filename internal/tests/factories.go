package tests

import (
	"fmt"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/commonvo"

	"github.com/pet-sitter/pets-next-door-api/internal/datatype"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sospost"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/shopspring/decimal"
)

func randomUUID() string {
	uuid, _ := datatype.NewV7()
	return uuid.String()
}

func NewDummyRegisterUserRequest(profileImageID *int64) *user.RegisterUserRequest {
	return &user.RegisterUserRequest{
		Email:                randomUUID() + "@example.com",
		Nickname:             randomUUID()[0:20],
		Fullname:             randomUUID()[0:20],
		ProfileImageID:       profileImageID,
		FirebaseProviderType: user.FirebaseProviderTypeKakao,
		FirebaseUID:          randomUUID(),
	}
}

func NewDummyAddPetRequest(
	profileImageID *int64, petType commonvo.PetType, gender pet.Gender, breed string,
) *pet.AddPetRequest {
	birthDate, _ := datatype.ParseDate("2020-01-01")
	return &pet.AddPetRequest{
		Name:           randomUUID()[0:20],
		PetType:        petType,
		Sex:            gender,
		Neutered:       true,
		Breed:          breed,
		BirthDate:      birthDate.String(),
		WeightInKg:     decimal.NewFromFloat(10.0),
		ProfileImageID: profileImageID,
	}
}

func NewDummyWriteSOSPostRequest(imageID, petIDs []int64, sosPostCnt int) *sospost.WriteSOSPostRequest {
	return &sospost.WriteSOSPostRequest{
		Title:    fmt.Sprintf("Title%d", sosPostCnt),
		Content:  fmt.Sprintf("Content%d", sosPostCnt),
		ImageIDs: imageID,
		Reward:   "Reward",
		Dates: []sospost.SOSDateView{
			{DateStartAt: fmt.Sprintf("2024-04-1%d", sosPostCnt), DateEndAt: fmt.Sprintf("2024-04-2%d", sosPostCnt)},
			{DateStartAt: fmt.Sprintf("2024-05-1%d", sosPostCnt), DateEndAt: fmt.Sprintf("2024-05-2%d", sosPostCnt)},
		},
		CareType:     commonvo.CareTypeFoster,
		CarerGender:  commonvo.CarerGenderMale,
		RewardType:   commonvo.RewardTypeFee,
		ConditionIDs: []int{1, 2},
		PetIDs:       petIDs,
	}
}
