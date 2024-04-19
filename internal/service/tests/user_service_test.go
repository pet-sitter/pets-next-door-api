package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
	"github.com/pet-sitter/pets-next-door-api/internal/tests"
)

func TestUserService(t *testing.T) {
	setUp := func(t *testing.T) (*database.DB, func(t *testing.T)) {
		t.Helper()

		db, _ := database.Open(tests.TestDatabaseURL)
		db.Flush()

		return db, func(t *testing.T) {
			t.Helper()

			db.Close()
		}
	}

	t.Run("RegisterUser", func(t *testing.T) {
		t.Run("사용자를 새로 생성한다", func(t *testing.T) {
			db, tearDown := setUp(t)
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
			db, tearDown := setUp(t)
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
			db, tearDown := setUp(t)
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
	})

	t.Run("FindUsers", func(t *testing.T) {
		t.Run("사용자를 닉네임으로 검색한다", func(t *testing.T) {
			db, tearDown := setUp(t)
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
			found, _ := userService.FindUsers(ctx, 1, 20, &targetNickname)

			// Then
			if len(found.Items) != 1 {
				t.Errorf("got %v want %v", len(found.Items), 1)
			}
		})
	})

	t.Run("FindUserByEmail", func(t *testing.T) {
		t.Run("사용자를 이메일로 찾는다", func(t *testing.T) {
			db, tearDown := setUp(t)
			defer tearDown(t)
			ctx := context.Background()

			// Given
			mediaService := service.NewMediaService(db, nil)
			profileImage := tests.AddDummyMedia(t, ctx, mediaService)

			userService := service.NewUserService(db, mediaService)
			userRequest := tests.GenerateDummyRegisterUserRequest(&profileImage.ID)
			created, _ := userService.RegisterUser(ctx, userRequest)

			// When
			found, _ := userService.FindUserByEmail(ctx, created.Email)

			// Then
			if found.Email != userRequest.Email {
				t.Errorf("got %v want %v", found.Email, userRequest.Email)
			}
		})

		t.Run("사용자가 존재하지 않을 경우 에러를 반환한다", func(t *testing.T) {
			db, tearDown := setUp(t)
			defer tearDown(t)
			ctx := context.Background()

			// Given
			userService := service.NewUserService(db, nil)

			// When
			_, err := userService.FindUserByEmail(ctx, "non-existent@example.com")

			// Then
			if err == nil {
				t.Errorf("got %v want %v", err, nil)
			}
		})
	})

	t.Run("FindUserByUID", func(t *testing.T) {
		t.Run("사용자를 UID로 찾는다", func(t *testing.T) {
			db, tearDown := setUp(t)
			defer tearDown(t)
			ctx := context.Background()

			// Given
			mediaService := service.NewMediaService(db, nil)
			profileImage := tests.AddDummyMedia(t, ctx, mediaService)

			userService := service.NewUserService(db, mediaService)
			userRequest := tests.GenerateDummyRegisterUserRequest(&profileImage.ID)
			created, _ := userService.RegisterUser(ctx, userRequest)

			// When
			found, _ := userService.FindUserByUID(ctx, created.FirebaseUID)

			// Then
			if found.FirebaseUID != userRequest.FirebaseUID {
				t.Errorf("got %v want %v", found.FirebaseUID, userRequest.FirebaseUID)
			}
		})

		t.Run("사용자가 존재하지 않을 경우 에러를 반환한다", func(t *testing.T) {
			db, tearDown := setUp(t)
			defer tearDown(t)
			ctx := context.Background()

			// Given
			userService := service.NewUserService(db, nil)

			// When
			_, err := userService.FindUserByUID(ctx, "non-existent")

			// Then
			if err == nil {
				t.Errorf("got %v want %v", err, nil)
			}
		})
	})

	t.Run("ExistsByNickname", func(t *testing.T) {
		t.Run("사용자의 닉네임이 존재하지 않을 경우 false를 반환한다", func(t *testing.T) {
			db, tearDown := setUp(t)
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
			db, tearDown := setUp(t)
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
	})

	t.Run("FindUserStatusByEmail", func(t *testing.T) {
		t.Run("사용자의 상태를 반환한다", func(t *testing.T) {
			db, tearDown := setUp(t)
			defer tearDown(t)
			ctx := context.Background()

			// Given
			mediaService := service.NewMediaService(db, nil)
			profileImage, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "http://example.com",
			})

			userService := service.NewUserService(db, mediaService)
			userRequest := tests.GenerateDummyRegisterUserRequest(&profileImage.ID)
			created, _ := userService.RegisterUser(ctx, userRequest)

			// When
			status, _ := userService.FindUserStatusByEmail(ctx, created.Email)

			// Then
			if status.FirebaseProviderType != userRequest.FirebaseProviderType {
				t.Errorf("got %v want %v", status.FirebaseProviderType, userRequest.FirebaseProviderType)
			}
		})
	})

	t.Run("UpdateUserByUID", func(t *testing.T) {
		t.Run("사용자를 업데이트한다", func(t *testing.T) {
			db, tearDown := setUp(t)
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
			found, _ := userService.FindUserByUID(ctx, userRequest.FirebaseUID)
			if found.Nickname != updatedNickname {
				t.Errorf("got %v want %v", found.Nickname, updatedNickname)
			}

			if *found.ProfileImageURL != updatedProfileImage.URL {
				t.Errorf("got %v want %v", *found.ProfileImageURL, updatedProfileImage.URL)
			}
		})
	})

	t.Run("AddPetsToOwner", func(t *testing.T) {
		t.Run("사용자에게 반려동물을 추가한다", func(t *testing.T) {
			db, tearDown := setUp(t)
			defer tearDown(t)
			ctx := context.Background()

			// Given
			mediaService := service.NewMediaService(db, nil)
			userService := service.NewUserService(db, mediaService)

			owner := tests.RegisterDummyUser(t, ctx, userService, mediaService)
			profileImage := tests.AddDummyMedia(t, ctx, mediaService)
			pets := pet.AddPetsToOwnerRequest{Pets: []pet.AddPetRequest{*tests.GenerateDummyAddPetRequest(&profileImage.ID)}}

			// When
			userService.AddPetsToOwner(ctx, owner.FirebaseUID, pets)

			// Then
			found, _ := userService.FindPetsByOwnerUID(ctx, owner.FirebaseUID)
			if len(found.Pets) != 1 {
				t.Errorf("got %v want %v", len(found.Pets), 1)
			}

			for _, expected := range pets.Pets {
				for _, found := range found.Pets {
					assertPetRequestAndViewEquals(t, expected, found)
				}
			}
		})
	})

	t.Run("UpdatePet", func(t *testing.T) {
		t.Run("반려동물을 업데이트한다", func(t *testing.T) {
			db, tearDown := setUp(t)
			defer tearDown(t)
			ctx := context.Background()

			// Given
			mediaService := service.NewMediaService(db, nil)
			userService := service.NewUserService(db, mediaService)
			userRequest := tests.RegisterDummyUser(t, ctx, userService, mediaService)

			petProfileImage := tests.AddDummyMedia(t, ctx, mediaService)
			petRequest := tests.GenerateDummyAddPetRequest(&petProfileImage.ID)
			createdPets, _ := userService.AddPetsToOwner(
				ctx, userRequest.FirebaseUID, pet.AddPetsToOwnerRequest{Pets: []pet.AddPetRequest{*petRequest}})
			createdPet := createdPets[0]

			// When
			updatedPetProfileImage := tests.AddDummyMedia(t, ctx, mediaService)
			updatedPetRequest := pet.UpdatePetRequest{
				Name:           "updated",
				Neutered:       true,
				Breed:          "updated",
				BirthDate:      "2021-01-01",
				WeightInKg:     10.0,
				Remarks:        "updated",
				ProfileImageID: &updatedPetProfileImage.ID,
			}

			userService.UpdatePet(ctx, userRequest.FirebaseUID, createdPet.ID, updatedPetRequest)

			// Then
			found, _ := userService.FindPetsByOwnerUID(ctx, userRequest.FirebaseUID)
			if len(found.Pets) != 1 {
				t.Errorf("got %v want %v", len(found.Pets), 1)
			}

			assertUpdatedPetEquals(t, updatedPetRequest, found.Pets[0])
		})
	})

	t.Run("DeletePet", func(t *testing.T) {
		t.Run("반려동물을 삭제한다", func(t *testing.T) {
			db, tearDown := setUp(t)
			defer tearDown(t)
			ctx := context.Background()

			// Given
			mediaService := service.NewMediaService(db, nil)
			userService := service.NewUserService(db, mediaService)
			userRequest := tests.RegisterDummyUser(t, ctx, userService, mediaService)

			petProfileImage := tests.AddDummyMedia(t, ctx, mediaService)
			petRequest := tests.GenerateDummyAddPetRequest(&petProfileImage.ID)
			createdPets, _ := userService.AddPetsToOwner(
				ctx,
				userRequest.FirebaseUID,
				pet.AddPetsToOwnerRequest{Pets: []pet.AddPetRequest{*petRequest}},
			)
			createdPet := createdPets[0]

			// When
			userService.DeletePet(ctx, userRequest.FirebaseUID, createdPet.ID)

			// Then
			found, _ := userService.FindPetsByOwnerUID(ctx, userRequest.FirebaseUID)
			if len(found.Pets) != 0 {
				t.Errorf("got %v want %v", len(found.Pets), 0)
			}
		})
	})
}

func assertPetRequestAndViewEquals(t *testing.T, expected pet.AddPetRequest, found pet.PetView) {
	t.Helper()

	if expected.Name != found.Name {
		t.Errorf("got %v want %v", expected.Name, found.Name)
	}

	if expected.PetType != found.PetType {
		t.Errorf("got %v want %v", expected.PetType, found.PetType)
	}

	if expected.Sex != found.Sex {
		t.Errorf("got %v want %v", expected.Sex, found.PetType)
	}

	if expected.Neutered != found.Neutered {
		t.Errorf("got %v want %v", expected.Neutered, found.Neutered)
	}

	if expected.Breed != found.Breed {
		t.Errorf("got %v want %v", expected.Breed, found.Breed)
	}

	if expected.BirthDate != found.BirthDate {
		t.Errorf("got %v want %v", expected.BirthDate, found.BirthDate)
	}

	if expected.WeightInKg != found.WeightInKg {
		t.Errorf("got %v want %v", expected.WeightInKg, found.WeightInKg)
	}
}

func assertUpdatedPetEquals(t *testing.T, expected pet.UpdatePetRequest, found pet.PetView) {
	t.Helper()

	if expected.Name != found.Name {
		t.Errorf("got %v want %v", expected.Name, found.Name)
	}

	if expected.Neutered != found.Neutered {
		t.Errorf("got %v want %v", expected.Neutered, found.Neutered)
	}

	if expected.Breed != found.Breed {
		t.Errorf("got %v want %v", expected.Breed, found.Breed)
	}

	if expected.BirthDate != found.BirthDate {
		t.Errorf("got %v want %v", expected.BirthDate, found.BirthDate)
	}

	if expected.WeightInKg != found.WeightInKg {
		t.Errorf("got %v want %v", expected.WeightInKg, found.WeightInKg)
	}
}
