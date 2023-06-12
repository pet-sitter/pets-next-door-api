package user

import "github.com/pet-sitter/pets-next-door-api/internal/models"

type UserService struct {
	userRepo UserRepo
}

func NewUserService(userRepo UserRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

type UserServicer interface {
	CreateUser(user *models.UserModel) (*models.UserModel, error)
	FindUserByEmail(email string) (*models.UserModel, error)
	FindUserByUID(uid string) (*models.UserModel, error)
	UpdateUserByUID(uid string, nickname string) (*models.UserModel, error)
}

func (service *UserService) CreateUser(user *models.UserModel) (*models.UserModel, error) {
	created, err := service.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (service *UserService) FindUserByEmail(email string) (*models.UserModel, error) {
	user, err := service.userRepo.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *UserService) FindUserByUID(uid string) (*models.UserModel, error) {
	user, err := service.userRepo.FindUserByUID(uid)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *UserService) UpdateUserByUID(uid string, nickname string) (*models.UserModel, error) {
	updated, err := service.userRepo.UpdateUserByUID(uid, nickname)
	if err != nil {
		return nil, err
	}

	return updated, nil
}
