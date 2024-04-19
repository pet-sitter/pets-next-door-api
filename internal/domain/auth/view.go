package auth

import (
	"fmt"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	kakaoinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/kakao"
)

type KakaoCallbackView struct {
	AuthToken            string                    `json:"authToken"`
	FirebaseProviderType user.FirebaseProviderType `json:"fbProviderType"`
	FirebaseUID          string                    `json:"fbUid"`
	Email                string                    `json:"email"`
	PhotoURL             string                    `json:"photoURL"`
}

func NewKakaoCallbackView(authToken string, kakaoUserProfile *kakaoinfra.KakaoUserProfile) KakaoCallbackView {
	return KakaoCallbackView{
		AuthToken:            authToken,
		FirebaseProviderType: user.FirebaseProviderTypeKakao,
		FirebaseUID:          fmt.Sprintf("%d", kakaoUserProfile.ID),
		Email:                kakaoUserProfile.KakaoAccount.Email,
		PhotoURL:             kakaoUserProfile.Properties.ProfileImage,
	}
}

// GenerateFBCustomTokenRequest 는 OAuth 토큰 정보를 기반으로 Firebase Custom Token을 생성하기 위한 요청이다.
type GenerateFBCustomTokenRequest struct {
	OAuthToken string `json:"oauthToken"`
}

// GenerateFBCustomTokenResponse 는 Firebase Custom Token을 생성하기 위한 응답이다.
type GenerateFBCustomTokenResponse struct {
	AuthToken            string                    `json:"authToken"`
	FirebaseProviderType user.FirebaseProviderType `json:"fbProviderType"`
	FirebaseUID          string                    `json:"fbUid"`
	Email                string                    `json:"email"`
	PhotoURL             string                    `json:"photoURL"`
}

func NewGenerateFBCustomTokenResponse(authToken string, kakaoUserProfile *kakaoinfra.KakaoUserProfile) GenerateFBCustomTokenResponse {
	return GenerateFBCustomTokenResponse{
		AuthToken:            authToken,
		FirebaseProviderType: user.FirebaseProviderTypeKakao,
		FirebaseUID:          fmt.Sprintf("%d", kakaoUserProfile.ID),
		Email:                kakaoUserProfile.KakaoAccount.Email,
		PhotoURL:             kakaoUserProfile.Properties.ProfileImage,
	}
}
