package user

import (
	"testing"

	"github.com/pet-sitter/pets-next-door-api/internal/database"
	"github.com/pet-sitter/pets-next-door-api/internal/models"
	"github.com/pet-sitter/pets-next-door-api/internal/tests"
)

var db *database.DB

func setUp(t *testing.T) func(t *testing.T) {
	db, _ = database.Open(tests.TestDatabaseURL)
	db.Flush()

	return func(t *testing.T) {
		db.Close()
	}
}

func TestUserService(t *testing.T) {

	t.Run("CreateUser", func(t *testing.T) {
		t.Run("사용자를 새로 생성한다", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			service := NewUserService(db)

			user := &models.User{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				Password:             "password",
				FirebaseProviderType: "kakao",
				FirebaseUID:          "uid",
			}

			created, _ := service.CreateUser(user)

			if created.Email != user.Email {
				t.Errorf("got %v want %v", created.ID, user.ID)
			}
		})

		t.Run("사용자가 이미 존재할 경우 에러를 반환한다", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			service := NewUserService(db)

			user := &models.User{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				Password:             "password",
				FirebaseProviderType: "kakao",
				FirebaseUID:          "uid",
			}

			_, _ = service.CreateUser(user)
			_, err := service.CreateUser(user)

			if err == nil {
				t.Errorf("got %v want %v", err, nil)
			}
		})
	})

	t.Run("FindUserByEmail", func(t *testing.T) {
		t.Run("사용자를 이메일로 찾는다", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			service := NewUserService(db)

			user := &models.User{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				Password:             "password",
				FirebaseProviderType: "kakao",
				FirebaseUID:          "uid",
			}

			created, _ := service.CreateUser(user)

			found, err := service.FindUserByEmail(created.Email)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			if found.Email != user.Email {
				t.Errorf("got %v want %v", found.Email, user.Email)
			}
		})

		t.Run("사용자가 존재하지 않을 경우 에러를 반환한다", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			service := NewUserService(db)

			_, err := service.FindUserByEmail("non-existent@example.com")
			if err == nil {
				t.Errorf("got %v want %v", err, nil)
			}
		})
	})

	t.Run("FindUserByUID", func(t *testing.T) {
		t.Run("사용자를 UID로 찾는다", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			service := NewUserService(db)

			user := &models.User{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				Password:             "password",
				FirebaseProviderType: "kakao",
				FirebaseUID:          "uid",
			}

			created, _ := service.CreateUser(user)

			found, err := service.FindUserByUID(created.FirebaseUID)

			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			if found.FirebaseUID != user.FirebaseUID {
				t.Errorf("got %v want %v", found.FirebaseUID, user.FirebaseUID)
			}
		})

		t.Run("사용자가 존재하지 않을 경우 에러를 반환한다", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			service := NewUserService(db)

			_, err := service.FindUserByUID("non-existent")
			if err == nil {
				t.Errorf("got %v want %v", err, nil)
			}
		})
	})

	t.Run("FindUserStatusByEmail", func(t *testing.T) {
		t.Run("사용자의 상태를 반환한다", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			service := NewUserService(db)

			user := &models.User{
				Email:                "pnd@example.com",
				Password:             "password",
				Nickname:             "nickname",
				Fullname:             "fullname",
				FirebaseProviderType: "kakao",
				FirebaseUID:          "uid",
			}

			created, err := service.CreateUser(user)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			status, err := service.FindUserStatusByEmail(created.Email)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			if status.FirebaseProviderType != user.FirebaseProviderType {
				t.Errorf("got %v want %v", status.FirebaseProviderType, user.FirebaseProviderType)
			}
		})
	})

	t.Run("UpdateUserByUID", func(t *testing.T) {
		t.Run("사용자를 업데이트한다", func(t *testing.T) {
			tearDown := setUp(t)
			defer tearDown(t)

			service := NewUserService(db)

			user := &models.User{
				Email:                "test@example.com",
				Nickname:             "nickname",
				Fullname:             "fullname",
				Password:             "password",
				FirebaseProviderType: "kakao",
				FirebaseUID:          "uid",
			}

			created, _ := service.CreateUser(user)

			updatedNickname := "updated"

			_, err := service.UpdateUserByUID(created.FirebaseUID, updatedNickname)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			found, err := service.FindUserByUID(created.FirebaseUID)
			if err != nil {
				t.Errorf("got %v want %v", err, nil)
			}

			if found.Nickname != updatedNickname {
				t.Errorf("got %v want %v", found.Nickname, updatedNickname)
			}
		})
	})
}
