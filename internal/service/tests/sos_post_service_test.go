package service_test

import (
	"context"
	"fmt"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sos_post"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database/sql"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
	"github.com/pet-sitter/pets-next-door-api/internal/tests"
	"reflect"
	"testing"
)

func TestSosPostService(t *testing.T) {
	setUp := func(ctx context.Context, t *testing.T) (*database.DB, func(t *testing.T)) {
		db, _ := sql.OpenSqlDB(tests.TestDatabaseURL)
		db.Flush()

		if err := sql.WithTransaction(ctx, &db, func(tx *database.Tx) *pnd.AppError {
			postgres.InitConditions(ctx, *tx, sos_post.ConditionName)
			return nil
		}); err != nil {
			t.Errorf("InitConditions failed: %v", err)
		}

		return &db, func(t *testing.T) {
			db.Close()
		}
	}

	t.Run("CreateSosPost", func(t *testing.T) {
		t.Run("돌봄 급구 게시글을 작성합니다.", func(t *testing.T) {
			ctx := context.Background()
			db, tearDown := setUp(ctx, t)
			defer tearDown(t)

			mediaService := service.NewMediaService(*db, nil)
			profileImage, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test.com",
			})
			sosPostImage, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test2.com",
			})
			sosPostImage2, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test3.com",
			})

			userService := service.NewUserService(*db, mediaService)

			owner, err := userService.RegisterUser(ctx, &user.RegisterUserRequest{
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
						Name:           "name",
						PetType:        "dog",
						Sex:            "male",
						Neutered:       true,
						Breed:          "poodle",
						BirthDate:      "2020-01-01T00:00:00Z",
						WeightInKg:     10.0,
						Remarks:        "",
						ProfileImageID: nil,
					},
				},
			}

			addPets, err := userService.AddPetsToOwner(ctx, uid, pets)
			if err != nil {
				t.Errorf(err.Err.Error())
			}

			sosPostService := service.NewSosPostService(*db)

			conditionIDs := []int{1, 2}

			writeSosPostRequest := &sos_post.WriteSosPostRequest{
				Title:    "Test Title",
				Content:  "Test Content",
				ImageIDs: []int{sosPostImage.ID, sosPostImage2.ID},
				Reward:   "Test Reward",
				Dates: []sos_post.SosDateView{{"2024-03-30", "2024-03-30"},
					{"2024-04-01", "2024-04-02"}},
				CareType:     sos_post.CareTypeFoster,
				CarerGender:  sos_post.CarerGenderMale,
				RewardType:   sos_post.RewardTypeFee,
				ConditionIDs: conditionIDs,
				PetIDs:       []int{addPets[0].ID},
			}

			sosPost, err := sosPostService.WriteSosPost(ctx, uid, writeSosPostRequest)
			if err != nil {
				t.Errorf(err.Err.Error())
			}

			assertConditionEquals(t, sosPost.Conditions, conditionIDs)
			assertPetEquals(t, sosPost.Pets[0], addPets[0])
			assertMediaEquals(t, sosPost.Media, (&media.MediaList{sosPostImage, sosPostImage2}).ToMediaViewList())
			assertDatesEquals(t, sosPost.Dates, writeSosPostRequest.Dates)

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
			if sosPost.CareType != sos_post.CareTypeFoster {
				t.Errorf("got %v want %v", sosPost.CareType, sos_post.CareTypeFoster)
			}
			if sosPost.CarerGender != sos_post.CarerGenderMale {
				t.Errorf("got %v want %v", sosPost.CarerGender, sos_post.CarerGenderMale)
			}
			if sosPost.RewardType != sos_post.RewardTypeFee {
				t.Errorf("got %v want %v", sosPost.RewardType, sos_post.RewardTypeFee)
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
			ctx := context.Background()
			db, tearDown := setUp(ctx, t)
			defer tearDown(t)

			mediaService := service.NewMediaService(*db, nil)
			profileImage, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test.com",
			})
			sosPostImage, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test2.com",
			})
			sosPostImage2, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test3.com",
			})

			userService := service.NewUserService(*db, mediaService)
			owner, err := userService.RegisterUser(ctx, &user.RegisterUserRequest{
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
						Name:           "name",
						PetType:        "dog",
						Sex:            "male",
						Neutered:       true,
						Breed:          "poodle",
						BirthDate:      "2020-01-01T00:00:00Z",
						WeightInKg:     10.0,
						Remarks:        "",
						ProfileImageID: nil,
					},
				},
			}

			addPets, err := userService.AddPetsToOwner(ctx, uid, pets)
			if err != nil {
				t.Errorf(err.Err.Error())
			}

			sosPostService := service.NewSosPostService(*db)

			conditionIDs := []int{1, 2}

			var sosPosts []sos_post.WriteSosPostView

			for i := 1; i < 4; i++ {
				sosPost, err := sosPostService.WriteSosPost(ctx, uid, &sos_post.WriteSosPostRequest{
					Title:    fmt.Sprintf("Title%d", i),
					Content:  fmt.Sprintf("Test Content%d", i),
					ImageIDs: []int{sosPostImage.ID, sosPostImage2.ID},
					Reward:   fmt.Sprintf("Test Reward%d", i),
					Dates: []sos_post.SosDateView{{"2024-03-30", "2024-03-30"},
						{"2024-04-01", "2024-04-02"}},
					CareType:     sos_post.CareTypeFoster,
					CarerGender:  sos_post.CarerGenderMale,
					RewardType:   sos_post.RewardTypeFee,
					ConditionIDs: conditionIDs,
					PetIDs:       []int{addPets[0].ID},
				})

				if err != nil {
					t.Errorf(err.Err.Error())
				}

				sosPosts = append(sosPosts, *sosPost)
			}

			author := &user.UserWithoutPrivateInfo{
				ID:              owner.ID,
				ProfileImageURL: owner.ProfileImageURL,
				Nickname:        owner.Nickname,
			}

			sosPostList, err := sosPostService.FindSosPosts(ctx, 1, 3, "newest")
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			for i, sosPost := range sosPostList.Items {
				assertConditionEquals(t, sosPost.Conditions, conditionIDs)
				assertPetEquals(t, sosPost.Pets[0], addPets[0])
				assertMediaEquals(t, sosPost.Media, (&media.MediaList{sosPostImage, sosPostImage2}).ToMediaViewList())
				assertAuthorEquals(t, sosPost.Author, author)

				idx := len(sosPostList.Items) - i - 1

				assertDatesEquals(t, sosPost.Dates, sosPosts[idx].Dates)

				if sosPost.Title != sosPosts[idx].Title {
					t.Errorf("got %v want %v", sosPost.Title, sosPosts[idx].Title)
				}
				if sosPost.Content != sosPosts[idx].Content {
					t.Errorf("got %v want %v", sosPost.Content, sosPosts[idx].Content)
				}
				if sosPost.Reward != sosPosts[idx].Reward {
					t.Errorf("got %v want %v", sosPost.Reward, sosPosts[idx].Reward)
				}
				if sosPost.CareType != sosPosts[idx].CareType {
					t.Errorf("got %v want %v", sosPost.CareType, sosPosts[idx].CareType)
				}
				if sosPost.CarerGender != sosPosts[idx].CarerGender {
					t.Errorf("got %v want %v", sosPost.CarerGender, sosPosts[idx].CarerGender)
				}
				if sosPost.RewardType != sosPosts[idx].RewardType {
					t.Errorf("got %v want %v", sosPost.RewardType, sosPosts[idx].RewardType)
				}
				if sosPost.ThumbnailID != sosPostImage.ID {
					t.Errorf("got %v want %v", sosPost.ThumbnailID, sosPostImage.ID)
				}
			}
		})
		t.Run("작성자 ID로 돌봄 급구 게시글을 조회합니다.", func(t *testing.T) {
			ctx := context.Background()
			db, tearDown := setUp(ctx, t)
			defer tearDown(t)

			mediaService := service.NewMediaService(*db, nil)
			profileImage, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test.com",
			})
			sosPostImage, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test2.com",
			})
			sosPostImage2, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test3.com",
			})

			userService := service.NewUserService(*db, mediaService)

			owner, err := userService.RegisterUser(ctx, &user.RegisterUserRequest{
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
						Name:           "name",
						PetType:        "dog",
						Sex:            "male",
						Neutered:       true,
						Breed:          "poodle",
						BirthDate:      "2020-01-01T00:00:00Z",
						WeightInKg:     10.0,
						Remarks:        "",
						ProfileImageID: nil,
					},
				},
			}

			addPets, err := userService.AddPetsToOwner(ctx, uid, pets)
			if err != nil {
				t.Errorf(err.Err.Error())
			}

			sosPostService := service.NewSosPostService(*db)

			conditionIDs := []int{1, 2}
			//krLocation, _ := time.LoadLocation("Asia/Seoul")

			sosPosts := make([]sos_post.WriteSosPostView, 0)
			for i := 1; i < 4; i++ {
				sosPost, err := sosPostService.WriteSosPost(ctx, uid, &sos_post.WriteSosPostRequest{
					Title:        fmt.Sprintf("Title%d", i),
					Content:      fmt.Sprintf("Test Content%d", i),
					ImageIDs:     []int{sosPostImage.ID, sosPostImage2.ID},
					Reward:       fmt.Sprintf("Test Reward%d", i),
					CareType:     sos_post.CareTypeFoster,
					CarerGender:  sos_post.CarerGenderMale,
					RewardType:   sos_post.RewardTypeFee,
					ConditionIDs: conditionIDs,
					PetIDs:       []int{addPets[0].ID},
				})

				if err != nil {
					t.Errorf(err.Err.Error())
				}

				sosPosts = append(sosPosts, *sosPost)
			}

			sosPostListByAuthorID, err := sosPostService.FindSosPostsByAuthorID(ctx, owner.ID, 1, 3, "newest")
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			author := &user.UserWithoutPrivateInfo{
				ID:              owner.ID,
				ProfileImageURL: owner.ProfileImageURL,
				Nickname:        owner.Nickname,
			}

			for i, sosPost := range sosPostListByAuthorID.Items {
				assertConditionEquals(t, sosPost.Conditions, conditionIDs)
				assertPetEquals(t, sosPost.Pets[0], addPets[0])
				assertMediaEquals(t, sosPost.Media, (&media.MediaList{sosPostImage, sosPostImage2}).ToMediaViewList())
				assertAuthorEquals(t, sosPost.Author, author)

				idx := len(sosPostListByAuthorID.Items) - i - 1

				assertDatesEquals(t, sosPost.Dates, sosPosts[idx].Dates)

				if sosPost.Title != sosPosts[idx].Title {
					t.Errorf("got %v want %v", sosPost.Title, sosPosts[idx].Title)
				}
				if sosPost.Content != sosPosts[idx].Content {
					t.Errorf("got %v want %v", sosPost.Content, sosPosts[idx].Content)
				}
				if sosPost.Reward != sosPosts[idx].Reward {
					t.Errorf("got %v want %v", sosPost.Reward, sosPosts[idx].Reward)
				}
				if sosPost.CareType != sosPosts[idx].CareType {
					t.Errorf("got %v want %v", sosPost.CareType, sosPosts[idx].CareType)
				}
				if sosPost.CarerGender != sosPosts[idx].CarerGender {
					t.Errorf("got %v want %v", sosPost.CarerGender, sosPosts[idx].CarerGender)
				}
				if sosPost.RewardType != sosPosts[idx].RewardType {
					t.Errorf("got %v want %v", sosPost.RewardType, sosPosts[idx].RewardType)
				}
				if sosPost.ThumbnailID != sosPostImage.ID {
					t.Errorf("got %v want %v", sosPost.ThumbnailID, sosPostImage.ID)
				}
			}
		})
	})

	t.Run("FindSosPostByID", func(t *testing.T) {
		t.Run("게시글 ID로 돌봄 급구 게시글을 조회합니다.", func(t *testing.T) {
			ctx := context.Background()
			db, tearDown := setUp(ctx, t)
			defer tearDown(t)

			mediaService := service.NewMediaService(*db, nil)
			profileImage, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test.com",
			})
			sosPostImage, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test2.com",
			})
			sosPostImage2, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test3.com",
			})

			userService := service.NewUserService(*db, mediaService)

			owner, err := userService.RegisterUser(ctx, &user.RegisterUserRequest{
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
						Name:           "name",
						PetType:        "dog",
						Sex:            "male",
						Neutered:       true,
						Breed:          "poodle",
						BirthDate:      "2020-01-01T00:00:00Z",
						WeightInKg:     10.0,
						Remarks:        "",
						ProfileImageID: nil,
					},
				},
			}

			addPets, err := userService.AddPetsToOwner(ctx, uid, pets)
			if err != nil {
				t.Errorf(err.Err.Error())
			}

			sosPostService := service.NewSosPostService(*db)

			conditionIDs := []int{1, 2}

			sosPosts := make([]sos_post.WriteSosPostView, 0)
			for i := 1; i < 4; i++ {
				sosPost, err := sosPostService.WriteSosPost(ctx, uid, &sos_post.WriteSosPostRequest{
					Title:    fmt.Sprintf("Title%d", i),
					Content:  fmt.Sprintf("Test Content%d", i),
					ImageIDs: []int{sosPostImage.ID, sosPostImage2.ID},
					Reward:   fmt.Sprintf("Test Reward%d", i),
					Dates: []sos_post.SosDateView{{"2024-03-30", "2024-03-30"},
						{"2024-04-01", "2024-04-02"}},
					CareType:     sos_post.CareTypeFoster,
					CarerGender:  sos_post.CarerGenderMale,
					RewardType:   sos_post.RewardTypeFee,
					ConditionIDs: conditionIDs,
					PetIDs:       []int{addPets[0].ID},
				})

				if err != nil {
					t.Errorf(err.Err.Error())
				}

				sosPosts = append(sosPosts, *sosPost)
			}

			author := &user.UserWithoutPrivateInfo{
				ID:              owner.ID,
				ProfileImageURL: owner.ProfileImageURL,
				Nickname:        owner.Nickname,
			}

			findSosPostByID, err := sosPostService.FindSosPostByID(ctx, sosPosts[0].ID)

			assertConditionEquals(t, sosPosts[0].Conditions, conditionIDs)
			assertPetEquals(t, sosPosts[0].Pets[0], addPets[0])
			assertMediaEquals(t, findSosPostByID.Media, (&media.MediaList{sosPostImage, sosPostImage2}).ToMediaViewList())
			assertAuthorEquals(t, findSosPostByID.Author, author)
			assertDatesEquals(t, findSosPostByID.Dates, sosPosts[0].Dates)
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
			if findSosPostByID.CareType != sosPosts[0].CareType {
				t.Errorf("got %v want %v", findSosPostByID.CareType, sosPosts[0].CareType)
			}
			if findSosPostByID.CarerGender != sosPosts[0].CarerGender {
				t.Errorf("got %v want %v", findSosPostByID.CarerGender, sosPosts[0].CarerGender)
			}
			if findSosPostByID.RewardType != sosPosts[0].RewardType {
				t.Errorf("got %v want %v", findSosPostByID.RewardType, sosPosts[0].RewardType)
			}
			if findSosPostByID.ThumbnailID != sosPostImage.ID {
				t.Errorf("got %v want %v", findSosPostByID.ThumbnailID, sosPostImage.ID)
			}
		})
	})

	t.Run("UpdateSosPost", func(t *testing.T) {
		t.Run("돌봄 급구 게시글을 수정합니다.", func(t *testing.T) {
			ctx := context.Background()
			db, tearDown := setUp(ctx, t)
			defer tearDown(t)

			mediaService := service.NewMediaService(*db, nil)
			profileImage, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test.com",
			})
			sosPostImage, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test2.com",
			})
			sosPostImage2, _ := mediaService.CreateMedia(ctx, &media.Media{
				MediaType: media.IMAGE_MEDIA_TYPE,
				URL:       "https://test3.com",
			})

			userService := service.NewUserService(*db, mediaService)

			owner, err := userService.RegisterUser(ctx, &user.RegisterUserRequest{
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
						Name:           "name",
						PetType:        "dog",
						Sex:            "male",
						Neutered:       true,
						Breed:          "poodle",
						BirthDate:      "2020-01-01T00:00:00Z",
						WeightInKg:     10.0,
						Remarks:        "",
						ProfileImageID: nil,
					},
				},
			}

			addPets, err := userService.AddPetsToOwner(ctx, uid, pets)
			if err != nil {
				t.Errorf(err.Err.Error())
			}

			sosPostService := service.NewSosPostService(*db)

			conditionIDs := []int{1, 2}

			sosPost, err := sosPostService.WriteSosPost(ctx, uid, &sos_post.WriteSosPostRequest{
				Title:        "Title1",
				Content:      "Test Content1",
				ImageIDs:     []int{sosPostImage.ID, sosPostImage2.ID},
				Reward:       "Test Reward1",
				CareType:     sos_post.CareTypeFoster,
				CarerGender:  sos_post.CarerGenderMale,
				RewardType:   sos_post.RewardTypeFee,
				ConditionIDs: conditionIDs,
				PetIDs:       []int{addPets[0].ID},
			})

			if err != nil {
				t.Errorf(err.Err.Error())
			}

			updateSosPostData := &sos_post.UpdateSosPostRequest{
				ID:       sosPost.ID,
				Title:    "Title2",
				Content:  "Test Content2",
				ImageIDs: []int{sosPostImage.ID, sosPostImage2.ID},
				Reward:   "Test Reward2",
				Dates: []sos_post.SosDateView{{"2024-03-30", "2024-03-30"},
					{"2024-04-01", "2024-04-02"}},
				CareType:     sos_post.CareTypeFoster,
				CarerGender:  sos_post.CarerGenderMale,
				RewardType:   sos_post.RewardTypeFee,
				ConditionIDs: []int{1, 2, 3},
				PetIDs:       []int{addPets[0].ID},
			}

			updateSosPost, err := sosPostService.UpdateSosPost(ctx, updateSosPostData)

			assertConditionEquals(t, sosPost.Conditions, conditionIDs)
			assertPetEquals(t, sosPost.Pets[0], addPets[0])
			assertMediaEquals(t, updateSosPost.Media, (&media.MediaList{sosPostImage, sosPostImage2}).ToMediaViewList())
			assertDatesEquals(t, updateSosPost.Dates, updateSosPostData.Dates)

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
			if updateSosPost.CareType != updateSosPostData.CareType {
				t.Errorf("got %v want %v", updateSosPost.CareType, updateSosPostData.CareType)
			}
			if updateSosPost.CarerGender != updateSosPostData.CarerGender {
				t.Errorf("got %v want %v", updateSosPost.CarerGender, updateSosPostData.CarerGender)
			}
			if updateSosPost.RewardType != updateSosPostData.RewardType {
				t.Errorf("got %v want %v", updateSosPost.RewardType, updateSosPostData.RewardType)
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
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertMediaEquals(t *testing.T, got media.MediaViewList, want media.MediaViewList) {
	for i, media := range want {
		if !reflect.DeepEqual(got[i], media) {
			t.Errorf("got %v want %v", got[i], media)
		}
	}
}

func assertAuthorEquals(t *testing.T, got *user.UserWithoutPrivateInfo, want *user.UserWithoutPrivateInfo) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertDatesEquals(t *testing.T, got []sos_post.SosDateView, want []sos_post.SosDateView) {
	for i, date := range want {
		if got[i].DateStartAt != date.DateStartAt {
			t.Errorf("got %v want %v", got[i].DateStartAt, date.DateStartAt)
		}
		if got[i].DateEndAt != date.DateEndAt {
			t.Errorf("got %v want %v", got[i].DateEndAt, date.DateEndAt)
		}
	}
}
