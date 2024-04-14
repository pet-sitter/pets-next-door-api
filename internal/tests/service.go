package tests

import (
	"context"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
	"testing"
)

func AddDummyMedia(t *testing.T, mediaService *service.MediaService) *media.Media {
	media, err := mediaService.CreateMedia(context.Background(), &media.Media{
		MediaType: media.IMAGE_MEDIA_TYPE,
		URL:       "http://example.com",
	})
	if err != nil {
		t.Errorf("got %v want %v", err, nil)
	}

	return media
}

func RegisterDummyUser(
	t *testing.T,
	ctx context.Context,
	userService *service.UserService,
	mediaService *service.MediaService,
) *user.RegisterUserView {
	profileImage := AddDummyMedia(t, mediaService)
	user, err := userService.RegisterUser(ctx, GenerateDummyRegisterUserRequest(&profileImage.ID))
	if err != nil {
		t.Errorf("got %v want %v", err, nil)
	}

	return user
}

func AddDummyPets(
	t *testing.T,
	ctx context.Context,
	userService *service.UserService,
	ownerUID string,
	profileImageID *int,
) *pet.PetView {
	pets, err := userService.AddPetsToOwner(ctx, ownerUID, pet.AddPetsToOwnerRequest{
		Pets: []pet.AddPetRequest{*GenerateDummyAddPetRequest(profileImageID)},
	})
	if err != nil {
		t.Errorf("got %v want %v", err, nil)
	}

	return &pets[0]
}
