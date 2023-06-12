package user

import (
	"testing"

	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"github.com/pet-sitter/pets-next-door-api/internal/models"
)

func TestUserInMemoryRepo(t *testing.T) {
	t.Run("CreateUser", func(t *testing.T) {
		t.Run("should create user", func(t *testing.T) {
			repo := NewUserInMemoryRepo(database.NewInMemoryDB())

			expected := &models.UserModel{
				ID:       1,
				UID:      "uid",
				Email:    "test@example.com",
				Password: "password",
			}

			created, err := repo.CreateUser(expected)
			if err != nil {
				t.Errorf("expected nil, got %s", err.Error())
			}

			if created == nil {
				t.Errorf("expected user, got nil")
			}

			actual, err := repo.FindUserByEmail(expected.Email)
			if err != nil {
				t.Errorf("expected nil, got %s", err.Error())
			}
			if actual == nil {
				t.Errorf("expected user, got nil")
			}

			assertUserEquals(t, expected, actual)
		})
	})

	t.Run("FindUserByEmail", func(t *testing.T) {
		t.Run("should return user when user exists", func(t *testing.T) {
			repo := NewUserInMemoryRepo(database.NewInMemoryDB())

			email := "test@example.com"
			expected := &models.UserModel{
				ID:       1,
				UID:      "1234",
				Email:    email,
				Password: "password",
			}
			repo.db.Users = append(repo.db.Users, *expected)

			actual, err := repo.FindUserByEmail(email)
			if err != nil {
				t.Errorf("expected nil, got %s", err.Error())
			}
			if actual == nil {
				t.Errorf("expected user, got nil")
			}

			assertUserEquals(t, expected, actual)
		})
	})

	t.Run("FindUserByUID", func(t *testing.T) {
		t.Run("should return user when user exists", func(t *testing.T) {
			repo := NewUserInMemoryRepo(database.NewInMemoryDB())

			uid := "uid"
			expected := &models.UserModel{
				ID:       1,
				UID:      uid,
				Email:    "test@example.com",
				Password: "password",
			}
			repo.db.Users = append(repo.db.Users, *expected)

			actual, err := repo.FindUserByUID(uid)
			if err != nil {
				t.Errorf("expected nil, got %s", err.Error())
			}
			if actual == nil {
				t.Errorf("expected user, got nil")
			}

			assertUserEquals(t, expected, actual)
		})
	})

	t.Run("UpdateUserByUID", func(t *testing.T) {
		t.Run("should update user by UID", func(t *testing.T) {
			repo := NewUserInMemoryRepo(database.NewInMemoryDB())

			nickname := "test"
			expected := &models.UserModel{
				ID:       1,
				UID:      "1234",
				Email:    "test@example.com",
				Nickname: nickname,
				Password: "password",
			}
			repo.db.Users = append(repo.db.Users, *expected)

			expected.Email = "updated@example.com"
			updated, err := repo.UpdateUserByUID(expected.UID, nickname)
			if err != nil {
				t.Errorf("expected nil, got %s", err.Error())
			}
			if updated == nil {
				t.Errorf("expected user, got nil")
			}

			actual, err := repo.FindUserByUID(expected.UID)
			if err != nil {
				t.Errorf("expected nil, got %s", err.Error())
			}
			if actual == nil {
				t.Errorf("expected user, got nil")
			}

			if actual.Nickname != nickname {
				t.Errorf("expected %s, got %s", nickname, actual.Nickname)
			}
		})
	})
}

func assertUserEquals(t testing.TB, expected *models.UserModel, actual *models.UserModel) {
	t.Helper()

	if actual.ID != expected.ID {
		t.Errorf("expected %d, got %d", expected.ID, actual.ID)
	}

	if actual.UID != expected.UID {
		t.Errorf("expected %s, got %s", expected.UID, actual.UID)
	}

	if actual.Email != expected.Email {
		t.Errorf("expected %s, got %s", expected.Email, actual.Email)
	}
}
