package service_test

import (
	"context"
	"reflect"
	"testing"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sospost"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
	"github.com/pet-sitter/pets-next-door-api/internal/tests"
)

//nolint:gocognit
func TestSOSPostService(t *testing.T) {
	setUp := func(ctx context.Context, t *testing.T) (*database.DB, func(t *testing.T)) {
		t.Helper()
		db, _ := database.Open(tests.TestDatabaseURL)
		db.Flush()

		if err := database.WithTransaction(ctx, db, func(tx *database.Tx) *pnd.AppError {
			postgres.InitConditions(ctx, tx, sospost.ConditionName)
			return nil
		}); err != nil {
			t.Errorf("InitConditions failed: %v", err)
		}

		return db, func(t *testing.T) {
			t.Helper()
			db.Close()
		}
	}
	t.Run("CreateSOSPost", func(t *testing.T) {
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
			sosPostService := service.NewSOSPostService(db)
			imageIDs := []int64{int64(sosPostImage.ID), int64(sosPostImage2.ID)}
			petIDs := []int64{addPets.ID}

			sosPostData := tests.GenerateDummyWriteSOSPostRequest(imageIDs, petIDs, 0)
			sosPost, err := sosPostService.WriteSOSPost(ctx, uid, sosPostData)
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
			if int64(sosPost.ThumbnailID) != sosPostData.ImageIDs[0] {
				t.Errorf("got %v want %v", sosPost.ThumbnailID, sosPostData.ImageIDs[0])
			}
			if int64(sosPost.AuthorID) != owner.ID {
				t.Errorf("got %v want %v", sosPost.AuthorID, owner.ID)
			}
		})
	})
	t.Run("FindSOSPosts", func(t *testing.T) {
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
			author := &user.WithoutPrivateInfo{
				ID:              owner.ID,
				ProfileImageURL: owner.ProfileImageURL,
				Nickname:        owner.Nickname,
			}
			uid := owner.FirebaseUID
			addPets := tests.AddDummyPet(t, ctx, userService, uid, &profileImage.ID)

			sosPostService := service.NewSOSPostService(db)
			imageIDs := []int64{int64(sosPostImage.ID), int64(sosPostImage2.ID)}
			petIDs := []int64{addPets.ID}
			conditionIDs := []int{1, 2}

			var sosPosts []sospost.WriteSOSPostView
			for i := 1; i < 4; i++ {
				sosPost := tests.WriteDummySOSPosts(t, ctx, sosPostService, uid, imageIDs, petIDs, i)
				sosPosts = append(sosPosts, *sosPost)
			}

			// when
			sosPostList, err := sosPostService.FindSOSPosts(ctx, 1, 3, "newest", "all")
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
				assertFindSOSPostEquals(t, sosPost, sosPosts[idx])
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
			author := &user.WithoutPrivateInfo{
				ID:              owner.ID,
				ProfileImageURL: owner.ProfileImageURL,
				Nickname:        owner.Nickname,
			}
			uid := owner.FirebaseUID
			addPets := tests.AddDummyPets(t, ctx, userService, uid, &profileImage.ID)

			sosPostService := service.NewSOSPostService(db)
			imageIDs := []int64{int64(sosPostImage.ID), int64(sosPostImage2.ID)}
			conditionIDs := []int{1, 2}

			var sosPosts []sospost.WriteSOSPostView
			for i := 1; i < 4; i++ {
				sosPost := tests.WriteDummySOSPosts(t, ctx, sosPostService, uid, imageIDs, []int64{addPets.Pets[i-1].ID}, i)
				sosPosts = append(sosPosts, *sosPost)
			}

			// when
			sosPostList, err := sosPostService.FindSOSPosts(ctx, 1, 3, "newest", "all")
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			// then
			for i, sosPost := range sosPostList.Items {
				idx := len(sosPostList.Items) - i - 1
				assertConditionEquals(t, sosPost.Conditions, conditionIDs)
				assertPetEquals(t, sosPost.Pets[i-1], addPets.Pets[i-1])
				assertMediaEquals(t, sosPost.Media, (&media.MediaList{sosPostImage, sosPostImage2}).ToMediaViewList())
				assertAuthorEquals(t, sosPost.Author, author)
				assertDatesEquals(t, sosPost.Dates, sosPosts[idx].Dates)
				assertFindSOSPostEquals(t, sosPost, sosPosts[idx])
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
			author := &user.WithoutPrivateInfo{
				ID:              owner.ID,
				ProfileImageURL: owner.ProfileImageURL,
				Nickname:        owner.Nickname,
			}
			uid := owner.FirebaseUID
			addPets := tests.AddDummyPets(t, ctx, userService, uid, &profileImage.ID)

			sosPostService := service.NewSOSPostService(db)
			imageIDs := []int64{int64(sosPostImage.ID), int64(sosPostImage2.ID)}
			conditionIDs := []int{1, 2}

			var sosPosts []sospost.WriteSOSPostView
			// 강아지인 경우
			for i := 1; i < 3; i++ {
				sosPost := tests.WriteDummySOSPosts(t, ctx, sosPostService, uid, imageIDs, []int64{addPets.Pets[i-1].ID}, i)
				sosPosts = append(sosPosts, *sosPost)
			}

			// 고양이인 경우
			sosPosts = append(sosPosts,
				*tests.WriteDummySOSPosts(t, ctx, sosPostService, uid, imageIDs, []int64{addPets.Pets[2].ID}, 3))

			// 강아지, 고양이인 경우
			sosPosts = append(sosPosts,
				*tests.WriteDummySOSPosts(t, ctx,
					sosPostService, uid, imageIDs, []int64{addPets.Pets[1].ID, addPets.Pets[2].ID},
					4,
				),
			)

			// when
			sosPostList, err := sosPostService.FindSOSPosts(ctx, 1, 3, "newest", "all")
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			// then
			for i, sosPost := range sosPostList.Items {
				idx := len(sosPostList.Items) - i - 1
				assertConditionEquals(t, sosPost.Conditions, conditionIDs)
				assertPetEquals(t, sosPost.Pets[i-1], addPets.Pets[i-1])
				assertMediaEquals(t, sosPost.Media, (&media.MediaList{sosPostImage, sosPostImage2}).ToMediaViewList())
				assertAuthorEquals(t, sosPost.Author, author)
				assertDatesEquals(t, sosPost.Dates, sosPosts[idx].Dates)
				assertFindSOSPostEquals(t, sosPost, sosPosts[idx])
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
			author := &user.WithoutPrivateInfo{
				ID:              owner.ID,
				ProfileImageURL: owner.ProfileImageURL,
				Nickname:        owner.Nickname,
			}
			uid := owner.FirebaseUID
			addPet := tests.AddDummyPet(t, ctx, userService, uid, &profileImage.ID)

			sosPostService := service.NewSOSPostService(db)
			imageIDs := []int64{int64(sosPostImage.ID), int64(sosPostImage2.ID)}
			conditionIDs := []int{1, 2}

			sosPosts := make([]sospost.WriteSOSPostView, 0)
			for i := 1; i < 4; i++ {
				sosPost := tests.WriteDummySOSPosts(t, ctx, sosPostService, uid, imageIDs, []int64{addPet.ID}, i)
				sosPosts = append(sosPosts, *sosPost)
			}

			// when
			sosPostListByAuthorID, err := sosPostService.FindSOSPostsByAuthorID(ctx, int(owner.ID), 1, 3, "newest", "all")
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
				assertFindSOSPostEquals(t, sosPost, sosPosts[idx])
			}
		})
	})
	t.Run("FindSOSPostByID", func(t *testing.T) {
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
			author := &user.WithoutPrivateInfo{
				ID:              owner.ID,
				ProfileImageURL: owner.ProfileImageURL,
				Nickname:        owner.Nickname,
			}
			uid := owner.FirebaseUID
			addPet := tests.AddDummyPet(t, ctx, userService, uid, &profileImage.ID)

			sosPostService := service.NewSOSPostService(db)
			imageIDs := []int64{int64(sosPostImage.ID), int64(sosPostImage2.ID)}
			conditionIDs := []int{1, 2}

			sosPosts := make([]sospost.WriteSOSPostView, 0)
			for i := 1; i < 4; i++ {
				sosPost := tests.WriteDummySOSPosts(t, ctx, sosPostService, uid, imageIDs, []int64{addPet.ID}, i)
				sosPosts = append(sosPosts, *sosPost)
			}

			// when
			findSOSPostByID, err := sosPostService.FindSOSPostByID(ctx, sosPosts[0].ID)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			// then
			assertConditionEquals(t, sosPosts[0].Conditions, conditionIDs)
			assertPetEquals(t, sosPosts[0].Pets[0], *addPet)
			assertMediaEquals(t, findSOSPostByID.Media, (&media.MediaList{sosPostImage, sosPostImage2}).ToMediaViewList())
			assertAuthorEquals(t, findSOSPostByID.Author, author)
			assertDatesEquals(t, findSOSPostByID.Dates, sosPosts[0].Dates)
			assertFindSOSPostEquals(t, *findSOSPostByID, sosPosts[0])
		})
	})

	t.Run("UpdateSOSPost", func(t *testing.T) {
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

			sosPostService := service.NewSOSPostService(db)
			sosPost := tests.WriteDummySOSPosts(t, ctx,
				sosPostService, uid, []int64{int64(sosPostImage.ID)}, []int64{addPet.ID},
				1,
			)

			updateSOSPostData := &sospost.UpdateSOSPostRequest{
				ID:       sosPost.ID,
				Title:    "Title2",
				Content:  "Content2",
				ImageIDs: []int{sosPostImage.ID, sosPostImage2.ID},
				Reward:   "Reward2",
				Dates: []sospost.SOSDateView{
					{"2024-04-10", "2024-04-20"},
					{"2024-05-10", "2024-05-20"},
				},
				CareType:     sospost.CareTypeFoster,
				CarerGender:  sospost.CarerGenderMale,
				RewardType:   sospost.RewardTypeFee,
				ConditionIDs: []int{1, 2},
				PetIDs:       []int{int(addPet.ID)},
			}

			// when
			updateSOSPost, err := sosPostService.UpdateSOSPost(ctx, updateSOSPostData)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			assertConditionEquals(t, sosPost.Conditions, updateSOSPostData.ConditionIDs)
			assertPetEquals(t, sosPost.Pets[0], *addPet)
			assertMediaEquals(t, updateSOSPost.Media, (&media.MediaList{sosPostImage, sosPostImage2}).ToMediaViewList())
			assertDatesEquals(t, updateSOSPost.Dates, updateSOSPostData.Dates)

			if updateSOSPost.Title != updateSOSPostData.Title {
				t.Errorf("got %v want %v", updateSOSPost.Title, updateSOSPostData.Title)
			}
			if updateSOSPost.Content != updateSOSPostData.Content {
				t.Errorf("got %v want %v", updateSOSPost.Content, updateSOSPostData.Content)
			}
			if updateSOSPost.Reward != updateSOSPostData.Reward {
				t.Errorf("got %v want %v", updateSOSPost.Reward, updateSOSPostData.Reward)
			}
			if updateSOSPost.CareType != updateSOSPostData.CareType {
				t.Errorf("got %v want %v", updateSOSPost.CareType, updateSOSPostData.CareType)
			}
			if updateSOSPost.CarerGender != updateSOSPostData.CarerGender {
				t.Errorf("got %v want %v", updateSOSPost.CarerGender, updateSOSPostData.CarerGender)
			}
			if updateSOSPost.RewardType != updateSOSPostData.RewardType {
				t.Errorf("got %v want %v", updateSOSPost.RewardType, updateSOSPostData.RewardType)
			}
			if updateSOSPost.ThumbnailID != sosPostImage.ID {
				t.Errorf("got %v want %v", updateSOSPost.ThumbnailID, sosPostImage.ID)
			}
			if int64(updateSOSPost.AuthorID) != owner.ID {
				t.Errorf("got %v want %v", updateSOSPost.AuthorID, owner.ID)
			}
		})
	})
}

func assertFindSOSPostEquals(t *testing.T, got sospost.FindSOSPostView, want sospost.WriteSOSPostView) {
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

func assertConditionEquals(t *testing.T, got []sospost.ConditionView, want []int) {
	t.Helper()

	for i := range want {
		if i+1 != got[i].ID {
			t.Errorf("got %v want %v", got[i].ID, i+1)
		}
	}
}

func assertPetEquals(t *testing.T, got pet.PetView, want pet.DetailView) {
	t.Helper()

	if int64(got.ID) != want.ID {
		t.Errorf("got %v want %v", got.ID, want.ID)
	}

	if got.Name != want.Name {
		t.Errorf("got %v want %v", got.Name, want.Name)
	}

	if got.PetType != want.PetType {
		t.Errorf("got %v want %v", got.PetType, want.PetType)
	}

	if got.Sex != want.Sex {
		t.Errorf("got %v want %v", got.Sex, want.Sex)
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

	if got.WeightInKg.String() != want.WeightInKg.String() {
		t.Errorf("got %v want %v", got.WeightInKg, want.WeightInKg)
	}

	if got.Remarks != want.Remarks {
		t.Errorf("got %v want %v", got.Remarks, want.Remarks)
	}

	switch {
	case got.ProfileImageURL == nil && want.ProfileImageURL != nil:
		t.Errorf("got %v want %v", got.ProfileImageURL, want.ProfileImageURL)
	case got.ProfileImageURL != nil && want.ProfileImageURL == nil:
		t.Errorf("got %v want %v", got.ProfileImageURL, want.ProfileImageURL)
	case *got.ProfileImageURL != *want.ProfileImageURL:
		t.Errorf("got %v want %v", *got.ProfileImageURL, *want.ProfileImageURL)
	}
}

func assertMediaEquals(t *testing.T, got, want media.MediaViewList) {
	t.Helper()

	for i, mediaData := range want {
		if got[i].ID != mediaData.ID {
			t.Errorf("got %v want %v", got[i].ID, mediaData.ID)
		}
		if got[i].MediaType != mediaData.MediaType {
			t.Errorf("got %v want %v", got[i].MediaType, mediaData.MediaType)
		}
		if got[i].URL != mediaData.URL {
			t.Errorf("got %v want %v", got[i].URL, mediaData.URL)
		}
	}
}

func assertAuthorEquals(t *testing.T, got, want *user.WithoutPrivateInfo) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertDatesEquals(t *testing.T, got, want []sospost.SOSDateView) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
