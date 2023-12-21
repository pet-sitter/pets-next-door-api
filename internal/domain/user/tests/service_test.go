package user_test

import (
	"testing"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"
	"github.com/pet-sitter/pets-next-door-api/internal/tests"
)

var db *database.DB

func setUp(t *testing.T) func(t *testing.T) {
	db, _ = database.Open(tests.TestDatabaseURL)
	db.Flush()

	return func(t *testing.T) {
		db.Close()
	}
}

func TestUserService(t *testing.T) {

	t.Run("RegisterUser", func(t *testing.T) {
		t.Run("사용자를 새로 생성한다", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			media_service := media.NewMediaService(postgres.NewMediaPostgresStore(db), nil)
			profile_image, _ := media_service.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "http://example.com",
			})

			service := user.NewUserService(postgres.NewUserPostgresStore(db), postgres.NewPetPostgresStore(db), media_service)

			user := &user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       profile_image.ID,
				FirebaseProviderType: "kakao",
				FirebaseUID:          "uid",
			}

			created, _ := service.RegisterUser(user)

			if created.Email != user.Email {
				t.Errorf("got %v want %v", created.Email, user.Email)
			}
		})

		t.Run("사용자가 이미 존재할 경우 에러를 반환한다", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			media_service := media.NewMediaService(postgres.NewMediaPostgresStore(db), nil)
			profile_image, _ := media_service.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "http://example.com",
			})
			service := user.NewUserService(postgres.NewUserPostgresStore(db), postgres.NewPetPostgresStore(db), media_service)

			user := &user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       profile_image.ID,
				FirebaseProviderType: "kakao",
				FirebaseUID:          "uid",
			}

			_, _ = service.RegisterUser(user)
			_, err := service.RegisterUser(user)

			if err == nil {
				t.Errorf("got %v want %v", err, nil)
			}
		})
	})

	t.Run("FindUsers", func(t *testing.T) {
		t.Run("사용자를 닉네임으로 검색한다", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			media_service := media.NewMediaService(postgres.NewMediaPostgresStore(db), nil)
			profile_image, _ := media_service.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "http://example.com",
			})

			service := user.NewUserService(postgres.NewUserPostgresStore(db), postgres.NewPetPostgresStore(db), media_service)

			targetNickname := "target"
			targetUserRequest := &user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             targetNickname,
				Fullname:             "fullname",
				ProfileImageID:       profile_image.ID,
				FirebaseProviderType: "kakao",
				FirebaseUID:          "uid",
			}
			service.RegisterUser(targetUserRequest)
			for i := 0; i < 2; i++ {
				service.RegisterUser(&user.RegisterUserRequest{
					Email:                "test" + string(rune(i)) + "@example.com",
					Nickname:             "nickname" + string(rune(i)),
					Fullname:             "fullname" + string(rune(i)),
					ProfileImageID:       profile_image.ID,
					FirebaseProviderType: "kakao",
					FirebaseUID:          "uid" + string(rune(i)),
				})
			}

			found, _ := service.FindUsers(1, 20, &targetNickname)

			if len(found) != 1 {
				t.Errorf("got %v want %v", len(found), 1)
			}
		})
	})

	t.Run("FindUserByEmail", func(t *testing.T) {
		t.Run("사용자를 이메일로 찾는다", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			media_service := media.NewMediaService(postgres.NewMediaPostgresStore(db), nil)
			profile_image, _ := media_service.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "http://example.com",
			})

			service := user.NewUserService(postgres.NewUserPostgresStore(db), postgres.NewPetPostgresStore(db), media_service)

			user := &user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       profile_image.ID,
				FirebaseProviderType: "kakao",
				FirebaseUID:          "uid",
			}

			created, _ := service.RegisterUser(user)

			found, err := service.FindUserByEmail(created.Email)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			if found.Email != user.Email {
				t.Errorf("got %v want %v", found.Email, user.Email)
			}
		})

		t.Run("사용자가 존재하지 않을 경우 에러를 반환한다", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			media_service := media.NewMediaService(postgres.NewMediaPostgresStore(db), nil)

			service := user.NewUserService(postgres.NewUserPostgresStore(db), postgres.NewPetPostgresStore(db), media_service)

			_, err := service.FindUserByEmail("non-existent@example.com")
			if err == nil {
				t.Errorf("got %v want %v", err, nil)
			}
		})
	})

	t.Run("FindUserByUID", func(t *testing.T) {
		t.Run("사용자를 UID로 찾는다", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			media_service := media.NewMediaService(postgres.NewMediaPostgresStore(db), nil)
			profile_image, _ := media_service.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "http://example.com",
			})

			service := user.NewUserService(postgres.NewUserPostgresStore(db), postgres.NewPetPostgresStore(db), media_service)

			user := &user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       profile_image.ID,
				FirebaseProviderType: "kakao",
				FirebaseUID:          "uid",
			}

			created, _ := service.RegisterUser(user)

			found, err := service.FindUserByUID(created.FirebaseUID)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			if found.FirebaseUID != user.FirebaseUID {
				t.Errorf("got %v want %v", found.FirebaseUID, user.FirebaseUID)
			}
		})

		t.Run("사용자가 존재하지 않을 경우 에러를 반환한다", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			media_service := media.NewMediaService(postgres.NewMediaPostgresStore(db), nil)

			service := user.NewUserService(postgres.NewUserPostgresStore(db), postgres.NewPetPostgresStore(db), media_service)

			_, err := service.FindUserByUID("non-existent")
			if err == nil {
				t.Errorf("got %v want %v", err, nil)
			}
		})
	})

	t.Run("ExistsByNickname", func(t *testing.T) {
		t.Run("사용자의 닉네임이 존재하지 않을 경우 false를 반환한다", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			media_service := media.NewMediaService(postgres.NewMediaPostgresStore(db), nil)

			service := user.NewUserService(postgres.NewUserPostgresStore(db), postgres.NewPetPostgresStore(db), media_service)

			exists, _ := service.ExistsByNickname("non-existent")
			if exists {
				t.Errorf("got %v want %v", exists, false)
			}
		})

		t.Run("사용자의 닉네임이 존재할 경우 true를 반환한다", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			media_service := media.NewMediaService(postgres.NewMediaPostgresStore(db), nil)
			profile_image, _ := media_service.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "http://example.com",
			})

			service := user.NewUserService(postgres.NewUserPostgresStore(db), postgres.NewPetPostgresStore(db), media_service)

			user := &user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       profile_image.ID,
				FirebaseProviderType: "kakao",
				FirebaseUID:          "uid",
			}

			_, _ = service.RegisterUser(user)

			exists, _ := service.ExistsByNickname(user.Nickname)

			if !exists {
				t.Errorf("got %v want %v", exists, true)
			}
		})
	})

	t.Run("FindUserStatusByEmail", func(t *testing.T) {
		t.Run("사용자의 상태를 반환한다", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			media_service := media.NewMediaService(postgres.NewMediaPostgresStore(db), nil)
			profile_image, _ := media_service.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "http://example.com",
			})

			service := user.NewUserService(postgres.NewUserPostgresStore(db), postgres.NewPetPostgresStore(db), media_service)

			user := &user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       profile_image.ID,
				FirebaseProviderType: "kakao",
				FirebaseUID:          "uid",
			}

			created, err := service.RegisterUser(user)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			status, err := service.FindUserStatusByEmail(created.Email)
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
			tearDown := setUp(t)
			defer tearDown(t)

			media_service := media.NewMediaService(postgres.NewMediaPostgresStore(db), nil)
			profile_image, _ := media_service.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "http://example.com",
			})

			service := user.NewUserService(postgres.NewUserPostgresStore(db), postgres.NewPetPostgresStore(db), media_service)

			user := &user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       profile_image.ID,
				FirebaseProviderType: "kakao",
				FirebaseUID:          "uid",
			}

			created, _ := service.RegisterUser(user)

			updatedNickname := "updated"

			_, err := service.UpdateUserByUID(created.FirebaseUID, updatedNickname, profile_image.ID)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			found, err := service.FindUserByUID(created.FirebaseUID)
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
			tearDown := setUp(t)
			defer tearDown(t)

			media_service := media.NewMediaService(postgres.NewMediaPostgresStore(db), nil)
			profile_image, _ := media_service.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "http://example.com",
			})

			service := user.NewUserService(postgres.NewUserPostgresStore(db), postgres.NewPetPostgresStore(db), media_service)

			owner, _ := service.RegisterUser(&user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       profile_image.ID,
				FirebaseProviderType: "kakao",
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
						BirthDate:  "2020-01-01T00:00:00Z",
						WeightInKg: 10.0,
					},
				},
			}

			_, err := service.AddPetsToOwner(owner.FirebaseUID, pets)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			found, err := service.FindPetsByOwnerUID(owner.FirebaseUID)
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

func assertPetEquals(t *testing.T, expected pet.AddPetRequest, found pet.PetView) {
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
