package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/pet-sitter/pets-next-door-api/internal/tests/assert"

	"github.com/pet-sitter/pets-next-door-api/internal/datatype"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
	"github.com/pet-sitter/pets-next-door-api/internal/tests"
	"github.com/shopspring/decimal"
)

func TestRegisterUser(t *testing.T) {
	t.Run("사용자를 새로 생성한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()

		// Given
		mediaService := service.NewMediaService(db, nil)
		profileImage := tests.AddDummyMedia(t, ctx, mediaService)

		userService := service.NewUserService(db, mediaService)
		userRequest := tests.GenerateDummyRegisterUserRequest(&profileImage.ID)

		// When
		created, _ := userService.RegisterUser(ctx, userRequest)

		// Then
		if created.Email != userRequest.Email {
			t.Errorf("got %v want %v", created.Email, userRequest.Email)
		}
	})

	t.Run("사용자의 프로필 이미지가 존재하지 않아도 생성한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()

		// Given
		userService := service.NewUserService(db, nil)
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
		if created.Email != userRequest.Email {
			t.Errorf("got %v want %v", created.Email, userRequest.Email)
		}
	})

	t.Run("사용자가 이미 존재할 경우 에러를 반환한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()

		// Given
		mediaService := service.NewMediaService(db, nil)
		profileImage := tests.AddDummyMedia(t, ctx, mediaService)

		userService := service.NewUserService(db, mediaService)
		userRequest := tests.GenerateDummyRegisterUserRequest(&profileImage.ID)

		userService.RegisterUser(ctx, userRequest)

		// When & Then
		if _, err := userService.RegisterUser(ctx, userRequest); err == nil {
			t.Errorf("got %v want %v", err, nil)
		}
	})
}

func TestFindUsers(t *testing.T) {
	t.Run("사용자를 닉네임으로 검색한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()

		// Given
		mediaService := service.NewMediaService(db, nil)
		profileImage := tests.AddDummyMedia(t, ctx, mediaService)

		userService := service.NewUserService(db, mediaService)
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
		if len(found.Items) != 1 {
			t.Errorf("got %v want %v", len(found.Items), 1)
		}
	})
}

func TestFindUser(t *testing.T) {
	t.Run("사용자를 이메일로 찾는다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()

		// Given
		mediaService := service.NewMediaService(db, nil)
		profileImage := tests.AddDummyMedia(t, ctx, mediaService)

		userService := service.NewUserService(db, mediaService)
		userRequest := tests.GenerateDummyRegisterUserRequest(&profileImage.ID)
		created, _ := userService.RegisterUser(ctx, userRequest)

		// When
		found, _ := userService.FindUser(ctx, user.FindUserParams{Email: &created.Email})

		// Then
		if found.Email != userRequest.Email {
			t.Errorf("got %v want %v", found.Email, userRequest.Email)
		}
	})

	t.Run("사용자를 UID로 찾는다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()

		// Given
		mediaService := service.NewMediaService(db, nil)
		profileImage := tests.AddDummyMedia(t, ctx, mediaService)

		userService := service.NewUserService(db, mediaService)
		userRequest := tests.GenerateDummyRegisterUserRequest(&profileImage.ID)
		created, _ := userService.RegisterUser(ctx, userRequest)

		// When
		found, _ := userService.FindUser(ctx, user.FindUserParams{FbUID: &created.FirebaseUID})

		// Then
		if found.FirebaseUID != userRequest.FirebaseUID {
			t.Errorf("got %v want %v", found.FirebaseUID, userRequest.FirebaseUID)
		}
	})

	t.Run("사용자가 존재하지 않을 경우 에러를 반환한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()

		// Given
		userService := service.NewUserService(db, nil)

		// When
		email := "non-existent@example.com"
		_, err := userService.FindUser(ctx, user.FindUserParams{Email: &email})

		// Then
		if err == nil {
			t.Errorf("got %v want %v", err, nil)
		}
	})
}

func TestExistsByEmail(t *testing.T) {
	t.Run("사용자의 닉네임이 존재하지 않을 경우 false를 반환한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()

		// Given
		userService := service.NewUserService(db, nil)

		// When
		exists, _ := userService.ExistsByNickname(ctx, "non-existent")

		// Then
		if exists {
			t.Errorf("got %v want %v", exists, false)
		}
	})

	t.Run("사용자의 닉네임이 존재할 경우 true를 반환한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()

		// Given
		mediaService := service.NewMediaService(db, nil)
		profileImage := tests.AddDummyMedia(t, ctx, mediaService)

		userService := service.NewUserService(db, mediaService)
		userRequest := tests.GenerateDummyRegisterUserRequest(&profileImage.ID)
		userService.RegisterUser(ctx, userRequest)

		// When
		exists, _ := userService.ExistsByNickname(ctx, userRequest.Nickname)

		// Then
		if !exists {
			t.Errorf("got %v want %v", exists, true)
		}
	})
}

func TestUpdateUserByUID(t *testing.T) {
	t.Run("사용자를 업데이트한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()

		// Given
		mediaService := service.NewMediaService(db, nil)

		userService := service.NewUserService(db, mediaService)
		userRequest := tests.RegisterDummyUser(t, ctx, userService, mediaService)

		// When
		updatedNickname := "updated"
		updatedProfileImage := tests.AddDummyMedia(t, ctx, mediaService)
		userService.UpdateUserByUID(ctx, userRequest.FirebaseUID, updatedNickname, &updatedProfileImage.ID)

		// Then
		found, _ := userService.FindUser(ctx, user.FindUserParams{FbUID: &userRequest.FirebaseUID})
		if found.Nickname != updatedNickname {
			t.Errorf("got %v want %v", found.Nickname, updatedNickname)
		}

		if *found.ProfileImageURL != updatedProfileImage.URL {
			t.Errorf("got %v want %v", *found.ProfileImageURL, updatedProfileImage.URL)
		}
	})
}

func TestAddPetsToOwner(t *testing.T) {
	t.Run("사용자에게 반려동물을 추가한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()

		// Given
		mediaService := service.NewMediaService(db, nil)
		userService := service.NewUserService(db, mediaService)

		owner := tests.RegisterDummyUser(t, ctx, userService, mediaService)
		profileImage := tests.AddDummyMedia(t, ctx, mediaService)
		petsToAdd := pet.AddPetsToOwnerRequest{
			Pets: []pet.AddPetRequest{
				*tests.GenerateDummyAddPetRequest(&profileImage.ID),
			},
		}

		// When
		created, _ := userService.AddPetsToOwner(ctx, owner.FirebaseUID, petsToAdd)

		// Then
		if len(created.Pets) != 1 {
			t.Errorf("got %v want %v", len(created.Pets), 1)
		}

		for _, expected := range petsToAdd.Pets {
			for _, found := range created.Pets {
				assert.PetRequestAndViewEquals(t, expected, found)
			}
		}
	})
}

func TestUpdatePet(t *testing.T) {
	t.Run("반려동물을 업데이트한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()

		// Given
		mediaService := service.NewMediaService(db, nil)
		userService := service.NewUserService(db, mediaService)
		userData := tests.RegisterDummyUser(t, ctx, userService, mediaService)

		petProfileImage := tests.AddDummyMedia(t, ctx, mediaService)
		petRequest := tests.GenerateDummyAddPetRequest(&petProfileImage.ID)
		createdPets, _ := userService.AddPetsToOwner(
			ctx, userData.FirebaseUID, pet.AddPetsToOwnerRequest{Pets: []pet.AddPetRequest{*petRequest}})
		createdPet := createdPets.Pets[0]

		// When
		updatedPetProfileImage := tests.AddDummyMedia(t, ctx, mediaService)
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

		userService.UpdatePet(ctx, userData.FirebaseUID, createdPet.ID, updatedPetRequest)

		// Then
		found, _ := userService.FindPets(ctx, pet.FindPetsParams{OwnerID: &userData.ID})
		if len(found.Pets) != 1 {
			t.Errorf("got %v want %v", len(found.Pets), 1)
		}

		assert.UpdatedPetEquals(t, updatedPetRequest, found.Pets[0])
	})
}

func TestDeletePet(t *testing.T) {
	t.Run("반려동물을 삭제한다", func(t *testing.T) {
		db, tearDown := tests.SetUp(t)
		defer tearDown(t)
		ctx := context.Background()

		// Given
		mediaService := service.NewMediaService(db, nil)
		userService := service.NewUserService(db, mediaService)
		userData := tests.RegisterDummyUser(t, ctx, userService, mediaService)

		petProfileImage := tests.AddDummyMedia(t, ctx, mediaService)
		petRequest := tests.GenerateDummyAddPetRequest(&petProfileImage.ID)
		createdPets, err := userService.AddPetsToOwner(
			ctx,
			userData.FirebaseUID,
			pet.AddPetsToOwnerRequest{Pets: []pet.AddPetRequest{*petRequest}},
		)
		if err != nil {
			t.Fatalf("got %v want %v", err, nil)
		}
		createdPet := createdPets.Pets[0]

		// When
		err = userService.DeletePet(ctx, userData.FirebaseUID, createdPet.ID)
		if err != nil {
			t.Fatalf("got %v want %v", err, nil)
		}

		// Then
		found, _ := userService.FindPets(ctx, pet.FindPetsParams{
			OwnerID: &userData.ID,
		})
		if len(found.Pets) != 0 {
			t.Fatalf("got %v want %v", len(found.Pets), 0)
		}
	})
}
