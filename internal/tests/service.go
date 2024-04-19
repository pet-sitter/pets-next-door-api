package tests

import (
	"context"
	"testing"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sospost"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

func AddDummyMedia(t *testing.T, ctx context.Context, mediaService *service.MediaService) *media.Media {
	t.Helper()
	mediaData, err := mediaService.CreateMedia(ctx, &media.Media{
		MediaType: media.MediaTypeImage,
		URL:       "http://example.com",
	})
	if err != nil {
		t.Errorf("got %v want %v", err, nil)
	}

	return mediaData
}

func RegisterDummyUser(
	t *testing.T,
	ctx context.Context,
	userService *service.UserService,
	mediaService *service.MediaService,
) *user.RegisterUserView {
	t.Helper()
	profileImage := AddDummyMedia(t, ctx, mediaService)
	userRequest := GenerateDummyRegisterUserRequest(&profileImage.ID)
	registeredUser, err := userService.RegisterUser(ctx, userRequest)
	if err != nil {
		t.Errorf("got %v want %v", err, nil)
	}

	return registeredUser
}

func AddDummyPet(
	t *testing.T,
	ctx context.Context,
	userService *service.UserService,
	ownerUID string,
	profileImageID *int,
) *pet.PetView {
	t.Helper()
	pets, err := userService.AddPetsToOwner(ctx, ownerUID, pet.AddPetsToOwnerRequest{
		Pets: []pet.AddPetRequest{*GenerateDummyAddPetRequest(profileImageID)},
	})
	if err != nil {
		t.Errorf("got %v want %v", err, nil)
	}

	return &pets[0]
}

func AddDummyPets(
	t *testing.T,
	ctx context.Context,
	userService *service.UserService,
	ownerUID string,
	profileImageID *int,
) []pet.PetView {
	t.Helper()
	pets, err := userService.AddPetsToOwner(ctx, ownerUID, pet.AddPetsToOwnerRequest{
		Pets: GenerateDummyAddPetsRequest(profileImageID),
	})
	if err != nil {
		t.Errorf("got %v want %v", err, nil)
	}

	return pets
}

func WriteDummySOSPosts(
	t *testing.T,
	ctx context.Context,
	sosPostService *service.SOSPostService,
	uid string,
	imageID []int,
	petIDs []int,
	sosPostCnt int,
) *sospost.WriteSOSPostView {
	t.Helper()
	sosPostRequest := GenerateDummyWriteSOSPostRequest(imageID, petIDs, sosPostCnt)
	sosPost, err := sosPostService.WriteSOSPost(ctx, uid, sosPostRequest)
	if err != nil {
		t.Errorf("got %v want %v", err, nil)
	}
	return sosPost
}
