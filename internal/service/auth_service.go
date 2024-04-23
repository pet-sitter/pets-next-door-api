package service

import (
	"context"
	"errors"
	"strings"

	"firebase.google.com/go/auth"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
)

type AuthService interface {
	VerifyAuthAndGetUser(ctx context.Context, authHeader string) (*user.InternalView, *pnd.AppError)
	CustomToken(ctx context.Context, uid string) (*string, *pnd.AppError)
}

type FirebaseBearerAuthService struct {
	authClient  *auth.Client
	userService *UserService
}

func NewFirebaseBearerAuthService(authClient *auth.Client, userService *UserService) *FirebaseBearerAuthService {
	return &FirebaseBearerAuthService{
		authClient:  authClient,
		userService: userService,
	}
}

func (service *FirebaseBearerAuthService) verifyAuth(
	ctx context.Context, authHeader string,
) (*auth.Token, *pnd.AppError) {
	idToken, err := service.stripBearerToken(authHeader)
	if err != nil {
		return nil, err
	}

	authToken, err2 := service.authClient.VerifyIDToken(ctx, idToken)
	if err2 != nil {
		return nil, pnd.ErrInvalidFBToken(err2)
	}

	return authToken, nil
}

func (service *FirebaseBearerAuthService) VerifyAuthAndGetUser(
	ctx context.Context, authHeader string,
) (*user.InternalView, *pnd.AppError) {
	authToken, err := service.verifyAuth(ctx, authHeader)
	if err != nil {
		return nil, err
	}

	foundUser, err := service.userService.FindUser(ctx, user.FindUserParams{FbUID: &authToken.UID})
	if err != nil {
		return nil, pnd.ErrUserNotRegistered(errors.New("가입되지 않은 사용자입니다"))
	}

	return foundUser.ToInternalView(), nil
}

func (service *FirebaseBearerAuthService) CustomToken(ctx context.Context, uid string) (*string, *pnd.AppError) {
	customToken, err := service.authClient.CustomToken(ctx, uid)
	if err != nil {
		return nil, pnd.ErrUnknown(err)
	}

	return &customToken, nil
}

func (service *FirebaseBearerAuthService) stripBearerToken(authHeader string) (string, *pnd.AppError) {
	if len(authHeader) > 6 && strings.ToUpper(authHeader[0:7]) == "BEARER " {
		return authHeader[7:], nil
	}

	return "", pnd.ErrInvalidBearerToken(errors.New("올바른 Bearer 토큰이 아닙니다"))
}
