package server

import (
	"encoding/json"
	"log"
	"net/http"

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

// Kakao 로그인 콜백을 처리한다.
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

	log.Println(userProfile)

	// TODO: DB에 저장 및 Firebase custom token 발급
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
