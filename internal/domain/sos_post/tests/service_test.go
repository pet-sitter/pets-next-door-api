package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sos_post"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"
	"github.com/pet-sitter/pets-next-door-api/internal/tests"
)

var db *database.DB

func setUp(t *testing.T) func(t *testing.T) {
	db, _ = database.Open(tests.TestDatabaseURL)
	db.Flush()
	postgres.NewConditionPostgresStore(db).InitConditions(sos_post.ConditionName)
	return func(t *testing.T) {
		db.Close()
	}
}

func TestSosPostService(t *testing.T) {

	t.Run("CreateSosPost", func(t *testing.T) {
		t.Run("돌봄 급구 게시글을 작성합니다.", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			mediaService := media.NewMediaService(postgres.NewMediaPostgresStore(db), nil)
			profileImage, _ := mediaService.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test.com",
			})
			sosPostImage, _ := mediaService.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test2.com",
			})
			sosPostImage2, _ := mediaService.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test3.com",
			})
			sosPostMedia := []media.MediaView{
				{
					ID:        sosPostImage.ID,
					MediaType: sosPostImage.MediaType,
					URL:       sosPostImage.URL,
					CreatedAt: sosPostImage.CreatedAt,
				},
				{
					ID:        sosPostImage2.ID,
					MediaType: sosPostImage2.MediaType,
					URL:       sosPostImage2.URL,
					CreatedAt: sosPostImage2.CreatedAt,
				},
			}

			userService := user.NewUserService(postgres.NewUserPostgresStore(db), postgres.NewPetPostgresStore(db), *mediaService)

			owner, err := userService.RegisterUser(&user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       &profileImage.ID,
				FirebaseProviderType: user.FirebaseProviderTypeKakao,
				FirebaseUID:          "1234",
			})
			if err != nil {
				t.Errorf("RegisterUser failed: %v", err)
				return
			}

			uid := owner.FirebaseUID

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

			addPets, err := userService.AddPetsToOwner(uid, pets)
			if err != nil {
				t.Errorf(err.Err.Error())
			}

			sosPostService := sos_post.NewSosPostService(postgres.NewSosPostPostgresStore(db), postgres.NewResourceMediaPostgresStore(db), postgres.NewUserPostgresStore(db))

			conditionIDs := []int{1, 2}
			krLocation, _ := time.LoadLocation("Asia/Seoul")

			writeSosPostRequest := &sos_post.WriteSosPostRequest{
				Title:        "Test Title",
				Content:      "Test Content",
				ImageIDs:     []int{sosPostImage.ID, sosPostImage2.ID},
				Reward:       "Test Reward",
				DateStartAt:  time.Date(2023, time.December, 18, 8, 00, 0, 0, krLocation),
				DateEndAt:    time.Date(2023, time.December, 20, 18, 00, 0, 0, krLocation),
				CareType:     sos_post.CareTypeFoster,
				CarerGender:  sos_post.CarerGenderMale,
				RewardAmount: sos_post.RewardAmountHour,
				ConditionIDs: conditionIDs,
				PetIDs:       []int{addPets[0].ID},
			}

			sosPost, err := sosPostService.WriteSosPost(uid, writeSosPostRequest)
			if err != nil {
				t.Errorf(err.Err.Error())
			}

			assertConditionEquals(t, sosPost.Conditions, conditionIDs)
			assertPetEquals(t, sosPost.Pets[0], addPets[0])
			assertMediaEquals(t, sosPost.Media, sosPostMedia)

			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}
			if sosPost.Title != writeSosPostRequest.Title {
				t.Errorf("got %v want %v", sosPost.Title, writeSosPostRequest.Title)
			}
			if sosPost.Content != writeSosPostRequest.Content {
				t.Errorf("got %v want %v", sosPost.Content, writeSosPostRequest.Content)
			}
			if sosPost.Reward != writeSosPostRequest.Reward {
				t.Errorf("got %v want %v", sosPost.Reward, writeSosPostRequest.Reward)
			}
			if sosPost.DateStartAt != "2023-12-17T00:00:00Z" {
				t.Errorf("got %v want %v", sosPost.DateStartAt, writeSosPostRequest.DateStartAt)
			}
			if sosPost.DateEndAt != "2023-12-20T00:00:00Z" {
				t.Errorf("got %v want %v", sosPost.DateEndAt, writeSosPostRequest.DateEndAt)
			}
			if sosPost.CareType != sos_post.CareTypeFoster {
				t.Errorf("got %v want %v", sosPost.CareType, sos_post.CareTypeFoster)
			}
			if sosPost.CarerGender != sos_post.CarerGenderMale {
				t.Errorf("got %v want %v", sosPost.CarerGender, sos_post.CarerGenderMale)
			}
			if sosPost.RewardAmount != sos_post.RewardAmountHour {
				t.Errorf("got %v want %v", sosPost.RewardAmount, sos_post.RewardAmountHour)
			}
			if sosPost.ThumbnailID != sosPostImage.ID {
				t.Errorf("got %v want %v", sosPost.ThumbnailID, sosPostImage.ID)
			}
			if sosPost.AuthorID != owner.ID {
				t.Errorf("got %v want %v", sosPost.AuthorID, owner.ID)
			}
		})
	})

	t.Run("FindSosPosts", func(t *testing.T) {
		t.Run("전체 돌봄 급구 게시글을 조회합니다.", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			mediaService := media.NewMediaService(postgres.NewMediaPostgresStore(db), nil)
			profileImage, _ := mediaService.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test.com",
			})
			sosPostImage, _ := mediaService.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test2.com",
			})
			sosPostImage2, _ := mediaService.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test3.com",
			})
			sosPostMedia := []media.MediaView{
				{
					ID:        sosPostImage.ID,
					MediaType: sosPostImage.MediaType,
					URL:       sosPostImage.URL,
					CreatedAt: sosPostImage.CreatedAt,
				},
				{
					ID:        sosPostImage2.ID,
					MediaType: sosPostImage2.MediaType,
					URL:       sosPostImage2.URL,
					CreatedAt: sosPostImage2.CreatedAt,
				},
			}

			userService := user.NewUserService(postgres.NewUserPostgresStore(db), postgres.NewPetPostgresStore(db), *mediaService)

			owner, err := userService.RegisterUser(&user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       &profileImage.ID,
				FirebaseProviderType: user.FirebaseProviderTypeKakao,
				FirebaseUID:          "1234",
			})
			if err != nil {
				t.Errorf("RegisterUser failed: %v", err)
				return
			}

			uid := owner.FirebaseUID

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

			addPets, err := userService.AddPetsToOwner(uid, pets)
			if err != nil {
				t.Errorf(err.Err.Error())
			}

			sosPostService := sos_post.NewSosPostService(postgres.NewSosPostPostgresStore(db), postgres.NewResourceMediaPostgresStore(db), postgres.NewUserPostgresStore(db))

			conditionIDs := []int{1, 2}
			krLocation, _ := time.LoadLocation("Asia/Seoul")

			var sosPosts []sos_post.WriteSosPostView

			for i := 1; i < 4; i++ {
				sosPost, err := sosPostService.WriteSosPost(uid, &sos_post.WriteSosPostRequest{
					Title:        fmt.Sprintf("Title%d", i),
					Content:      fmt.Sprintf("Test Content%d", i),
					ImageIDs:     []int{sosPostImage.ID, sosPostImage2.ID},
					Reward:       fmt.Sprintf("Test Reward%d", i),
					DateStartAt:  time.Date(2023, time.December, i, 8, 00, 0, 0, krLocation),
					DateEndAt:    time.Date(2023, time.December, i, 18, 00, 0, 0, krLocation),
					CareType:     sos_post.CareTypeFoster,
					CarerGender:  sos_post.CarerGenderMale,
					RewardAmount: sos_post.RewardAmountHour,
					ConditionIDs: conditionIDs,
					PetIDs:       []int{addPets[0].ID},
				})

				if err != nil {
					t.Errorf(err.Err.Error())
				}

				sosPosts = append(sosPosts, *sosPost)
			}

			sosPostList, err := sosPostService.FindSosPosts(1, 3, "newest")
			for i, sosPost := range sosPostList.Items {
				assertConditionEquals(t, sosPost.Conditions, conditionIDs)
				assertPetEquals(t, sosPost.Pets[0], addPets[0])
				assertMediaEquals(t, sosPost.Media, sosPostMedia)

				idx := len(sosPostList.Items) - i - 1

				if err != nil {
					t.Errorf("got %v want %v", err, nil)
				}
				if sosPost.Title != sosPosts[idx].Title {
					t.Errorf("got %v want %v", sosPost.Title, sosPosts[idx].Title)
				}
				if sosPost.Content != sosPosts[idx].Content {
					t.Errorf("got %v want %v", sosPost.Content, sosPosts[idx].Content)
				}
				if sosPost.Reward != sosPosts[idx].Reward {
					t.Errorf("got %v want %v", sosPost.Reward, sosPosts[idx].Reward)
				}
				if sosPost.DateStartAt != sosPosts[idx].DateStartAt {
					t.Errorf("got %v want %v", sosPost.DateStartAt, sosPosts[idx].DateStartAt)
				}
				if sosPost.DateEndAt != sosPosts[idx].DateEndAt {
					t.Errorf("got %v want %v", sosPost.DateEndAt, sosPosts[idx].DateEndAt)
				}
				if sosPost.CareType != sosPosts[idx].CareType {
					t.Errorf("got %v want %v", sosPost.CareType, sosPosts[idx].CareType)
				}
				if sosPost.CarerGender != sosPosts[idx].CarerGender {
					t.Errorf("got %v want %v", sosPost.CarerGender, sosPosts[idx].CarerGender)
				}
				if sosPost.RewardAmount != sosPosts[idx].RewardAmount {
					t.Errorf("got %v want %v", sosPost.RewardAmount, sosPosts[idx].RewardAmount)
				}
				if sosPost.ThumbnailID != sosPostImage.ID {
					t.Errorf("got %v want %v", sosPost.ThumbnailID, sosPostImage.ID)
				}
				if sosPost.AuthorID != owner.ID {
					t.Errorf("got %v want %v", sosPost.AuthorID, owner.ID)
				}
			}
		})
		t.Run("작성자 ID로 돌봄 급구 게시글을 조회합니다.", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			mediaService := media.NewMediaService(postgres.NewMediaPostgresStore(db), nil)
			profileImage, _ := mediaService.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test.com",
			})
			sosPostImage, _ := mediaService.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test2.com",
			})
			sosPostImage2, _ := mediaService.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test3.com",
			})
			sosPostMedia := []media.MediaView{
				{
					ID:        sosPostImage.ID,
					MediaType: sosPostImage.MediaType,
					URL:       sosPostImage.URL,
					CreatedAt: sosPostImage.CreatedAt,
				},
				{
					ID:        sosPostImage2.ID,
					MediaType: sosPostImage2.MediaType,
					URL:       sosPostImage2.URL,
					CreatedAt: sosPostImage2.CreatedAt,
				},
			}

			userService := user.NewUserService(postgres.NewUserPostgresStore(db), postgres.NewPetPostgresStore(db), *mediaService)

			owner, err := userService.RegisterUser(&user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       &profileImage.ID,
				FirebaseProviderType: user.FirebaseProviderTypeKakao,
				FirebaseUID:          "1234",
			})
			if err != nil {
				t.Errorf("RegisterUser failed: %v", err)
				return
			}

			uid := owner.FirebaseUID
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

			addPets, err := userService.AddPetsToOwner(uid, pets)
			if err != nil {
				t.Errorf(err.Err.Error())
			}

			sosPostService := sos_post.NewSosPostService(postgres.NewSosPostPostgresStore(db), postgres.NewResourceMediaPostgresStore(db), postgres.NewUserPostgresStore(db))

			conditionIDs := []int{1, 2}
			krLocation, _ := time.LoadLocation("Asia/Seoul")

			sosPosts := make([]sos_post.WriteSosPostView, 0)
			for i := 1; i < 4; i++ {
				sosPost, err := sosPostService.WriteSosPost(uid, &sos_post.WriteSosPostRequest{
					Title:        fmt.Sprintf("Title%d", i),
					Content:      fmt.Sprintf("Test Content%d", i),
					ImageIDs:     []int{sosPostImage.ID, sosPostImage2.ID},
					Reward:       fmt.Sprintf("Test Reward%d", i),
					DateStartAt:  time.Date(2023, time.December, i, 8, 00, 0, 0, krLocation),
					DateEndAt:    time.Date(2023, time.December, i, 18, 00, 0, 0, krLocation),
					CareType:     sos_post.CareTypeFoster,
					CarerGender:  sos_post.CarerGenderMale,
					RewardAmount: sos_post.RewardAmountHour,
					ConditionIDs: conditionIDs,
					PetIDs:       []int{addPets[0].ID},
				})

				if err != nil {
					t.Errorf(err.Err.Error())
				}

				sosPosts = append(sosPosts, *sosPost)
			}

			sosPostListByAuthorID, err := sosPostService.FindSosPostsByAuthorID(owner.ID, 1, 3, "newest")
			for i, sosPost := range sosPostListByAuthorID.Items {
				assertConditionEquals(t, sosPost.Conditions, conditionIDs)
				assertPetEquals(t, sosPost.Pets[0], addPets[0])
				assertMediaEquals(t, sosPost.Media, sosPostMedia)

				idx := len(sosPostListByAuthorID.Items) - i - 1

				if err != nil {
					t.Errorf("got %v want %v", err, nil)
				}
				if sosPost.Title != sosPosts[idx].Title {
					t.Errorf("got %v want %v", sosPost.Title, sosPosts[idx].Title)
				}
				if sosPost.Content != sosPosts[idx].Content {
					t.Errorf("got %v want %v", sosPost.Content, sosPosts[idx].Content)
				}
				if sosPost.Reward != sosPosts[idx].Reward {
					t.Errorf("got %v want %v", sosPost.Reward, sosPosts[idx].Reward)
				}
				if sosPost.DateStartAt != sosPosts[idx].DateStartAt {
					t.Errorf("got %v want %v", sosPost.DateStartAt, sosPosts[idx].DateStartAt)
				}
				if sosPost.DateEndAt != sosPosts[idx].DateEndAt {
					t.Errorf("got %v want %v", sosPost.DateEndAt, sosPosts[idx].DateEndAt)
				}
				if sosPost.CareType != sosPosts[idx].CareType {
					t.Errorf("got %v want %v", sosPost.CareType, sosPosts[idx].CareType)
				}
				if sosPost.CarerGender != sosPosts[idx].CarerGender {
					t.Errorf("got %v want %v", sosPost.CarerGender, sosPosts[idx].CarerGender)
				}
				if sosPost.RewardAmount != sosPosts[idx].RewardAmount {
					t.Errorf("got %v want %v", sosPost.RewardAmount, sosPosts[idx].RewardAmount)
				}
				if sosPost.ThumbnailID != sosPostImage.ID {
					t.Errorf("got %v want %v", sosPost.ThumbnailID, sosPostImage.ID)
				}
				if sosPost.AuthorID != owner.ID {
					t.Errorf("got %v want %v", sosPost.AuthorID, owner.ID)
				}
			}
		})
	})

	t.Run("FindSosPostByID", func(t *testing.T) {
		t.Run("게시글 ID로 돌봄 급구 게시글을 조회합니다.", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			mediaService := media.NewMediaService(postgres.NewMediaPostgresStore(db), nil)
			profileImage, _ := mediaService.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test.com",
			})
			sosPostImage, _ := mediaService.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test2.com",
			})
			sosPostImage2, _ := mediaService.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test3.com",
			})
			sosPostMedia := []media.MediaView{
				{
					ID:        sosPostImage.ID,
					MediaType: sosPostImage.MediaType,
					URL:       sosPostImage.URL,
					CreatedAt: sosPostImage.CreatedAt,
				},
				{
					ID:        sosPostImage2.ID,
					MediaType: sosPostImage2.MediaType,
					URL:       sosPostImage2.URL,
					CreatedAt: sosPostImage2.CreatedAt,
				},
			}

			userService := user.NewUserService(postgres.NewUserPostgresStore(db), postgres.NewPetPostgresStore(db), *mediaService)

			owner, err := userService.RegisterUser(&user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       &profileImage.ID,
				FirebaseProviderType: user.FirebaseProviderTypeKakao,
				FirebaseUID:          "1234",
			})
			if err != nil {
				t.Errorf("RegisterUser failed: %v", err)
				return
			}

			uid := owner.FirebaseUID

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

			addPets, err := userService.AddPetsToOwner(uid, pets)
			if err != nil {
				t.Errorf(err.Err.Error())
			}

			sosPostService := sos_post.NewSosPostService(postgres.NewSosPostPostgresStore(db), postgres.NewResourceMediaPostgresStore(db), postgres.NewUserPostgresStore(db))

			conditionIDs := []int{1, 2}
			krLocation, _ := time.LoadLocation("Asia/Seoul")

			sosPosts := make([]sos_post.WriteSosPostView, 0)
			for i := 1; i < 4; i++ {
				sosPost, err := sosPostService.WriteSosPost(uid, &sos_post.WriteSosPostRequest{
					Title:        fmt.Sprintf("Title%d", i),
					Content:      fmt.Sprintf("Test Content%d", i),
					ImageIDs:     []int{sosPostImage.ID, sosPostImage2.ID},
					Reward:       fmt.Sprintf("Test Reward%d", i),
					DateStartAt:  time.Date(2023, time.December, i, 8, 00, 0, 0, krLocation),
					DateEndAt:    time.Date(2023, time.December, i, 18, 00, 0, 0, krLocation),
					CareType:     sos_post.CareTypeFoster,
					CarerGender:  sos_post.CarerGenderMale,
					RewardAmount: sos_post.RewardAmountHour,
					ConditionIDs: conditionIDs,
					PetIDs:       []int{addPets[0].ID},
				})

				if err != nil {
					t.Errorf(err.Err.Error())
				}

				sosPosts = append(sosPosts, *sosPost)
			}

			findSosPostByID, err := sosPostService.FindSosPostByID(sosPosts[0].ID)

			assertConditionEquals(t, sosPosts[0].Conditions, conditionIDs)
			assertPetEquals(t, sosPosts[0].Pets[0], addPets[0])
			assertMediaEquals(t, findSosPostByID.Media, sosPostMedia)

			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}
			if findSosPostByID.Title != sosPosts[0].Title {
				t.Errorf("got %v want %v", findSosPostByID.Title, sosPosts[0].Title)
			}
			if findSosPostByID.Content != sosPosts[0].Content {
				t.Errorf("got %v want %v", findSosPostByID.Content, sosPosts[0].Content)
			}
			if findSosPostByID.Reward != sosPosts[0].Reward {
				t.Errorf("got %v want %v", findSosPostByID.Reward, sosPosts[0].Reward)
			}
			if findSosPostByID.DateStartAt != sosPosts[0].DateStartAt {
				t.Errorf("got %v want %v", findSosPostByID.DateStartAt, sosPosts[0].DateStartAt)
			}
			if findSosPostByID.DateEndAt != sosPosts[0].DateEndAt {
				t.Errorf("got %v want %v", findSosPostByID.DateEndAt, sosPosts[0].DateEndAt)
			}
			if findSosPostByID.CareType != sosPosts[0].CareType {
				t.Errorf("got %v want %v", findSosPostByID.CareType, sosPosts[0].CareType)
			}
			if findSosPostByID.CarerGender != sosPosts[0].CarerGender {
				t.Errorf("got %v want %v", findSosPostByID.CarerGender, sosPosts[0].CarerGender)
			}
			if findSosPostByID.RewardAmount != sosPosts[0].RewardAmount {
				t.Errorf("got %v want %v", findSosPostByID.RewardAmount, sosPosts[0].RewardAmount)
			}
			if findSosPostByID.ThumbnailID != sosPostImage.ID {
				t.Errorf("got %v want %v", findSosPostByID.ThumbnailID, sosPostImage.ID)
			}
			if findSosPostByID.AuthorID != owner.ID {
				t.Errorf("got %v want %v", findSosPostByID.AuthorID, owner.ID)
			}
		})
	})

	t.Run("UpdateSosPost", func(t *testing.T) {
		t.Run("돌봄 급구 게시글을 수정합니다.", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			mediaService := media.NewMediaService(postgres.NewMediaPostgresStore(db), nil)
			profileImage, _ := mediaService.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test.com",
			})
			sosPostImage, _ := mediaService.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test2.com",
			})
			sosPostImage2, _ := mediaService.CreateMedia(&media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test3.com",
			})
			sosPostMedia := []media.MediaView{
				{
					ID:        sosPostImage.ID,
					MediaType: sosPostImage.MediaType,
					URL:       sosPostImage.URL,
					CreatedAt: sosPostImage.CreatedAt,
				},
				{
					ID:        sosPostImage2.ID,
					MediaType: sosPostImage2.MediaType,
					URL:       sosPostImage2.URL,
					CreatedAt: sosPostImage2.CreatedAt,
				},
			}

			userService := user.NewUserService(postgres.NewUserPostgresStore(db), postgres.NewPetPostgresStore(db), *mediaService)

			owner, err := userService.RegisterUser(&user.RegisterUserRequest{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				ProfileImageID:       &profileImage.ID,
				FirebaseProviderType: user.FirebaseProviderTypeKakao,
				FirebaseUID:          "1234",
			})
			if err != nil {
				t.Errorf("RegisterUser failed: %v", err)
				return
			}

			uid := owner.FirebaseUID

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

			addPets, err := userService.AddPetsToOwner(uid, pets)
			if err != nil {
				t.Errorf(err.Err.Error())
			}

			sosPostService := sos_post.NewSosPostService(postgres.NewSosPostPostgresStore(db), postgres.NewResourceMediaPostgresStore(db), postgres.NewUserPostgresStore(db))

			conditionIDs := []int{1, 2}
			krLocation, _ := time.LoadLocation("Asia/Seoul")

			sosPost, err := sosPostService.WriteSosPost(uid, &sos_post.WriteSosPostRequest{
				Title:        "Title1",
				Content:      "Test Content1",
				ImageIDs:     []int{sosPostImage.ID, sosPostImage2.ID},
				Reward:       "Test Reward1",
				DateStartAt:  time.Date(2023, time.December, 0, 8, 00, 0, 0, krLocation),
				DateEndAt:    time.Date(2023, time.December, 0, 18, 00, 0, 0, krLocation),
				CareType:     sos_post.CareTypeFoster,
				CarerGender:  sos_post.CarerGenderMale,
				RewardAmount: sos_post.RewardAmountHour,
				ConditionIDs: conditionIDs,
				PetIDs:       []int{addPets[0].ID},
			})

			if err != nil {
				t.Errorf(err.Err.Error())
			}

			updateSosPostData := &sos_post.UpdateSosPostRequest{
				ID:           sosPost.ID,
				Title:        "Title2",
				Content:      "Test Content2",
				ImageIDs:     []int{sosPostImage.ID, sosPostImage2.ID},
				Reward:       "Test Reward2",
				DateStartAt:  "2023-12-01T00:00:00Z",
				DateEndAt:    "2023-12-05T00:00:00Z",
				CareType:     sos_post.CareTypeFoster,
				CarerGender:  sos_post.CarerGenderMale,
				RewardAmount: sos_post.RewardAmountHour,
				ConditionIDs: []int{1, 2, 3},
				PetIDs:       []int{addPets[0].ID},
			}

			updateSosPost, err := sosPostService.UpdateSosPost(updateSosPostData)

			assertConditionEquals(t, sosPost.Conditions, conditionIDs)
			assertPetEquals(t, sosPost.Pets[0], addPets[0])
			assertMediaEquals(t, updateSosPost.Media, sosPostMedia)

			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}
			if updateSosPost.Title != updateSosPostData.Title {
				t.Errorf("got %v want %v", updateSosPost.Title, updateSosPostData.Title)
			}
			if updateSosPost.Content != updateSosPostData.Content {
				t.Errorf("got %v want %v", updateSosPost.Content, updateSosPostData.Content)
			}
			if updateSosPost.Reward != updateSosPostData.Reward {
				t.Errorf("got %v want %v", updateSosPost.Reward, updateSosPostData.Reward)
			}
			if updateSosPost.DateStartAt != updateSosPostData.DateStartAt {
				t.Errorf("got %v want %v", updateSosPost.DateStartAt, updateSosPostData.DateStartAt)
			}
			if updateSosPost.DateEndAt != updateSosPostData.DateEndAt {
				t.Errorf("got %v want %v", updateSosPost.DateEndAt, updateSosPostData.DateEndAt)
			}
			if updateSosPost.CareType != updateSosPostData.CareType {
				t.Errorf("got %v want %v", updateSosPost.CareType, updateSosPostData.CareType)
			}
			if updateSosPost.CarerGender != updateSosPostData.CarerGender {
				t.Errorf("got %v want %v", updateSosPost.CarerGender, updateSosPostData.CarerGender)
			}
			if updateSosPost.RewardAmount != updateSosPostData.RewardAmount {
				t.Errorf("got %v want %v", updateSosPost.RewardAmount, updateSosPostData.RewardAmount)
			}
			if updateSosPost.ThumbnailID != sosPostImage.ID {
				t.Errorf("got %v want %v", updateSosPost.ThumbnailID, sosPostImage.ID)
			}
			if updateSosPost.AuthorID != owner.ID {
				t.Errorf("got %v want %v", updateSosPost.AuthorID, owner.ID)
			}
		})
	})
}
func assertConditionEquals(t *testing.T, got []sos_post.ConditionView, want []int) {
	for i := range want {
		if i+1 != got[i].ID {
			t.Errorf("got %v want %v", got[i].ID, i+1)
		}
	}
}

func assertPetEquals(t *testing.T, got pet.PetView, want pet.PetView) {
	if got.Name != want.Name {
		t.Errorf("got %v want %v", got.Name, want.Name)
	}

	if got.PetType != want.PetType {
		t.Errorf("got %v want %v", got.PetType, want.PetType)
	}

	if got.Sex != want.Sex {
		t.Errorf("got %v want %v", got.Sex, want.PetType)
	}

	if got.Neutered != want.Neutered {
		t.Errorf("got %v want %v", got.Neutered, want.Neutered)
	}

	if got.Breed != want.Breed {
		t.Errorf("got %v want %v", got.Breed, want.Breed)
	}

	if got.BirthDate != want.BirthDate {
		t.Errorf("got %v want %v", got.BirthDate, want.BirthDate)
	}

	if got.WeightInKg != want.WeightInKg {
		t.Errorf("got %v want %v", got.WeightInKg, want.WeightInKg)
	}
}

func assertMediaEquals(t *testing.T, got []media.MediaView, want []media.MediaView) {
	for i, media := range want {
		if got[i].ID != media.ID {
			t.Errorf("got %v want %v", got[i].ID, media.ID)
		}
		if got[i].MediaType != media.MediaType {
			t.Errorf("got %v want %v", got[i].MediaType, media.MediaType)
		}
		if got[i].URL != media.URL {
			t.Errorf("got %v want %v", got[i].URL, media.URL)
		}
		if got[i].CreatedAt != media.CreatedAt {
			t.Errorf("got %v want %v", got[i].CreatedAt, media.CreatedAt)
		}
	}
}
