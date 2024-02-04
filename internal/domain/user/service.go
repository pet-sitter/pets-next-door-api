package user

import (
	"context"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
)

type UserService struct {
	userStore    UserStore
	petStore     pet.PetStore
	mediaService media.MediaService
}

func NewUserService(userStore UserStore, petStore pet.PetStore, mediaService media.MediaService) *UserService {
	return &UserService{
		userStore:    userStore,
		petStore:     petStore,
		mediaService: mediaService,
	}
}

func (service *UserService) RegisterUser(ctx context.Context, registerUserRequest *RegisterUserRequest) (*RegisterUserView, *pnd.AppError) {
	var media *media.Media
	var err *pnd.AppError
	if registerUserRequest.ProfileImageID != nil {
		media, err = service.mediaService.FindMediaByID(ctx, *registerUserRequest.ProfileImageID)
		if err != nil {
			return nil, err
		}
	}

	created, err := service.userStore.CreateUser(ctx, registerUserRequest)
	if err != nil {
		return nil, pnd.ErrUnknown(err)
	}

	var profileImageURL *string
	if media != nil {
		profileImageURL = &media.URL
	}

	return &RegisterUserView{
		ID:                   created.ID,
		Email:                created.Email,
		Nickname:             created.Nickname,
		Fullname:             created.Fullname,
		ProfileImageURL:      profileImageURL,
		FirebaseProviderType: created.FirebaseProviderType,
		FirebaseUID:          created.FirebaseUID,
	}, nil
}

func (service *UserService) FindUsers(ctx context.Context, page int, size int, nickname *string) (*UserWithoutPrivateInfoList, *pnd.AppError) {
	userList, err := service.userStore.FindUsers(ctx, page, size, nickname)
	if err != nil {
		return nil, err
	}

	return userList, nil
}

func (service *UserService) FindUserByEmail(ctx context.Context, email string) (*UserWithProfileImage, *pnd.AppError) {
	user, err := service.userStore.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *UserService) FindUserByUID(ctx context.Context, uid string) (*FindUserView, *pnd.AppError) {
	user, err := service.userStore.FindUserByUID(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &FindUserView{
		ID:                   user.ID,
		Email:                user.Email,
		Nickname:             user.Nickname,
		Fullname:             user.Fullname,
		ProfileImageURL:      user.ProfileImageURL,
		FirebaseProviderType: user.FirebaseProviderType,
		FirebaseUID:          user.FirebaseUID,
	}, nil
}

func (service *UserService) ExistsByNickname(ctx context.Context, nickname string) (bool, *pnd.AppError) {
	existsByNickname, err := service.userStore.ExistsByNickname(ctx, nickname)
	if err != nil {
		return false, pnd.ErrUnknown(err)
	}

	return existsByNickname, nil
}

func (service *UserService) FindUserStatusByEmail(ctx context.Context, email string) (*UserStatus, *pnd.AppError) {
	userStatus, err := service.userStore.FindUserStatusByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &UserStatus{
		FirebaseProviderType: userStatus.FirebaseProviderType,
	}, nil
}

func (service *UserService) UpdateUserByUID(ctx context.Context, uid string, nickname string, profileImageID *int) (*UserWithProfileImage, *pnd.AppError) {
	updated, err := service.userStore.UpdateUserByUID(ctx, uid, nickname, profileImageID)
	if err != nil {
		return nil, err
	}

	var profileImage *media.Media
	var err2 *pnd.AppError
	if updated.ProfileImageID != nil {
		profileImage, err2 = service.mediaService.FindMediaByID(ctx, *updated.ProfileImageID)
		if err != nil {
			return nil, err2
		}
	}

	var profileImageURL *string
	if profileImage != nil {
		profileImageURL = &profileImage.URL
	}
	return &UserWithProfileImage{
		ID:                   updated.ID,
		Email:                updated.Email,
		Nickname:             updated.Nickname,
		Fullname:             updated.Fullname,
		ProfileImageURL:      profileImageURL,
		FirebaseProviderType: updated.FirebaseProviderType,
		FirebaseUID:          updated.FirebaseUID,
	}, nil
}

func (service *UserService) AddPetsToOwner(ctx context.Context, uid string, addPetsRequest pet.AddPetsToOwnerRequest) ([]pet.PetView, *pnd.AppError) {
	pets := make([]pet.Pet, len(addPetsRequest.Pets))

	user, err := service.userStore.FindUserByUID(ctx, uid)
	if err != nil {
		return nil, err
	}

	for i, item := range addPetsRequest.Pets {
		pets[i] = pet.Pet{
			OwnerID:    user.ID,
			Name:       item.Name,
			PetType:    item.PetType,
			Sex:        item.Sex,
			Neutered:   item.Neutered,
			Breed:      item.Breed,
			BirthDate:  item.BirthDate,
			WeightInKg: item.WeightInKg,
		}

		if _, err := service.petStore.CreatePet(ctx, &pets[i]); err != nil {
			return nil, pnd.ErrUnknown(err)
		}
	}

	petViews := make([]pet.PetView, len(addPetsRequest.Pets))
	for i, item := range pets {
		petViews[i] = pet.PetView{
			ID:         item.ID,
			Name:       item.Name,
			PetType:    item.PetType,
			Sex:        item.Sex,
			Neutered:   item.Neutered,
			Breed:      item.Breed,
			BirthDate:  item.BirthDate,
			WeightInKg: item.WeightInKg,
		}
	}

	return petViews, nil
}

func (service *UserService) FindPetsByOwnerUID(ctx context.Context, uid string) (*pet.FindMyPetsView, *pnd.AppError) {
	user, err := service.userStore.FindUserByUID(ctx, uid)
	if err != nil {
		return nil, err
	}

	pets, err := service.petStore.FindPetsByOwnerID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	petViews := make([]pet.PetView, len(pets))
	for i, item := range pets {
		petViews[i] = pet.PetView{
			ID:         item.ID,
			Name:       item.Name,
			PetType:    item.PetType,
			Sex:        item.Sex,
			Neutered:   item.Neutered,
			Breed:      item.Breed,
			BirthDate:  item.BirthDate,
			WeightInKg: item.WeightInKg,
		}
	}

	return &pet.FindMyPetsView{
		Pets: petViews,
	}, nil
}
