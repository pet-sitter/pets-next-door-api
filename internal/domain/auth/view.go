package auth

import "github.com/pet-sitter/pets-next-door-api/internal/domain/user"

type KakaoCallbackView struct {
	AuthToken            string                    `json:"authToken"`
	FirebaseProviderType user.FirebaseProviderType `json:"fbProviderType"`
	FirebaseUID          string                    `json:"fbUid"`
	Email                string                    `json:"email"`
	PhotoURL             string                    `json:"photoURL"`
}

// GenerateFBCustomTokenRequest는 OAuth 토큰 정보를 기반으로 Firebase Custom Token을 생성하기 위한 요청이다.
type GenerateFBCustomTokenRequest struct {
	OAuthToken string `json:"oauthToken"`
}

// GenerateFBCustomTokenResponse는 Firebase Custom Token을 생성하기 위한 응답이다.
type GenerateFBCustomTokenResponse struct {
	AuthToken            string                    `json:"authToken"`
	FirebaseProviderType user.FirebaseProviderType `json:"fbProviderType"`
	FirebaseUID          string                    `json:"fbUid"`
	Email                string                    `json:"email"`
	PhotoURL             string                    `json:"photoURL"`
}
