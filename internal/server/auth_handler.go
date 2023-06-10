package server

import (
	"net/http"

	"github.com/pet-sitter/pets-next-door-api/internal/configs"
)

type authHandler struct{}

func newAuthHandler() *authHandler {
	return &authHandler{}
}

// Kakao 로그인 페이지로 redirect한다.
func (h *authHandler) kakaoLogin(w http.ResponseWriter, r *http.Request) {
	redirectUri := "http://localhost:8080/api/auth/callback/kakao"

	http.Redirect(w, r, "https://kauth.kakao.com/oauth/authorize?"+
		"client_id="+configs.KakaoRestAPIKey+
		"&redirect_uri="+redirectUri+
		"&response_type=code"+
		"&scope=profile_nickname,profile_image,account_email,gender,age_range",
		http.StatusTemporaryRedirect,
	)
}
