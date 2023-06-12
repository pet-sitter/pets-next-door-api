package user

import (
	"fmt"

	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"github.com/pet-sitter/pets-next-door-api/internal/models"
)

type UserInMemoryRepo struct {
	db *database.InMemoryDB
}

func NewUserInMemoryRepo(db *database.InMemoryDB) *UserInMemoryRepo {
	return &UserInMemoryRepo{
		db: db,
	}
}

func (repo *UserInMemoryRepo) CreateUser(user *models.UserModel) (*models.UserModel, error) {
	nextID := len(repo.db.Users) + 1
	user.ID = nextID

	repo.db.Users = append(repo.db.Users, *user)

	return user, nil
}

func (repo *UserInMemoryRepo) FindUserByEmail(email string) (*models.UserModel, error) {
	var user *models.UserModel

	for _, v := range repo.db.Users {
		if v.Email == email {
			user = &v
		}
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (repo *UserInMemoryRepo) FindUserByUID(uid string) (*models.UserModel, error) {
	var user *models.UserModel

	for _, v := range repo.db.Users {
		if v.UID == uid {
			user = &v
		}
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (repo *UserInMemoryRepo) UpdateUserByUID(uid string, nickname string) (*models.UserModel, error) {
	var user *models.UserModel
	for i, v := range repo.db.Users {
		if v.UID == uid {
			repo.db.Users[i].Nickname = nickname
			user = &repo.db.Users[i]
		}
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}
