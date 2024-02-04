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
		db, _ := database.Open(tests.TestDatabaseURL)
		db.Flush()

		return db, func(t *testing.T) {
			db.Close()
		}
	}

	assertPetEquals := func(t *testing.T, expected pet.AddPetRequest, found pet.PetView) {
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

	t.Run("RegisterUser", func(t *testing.T) {
		t.Run("사용자를 새로 생성한다", func(t *testing.T) {
			db, tearDown := setUp(t)
			defer tearDown(t)
			ctx := context.Background()

			mediaService := service.NewMediaService(db, nil)
			profile_image, err := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "http://example.com",
			})
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			service := service.NewUserService(db, nil)

			user := &user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       &profile_image.ID,
				FirebaseProviderType: user.FirebaseProviderTypeKakao,
				FirebaseUID:          "uid",
			}
			created, err := service.RegisterUser(ctx, user)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			if created.Email != user.Email {
				t.Errorf("got %v want %v", created.Email, user.Email)
			}
		})

		t.Run("사용자의 프로필 이미지가 존재하지 않아도 생성한다", func(t *testing.T) {
			db, tearDown := setUp(t)
			defer tearDown(t)
			ctx := context.Background()

			service := service.NewUserService(db, nil)

			user := &user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       nil,
				FirebaseProviderType: user.FirebaseProviderTypeKakao,
				FirebaseUID:          "uid",
			}

			created, err := service.RegisterUser(ctx, user)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			if created.Email != user.Email {
				t.Errorf("got %v want %v", created.Email, user.Email)
			}
		})

		t.Run("사용자가 이미 존재할 경우 에러를 반환한다", func(t *testing.T) {
			db, tearDown := setUp(t)
			defer tearDown(t)
			ctx := context.Background()

			mediaService := service.NewMediaService(db, nil)
			profile_image, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "http://example.com",
			})
			userService := service.NewUserService(db, nil)

			user := &user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       &profile_image.ID,
				FirebaseProviderType: user.FirebaseProviderTypeKakao,
				FirebaseUID:          "uid",
			}

			userService.RegisterUser(ctx, user)
			if _, err := userService.RegisterUser(ctx, user); err == nil {
				t.Errorf("got %v want %v", err, nil)
			}
		})
	})

	t.Run("FindUsers", func(t *testing.T) {
		t.Run("사용자를 닉네임으로 검색한다", func(t *testing.T) {
			db, tearDown := setUp(t)
			defer tearDown(t)
			ctx := context.Background()

			mediaService := service.NewMediaService(db, nil)
			profileImage, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "http://example.com",
			})

			userService := service.NewUserService(db, nil)

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

			found, _ := userService.FindUsers(ctx, 1, 20, &targetNickname)
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

			mediaService := service.NewMediaService(db, nil)
			profile_image, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "http://example.com",
			})

			userService := service.NewUserService(db, nil)

			user := &user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       &profile_image.ID,
				FirebaseProviderType: user.FirebaseProviderTypeKakao,
				FirebaseUID:          "uid",
			}
			created, _ := userService.RegisterUser(ctx, user)

			found, err := userService.FindUserByEmail(ctx, created.Email)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			if found.Email != user.Email {
				t.Errorf("got %v want %v", found.Email, user.Email)
			}
		})

		t.Run("사용자가 존재하지 않을 경우 에러를 반환한다", func(t *testing.T) {
			db, tearDown := setUp(t)
			defer tearDown(t)
			ctx := context.Background()

			userService := service.NewUserService(db, nil)

			_, err := userService.FindUserByEmail(ctx, "non-existent@example.com")
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

			mediaService := service.NewMediaService(db, nil)
			profile_image, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "http://example.com",
			})

			userService := service.NewUserService(db, nil)
			user := &user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       &profile_image.ID,
				FirebaseProviderType: user.FirebaseProviderTypeKakao,
				FirebaseUID:          "uid",
			}
			created, _ := userService.RegisterUser(ctx, user)

			found, err := userService.FindUserByUID(ctx, created.FirebaseUID)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			if found.FirebaseUID != user.FirebaseUID {
				t.Errorf("got %v want %v", found.FirebaseUID, user.FirebaseUID)
			}
		})

		t.Run("사용자가 존재하지 않을 경우 에러를 반환한다", func(t *testing.T) {
			db, tearDown := setUp(t)
			defer tearDown(t)
			ctx := context.Background()

			userService := service.NewUserService(db, nil)

			_, err := userService.FindUserByUID(ctx, "non-existent")
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

			userService := service.NewUserService(db, nil)

			exists, _ := userService.ExistsByNickname(ctx, "non-existent")
			if exists {
				t.Errorf("got %v want %v", exists, false)
			}
		})

		t.Run("사용자의 닉네임이 존재할 경우 true를 반환한다", func(t *testing.T) {
			db, tearDown := setUp(t)
			defer tearDown(t)
			ctx := context.Background()

			mediaService := service.NewMediaService(db, nil)
			profile_image, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "http://example.com",
			})

			userService := service.NewUserService(db, nil)

			user := &user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       &profile_image.ID,
				FirebaseProviderType: user.FirebaseProviderTypeKakao,
				FirebaseUID:          "uid",
			}
			userService.RegisterUser(ctx, user)

			exists, _ := userService.ExistsByNickname(ctx, user.Nickname)
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

			mediaService := service.NewMediaService(db, nil)
			profile_image, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "http://example.com",
			})

			userService := service.NewUserService(db, nil)

			user := &user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       &profile_image.ID,
				FirebaseProviderType: user.FirebaseProviderTypeKakao,
				FirebaseUID:          "uid",
			}
			created, err := userService.RegisterUser(ctx, user)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			status, err := userService.FindUserStatusByEmail(ctx, created.Email)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			if status.FirebaseProviderType != user.FirebaseProviderType {
				t.Errorf("got %v want %v", status.FirebaseProviderType, user.FirebaseProviderType)
			}
		})
	})

	t.Run("UpdateUserByUID", func(t *testing.T) {
		t.Run("사용자를 업데이트한다", func(t *testing.T) {
			db, tearDown := setUp(t)
			defer tearDown(t)
			ctx := context.Background()

			mediaService := service.NewMediaService(db, nil)
			profile_image, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "http://example.com",
			})

			userService := service.NewUserService(db, nil)

			user := &user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       &profile_image.ID,
				FirebaseProviderType: user.FirebaseProviderTypeKakao,
				FirebaseUID:          "uid",
			}
			created, _ := userService.RegisterUser(ctx, user)

			updatedNickname := "updated"
			_, err := userService.UpdateUserByUID(ctx, created.FirebaseUID, updatedNickname, &profile_image.ID)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			found, err := userService.FindUserByUID(ctx, created.FirebaseUID)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			if found.Nickname != updatedNickname {
				t.Errorf("got %v want %v", found.Nickname, updatedNickname)
			}
		})
	})

	t.Run("AddPetsToOwner", func(t *testing.T) {
		t.Run("사용자에게 반려동물을 추가한다", func(t *testing.T) {
			db, tearDown := setUp(t)
			defer tearDown(t)
			ctx := context.Background()

			mediaService := service.NewMediaService(db, nil)
			profile_image, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "http://example.com",
			})

			userService := service.NewUserService(db, nil)
			owner, _ := userService.RegisterUser(ctx, &user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       &profile_image.ID,
				FirebaseProviderType: user.FirebaseProviderTypeKakao,
				FirebaseUID:          "uid",
			})

			pets := pet.AddPetsToOwnerRequest{
				Pets: []pet.AddPetRequest{
					{
						Name:       "name",
						PetType:    "dog",
						Sex:        "male",
						Neutered:   true,
						Breed:      "poodle",
						BirthDate:  "2020-01-01",
						WeightInKg: 10.0,
					},
				},
			}

			_, err := userService.AddPetsToOwner(ctx, owner.FirebaseUID, pets)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			found, err := userService.FindPetsByOwnerUID(ctx, owner.FirebaseUID)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			if len(found.Pets) != 1 {
				t.Errorf("got %v want %v", len(found.Pets), 1)
			}

			for _, expected := range pets.Pets {
				for _, found := range found.Pets {
					assertPetEquals(t, expected, found)
				}
			}
		})
	})
}
