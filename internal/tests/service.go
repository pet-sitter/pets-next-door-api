package tests

import (
	"context"
	"io"
	"testing"

	"github.com/google/uuid"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/commonvo"

	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	bucketinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/bucket"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type StubUploader struct{}

func (u StubUploader) UploadFile(_ io.ReadSeeker, fileName string) (url string, appError *pnd.AppError) {
	return "https://example.com/files/" + fileName, nil
}

func NewStubFileUploader() bucketinfra.FileUploader {
	return StubUploader{}
}

func NewMockMediaService(db *database.DB) *service.MediaService {
	return service.NewMediaService(db, NewStubFileUploader())
}

func NewMockUserService(db *database.DB) *service.UserService {
	return service.NewUserService(db, NewMockMediaService(db))
}

// func NewMockSOSPostService(db *database.DB) *service.SOSPostService {
// 	return service.NewSOSPostService(db)
// }

func AddDummyPet(
	t *testing.T,
	ctx context.Context,
	userService *service.UserService,
	ownerUID string,
	profileImageID uuid.NullUUID,
) *pet.DetailView {
	t.Helper()
	petList, err := userService.AddPetsToOwner(ctx, ownerUID, pet.AddPetsToOwnerRequest{
		Pets: []pet.AddPetRequest{*NewDummyAddPetRequest(profileImageID, commonvo.PetTypeDog, pet.GenderMale, "poodle")},
	})
	if err != nil {
		t.Errorf("got %v want %v", err, nil)
	}

	return &petList.Pets[0]
}
