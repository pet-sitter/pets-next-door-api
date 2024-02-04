package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	s3infra "github.com/pet-sitter/pets-next-door-api/internal/infra/s3"

	"firebase.google.com/go/auth"
)

type AuthService interface {
	VerifyAuthAndGetUser(ctx context.Context, r *http.Request) (*user.FindUserView, *pnd.AppError)
	CustomToken(ctx context.Context, uid string) (*string, *pnd.AppError)
}

type FirebaseBearerAuthService struct {
	conn       *database.DB
	authClient *auth.Client
	s3Client   *s3infra.S3Client
}

func NewFirebaseBearerAuthService(conn *database.DB, authClient *auth.Client, s3Client *s3infra.S3Client) *FirebaseBearerAuthService {
	return &FirebaseBearerAuthService{
		conn:       conn,
		authClient: authClient,
		s3Client:   s3Client,
	}
}

func (s *FirebaseBearerAuthService) verifyAuth(ctx context.Context, authHeader string) (*auth.Token, error) {
	idToken, err := s.stripBearerToken(authHeader)
	if err != nil {
		return nil, err
	}

	authToken, err := s.authClient.VerifyIDToken(ctx, idToken)
	return authToken, err
}

func (s *FirebaseBearerAuthService) VerifyAuthAndGetUser(ctx context.Context, r *http.Request) (*user.FindUserView, *pnd.AppError) {
	authToken, err := s.verifyAuth(ctx, r.Header.Get("Authorization"))
	if err != nil {
		return nil, pnd.ErrInvalidFBToken(fmt.Errorf("유효하지 않은 인증 토큰입니다"))
	}

	var foundUser *user.FindUserView
	var err2 *pnd.AppError

	err2 = database.WithTransaction(ctx, s.conn, func(tx *database.Tx) *pnd.AppError {
		userService := NewUserService(s.conn, s.s3Client)

		foundUser, err2 = userService.FindUserByUID(ctx, authToken.UID)
		if err2 != nil {
			return pnd.ErrUserNotRegistered(fmt.Errorf("가입되지 않은 사용자입니다"))
		}

		return nil
	})
	if err2 != nil {
		return nil, err2
	}

	return foundUser, nil
}

func (s *FirebaseBearerAuthService) CustomToken(ctx context.Context, uid string) (*string, *pnd.AppError) {
	customToken, err := s.authClient.CustomToken(ctx, uid)
	if err != nil {
		return nil, pnd.ErrUnknown(err)
	}

	return &customToken, nil
}

func (s *FirebaseBearerAuthService) stripBearerToken(authHeader string) (string, error) {
	if len(authHeader) > 6 && strings.ToUpper(authHeader[0:7]) == "BEARER " {
		return authHeader[7:], nil
	}

	return authHeader, fmt.Errorf("유효하지 않은 인증 토큰입니다")
}
