package database

import "github.com/pet-sitter/pets-next-door-api/internal/models"

type InMemoryDB struct {
	Users []models.UserModel
}

func NewInMemoryDB() *InMemoryDB {
	users := []models.UserModel{}

	return &InMemoryDB{
		Users: users,
	}
}

func (s *InMemoryDB) Migrate() error {
	return nil
}
