package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type UserHandler struct {
	userService service.UserService
	authService service.AuthService
}

func NewUserHandler(userService service.UserService, authService service.AuthService) *UserHandler {
	return &UserHandler{
		userService: userService,
		authService: authService,
	}
}

// RegisterUser godoc
// @Summary 파이어베이스 가입 이후 정보를 입력 받아 유저를 생성합니다.
// @Description
// @Tags users
// @Accept  json
// @Produce  json
// @Param request body user.RegisterUserRequest true "사용자 회원가입 요청"
// @Success 201 {object} user.InternalView
// @Router /users [post]
func (h *UserHandler) RegisterUser(c echo.Context) error {
	var registerUserRequest user.RegisterUserRequest
	if err := pnd.ParseBody(c, &registerUserRequest); err != nil {
		return c.JSON(err.StatusCode, err)
	}

	res, err := h.userService.RegisterUser(c.Request().Context(), &registerUserRequest)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusCreated, res)
}

// CheckUserNickname godoc
// @Summary 닉네임 중복 여부를 조회합니다.
// @Description
// @Tags users
// @Accept  json
// @Produce  json
// @Param request body user.CheckNicknameRequest true "사용자 닉네임 중복 조회 요청"
// @Success 200 {object} user.CheckNicknameView
// @Router /users/check/nickname [post]
func (h *UserHandler) CheckUserNickname(c echo.Context) error {
	var checkUserNicknameRequest user.CheckNicknameRequest
	if err := pnd.ParseBody(c, &checkUserNicknameRequest); err != nil {
		return c.JSON(err.StatusCode, err)
	}

	exists, err := h.userService.ExistsByNickname(c.Request().Context(), checkUserNicknameRequest.Nickname)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusOK, user.CheckNicknameView{IsAvailable: !exists})
}

// FindUserStatusByEmail godoc
// @Summary 이메일로 유저의 가입 상태를 조회합니다.
// @Description
// @Tags users
// @Accept  json
// @Produce  json
// @Param request body user.UserStatusRequest true "사용자 가입 상태 조회 요청"
// @Success 200 {object} user.StatusView
// @Router /users/status [post]
func (h *UserHandler) FindUserStatusByEmail(c echo.Context) error {
	var providerRequest user.UserStatusRequest
	if err := pnd.ParseBody(c, &providerRequest); err != nil {
		return c.JSON(err.StatusCode, err)
	}

	userData, err := h.userService.FindUser(c.Request().Context(), user.FindUserParams{Email: &providerRequest.Email})
	if err != nil {
		return c.JSON(http.StatusOK, user.StatusView{
			Status: user.StatusNotRegistered,
		})
	}

	return c.JSON(http.StatusOK, user.NewStatusView(userData.FirebaseProviderType))
}

// FindUsers godoc
// @Summary 사용자 목록을 조회합니다.
// @Description
// @Tags users
// @Produce  json
// @Security FirebaseAuth
// @Param page query int false "페이지 번호" default(1)
// @Param size query int false "페이지 사이즈" default(10)
// @Param nickname query string false "닉네임 (완전 일치)"
// @Success 200 {object} user.ListWithoutPrivateInfo
// @Router /users [get]
func (h *UserHandler) FindUsers(c echo.Context) error {
	_, err := h.authService.VerifyAuthAndGetUser(c.Request().Context(), c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	nickname := pnd.ParseOptionalStringQuery(c, "nickname")
	page, size, err := pnd.ParsePaginationQueries(c, 1, 10)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	var res *user.ListWithoutPrivateInfo

	res, err = h.userService.FindUsers(c.Request().Context(),
		user.FindUsersParams{Page: page, Size: size, Nickname: nickname})
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusOK, res)
}

// FindUserByID godoc
// @Summary 단일 사용자 정보를 조회합니다.
// @Description
// @Tags users
// @Produce  json
// @Security FirebaseAuth
// @Param userID path int true "사용자 ID"
// @Success 200 {object} user.ProfileView
// @Router /users/{userID} [get]
func (h *UserHandler) FindUserByID(c echo.Context) error {
	_, err := h.authService.VerifyAuthAndGetUser(c.Request().Context(), c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	userID, err := pnd.ParseIDFromPath(c, "userID")
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	res, err := h.userService.FindUserProfile(c.Request().Context(), user.FindUserParams{ID: userID})
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusOK, res)
}

// FindMyProfile godoc
// @Summary 내 프로필 정보를 조회합니다.
// @Description
// @Tags users
// @Produce  json
// @Security FirebaseAuth
// @Success 200 {object} user.MyProfileView
// @Router /users/me [get]
func (h *UserHandler) FindMyProfile(c echo.Context) error {
	res, err := h.authService.VerifyAuthAndGetUser(c.Request().Context(), c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusOK, res.ToMyProfileView())
}

// UpdateMyProfile godoc
// @Summary 내 프로필 정보를 수정합니다.
// @Description
// @Tags users
// @Accept  json
// @Produce  json
// @Security FirebaseAuth
// @Param request body user.UpdateUserRequest true "사용자 프로필 수정 요청"
// @Success 200 {object} user.MyProfileView
// @Router /users/me [put]
func (h *UserHandler) UpdateMyProfile(c echo.Context) error {
	foundUser, err := h.authService.VerifyAuthAndGetUser(c.Request().Context(), c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	uid := foundUser.FirebaseUID

	var updateUserRequest user.UpdateUserRequest
	if err = pnd.ParseBody(c, &updateUserRequest); err != nil {
		return c.JSON(err.StatusCode, err)
	}

	view, err := h.userService.UpdateUserByUID(
		c.Request().Context(),
		uid,
		updateUserRequest.Nickname,
		updateUserRequest.ProfileImageID,
	)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusOK, view)
}

// DeleteMyAccount godoc
// @Summary 내 계정을 삭제합니다.
// @Description
// @Tags users
// @Security FirebaseAuth
// @Success 204
// @Router /users/me [delete]
func (h *UserHandler) DeleteMyAccount(c echo.Context) error {
	loggedInUser, err := h.authService.VerifyAuthAndGetUser(c.Request().Context(), c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	if err := h.userService.DeleteUserByUID(c.Request().Context(), loggedInUser.FirebaseUID); err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// AddMyPets godoc
// @Summary 내 반려동물을 등록합니다.
// @Description
// @Tags users,pets
// @Accept json
// @Produce json
// @Security FirebaseAuth
// @Param request body pet.AddPetsToOwnerRequest true "반려동물 등록 요청"
// @Success 200
// @Router /users/me/pets [put]
func (h *UserHandler) AddMyPets(c echo.Context) error {
	foundUser, err := h.authService.VerifyAuthAndGetUser(c.Request().Context(), c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	uid := foundUser.FirebaseUID

	var addPetsToOwnerRequest pet.AddPetsToOwnerRequest
	if err := pnd.ParseBody(c, &addPetsToOwnerRequest); err != nil {
		return c.JSON(err.StatusCode, err)
	}

	if _, err := h.userService.AddPetsToOwner(c.Request().Context(), uid, addPetsToOwnerRequest); err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.NoContent(http.StatusOK)
}

// FindMyPets godoc
// @Summary 내 반려동물 목록을 조회합니다.
// @Description
// @Tags users,pets
// @Produce json
// @Security FirebaseAuth
// @Success 200 {object} pet.ListView
// @Router /users/me/pets [get]
func (h *UserHandler) FindMyPets(c echo.Context) error {
	foundUser, err := h.authService.VerifyAuthAndGetUser(c.Request().Context(), c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	res, err := h.userService.FindPets(c.Request().Context(), pet.FindPetsParams{
		Page:    1,
		Size:    100,
		OwnerID: &foundUser.ID,
	})
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusOK, res)
}

// UpdateMyPet godoc
// @Summary 내 반려동물 정보를 수정합니다.
// @Description
// @Tags users,pets
// @Accept json
// @Produce json
// @Security FirebaseAuth
// @Param petID path int true "반려동물 ID"
// @Param request body pet.UpdatePetRequest true "반려동물 수정 요청"
// @Success 200 {object} pet.DetailView
// @Router /users/me/pets/{petID} [put]
func (h *UserHandler) UpdateMyPet(c echo.Context) error {
	foundUser, err := h.authService.VerifyAuthAndGetUser(c.Request().Context(), c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}
	uid := foundUser.FirebaseUID

	petID, err := pnd.ParseIDFromPath(c, "petID")
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	var updatePetRequest pet.UpdatePetRequest
	if err = pnd.ParseBody(c, &updatePetRequest); err != nil {
		return c.JSON(err.StatusCode, err)
	}

	res, err := h.userService.UpdatePet(c.Request().Context(), uid, int64(*petID), updatePetRequest)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusOK, res)
}

// DeleteMyPet godoc
// @Summary 내 반려동물을 삭제합니다.
// @Description
// @Tags users,pets
// @Security FirebaseAuth
// @Param petID path int true "반려동물 ID"
// @Success 204
// @Router /users/me/pets/{petID} [delete]
func (h *UserHandler) DeleteMyPet(c echo.Context) error {
	foundUser, err := h.authService.VerifyAuthAndGetUser(c.Request().Context(), c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}
	uid := foundUser.FirebaseUID

	petID, err := pnd.ParseIDFromPath(c, "petID")
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	if err := h.userService.DeletePet(c.Request().Context(), uid, int64(*petID)); err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.NoContent(http.StatusNoContent)
}
