package service

import (
	"context"
	"net/http"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"
)

type UserDeleter struct {
	conn        *database.DB
	authService AuthService
}

func NewUserDeleter(conn *database.DB, authService AuthService) *UserDeleter {
	return &UserDeleter{
		conn:        conn,
		authService: authService,
	}
}

func (s *UserDeleter) HardDeleteUserByUID(ctx context.Context, r *http.Request, uid string) *pnd.AppError {
	err := database.WithTransaction(ctx, s.conn, func(tx *database.Tx) *pnd.AppError {
		userStore := postgres.NewUserPostgresStore(tx)

		// PND 계정 삭제
		if err := userStore.HardDeleteUserByUID(ctx, uid); err != nil {
			// TODO: 유저 데이터 삭제
			return err
		}
		// Firebase 계정 삭제
		if err := s.authService.DeleteMyAccount(ctx, r); err != nil {
			return err
		}

		return nil
	})

	return err
}
