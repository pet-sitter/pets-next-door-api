package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	kakaoinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/kakao"
)

type authHandler struct{}

func newAuthHandler() *authHandler {
	return &authHandler{}
}

// Kakao 로그인 페이지로 redirect한다.
func (h *authHandler) kakaoLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://kauth.kakao.com/oauth/authorize?"+
		"client_id="+configs.KakaoRestAPIKey+
		"&redirect_uri="+configs.KakaoRedirectURI+
		"&response_type=code"+
		"&scope=profile_nickname,profile_image,account_email,gender,age_range",
		http.StatusTemporaryRedirect,
	)
}

// Kakao 로그인 콜백을 처리하고, 사용자 기본 정보를 채워 사용자를 생성하고, Firebase Custom Token을 발급한다.
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

	fbUser, err := authClient.GetUserByEmail(ctx, userProfile.KakaoAccount.Email)
	if err != nil {
		newFBUserParams := (&auth.UserToCreate{}).
			UID(fmt.Sprintf("%d", userProfile.ID)).
			Email(userProfile.KakaoAccount.Email).
			EmailVerified(true).
			PhotoURL(userProfile.Properties.ProfileImage).
			Disabled(false)

		newFBUser, err := authClient.CreateUser(ctx, newFBUserParams)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		claims := map[string]interface{}{"provider": "kakao"}
		authClient.SetCustomUserClaims(ctx, newFBUser.UID, claims)
	}
	if provider, ok := fbUser.CustomClaims["provider"]; ok {
		if provider != "kakao" {
			http.Error(w, fmt.Sprintf("user already registered with another provider: %s", fbUser.ProviderID), http.StatusConflict)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"authToken": customToken,
		"provider":  "kakao.com",
		"uid":       fmt.Sprintf("%d", userProfile.ID),
		"email":     userProfile.KakaoAccount.Email,
	})
}
