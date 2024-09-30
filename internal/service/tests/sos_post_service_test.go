package service_test

// func setUp(ctx context.Context, t *testing.T) (*database.DB, func(t *testing.T)) {
// 	t.Helper()
// 	db, tearDown := tests.SetUp(t)
//
// 	conditionService := service.NewSOSConditionService(db)
// 	if _, err2 := conditionService.InitConditions(ctx); err2 != nil {
// 		t.Errorf("InitConditions failed: %v", err2)
// 	}
//
// 	return db, tearDown
// }
//
// func TestCreateSOSPost(t *testing.T) {
// 	t.Run("돌봄 급구 게시글을 작성한다", func(t *testing.T) {
// 		ctx := context.Background()
// 		db, tearDown := setUp(ctx, t)
// 		defer tearDown(t)
// 		mediaService := tests.NewMockMediaService(db)
// 		userService := tests.NewMockUserService(db)
// 		sosPostService := tests.NewMockSOSPostService(db)
//
// 		// given
// 		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")
// 		owner, _ := userService.RegisterUser(ctx, tests.NewDummyRegisterUserRequest(&profileImage.ID))
// 		addPets, _ := userService.AddPetsToOwner(
// 			ctx,
// 			owner.FirebaseUID,
// 			pet.AddPetsToOwnerRequest{Pets: []pet.AddPetRequest{
// 				*tests.NewDummyAddPetRequest(&profileImage.ID, commonvo.PetTypeDog, pet.GenderMale, "poodle"),
// 			}},
// 		)
//
// 		// when
// 		sosPostImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "sos_post_image.jpg")
// 		sosPostImage2, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "sos_post_image2.jpg")
// 		sosPostData := tests.NewDummyWriteSOSPostRequest(
// 			[]int64{sosPostImage.ID, sosPostImage2.ID},
// 			[]int64{addPets.Pets[0].ID},
// 			0,
// 		)
// 		created, _ := sosPostService.WriteSOSPost(ctx, owner.FirebaseUID, tests.NewDummyWriteSOSPostRequest(
// 			[]int64{sosPostImage.ID, sosPostImage2.ID},
// 			[]int64{addPets.Pets[0].ID},
// 			0,
// 		))
//
// 		// then
// 		found, _ := sosPostService.FindSOSPostByID(ctx, created.ID)
// 		asserts.ConditionIDEquals(t, sosPostData.ConditionIDs, found.Conditions)
// 		assertPetEquals(t, addPets.Pets[0], found.Pets[0])
// 		asserts.MediaEquals(t, media.ListView{sosPostImage, sosPostImage2}, found.Media)
// 		asserts.DatesEquals(t, sosPostData.Dates, found.Dates)
// 		writtenAndFoundSOSPostEquals(t, *sosPostData, *found)
// 		assert.Equal(t, owner.ID, int64(created.AuthorID))
// 	})
// }
//
// func TestFindSOSPosts(t *testing.T) {
// 	t.Run("전체 돌봄 급구 게시글을 조회한다", func(t *testing.T) {
// 		ctx := context.Background()
// 		db, tearDown := setUp(ctx, t)
// 		defer tearDown(t)
// 		mediaService := tests.NewMockMediaService(db)
// 		userService := tests.NewMockUserService(db)
// 		sosPostService := tests.NewMockSOSPostService(db)
//
// 		// given
// 		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")
// 		sosPostImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "sos_post_image.jpg")
// 		sosPostImage2, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "sos_post_image2.jpg")
//
// 		userRequest := tests.NewDummyRegisterUserRequest(&profileImage.ID)
// 		owner, _ := userService.RegisterUser(ctx, userRequest)
// 		addPets, _ := userService.AddPetsToOwner(
// 			ctx,
// 			owner.FirebaseUID,
// 			pet.AddPetsToOwnerRequest{
// 				Pets: []pet.AddPetRequest{
// 					*tests.NewDummyAddPetRequest(&profileImage.ID, commonvo.PetTypeDog, pet.GenderMale, "poodle"),
// 				},
// 			},
// 		)
//
// 		imageIDs := []int64{sosPostImage.ID, sosPostImage2.ID}
// 		petIDs := []int64{addPets.Pets[0].ID}
// 		conditionIDs := []int{1, 2}
//
// 		sosPostRequests := make([]sospost.WriteSOSPostRequest, 0)
// 		var sosPosts []sospost.DetailView
// 		for i := 1; i < 4; i++ {
// 			request := tests.NewDummyWriteSOSPostRequest(imageIDs, petIDs, i)
// 			sosPost, _ := sosPostService.WriteSOSPost(
// 				ctx, owner.FirebaseUID, request,
// 			)
// 			sosPosts = append(sosPosts, *sosPost)
// 			sosPostRequests = append(sosPostRequests, *request)
// 		}
//
// 		// when
// 		foundList, err := sosPostService.FindSOSPosts(ctx, 1, 3, "newest", "all")
// 		if err != nil {
// 			t.Errorf("got %v want %v", err, nil)
// 		}
//
// 		// then
// 		for i, found := range foundList.Items {
// 			idx := len(foundList.Items) - i - 1
// 			asserts.ConditionIDEquals(t, conditionIDs, found.Conditions)
// 			assertPetEquals(t, addPets.Pets[0], found.Pets[0])
// 			asserts.MediaEquals(t, media.ListView{sosPostImage, sosPostImage2}, found.Media)
// 			assert.Equal(t, owner.ID, found.Author.ID)
// 			asserts.DatesEquals(t, sosPosts[idx].Dates, found.Dates)
// 			writtenAndFoundSOSPostEquals(t, sosPostRequests[idx], found)
// 		}
// 	})
//
// 	t.Run("전체 돌봄 급구 게시글의 정렬기준을 'deadline'으로 조회한다", func(t *testing.T) {
// 		ctx := context.Background()
// 		db, tearDown := setUp(ctx, t)
// 		defer tearDown(t)
// 		mediaService := tests.NewMockMediaService(db)
// 		userService := tests.NewMockUserService(db)
// 		sosPostService := tests.NewMockSOSPostService(db)
//
// 		// given
// 		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")
// 		sosPostImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "sos_post_image.jpg")
// 		sosPostImage2, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "sos_post_image2.jpg")
//
// 		userRequest := tests.NewDummyRegisterUserRequest(&profileImage.ID)
// 		owner, _ := userService.RegisterUser(ctx, userRequest)
// 		uid := owner.FirebaseUID
// 		addPetRequest, _ := userService.AddPetsToOwner(ctx, uid, pet.AddPetsToOwnerRequest{
// 			Pets: []pet.AddPetRequest{
// 				*tests.NewDummyAddPetRequest(&profileImage.ID, commonvo.PetTypeDog, pet.GenderMale, "poodle"),
// 				*tests.NewDummyAddPetRequest(&profileImage.ID, commonvo.PetTypeDog, pet.GenderMale, "poodle"),
// 				*tests.NewDummyAddPetRequest(&profileImage.ID, commonvo.PetTypeCat, pet.GenderMale, "munchkin"),
// 			},
// 		})
//
// 		imageIDs := []int64{sosPostImage.ID, sosPostImage2.ID}
// 		conditionIDs := []int{1, 2}
//
// 		sosPostRequests := make([]sospost.WriteSOSPostRequest, 0)
// 		for i := 1; i < 4; i++ {
// 			request := tests.NewDummyWriteSOSPostRequest(imageIDs, []int64{addPetRequest.Pets[i-1].ID}, i)
// 			sosPostService.WriteSOSPost(ctx, uid, request)
// 			sosPostRequests = append(sosPostRequests, *request)
// 		}
//
// 		// when
// 		sosPostList, _ := sosPostService.FindSOSPosts(ctx, 1, 3, "newest", "all")
//
// 		// then
// 		for i, sosPost := range sosPostList.Items {
// 			idx := len(sosPostList.Items) - i - 1
// 			asserts.ConditionIDEquals(t, conditionIDs, sosPost.Conditions)
// 			assertPetEquals(t, addPetRequest.Pets[i-1], sosPost.Pets[i-1])
// 			asserts.MediaEquals(t, media.ListView{sosPostImage, sosPostImage2}, sosPost.Media)
// 			assert.Equal(t, owner.ID, sosPost.Author.ID)
// 			asserts.DatesEquals(t, sosPostRequests[idx].Dates, sosPost.Dates)
// 			writtenAndFoundSOSPostEquals(t, sosPostRequests[idx], sosPost)
// 		}
// 	})
//
// 	t.Run("전체 돌봄 급구 게시글 중 반려동물이 'dog'인 경우만 조회한다", func(t *testing.T) {
// 		ctx := context.Background()
// 		db, tearDown := setUp(ctx, t)
// 		defer tearDown(t)
// 		mediaService := tests.NewMockMediaService(db)
// 		userService := tests.NewMockUserService(db)
// 		sosPostService := tests.NewMockSOSPostService(db)
//
// 		// given
// 		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")
// 		owner, _ := userService.RegisterUser(ctx, tests.NewDummyRegisterUserRequest(&profileImage.ID))
//
// 		sosPostImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "sos_post_image.jpg")
// 		sosPostImage2, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "sos_post_image2.jpg")
// 		uid := owner.FirebaseUID
// 		petList, _ := userService.AddPetsToOwner(ctx, uid, pet.AddPetsToOwnerRequest{
// 			Pets: []pet.AddPetRequest{
// 				*tests.NewDummyAddPetRequest(&profileImage.ID, commonvo.PetTypeDog, pet.GenderMale, "poodle"),
// 				*tests.NewDummyAddPetRequest(&profileImage.ID, commonvo.PetTypeDog, pet.GenderMale, "poodle"),
// 				*tests.NewDummyAddPetRequest(&profileImage.ID, commonvo.PetTypeCat, pet.GenderMale, "munchkin"),
// 			},
// 		})
//
// 		imageIDs := []int64{sosPostImage.ID, sosPostImage2.ID}
// 		conditionIDs := []int{1, 2}
//
// 		writeRequests := make([]sospost.WriteSOSPostRequest, 0)
// 		// 강아지인 경우
// 		for i := 1; i < 3; i++ {
// 			request := tests.NewDummyWriteSOSPostRequest(imageIDs, []int64{petList.Pets[i-1].ID}, i)
// 			sosPostService.WriteSOSPost(ctx, uid, request)
// 			writeRequests = append(writeRequests, *request)
// 		}
//
// 		// 고양이인 경우
// 		request := tests.NewDummyWriteSOSPostRequest(imageIDs, []int64{petList.Pets[2].ID}, 3)
// 		sosPostService.WriteSOSPost(ctx, uid, request)
// 		writeRequests = append(writeRequests, *request)
//
// 		// 강아지, 고양이인 경우
// 		request = tests.NewDummyWriteSOSPostRequest(imageIDs, []int64{petList.Pets[1].ID, petList.Pets[2].ID}, 4)
// 		sosPostService.WriteSOSPost(ctx, uid, request)
// 		writeRequests = append(writeRequests, *request)
//
// 		// when
// 		foundList, _ := sosPostService.FindSOSPosts(ctx, 1, 3, "newest", "all")
//
// 		// then
// 		for i, sosPost := range foundList.Items {
// 			idx := len(foundList.Items) - i - 1
// 			asserts.ConditionIDEquals(t, conditionIDs, sosPost.Conditions)
// 			assertPetEquals(t, petList.Pets[i-1], sosPost.Pets[i-1])
// 			asserts.MediaEquals(t, media.ListView{sosPostImage, sosPostImage2}, sosPost.Media)
// 			assert.Equal(t, owner.ID, sosPost.Author.ID)
// 			asserts.DatesEquals(t, writeRequests[idx].Dates, sosPost.Dates)
// 			writtenAndFoundSOSPostEquals(t, writeRequests[idx], sosPost)
// 		}
// 	})
//
// 	t.Run("작성자 ID로 돌봄 급구 게시글을 조회합니다.", func(t *testing.T) {
// 		ctx := context.Background()
// 		db, tearDown := setUp(ctx, t)
// 		defer tearDown(t)
// 		mediaService := tests.NewMockMediaService(db)
// 		userService := tests.NewMockUserService(db)
// 		sosPostService := tests.NewMockSOSPostService(db)
//
// 		// given
// 		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")
// 		sosPostImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "sos_post_image.jpg")
// 		sosPostImage2, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "sos_post_image2.jpg")
//
// 		userRequest := tests.NewDummyRegisterUserRequest(&profileImage.ID)
// 		owner, _ := userService.RegisterUser(ctx, userRequest)
// 		author := &user.WithoutPrivateInfo{
// 			ID:              owner.ID,
// 			ProfileImageURL: owner.ProfileImageURL,
// 			Nickname:        owner.Nickname,
// 		}
// 		petList, _ := userService.AddPetsToOwner(
// 			ctx,
// 			owner.FirebaseUID,
// 			pet.AddPetsToOwnerRequest{
// 				Pets: []pet.AddPetRequest{
// 					*tests.NewDummyAddPetRequest(&profileImage.ID, commonvo.PetTypeDog, pet.GenderMale, "poodle"),
// 				},
// 			},
// 		)
//
// 		imageIDs := []int64{sosPostImage.ID, sosPostImage2.ID}
// 		conditionIDs := []int{1, 2}
//
// 		writeRequests := make([]sospost.WriteSOSPostRequest, 0)
// 		for i := 1; i < 4; i++ {
// 			request := tests.NewDummyWriteSOSPostRequest(imageIDs, []int64{petList.Pets[0].ID}, i)
// 			sosPostService.WriteSOSPost(ctx, owner.FirebaseUID, request)
// 			writeRequests = append(writeRequests, *request)
// 		}
//
// 		// when
// 		foundList, _ := sosPostService.FindSOSPostsByAuthorID(ctx, int(owner.ID), 1, 3, "newest", "all")
//
// 		// then
// 		for i, sosPost := range foundList.Items {
// 			idx := len(foundList.Items) - i - 1
// 			asserts.ConditionIDEquals(t, conditionIDs, sosPost.Conditions)
// 			assertPetEquals(t, petList.Pets[0], sosPost.Pets[0])
// 			asserts.MediaEquals(t, media.ListView{sosPostImage, sosPostImage2}, sosPost.Media)
// 			assert.Equal(t, author.ID, sosPost.Author.ID)
// 			asserts.DatesEquals(t, writeRequests[idx].Dates, sosPost.Dates)
// 			writtenAndFoundSOSPostEquals(t, writeRequests[idx], sosPost)
// 		}
// 	})
// }
//
// func TestFindSOSPostByID(t *testing.T) {
// 	t.Run("게시글 ID로 돌봄 급구 게시글을 조회합니다.", func(t *testing.T) {
// 		ctx := context.Background()
// 		db, tearDown := setUp(ctx, t)
// 		defer tearDown(t)
// 		mediaService := tests.NewMockMediaService(db)
// 		userService := tests.NewMockUserService(db)
// 		sosPostService := tests.NewMockSOSPostService(db)
//
// 		// given
// 		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")
// 		sosPostImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "sos_post_image.jpg")
// 		sosPostImage2, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "sos_post_image2.jpg")
//
// 		userRequest := tests.NewDummyRegisterUserRequest(&profileImage.ID)
// 		owner, _ := userService.RegisterUser(ctx, userRequest)
// 		author := &user.WithoutPrivateInfo{
// 			ID:              owner.ID,
// 			ProfileImageURL: owner.ProfileImageURL,
// 			Nickname:        owner.Nickname,
// 		}
// 		addPets, _ := userService.AddPetsToOwner(
// 			ctx,
// 			owner.FirebaseUID,
// 			pet.AddPetsToOwnerRequest{Pets: []pet.AddPetRequest{
// 				*tests.NewDummyAddPetRequest(&profileImage.ID, commonvo.PetTypeDog, pet.GenderMale, "poodle"),
// 			}},
// 		)
//
// 		imageIDs := []int64{sosPostImage.ID, sosPostImage2.ID}
// 		conditionIDs := []int{1, 2}
//
// 		writeRequests := make([]sospost.WriteSOSPostRequest, 0)
// 		writtenIDs := make([]int64, 0)
// 		for i := 1; i < 4; i++ {
// 			request := tests.NewDummyWriteSOSPostRequest(imageIDs, []int64{addPets.Pets[0].ID}, i)
// 			sosPost, _ := sosPostService.WriteSOSPost(ctx, owner.FirebaseUID, request)
// 			writeRequests = append(writeRequests, *request)
// 			writtenIDs = append(writtenIDs, int64(sosPost.ID))
// 		}
//
// 		// when
// 		found, _ := sosPostService.FindSOSPostByID(ctx, int(writtenIDs[0]))
//
// 		// then
// 		asserts.ConditionIDEquals(t, conditionIDs, found.Conditions)
// 		assertPetEquals(t, addPets.Pets[0], found.Pets[0])
// 		asserts.MediaEquals(t, media.ListView{sosPostImage, sosPostImage2}, found.Media)
// 		assert.Equal(t, author.ID, found.Author.ID)
// 		asserts.DatesEquals(t, writeRequests[0].Dates, found.Dates)
// 		writtenAndFoundSOSPostEquals(t, writeRequests[0], *found)
// 	})
// }
//
// func TestUpdateSOSPost(t *testing.T) {
// 	t.Run("돌봄 급구 게시글을 수정합니다.", func(t *testing.T) {
// 		ctx := context.Background()
// 		db, tearDown := setUp(ctx, t)
// 		defer tearDown(t)
// 		mediaService := tests.NewMockMediaService(db)
// 		userService := tests.NewMockUserService(db)
// 		sosPostService := tests.NewMockSOSPostService(db)
//
// 		// given
// 		profileImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "profile_image.jpg")
// 		sosPostImage, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "sos_post_image.jpg")
// 		sosPostImage2, _ := mediaService.UploadMedia(ctx, nil, media.TypeImage, "sos_post_image2.jpg")
//
// 		userRequest := tests.NewDummyRegisterUserRequest(&profileImage.ID)
// 		owner, _ := userService.RegisterUser(ctx, userRequest)
// 		addPets, _ := userService.AddPetsToOwner(
// 			ctx,
// 			owner.FirebaseUID,
// 			pet.AddPetsToOwnerRequest{Pets: []pet.AddPetRequest{
// 				*tests.NewDummyAddPetRequest(&profileImage.ID, commonvo.PetTypeDog, pet.GenderMale, "poodle"),
// 			}},
// 		)
//
// 		sosPost, _ := sosPostService.WriteSOSPost(
// 			ctx, owner.FirebaseUID,
//			tests.NewDummyWriteSOSPostRequest([]int64{sosPostImage.ID}, []int64{addPets.Pets[0].ID}, 1),
// 		)
//
// 		// when
// 		updateRequest := &sospost.UpdateSOSPostRequest{
// 			ID:       sosPost.ID,
// 			Title:    "Title2",
// 			Content:  "Content2",
// 			ImageIDs: []int64{sosPostImage.ID, sosPostImage2.ID},
// 			Reward:   "Reward2",
// 			Dates: []sospost.SOSDateView{
// 				{"2024-04-10", "2024-04-20"},
// 				{"2024-05-10", "2024-05-20"},
// 			},
// 			CareType:     sospost.CareTypeFoster,
// 			CarerGender:  sospost.CarerGenderMale,
// 			RewardType:   sospost.RewardTypeFee,
// 			ConditionIDs: []int{1, 2},
// 			PetIDs:       []int64{addPets.Pets[0].ID},
// 		}
// 		updated, _ := sosPostService.UpdateSOSPost(ctx, updateRequest)
//
// 		// then
// 		found, _ := sosPostService.FindSOSPostByID(ctx, sosPost.ID)
// 		asserts.ConditionIDEquals(t, updateRequest.ConditionIDs, found.Conditions)
// 		assertPetEquals(t, sosPost.Pets[0], found.Pets[0])
// 		asserts.MediaEquals(t, media.ListView{sosPostImage, sosPostImage2}, found.Media)
// 		asserts.DatesEquals(t, updateRequest.Dates, found.Dates)
// 		assert.Equal(t, updateRequest.Title, found.Title)
// 		assert.Equal(t, updateRequest.Content, found.Content)
// 		assert.Equal(t, updateRequest.Reward, found.Reward)
// 		assert.Equal(t, updateRequest.CareType, found.CareType)
// 		assert.Equal(t, updateRequest.CarerGender, found.CarerGender)
// 		assert.Equal(t, updateRequest.RewardType, found.RewardType)
// 		assert.Equal(t, &updateRequest.ImageIDs[0], found.ThumbnailID)
// 		assert.Equal(t, int64(updated.AuthorID), owner.ID)
// 	})
// }
//
// func assertPetEquals(t *testing.T, want, got pet.DetailView) {
// 	t.Helper()
//
// 	assert.Equal(t, want.ID, got.ID)
// 	assert.Equal(t, want.Name, got.Name)
// 	assert.Equal(t, want.PetType, got.PetType)
// 	assert.Equal(t, want.Sex, got.Sex)
// 	assert.Equal(t, want.Neutered, got.Neutered)
// 	assert.Equal(t, want.Breed, got.Breed)
// 	assert.Equal(t, want.BirthDate, got.BirthDate)
// 	assert.Equal(t, want.WeightInKg.String(), got.WeightInKg.String())
// 	assert.Equal(t, want.Remarks, got.Remarks)
// 	assert.Equal(t, want.ProfileImageURL, got.ProfileImageURL)
// }
//
// func writtenAndFoundSOSPostEquals(t *testing.T, want sospost.WriteSOSPostRequest, got sospost.FindSOSPostView) {
// 	t.Helper()
//
// 	assert.Equal(t, want.Title, got.Title)
// 	assert.Equal(t, want.Content, got.Content)
// 	assert.Equal(t, want.Reward, got.Reward)
// 	assert.Equal(t, want.CareType, got.CareType)
// 	assert.Equal(t, want.CarerGender, got.CarerGender)
// 	assert.Equal(t, want.RewardType, got.RewardType)
// 	assert.Equal(t, &want.ImageIDs[0], got.ThumbnailID)
// }
//
