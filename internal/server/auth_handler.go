package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	kakaoinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/kakao"
	"github.com/pet-sitter/pets-next-door-api/internal/models"
)

type authHandler struct{}

func newAuthHandler() *authHandler {
	return &authHandler{}
}

// kakaoLogin godoc
// @Summary Kakao 로그인 페이지로 redirect 합니다.
// @Description
// @Tags auth
// @Success 302
// @Router /auth/login/kakao [get]
func (h *authHandler) kakaoLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://kauth.kakao.com/oauth/authorize?"+
		"client_id="+configs.KakaoRestAPIKey+
		"&redirect_uri="+configs.KakaoRedirectURI+
		"&response_type=code"+
		"&scope=profile_nickname,profile_image,account_email,gender,age_range",
		http.StatusTemporaryRedirect,
	)
}

type kakaoCallbackResponse struct {
	AuthToken            string                      `json:"authToken"`
	FirebaseProviderType models.FirebaseProviderType `json:"fbProviderType"`
	FirebaseUID          string                      `json:"fbUid"`
	Email                string                      `json:"email"`
	PhotoURL             string                      `json:"photoURL"`
}

// kakaoCallback godoc
// @Summary Kakao 회원가입 콜백 API
// @Description Kakao 로그인 콜백을 처리하고, 사용자 기본 정보와 함께 Firebase Custom Token을 발급합니다.
// @Tags auth
// @Success 200 {object} kakaoCallbackResponse
// @Router /auth/callback/kakao [get]
func (h *authHandler) kakaoCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	tokenView, err := kakaoinfra.FetchAccessToken(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userProfile, err := kakaoinfra.FetchUserProfile(tokenView.AccessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	authClient := ctx.Value(firebaseAuthClientKey).(*auth.Client)
	customToken, err := authClient.CustomToken(ctx, fmt.Sprintf("%d", userProfile.ID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fbUser, _ := authClient.GetUserByEmail(ctx, userProfile.KakaoAccount.Email)
	if provider, ok := fbUser.CustomClaims["provider"]; ok {
		if provider != models.FirebaseProviderTypeKakao {
			http.Error(w, fmt.Sprintf("user already registered with another provider: %s", fbUser.ProviderID), http.StatusConflict)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(kakaoCallbackResponse{
		AuthToken:            customToken,
		FirebaseProviderType: models.FirebaseProviderTypeKakao,
		FirebaseUID:          fmt.Sprintf("%d", userProfile.ID),
		Email:                userProfile.KakaoAccount.Email,
		PhotoURL:             userProfile.Properties.ProfileImage,
	})
}
