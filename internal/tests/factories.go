package tests

import (
	"github.com/google/uuid"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/commonvo"

	"github.com/pet-sitter/pets-next-door-api/internal/datatype"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/shopspring/decimal"
)

func randomUUID() uuid.UUID {
	newUUID, _ := datatype.NewV7()
	return newUUID
}

func randomUUIDStr() string {
	return randomUUID().String()
}

func NewDummyRegisterUserRequest(profileImageID uuid.NullUUID) *user.RegisterUserRequest {
	return &user.RegisterUserRequest{
		Email:                randomUUIDStr() + "@example.com",
		Nickname:             randomUUIDStr()[0:20],
		Fullname:             randomUUIDStr()[0:20],
		ProfileImageID:       profileImageID,
		FirebaseProviderType: user.FirebaseProviderTypeKakao,
		FirebaseUID:          randomUUIDStr(),
	}
}

func NewDummyAddPetRequest(
	profileImageID uuid.NullUUID, petType commonvo.PetType, gender pet.Gender, breed string,
) *pet.AddPetRequest {
	birthDate, _ := datatype.ParseDate("2020-01-01")
	return &pet.AddPetRequest{
		Name:           randomUUIDStr()[0:20],
		PetType:        petType,
		Sex:            gender,
		Neutered:       true,
		Breed:          breed,
		BirthDate:      birthDate.String(),
		WeightInKg:     decimal.NewFromFloat(10.0),
		ProfileImageID: profileImageID,
	}
}

// func NewDummyWriteSOSPostRequest(imageID, petIDs []uuid.UUID, sosPostCnt int) *sospost.WriteSOSPostRequest {
// 	return &sospost.WriteSOSPostRequest{
// 		Title:    fmt.Sprintf("Title%d", sosPostCnt),
// 		Content:  fmt.Sprintf("Content%d", sosPostCnt),
// 		ImageIDs: imageID,
// 		Reward:   "Reward",
// 		Dates: []sospost.SOSDateView{
// 			{DateStartAt: fmt.Sprintf("2024-04-1%d", sosPostCnt), DateEndAt: fmt.Sprintf("2024-04-2%d", sosPostCnt)},
// 			{DateStartAt: fmt.Sprintf("2024-05-1%d", sosPostCnt), DateEndAt: fmt.Sprintf("2024-05-2%d", sosPostCnt)},
// 		},
// 		CareType:     sospost.CareTypeFoster,
// 		CarerGender:  sospost.CarerGenderMale,
// 		RewardType:   sospost.RewardTypeFee,
// 		ConditionIDs: []uuid.UUID{randomUUID(), randomUUID()},
// 		PetIDs:       petIDs,
// 	}
// }
