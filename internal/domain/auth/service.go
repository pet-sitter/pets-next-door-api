package auth

import (
	"context"
	"fmt"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"strings"

	"firebase.google.com/go/auth"
)

type AuthService interface {
	VerifyAuthAndGetUser(ctx context.Context, authHeader string) (*user.FindUserResponse, error)
	CustomToken(ctx context.Context, uid string) (string, error)
}

type FirebaseBearerAuthService struct {
	authClient  *auth.Client
	userService user.UserServicer
}

func NewFirebaseBearerAuthService(authClient *auth.Client, userService user.UserServicer) *FirebaseBearerAuthService {
	return &FirebaseBearerAuthService{
		authClient:  authClient,
		userService: userService,
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

func (s *FirebaseBearerAuthService) VerifyAuthAndGetUser(ctx context.Context, authHeader string) (*user.FindUserResponse, error) {
	authToken, err := s.verifyAuth(ctx, authHeader)
	if err != nil {
		return nil, err
	}

	foundUser, err := s.userService.FindUserByUID(authToken.UID)
	if err != nil {
		return nil, err
	}

	return foundUser, nil
}

func (s *FirebaseBearerAuthService) CustomToken(ctx context.Context, uid string) (string, error) {
	return s.authClient.CustomToken(ctx, uid)
}

func (s *FirebaseBearerAuthService) stripBearerToken(authHeader string) (string, error) {
	if len(authHeader) > 6 && strings.ToUpper(authHeader[0:7]) == "BEARER " {
		return authHeader[7:], nil
	}

	return authHeader, fmt.Errorf("invalid auth header")
}
