package user

import "github.com/pet-sitter/pets-next-door-api/internal/models"

type UserRepo interface {
	CreateUser(user *models.UserModel) (*models.UserModel, error)
	FindUserByEmail(email string) (*models.UserModel, error)
	FindUserByUID(uid string) (*models.UserModel, error)
	UpdateUserByUID(uid string, nickname string) (*models.UserModel, error)
}
