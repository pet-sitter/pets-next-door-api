package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/configs"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/auth"
	kakaoinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/kakao"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
	kakaoClient kakaoinfra.KakaoClient
}

func NewAuthHandler(
	authService service.AuthService,
	kakaoClient kakaoinfra.KakaoClient,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		kakaoClient: kakaoClient,
	}
}

// KakaoLogin godoc
// @Summary Kakao 로그인 페이지로 redirect 합니다.
// @Description
// @Tags auth
// @Success 302
// @Router /auth/login/kakao [get]
func (h *AuthHandler) KakaoLogin(c echo.Context) error {
	return c.Redirect(
		http.StatusTemporaryRedirect,
		"https://kauth.kakao.com/oauth/authorize?"+
			"client_id="+configs.KakaoRestAPIKey+
			"&redirect_uri="+configs.KakaoRedirectURI+
			"&response_type=code"+
			"&scope=profile_nickname,profile_image,account_email,gender,age_range",
	)
}

// KakaoCallback godoc
// @Summary Kakao 회원가입 콜백 API
// @Description Kakao 로그인 콜백을 처리하고, 사용자 기본 정보와 함께 Firebase Custom Token을 발급합니다.
// @Tags auth
// @Success 200 {object} auth.KakaoCallbackView
// @Router /auth/callback/kakao [get]
func (h *AuthHandler) KakaoCallback(c echo.Context) error {
	code := pnd.ParseOptionalStringQuery(c, "code")
	tokenView, err := h.kakaoClient.FetchAccessToken(c.Request().Context(), *code)
	if err != nil {
		pndErr := pnd.ErrUnknown(err)
		return c.JSON(pndErr.StatusCode, pndErr)
	}

	userProfile, err := h.kakaoClient.FetchUserProfile(c.Request().Context(), tokenView.AccessToken)
	if err != nil {
		pndErr := pnd.ErrUnknown(err)
		return c.JSON(pndErr.StatusCode, pndErr)
	}

	customToken, err := h.authService.CustomToken(
		c.Request().Context(),
		strconv.FormatInt(userProfile.ID, 10),
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, auth.NewKakaoCallbackView(*customToken, userProfile))
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
func (h *AuthHandler) GenerateFBCustomTokenFromKakao(c echo.Context) error {
	var tokenRequest auth.GenerateFBCustomTokenRequest
	if err := pnd.ParseBody(c, &tokenRequest); err != nil {
		return err
	}

	userProfile, err := h.kakaoClient.FetchUserProfile(
		c.Request().Context(),
		tokenRequest.OAuthToken,
	)
	if err != nil {
		return pnd.ErrBadRequest(errors.New("유효하지 않은 Kakao 인증 정보입니다"))
	}

	customToken, err := h.authService.CustomToken(
		c.Request().Context(),
		strconv.FormatInt(userProfile.ID, 10),
	)
	if err != nil {
		return err
	}

	return c.JSON(
		http.StatusCreated,
		auth.NewGenerateFBCustomTokenResponse(*customToken, userProfile),
	)
}
