package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/commonvo"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"

	"github.com/pet-sitter/pets-next-door-api/internal/datatype"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/tests"
	"github.com/shopspring/decimal"
)

func TestRegisterUser(t *testing.T) {
	t.Run("사용자를 새로 생성한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		mediaService := tests.NewMockMediaService(db)
		userService := tests.NewMockUserService(db)

		// Given
		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")

		userRequest := tests.NewDummyRegisterUserRequest(&profileImage.ID)

		// When
		created, _ := userService.RegisterUser(ctx, userRequest)

		// Then
		assert.Equal(t, userRequest.Email, created.Email)
	})

	t.Run("사용자의 프로필 이미지가 존재하지 않아도 생성한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		userService := tests.NewMockUserService(db)

		// Given
		userRequest := &user.RegisterUserRequest{
			Email:                "test@example.com",
			Nickname:             "nickname",
			Fullname:             "fullname",
			ProfileImageID:       nil,
			FirebaseProviderType: user.FirebaseProviderTypeKakao,
			FirebaseUID:          "uid",
		}

		// When
		created, _ := userService.RegisterUser(ctx, userRequest)

		// Then
		assert.Equal(t, userRequest.Email, created.Email)
	})

	t.Run("사용자가 이미 존재할 경우 에러를 반환한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		mediaService := tests.NewMockMediaService(db)
		userService := tests.NewMockUserService(db)

		// Given
		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")
		userRequest := tests.NewDummyRegisterUserRequest(&profileImage.ID)
		userService.RegisterUser(ctx, userRequest)

		// When
		_, err := userService.RegisterUser(ctx, userRequest)

		// Then
		assert.NotNil(t, err)
	})
}

func TestFindUsers(t *testing.T) {
	t.Run("사용자를 닉네임으로 검색한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		mediaService := tests.NewMockMediaService(db)
		userService := tests.NewMockUserService(db)

		// Given
		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")

		targetNickname := "target"
		targetUserRequest := &user.RegisterUserRequest{
			Email:                "test@example.com",
			Nickname:             targetNickname,
			Fullname:             "fullname",
			ProfileImageID:       &profileImage.ID,
			FirebaseProviderType: user.FirebaseProviderTypeKakao,
			FirebaseUID:          "uid",
		}

		userService.RegisterUser(ctx, targetUserRequest)
		for i := 0; i < 2; i++ {
			userService.RegisterUser(ctx, &user.RegisterUserRequest{
				Email:                fmt.Sprintf("test%d@example.com", i),
				Nickname:             fmt.Sprintf("nickname%d", i),
				Fullname:             fmt.Sprintf("fullname%d", i),
				ProfileImageID:       &profileImage.ID,
				FirebaseProviderType: user.FirebaseProviderTypeKakao,
				FirebaseUID:          fmt.Sprintf("uid%d", i),
			})
		}

		// When
		found, _ := userService.FindUsers(ctx, user.FindUsersParams{Page: 1, Size: 20, Nickname: &targetNickname})

		// Then
		assert.Equal(t, 1, len(found.Items))
	})
}

func TestFindUser(t *testing.T) {
	t.Run("사용자를 이메일로 찾는다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		mediaService := tests.NewMockMediaService(db)
		userService := tests.NewMockUserService(db)

		// Given
		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")

		userRequest := tests.NewDummyRegisterUserRequest(&profileImage.ID)
		created, _ := userService.RegisterUser(ctx, userRequest)

		// When
		found, _ := userService.FindUser(ctx, user.FindUserParams{Email: &created.Email})

		// Then
		assert.Equal(t, created.Email, found.Email)
	})

	t.Run("사용자를 UID로 찾는다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		mediaService := tests.NewMockMediaService(db)
		userService := tests.NewMockUserService(db)

		// Given
		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")

		userRequest := tests.NewDummyRegisterUserRequest(&profileImage.ID)
		created, _ := userService.RegisterUser(ctx, userRequest)

		// When
		found, _ := userService.FindUser(ctx, user.FindUserParams{FbUID: &created.FirebaseUID})

		// Then
		assert.Equal(t, created.FirebaseUID, found.FirebaseUID)
	})

	t.Run("사용자가 존재하지 않을 경우 에러를 반환한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		userService := tests.NewMockUserService(db)

		// When
		email := "non-existent@example.com"
		_, err := userService.FindUser(ctx, user.FindUserParams{Email: &email})

		// Then
		assert.NotNil(t, err)
	})
}

func TestExistsByEmail(t *testing.T) {
	t.Run("사용자의 닉네임이 존재하지 않을 경우 false를 반환한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		userService := tests.NewMockUserService(db)

		// When
		exists, _ := userService.ExistsByNickname(ctx, "non-existent")

		// Then
		assert.False(t, exists)
	})

	t.Run("사용자의 닉네임이 존재할 경우 true를 반환한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		mediaService := tests.NewMockMediaService(db)
		userService := tests.NewMockUserService(db)

		// Given
		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")

		userRequest := tests.NewDummyRegisterUserRequest(&profileImage.ID)
		userService.RegisterUser(ctx, userRequest)

		// When
		exists, _ := userService.ExistsByNickname(ctx, userRequest.Nickname)

		// Then
		assert.True(t, exists)
	})
}

func TestUpdateUserByUID(t *testing.T) {
	t.Run("사용자를 업데이트한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		mediaService := tests.NewMockMediaService(db)
		userService := tests.NewMockUserService(db)

		// Given
		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")
		targetUser, _ := userService.RegisterUser(ctx, tests.NewDummyRegisterUserRequest(&profileImage.ID))

		// When
		updatedNickname := "updated"
		updatedProfileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "updated_profile_image.jpg")
		userService.UpdateUserByUID(ctx, targetUser.FirebaseUID, updatedNickname, &updatedProfileImage.ID)

		// Then
		found, _ := userService.FindUser(ctx, user.FindUserParams{FbUID: &targetUser.FirebaseUID})
		assert.Equal(t, updatedNickname, found.Nickname)
		assert.Equal(t, updatedProfileImage.URL, *found.ProfileImageURL)
	})
}

func TestAddPetsToOwner(t *testing.T) {
	t.Run("사용자에게 반려동물을 추가한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		mediaService := tests.NewMockMediaService(db)
		userService := tests.NewMockUserService(db)

		// Given
		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")
		owner, _ := userService.RegisterUser(ctx, tests.NewDummyRegisterUserRequest(&profileImage.ID))

		petProfileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "pet_profile_image.jpg")
		petsToAdd := pet.AddPetsToOwnerRequest{
			Pets: []pet.AddPetRequest{
				*tests.NewDummyAddPetRequest(&petProfileImage.ID, commonvo.PetTypeDog, pet.GenderMale, "poodle"),
			},
		}

		// When
		created, _ := userService.AddPetsToOwner(ctx, owner.FirebaseUID, petsToAdd)

		// Then
		assert.Equal(t, 1, len(created.Pets))

		for _, want := range petsToAdd.Pets {
			for _, got := range created.Pets {
				assert.Equal(t, want.Name, got.Name)
				assert.Equal(t, want.PetType, got.PetType)
				assert.Equal(t, want.Sex, got.Sex)
				assert.Equal(t, want.Neutered, got.Neutered)
				assert.Equal(t, want.Breed, got.Breed)
				assert.Equal(t, want.BirthDate, got.BirthDate)
				assert.Equal(t, want.WeightInKg.String(), got.WeightInKg.String())
			}
		}
	})
}

func TestUpdatePet(t *testing.T) {
	t.Run("반려동물을 업데이트한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		mediaService := tests.NewMockMediaService(db)
		userService := tests.NewMockUserService(db)

		// Given
		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")
		registeredUser, _ := userService.RegisterUser(ctx, tests.NewDummyRegisterUserRequest(&profileImage.ID))

		petProfileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "pet_profile_image.jpg")
		petRequest := tests.NewDummyAddPetRequest(&petProfileImage.ID, commonvo.PetTypeDog, pet.GenderMale, "poodle")
		createdPets, _ := userService.AddPetsToOwner(
			ctx, registeredUser.FirebaseUID, pet.AddPetsToOwnerRequest{Pets: []pet.AddPetRequest{*petRequest}})
		createdPet := createdPets.Pets[0]

		// When
		updatedPetProfileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "updated_pet_profile_image.jpg")
		birthData, _ := datatype.ParseDate("2021-01-01")
		updatedPetRequest := pet.UpdatePetRequest{
			Name:           "updated",
			Neutered:       true,
			Breed:          "updated",
			BirthDate:      birthData,
			WeightInKg:     decimal.NewFromFloat(10.0),
			Remarks:        "updated",
			ProfileImageID: &updatedPetProfileImage.ID,
		}

		_, err := userService.UpdatePet(ctx, registeredUser.FirebaseUID, createdPet.ID, updatedPetRequest)
		if err != nil {
			t.Errorf("got %v want %v", err, nil)
		}

		// Then
		found, _ := userService.FindPets(ctx, pet.FindPetsParams{OwnerID: &registeredUser.ID})
		assert.Equal(t, 1, len(found.Pets))

		want := updatedPetRequest
		got := found.Pets[0]
		assert.Equal(t, want.Name, got.Name)
		assert.Equal(t, want.Neutered, got.Neutered)
		assert.Equal(t, want.Breed, got.Breed)
		assert.Equal(t, want.BirthDate, got.BirthDate)
		assert.Equal(t, want.WeightInKg.String(), got.WeightInKg.String())
	})
}

func TestDeletePet(t *testing.T) {
	t.Run("반려동물을 삭제한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()
		mediaService := tests.NewMockMediaService(db)
		userService := tests.NewMockUserService(db)

		// Given
		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")
		registeredUser, _ := userService.RegisterUser(ctx, tests.NewDummyRegisterUserRequest(&profileImage.ID))

		petProfileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "pet_profile_image.jpg")
		petRequest := tests.NewDummyAddPetRequest(&petProfileImage.ID, commonvo.PetTypeDog, pet.GenderMale, "poodle")
		createdPets, _ := userService.AddPetsToOwner(
			ctx,
			registeredUser.FirebaseUID,
			pet.AddPetsToOwnerRequest{Pets: []pet.AddPetRequest{*petRequest}},
		)
		createdPet := createdPets.Pets[0]

		// When
		userService.DeletePet(ctx, registeredUser.FirebaseUID, createdPet.ID)

		// Then
		found, _ := userService.FindPets(ctx, pet.FindPetsParams{
			OwnerID: &registeredUser.ID,
		})
		assert.Equal(t, 0, len(found.Pets))
	})
}
