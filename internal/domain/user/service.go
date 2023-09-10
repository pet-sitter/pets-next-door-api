package user

import (
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/media"
)

type UserService struct {
	userStore    UserStore
	petStore     pet.PetStore
	mediaService media.MediaServicer
}

func NewUserService(userStore UserStore, petStore pet.PetStore, mediaService media.MediaServicer) *UserService {
	return &UserService{
		userStore:    userStore,
		petStore:     petStore,
		mediaService: mediaService,
	}
}

type UserServicer interface {
	RegisterUser(registerUserRequest *RegisterUserRequest) (*RegisterUserResponse, error)
	FindUserByEmail(email string) (*UserWithProfileImage, error)
	FindUserByUID(uid string) (*FindUserResponse, error)
	FindUserStatusByEmail(email string) (*UserStatus, error)
	UpdateUserByUID(uid string, nickname string, profileImageID int) (*UserWithProfileImage, error)
	AddPetsToOwner(uid string, addPetsRequest pet.AddPetsToOwnerRequest) ([]pet.PetView, error)
	FindPetsByOwnerUID(uid string) (*pet.FindMyPetsView, error)
}

func (service *UserService) RegisterUser(registerUserRequest *RegisterUserRequest) (*RegisterUserResponse, error) {
	media, err := service.mediaService.FindMediaByID(registerUserRequest.ProfileImageID)
	if err != nil {
		return nil, err
	}

	created, err := service.userStore.CreateUser(registerUserRequest)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &RegisterUserResponse{
		ID:                   created.ID,
		Email:                created.Email,
		Nickname:             created.Nickname,
		Fullname:             created.Fullname,
		ProfileImageURL:      media.URL,
		FirebaseProviderType: created.FirebaseProviderType,
		FirebaseUID:          created.FirebaseUID,
	}, nil
}

func (service *UserService) FindUserByEmail(email string) (*UserWithProfileImage, error) {
	user, err := service.userStore.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *UserService) FindUserByUID(uid string) (*FindUserResponse, error) {
	user, err := service.userStore.FindUserByUID(uid)
	if err != nil {
		return nil, err
	}

	return &FindUserResponse{
		ID:                   user.ID,
		Email:                user.Email,
		Nickname:             user.Nickname,
		Fullname:             user.Fullname,
		ProfileImageURL:      user.ProfileImageURL,
		FirebaseProviderType: user.FirebaseProviderType,
		FirebaseUID:          user.FirebaseUID,
	}, nil
}

func (service *UserService) FindUserStatusByEmail(email string) (*UserStatus, error) {
	userStatus, err := service.userStore.FindUserStatusByEmail(email)
	if err != nil {
		return nil, err
	}

	return &UserStatus{
		FirebaseProviderType: userStatus.FirebaseProviderType,
	}, nil
}

func (service *UserService) UpdateUserByUID(uid string, nickname string, profileImageID int) (*UserWithProfileImage, error) {
	updated, err := service.userStore.UpdateUserByUID(uid, nickname, profileImageID)
	if err != nil {
		return nil, err
	}

	profileImage, err := service.mediaService.FindMediaByID(updated.ProfileImageID)
	if err != nil {
		return nil, err
	}

	return &UserWithProfileImage{
		ID:                   updated.ID,
		Email:                updated.Email,
		Nickname:             updated.Nickname,
		Fullname:             updated.Fullname,
		ProfileImageURL:      profileImage.URL,
		FirebaseProviderType: updated.FirebaseProviderType,
		FirebaseUID:          updated.FirebaseUID,
	}, nil
}

func (service *UserService) AddPetsToOwner(uid string, addPetsRequest pet.AddPetsToOwnerRequest) ([]pet.PetView, error) {
	pets := make([]pet.Pet, len(addPetsRequest.Pets))

	user, err := service.userStore.FindUserByUID(uid)
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

		if _, err := service.petStore.CreatePet(&pets[i]); err != nil {
			return nil, err
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

func (service *UserService) FindPetsByOwnerUID(uid string) (*pet.FindMyPetsView, error) {
	user, err := service.userStore.FindUserByUID(uid)
	if err != nil {
		return nil, err
	}

	pets, err := service.petStore.FindPetsByOwnerID(user.ID)
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
