package tests

import (
	"context"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	bucketinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/bucket"
	"io"
	"testing"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sospost"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type DummyUploader struct{}

func (d DummyUploader) UploadFile(_ io.ReadSeeker, fileName string) (url string, appError *pnd.AppError) {
	return "https://example.com/files/" + fileName, nil
}

func NewDummyFileUploader() bucketinfra.FileUploader {
	return DummyUploader{}
}

func AddDummyMedia(t *testing.T, ctx context.Context, mediaService *service.MediaService) *media.DetailView {
	t.Helper()
	mediaData, err := mediaService.CreateMedia(ctx, media.TypeImage, "http://example.com")
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
) *user.InternalView {
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
	profileImageID *int64,
) *pet.DetailView {
	t.Helper()
	petList, err := userService.AddPetsToOwner(ctx, ownerUID, pet.AddPetsToOwnerRequest{
		Pets: []pet.AddPetRequest{*GenerateDummyAddPetRequest(profileImageID)},
	})
	if err != nil {
		t.Errorf("got %v want %v", err, nil)
	}

	return &petList.Pets[0]
}

func AddDummyPets(
	t *testing.T,
	ctx context.Context,
	userService *service.UserService,
	ownerUID string,
	profileImageID *int64,
) pet.ListView {
	t.Helper()
	petList, err := userService.AddPetsToOwner(ctx, ownerUID, pet.AddPetsToOwnerRequest{
		Pets: GenerateDummyAddPetsRequest(profileImageID),
	})
	if err != nil {
		t.Errorf("got %v want %v", err, nil)
	}

	return *petList
}

func WriteDummySOSPosts(
	t *testing.T,
	ctx context.Context,
	sosPostService *service.SOSPostService,
	uid string,
	imageID []int64,
	petIDs []int64,
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
