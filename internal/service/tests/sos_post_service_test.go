package service_test

import (
	"context"
	"reflect"
	"testing"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sos_post"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
	"github.com/pet-sitter/pets-next-door-api/internal/tests"
)

func TestSosPostService(t *testing.T) {
	setUp := func(ctx context.Context, t *testing.T) (*database.DB, func(t *testing.T)) {
		t.Helper()
		db, _ := database.Open(tests.TestDatabaseURL)
		db.Flush()

		if err := database.WithTransaction(ctx, db, func(tx *database.Tx) *pnd.AppError {
			postgres.InitConditions(ctx, tx, sos_post.ConditionName)
			return nil
		}); err != nil {
			t.Errorf("InitConditions failed: %v", err)
		}

		return db, func(t *testing.T) {
			t.Helper()
			db.Close()
		}
	}
	t.Run("CreateSosPost", func(t *testing.T) {
		t.Run("돌봄 급구 게시글을 작성한다", func(t *testing.T) {
			ctx := context.Background()
			db, tearDown := setUp(ctx, t)
			defer tearDown(t)

			// given
			mediaService := service.NewMediaService(db, nil)
			profileImage := tests.AddDummyMedia(t, ctx, mediaService)
			sosPostImage := tests.AddDummyMedia(t, ctx, mediaService)
			sosPostImage2 := tests.AddDummyMedia(t, ctx, mediaService)

			userService := service.NewUserService(db, mediaService)
			userRequest := tests.GenerateDummyRegisterUserRequest(&profileImage.ID)
			owner, _ := userService.RegisterUser(ctx, userRequest)
			uid := owner.FirebaseUID
			addPets := tests.AddDummyPet(t, ctx, userService, uid, &profileImage.ID)

			// when
			sosPostService := service.NewSosPostService(db)
			imageIDs := []int{sosPostImage.ID, sosPostImage2.ID}
			petIDs := []int{addPets.ID}

			sosPostData := tests.GenerateDummyWriteSosPostRequest(imageIDs, petIDs, 0)
			sosPost, err := sosPostService.WriteSosPost(ctx, uid, sosPostData)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			// then
			assertConditionEquals(t, sosPost.Conditions, sosPostData.ConditionIDs)
			assertPetEquals(t, sosPost.Pets[0], *addPets)
			assertMediaEquals(t, sosPost.Media, (&media.MediaList{sosPostImage, sosPostImage2}).ToMediaViewList())
			assertDatesEquals(t, sosPost.Dates, sosPostData.Dates)

			if sosPost.Title != sosPostData.Title {
				t.Errorf("got %v want %v", sosPost.Title, sosPostData.Title)
			}
			if sosPost.Content != sosPostData.Content {
				t.Errorf("got %v want %v", sosPost.Content, sosPostData.Content)
			}
			if sosPost.Reward != sosPostData.Reward {
				t.Errorf("got %v want %v", sosPost.Reward, sosPostData.Reward)
			}
			if sosPost.CareType != sosPostData.CareType {
				t.Errorf("got %v want %v", sosPost.CareType, sosPostData.CareType)
			}
			if sosPost.CarerGender != sosPostData.CarerGender {
				t.Errorf("got %v want %v", sosPost.CarerGender, sosPostData.CarerGender)
			}
			if sosPost.RewardType != sosPostData.RewardType {
				t.Errorf("got %v want %v", sosPost.RewardType, sosPostData.RewardType)
			}
			if sosPost.ThumbnailID != sosPostData.ImageIDs[0] {
				t.Errorf("got %v want %v", sosPost.ThumbnailID, sosPostData.ImageIDs[0])
			}
			if sosPost.AuthorID != owner.ID {
				t.Errorf("got %v want %v", sosPost.AuthorID, owner.ID)
			}
		})
	})
	t.Run("FindSosPosts", func(t *testing.T) {
		t.Run("전체 돌봄 급구 게시글을 조회한다", func(t *testing.T) {
			ctx := context.Background()
			db, tearDown := setUp(ctx, t)
			defer tearDown(t)

			// given
			mediaService := service.NewMediaService(db, nil)
			profileImage := tests.AddDummyMedia(t, ctx, mediaService)
			sosPostImage := tests.AddDummyMedia(t, ctx, mediaService)
			sosPostImage2 := tests.AddDummyMedia(t, ctx, mediaService)

			userService := service.NewUserService(db, mediaService)
			userRequest := tests.GenerateDummyRegisterUserRequest(&profileImage.ID)
			owner, _ := userService.RegisterUser(ctx, userRequest)
			author := &user.UserWithoutPrivateInfo{
				ID:              owner.ID,
				ProfileImageURL: owner.ProfileImageURL,
				Nickname:        owner.Nickname,
			}
			uid := owner.FirebaseUID
			addPets := tests.AddDummyPet(t, ctx, userService, uid, &profileImage.ID)

			sosPostService := service.NewSosPostService(db)
			imageIDs := []int{sosPostImage.ID, sosPostImage2.ID}
			petIDs := []int{addPets.ID}
			conditionIDs := []int{1, 2}

			var sosPosts []sos_post.WriteSosPostView
			for i := 1; i < 4; i++ {
				sosPost := tests.WriteDummySosPosts(t, ctx, sosPostService, uid, imageIDs, petIDs, i)
				sosPosts = append(sosPosts, *sosPost)
			}

			// when
			sosPostList, err := sosPostService.FindSosPosts(ctx, 1, 3, "newest", "all")
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			// then
			for i, sosPost := range sosPostList.Items {
				idx := len(sosPostList.Items) - i - 1
				assertConditionEquals(t, sosPost.Conditions, conditionIDs)
				assertPetEquals(t, sosPost.Pets[0], *addPets)
				assertMediaEquals(t, sosPost.Media, (&media.MediaList{sosPostImage, sosPostImage2}).ToMediaViewList())
				assertAuthorEquals(t, sosPost.Author, author)
				assertDatesEquals(t, sosPost.Dates, sosPosts[idx].Dates)
				assertFindSosPostEquals(t, sosPost, sosPosts[idx])
			}
		})
		t.Run("전체 돌봄 급구 게시글의 정렬기준을 'deadline'으로 조회한다", func(t *testing.T) {
			ctx := context.Background()
			db, tearDown := setUp(ctx, t)
			defer tearDown(t)

			// given
			mediaService := service.NewMediaService(db, nil)
			profileImage := tests.AddDummyMedia(t, ctx, mediaService)
			sosPostImage := tests.AddDummyMedia(t, ctx, mediaService)
			sosPostImage2 := tests.AddDummyMedia(t, ctx, mediaService)

			userService := service.NewUserService(db, mediaService)
			userRequest := tests.GenerateDummyRegisterUserRequest(&profileImage.ID)
			owner, _ := userService.RegisterUser(ctx, userRequest)
			author := &user.UserWithoutPrivateInfo{
				ID:              owner.ID,
				ProfileImageURL: owner.ProfileImageURL,
				Nickname:        owner.Nickname,
			}
			uid := owner.FirebaseUID
			addPets := tests.AddDummyPets(t, ctx, userService, uid, &profileImage.ID)

			sosPostService := service.NewSosPostService(db)
			imageIDs := []int{sosPostImage.ID, sosPostImage2.ID}
			conditionIDs := []int{1, 2}

			var sosPosts []sos_post.WriteSosPostView
			for i := 1; i < 4; i++ {
				sosPost := tests.WriteDummySosPosts(t, ctx, sosPostService, uid, imageIDs, []int{addPets[i-1].ID}, i)
				sosPosts = append(sosPosts, *sosPost)
			}

			// when
			sosPostList, err := sosPostService.FindSosPosts(ctx, 1, 3, "newest", "all")
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			// then
			for i, sosPost := range sosPostList.Items {
				idx := len(sosPostList.Items) - i - 1
				assertConditionEquals(t, sosPost.Conditions, conditionIDs)
				assertPetEquals(t, sosPost.Pets[i-1], addPets[i-1])
				assertMediaEquals(t, sosPost.Media, (&media.MediaList{sosPostImage, sosPostImage2}).ToMediaViewList())
				assertAuthorEquals(t, sosPost.Author, author)
				assertDatesEquals(t, sosPost.Dates, sosPosts[idx].Dates)
				assertFindSosPostEquals(t, sosPost, sosPosts[idx])
			}
		})
		t.Run("전체 돌봄 급구 게시글 중 반려동물이 'dog'인 경우만 조회한다", func(t *testing.T) {
			ctx := context.Background()
			db, tearDown := setUp(ctx, t)
			defer tearDown(t)

			// given
			mediaService := service.NewMediaService(db, nil)
			profileImage := tests.AddDummyMedia(t, ctx, mediaService)
			sosPostImage := tests.AddDummyMedia(t, ctx, mediaService)
			sosPostImage2 := tests.AddDummyMedia(t, ctx, mediaService)

			userService := service.NewUserService(db, mediaService)
			userRequest := tests.GenerateDummyRegisterUserRequest(&profileImage.ID)
			owner, _ := userService.RegisterUser(ctx, userRequest)
			author := &user.UserWithoutPrivateInfo{
				ID:              owner.ID,
				ProfileImageURL: owner.ProfileImageURL,
				Nickname:        owner.Nickname,
			}
			uid := owner.FirebaseUID
			addPets := tests.AddDummyPets(t, ctx, userService, uid, &profileImage.ID)

			sosPostService := service.NewSosPostService(db)
			imageIDs := []int{sosPostImage.ID, sosPostImage2.ID}
			conditionIDs := []int{1, 2}

			var sosPosts []sos_post.WriteSosPostView
			// 강아지인 경우
			for i := 1; i < 3; i++ {
				sosPost := tests.WriteDummySosPosts(t, ctx, sosPostService, uid, imageIDs, []int{addPets[i-1].ID}, i)
				sosPosts = append(sosPosts, *sosPost)
			}

			// 고양이인 경우
			sosPosts = append(sosPosts,
				*tests.WriteDummySosPosts(t, ctx, sosPostService, uid, imageIDs, []int{addPets[2].ID}, 3))

			// 강아지, 고양이인 경우
			sosPosts = append(sosPosts,
				*tests.WriteDummySosPosts(t, ctx, sosPostService, uid, imageIDs, []int{addPets[1].ID, addPets[2].ID}, 4))

			// when
			sosPostList, err := sosPostService.FindSosPosts(ctx, 1, 3, "newest", "all")
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			// then
			for i, sosPost := range sosPostList.Items {
				idx := len(sosPostList.Items) - i - 1
				assertConditionEquals(t, sosPost.Conditions, conditionIDs)
				assertPetEquals(t, sosPost.Pets[i-1], addPets[i-1])
				assertMediaEquals(t, sosPost.Media, (&media.MediaList{sosPostImage, sosPostImage2}).ToMediaViewList())
				assertAuthorEquals(t, sosPost.Author, author)
				assertDatesEquals(t, sosPost.Dates, sosPosts[idx].Dates)
				assertFindSosPostEquals(t, sosPost, sosPosts[idx])
			}
		})
		t.Run("작성자 ID로 돌봄 급구 게시글을 조회합니다.", func(t *testing.T) {
			ctx := context.Background()
			db, tearDown := setUp(ctx, t)
			defer tearDown(t)

			mediaService := service.NewMediaService(db, nil)
			profileImage := tests.AddDummyMedia(t, ctx, mediaService)
			sosPostImage := tests.AddDummyMedia(t, ctx, mediaService)
			sosPostImage2 := tests.AddDummyMedia(t, ctx, mediaService)

			userService := service.NewUserService(db, mediaService)
			userRequest := tests.GenerateDummyRegisterUserRequest(&profileImage.ID)
			owner, _ := userService.RegisterUser(ctx, userRequest)
			author := &user.UserWithoutPrivateInfo{
				ID:              owner.ID,
				ProfileImageURL: owner.ProfileImageURL,
				Nickname:        owner.Nickname,
			}
			uid := owner.FirebaseUID
			addPet := tests.AddDummyPet(t, ctx, userService, uid, &profileImage.ID)

			sosPostService := service.NewSosPostService(db)
			imageIDs := []int{sosPostImage.ID, sosPostImage2.ID}
			conditionIDs := []int{1, 2}

			sosPosts := make([]sos_post.WriteSosPostView, 0)
			for i := 1; i < 4; i++ {
				sosPost := tests.WriteDummySosPosts(t, ctx, sosPostService, uid, imageIDs, []int{addPet.ID}, i)
				sosPosts = append(sosPosts, *sosPost)
			}

			// when
			sosPostListByAuthorID, err := sosPostService.FindSosPostsByAuthorID(ctx, owner.ID, 1, 3, "newest", "all")
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			// then
			for i, sosPost := range sosPostListByAuthorID.Items {
				idx := len(sosPostListByAuthorID.Items) - i - 1
				assertConditionEquals(t, sosPost.Conditions, conditionIDs)
				assertPetEquals(t, sosPost.Pets[0], *addPet)
				assertMediaEquals(t, sosPost.Media, (&media.MediaList{sosPostImage, sosPostImage2}).ToMediaViewList())
				assertAuthorEquals(t, sosPost.Author, author)
				assertDatesEquals(t, sosPost.Dates, sosPosts[idx].Dates)
				assertFindSosPostEquals(t, sosPost, sosPosts[idx])
			}
		})
	})
	t.Run("FindSosPostByID", func(t *testing.T) {
		t.Run("게시글 ID로 돌봄 급구 게시글을 조회합니다.", func(t *testing.T) {
			ctx := context.Background()
			db, tearDown := setUp(ctx, t)
			defer tearDown(t)

			mediaService := service.NewMediaService(db, nil)
			profileImage := tests.AddDummyMedia(t, ctx, mediaService)
			sosPostImage := tests.AddDummyMedia(t, ctx, mediaService)
			sosPostImage2 := tests.AddDummyMedia(t, ctx, mediaService)

			userService := service.NewUserService(db, mediaService)
			userRequest := tests.GenerateDummyRegisterUserRequest(&profileImage.ID)
			owner, _ := userService.RegisterUser(ctx, userRequest)
			author := &user.UserWithoutPrivateInfo{
				ID:              owner.ID,
				ProfileImageURL: owner.ProfileImageURL,
				Nickname:        owner.Nickname,
			}
			uid := owner.FirebaseUID
			addPet := tests.AddDummyPet(t, ctx, userService, uid, &profileImage.ID)

			sosPostService := service.NewSosPostService(db)
			imageIDs := []int{sosPostImage.ID, sosPostImage2.ID}
			conditionIDs := []int{1, 2}

			sosPosts := make([]sos_post.WriteSosPostView, 0)
			for i := 1; i < 4; i++ {
				sosPost := tests.WriteDummySosPosts(t, ctx, sosPostService, uid, imageIDs, []int{addPet.ID}, i)
				sosPosts = append(sosPosts, *sosPost)
			}

			// when
			findSosPostByID, err := sosPostService.FindSosPostByID(ctx, sosPosts[0].ID)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			// then
			assertConditionEquals(t, sosPosts[0].Conditions, conditionIDs)
			assertPetEquals(t, sosPosts[0].Pets[0], *addPet)
			assertMediaEquals(t, findSosPostByID.Media, (&media.MediaList{sosPostImage, sosPostImage2}).ToMediaViewList())
			assertAuthorEquals(t, findSosPostByID.Author, author)
			assertDatesEquals(t, findSosPostByID.Dates, sosPosts[0].Dates)
			assertFindSosPostEquals(t, *findSosPostByID, sosPosts[0])
		})
	})

	t.Run("UpdateSosPost", func(t *testing.T) {
		t.Run("돌봄 급구 게시글을 수정합니다.", func(t *testing.T) {
			ctx := context.Background()
			db, tearDown := setUp(ctx, t)
			defer tearDown(t)

			// given
			mediaService := service.NewMediaService(db, nil)
			profileImage := tests.AddDummyMedia(t, ctx, mediaService)
			sosPostImage := tests.AddDummyMedia(t, ctx, mediaService)
			sosPostImage2 := tests.AddDummyMedia(t, ctx, mediaService)

			userService := service.NewUserService(db, mediaService)
			userRequest := tests.GenerateDummyRegisterUserRequest(&profileImage.ID)
			owner, _ := userService.RegisterUser(ctx, userRequest)
			uid := owner.FirebaseUID
			addPet := tests.AddDummyPet(t, ctx, userService, uid, &profileImage.ID)

			sosPostService := service.NewSosPostService(db)
			sosPost := tests.WriteDummySosPosts(t, ctx, sosPostService, uid, []int{sosPostImage.ID}, []int{addPet.ID}, 1)

			updateSosPostData := &sos_post.UpdateSosPostRequest{
				ID:       sosPost.ID,
				Title:    "Title2",
				Content:  "Content2",
				ImageIDs: []int{sosPostImage.ID, sosPostImage2.ID},
				Reward:   "Reward2",
				Dates: []sos_post.SosDateView{
					{"2024-04-10", "2024-04-20"},
					{"2024-05-10", "2024-05-20"},
				},
				CareType:     sos_post.CareTypeFoster,
				CarerGender:  sos_post.CarerGenderMale,
				RewardType:   sos_post.RewardTypeFee,
				ConditionIDs: []int{1, 2},
				PetIDs:       []int{addPet.ID},
			}

			// when
			updateSosPost, err := sosPostService.UpdateSosPost(ctx, updateSosPostData)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			assertConditionEquals(t, sosPost.Conditions, updateSosPostData.ConditionIDs)
			assertPetEquals(t, sosPost.Pets[0], *addPet)
			assertMediaEquals(t, updateSosPost.Media, (&media.MediaList{sosPostImage, sosPostImage2}).ToMediaViewList())
			assertDatesEquals(t, updateSosPost.Dates, updateSosPostData.Dates)

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

func assertFindSosPostEquals(t *testing.T, got sos_post.FindSosPostView, want sos_post.WriteSosPostView) {
	t.Helper()

	if got.Title != want.Title {
		t.Errorf("got %v want %v", got.Title, want.Title)
	}
	if got.Content != want.Content {
		t.Errorf("got %v want %v", got.Content, want.Content)
	}
	if got.Reward != want.Reward {
		t.Errorf("got %v want %v", got.Reward, want.Reward)
	}
	if got.CareType != want.CareType {
		t.Errorf("got %v want %v", got.CareType, want.CareType)
	}
	if got.CarerGender != want.CarerGender {
		t.Errorf("got %v want %v", got.CarerGender, want.CarerGender)
	}
	if got.RewardType != want.RewardType {
		t.Errorf("got %v want %v", got.RewardType, want.RewardType)
	}
	if got.ThumbnailID != want.ThumbnailID {
		t.Errorf("got %v want %v", got.ThumbnailID, want.ThumbnailID)
	}
}

func assertConditionEquals(t *testing.T, got []sos_post.ConditionView, want []int) {
	t.Helper()

	for i := range want {
		if i+1 != got[i].ID {
			t.Errorf("got %v want %v", got[i].ID, i+1)
		}
	}
}

func assertPetEquals(t *testing.T, got pet.PetView, want pet.PetView) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertMediaEquals(t *testing.T, got media.MediaViewList, want media.MediaViewList) {
	t.Helper()

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
	}
}

func assertAuthorEquals(t *testing.T, got *user.UserWithoutPrivateInfo, want *user.UserWithoutPrivateInfo) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertDatesEquals(t *testing.T, got []sos_post.SosDateView, want []sos_post.SosDateView) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
