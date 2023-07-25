package user

import (
	"github.com/pet-sitter/pets-next-door-api/internal/database"
	"github.com/pet-sitter/pets-next-door-api/internal/models"

	_ "github.com/lib/pq"
)

type UserService struct {
	db *database.DB
}

func NewUserService(db *database.DB) *UserService {
	return &UserService{
		db: db,
	}
}

type UserServicer interface {
	CreateUser(user *models.User) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	FindUserByUID(uid string) (*models.User, error)
	UpdateUserByUID(uid string, nickname string) (*models.User, error)
}

func (service *UserService) CreateUser(user *models.User) (*models.User, error) {
	tx, _ := service.db.Begin()

	created, err := tx.CreateUser(user)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return created, nil
}

func (service *UserService) FindUserByEmail(email string) (*models.User, error) {
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

func (service *UserService) FindUserByUID(uid string) (*models.User, error) {
	tx, _ := service.db.Begin()

	user, err := tx.FindUserByUID(uid)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (service *UserService) UpdateUserByUID(uid string, nickname string) (*models.User, error) {
	tx, _ := service.db.Begin()

	updated, err := tx.UpdateUserByUID(uid, nickname)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return updated, nil
}
