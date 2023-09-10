package user

import (
	"github.com/pet-sitter/pets-next-door-api/internal/database"
	"github.com/pet-sitter/pets-next-door-api/internal/media"
	"github.com/pet-sitter/pets-next-door-api/internal/models"
	"github.com/pet-sitter/pets-next-door-api/internal/views"

	_ "github.com/lib/pq"
)

type UserService struct {
	db           *database.DB
	mediaService media.MediaServicer
}

func NewUserService(db *database.DB, mediaService media.MediaServicer) *UserService {
	return &UserService{
		db:           db,
		mediaService: mediaService,
	}
}

type UserServicer interface {
	RegisterUser(registerUserRequest *views.RegisterUserRequest) (*views.RegisterUserResponse, error)
	FindUserByEmail(email string) (*models.UserWithProfileImage, error)
	FindUserByUID(uid string) (*views.FindUserResponse, error)
	FindUserStatusByEmail(email string) (*models.UserStatus, error)
	UpdateUserByUID(uid string, nickname string, profileImageID int) (*models.UserWithProfileImage, error)
	AddPetsToOwner(uid string, addPetsRequest views.AddPetsToOwnerRequest) ([]views.PetView, error)
	FindPetsByOwnerUID(uid string) (*views.FindMyPetsView, error)
}

func (service *UserService) RegisterUser(registerUserRequest *views.RegisterUserRequest) (*views.RegisterUserResponse, error) {
	media, err := service.mediaService.FindMediaByID(registerUserRequest.ProfileImageID)
	if err != nil {
		return nil, err
	}

	tx, _ := service.db.Begin()

	created, err := tx.CreateUser(registerUserRequest)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &views.RegisterUserResponse{
		ID:                   created.ID,
		Email:                created.Email,
		Nickname:             created.Nickname,
		Fullname:             created.Fullname,
		ProfileImageURL:      media.URL,
		FirebaseProviderType: created.FirebaseProviderType,
		FirebaseUID:          created.FirebaseUID,
	}, nil
}

func (service *UserService) FindUserByEmail(email string) (*models.UserWithProfileImage, error) {
	tx, _ := service.db.Begin()

	user, err := tx.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (service *UserService) FindUserByUID(uid string) (*views.FindUserResponse, error) {
	tx, _ := service.db.Begin()

	user, err := tx.FindUserByUID(uid)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &views.FindUserResponse{
		ID:                   user.ID,
		Email:                user.Email,
		Nickname:             user.Nickname,
		Fullname:             user.Fullname,
		ProfileImageURL:      user.ProfileImageURL,
		FirebaseProviderType: user.FirebaseProviderType,
		FirebaseUID:          user.FirebaseUID,
	}, nil
}

func (service *UserService) FindUserStatusByEmail(email string) (*models.UserStatus, error) {
	tx, _ := service.db.Begin()

	userStatus, err := tx.FindUserStatusByEmail(email)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &models.UserStatus{
		FirebaseProviderType: userStatus.FirebaseProviderType,
	}, nil
}

func (service *UserService) UpdateUserByUID(uid string, nickname string, profileImageID int) (*models.UserWithProfileImage, error) {
	tx, _ := service.db.Begin()

	updated, err := tx.UpdateUserByUID(uid, nickname, profileImageID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	profileImage, err := service.mediaService.FindMediaByID(updated.ProfileImageID)
	if err != nil {
		return nil, err
	}

	return &models.UserWithProfileImage{
		ID:                   updated.ID,
		Email:                updated.Email,
		Nickname:             updated.Nickname,
		Fullname:             updated.Fullname,
		ProfileImageURL:      profileImage.URL,
		FirebaseProviderType: updated.FirebaseProviderType,
		FirebaseUID:          updated.FirebaseUID,
	}, nil
}

func (service *UserService) AddPetsToOwner(uid string, addPetsRequest views.AddPetsToOwnerRequest) ([]views.PetView, error) {
	pets := make([]models.Pet, len(addPetsRequest.Pets))

	tx, _ := service.db.Begin()

	user, err := tx.FindUserByUID(uid)
	if err != nil {
		return nil, err
	}

	for i, pet := range addPetsRequest.Pets {
		pets[i] = models.Pet{
			OwnerID:    user.ID,
			Name:       pet.Name,
			PetType:    pet.PetType,
			Sex:        pet.Sex,
			Neutered:   pet.Neutered,
			Breed:      pet.Breed,
			BirthDate:  pet.BirthDate,
			WeightInKg: pet.WeightInKg,
		}

		if _, err := tx.CreatePet(&pets[i]); err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	petViews := make([]views.PetView, len(addPetsRequest.Pets))
	for i, pet := range pets {
		petViews[i] = views.PetView{
			ID:         pet.ID,
			Name:       pet.Name,
			PetType:    pet.PetType,
			Sex:        pet.Sex,
			Neutered:   pet.Neutered,
			Breed:      pet.Breed,
			BirthDate:  pet.BirthDate,
			WeightInKg: pet.WeightInKg,
		}
	}

	return petViews, nil
}

func (service *UserService) FindPetsByOwnerUID(uid string) (*views.FindMyPetsView, error) {
	tx, _ := service.db.Begin()

	user, err := tx.FindUserByUID(uid)
	if err != nil {
		return nil, err
	}

	pets, err := tx.FindPetsByOwnerID(user.ID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	petViews := make([]views.PetView, len(pets))
	for i, pet := range pets {
		petViews[i] = views.PetView{
			ID:         pet.ID,
			Name:       pet.Name,
			PetType:    pet.PetType,
			Sex:        pet.Sex,
			Neutered:   pet.Neutered,
			Breed:      pet.Breed,
			BirthDate:  pet.BirthDate,
			WeightInKg: pet.WeightInKg,
		}
	}

	return &views.FindMyPetsView{
		Pets: petViews,
	}, nil
}
