package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	app "github.com/pet-sitter/pets-next-door-api/api"
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/auth"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	kakaoinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/kakao"
)

type authHandler struct {
	authService auth.AuthService
	kakaoClient kakaoinfra.IKakaoClient
}

func NewAuthHandler(authService auth.AuthService, kakaoClient kakaoinfra.IKakaoClient) *authHandler {
	return &authHandler{
		authService: authService,
		kakaoClient: kakaoClient,
	}
}

// kakaoLogin godoc
// @Summary Kakao 로그인 페이지로 redirect 합니다.
// @Description
// @Tags auth
// @Success 302
// @Router /auth/login/kakao [get]
func (h *authHandler) KakaoLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://kauth.kakao.com/oauth/authorize?"+
		"client_id="+configs.KakaoRestAPIKey+
		"&redirect_uri="+configs.KakaoRedirectURI+
		"&response_type=code"+
		"&scope=profile_nickname,profile_image,account_email,gender,age_range",
		http.StatusTemporaryRedirect,
	)
}

// kakaoCallback godoc
// @Summary Kakao 회원가입 콜백 API
// @Description Kakao 로그인 콜백을 처리하고, 사용자 기본 정보와 함께 Firebase Custom Token을 발급합니다.
// @Tags auth
// @Success 200 {object} auth.KakaoCallbackResponse
// @Router /auth/callback/kakao [get]
func (h *authHandler) KakaoCallback(w http.ResponseWriter, r *http.Request) {
	code := utils.ParseOptionalStringQuery(r, "code")
	tokenView, err := h.kakaoClient.FetchAccessToken(*code)
	if err != nil {
		render.Render(w, r, app.ErrUnknown(err))
		return
	}

	userProfile, err := h.kakaoClient.FetchUserProfile(tokenView.AccessToken)
	if err != nil {
		render.Render(w, r, app.ErrUnknown(err))
		return
	}

	ctx := r.Context()
	customToken, err2 := h.authService.CustomToken(ctx, fmt.Sprintf("%d", userProfile.ID))
	if err != nil {
		render.Render(w, r, err2)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(auth.KakaoCallbackResponse{
		AuthToken:            *customToken,
		FirebaseProviderType: user.FirebaseProviderTypeKakao,
		FirebaseUID:          fmt.Sprintf("%d", userProfile.ID),
		Email:                userProfile.KakaoAccount.Email,
		PhotoURL:             userProfile.Properties.ProfileImage,
	})
}
