package auth

import (
	"context"
	"fmt"
	"strings"

	"firebase.google.com/go/auth"
)

type AuthService interface {
	VerifyAuth(ctx context.Context, authHeader string) (*auth.Token, error)
	CustomToken(ctx context.Context, uid string) (string, error)
}

type FirebaseBearerAuthService struct {
	authClient *auth.Client
}

func NewFirebaseBearerAuthService(authClient *auth.Client) *FirebaseBearerAuthService {
	return &FirebaseBearerAuthService{
		authClient: authClient,
	}
}

func (s *FirebaseBearerAuthService) VerifyAuth(ctx context.Context, authHeader string) (*auth.Token, error) {
	idToken, err := s.stripBearerToken(authHeader)
	if err != nil {
		return nil, err
	}

	authToken, err := s.authClient.VerifyIDToken(ctx, idToken)
	return authToken, err
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
