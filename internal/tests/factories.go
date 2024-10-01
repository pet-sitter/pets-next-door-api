package tests

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sospost"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/commonvo"

	"github.com/pet-sitter/pets-next-door-api/internal/datatype"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/shopspring/decimal"
)

func NewDummyRegisterUserRequest(profileImageID uuid.NullUUID) *user.RegisterUserRequest {
	return &user.RegisterUserRequest{
		Email:                datatype.NewUUIDV7().String() + "@example.com",
		Nickname:             datatype.NewUUIDV7().String()[0:20],
		Fullname:             datatype.NewUUIDV7().String()[0:20],
		ProfileImageID:       profileImageID,
		FirebaseProviderType: user.FirebaseProviderTypeKakao,
		FirebaseUID:          datatype.NewUUIDV7().String(),
	}
}

func NewDummyAddPetRequest(
	profileImageID uuid.NullUUID, petType commonvo.PetType, gender pet.Gender, breed string,
) *pet.AddPetRequest {
	birthDate, _ := datatype.ParseDate("2020-01-01")
	return &pet.AddPetRequest{
		Name:           datatype.NewUUIDV7().String()[0:20],
		PetType:        petType,
		Sex:            gender,
		Neutered:       true,
		Breed:          breed,
		BirthDate:      birthDate.String(),
		WeightInKg:     decimal.NewFromFloat(10.0),
		ProfileImageID: profileImageID,
	}
}

func NewDummyWriteSOSPostRequest(
	imageID,
	petIDs []uuid.UUID,
	sosPostCnt int,
	conditionIDs []uuid.UUID,
) *sospost.WriteSOSPostRequest {
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
		ConditionIDs: conditionIDs,
		PetIDs:       petIDs,
	}
}
