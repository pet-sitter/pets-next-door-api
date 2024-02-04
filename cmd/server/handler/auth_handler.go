package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/auth"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	kakaoinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/kakao"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type authHandler struct {
	authService service.AuthService
	kakaoClient kakaoinfra.KakaoClient
}

func NewAuthHandler(authService service.AuthService, kakaoClient kakaoinfra.KakaoClient) *authHandler {
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
// @Success 200 {object} auth.KakaoCallbackView
// @Router /auth/callback/kakao [get]
func (h *authHandler) KakaoCallback(w http.ResponseWriter, r *http.Request) {
	code := pnd.ParseOptionalStringQuery(r, "code")
	tokenView, err := h.kakaoClient.FetchAccessToken(*code)
	if err != nil {
		render.Render(w, r, pnd.ErrUnknown(err))
		return
	}

	userProfile, err := h.kakaoClient.FetchUserProfile(tokenView.AccessToken)
	if err != nil {
		render.Render(w, r, pnd.ErrUnknown(err))
		return
	}

	customToken, err2 := h.authService.CustomToken(r.Context(), fmt.Sprintf("%d", userProfile.ID))
	if err != nil {
		render.Render(w, r, err2)
		return
	}

	render.JSON(w, r, auth.KakaoCallbackView{
		AuthToken:            *customToken,
		FirebaseProviderType: user.FirebaseProviderTypeKakao,
		FirebaseUID:          fmt.Sprintf("%d", userProfile.ID),
		Email:                userProfile.KakaoAccount.Email,
		PhotoURL:             userProfile.Properties.ProfileImage,
	})
}

// GenerateFBCustomTokenFromKakao godoc
// @Summary Kakao OAuth 토큰 기반 Firebase Custom Token 생성 API
// @Description 주어진 카카오 토큰으로 사용자 기본 정보를 검증하고 Firebase Custom Token을 발급합니다.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.GenerateFBCustomTokenRequest true "Firebase Custom Token 생성 요청"
// @Success 201 {object} auth.GenerateFBCustomTokenResponse
// @Failure 400 {object} pnd.AppError
// @Router /auth/custom-tokens/kakao [post]
func (h *authHandler) GenerateFBCustomTokenFromKakao(w http.ResponseWriter, r *http.Request) {
	var tokenRequest auth.GenerateFBCustomTokenRequest
	if err := pnd.ParseBody(r, &tokenRequest); err != nil {
		render.Render(w, r, err)
		return
	}

	userProfile, err2 := h.kakaoClient.FetchUserProfile(tokenRequest.OAuthToken)
	if err2 != nil {
		render.Render(w, nil, pnd.ErrBadRequest(fmt.Errorf("유효하지 않은 Kakao 인증 정보입니다")))
		return
	}

	customToken, err := h.authService.CustomToken(r.Context(), fmt.Sprintf("%d", userProfile.ID))
	if err != nil {
		render.Render(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, auth.GenerateFBCustomTokenResponse{
		AuthToken:            *customToken,
		FirebaseProviderType: user.FirebaseProviderTypeKakao,
		FirebaseUID:          fmt.Sprintf("%d", userProfile.ID),
		Email:                userProfile.KakaoAccount.Email,
		PhotoURL:             userProfile.Properties.ProfileImage,
	})
}
