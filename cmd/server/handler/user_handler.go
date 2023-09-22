package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/pet-sitter/pets-next-door-api/api/commonviews"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/auth"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
)

type UserHandler struct {
	userService user.UserServicer
	authService auth.AuthService
}

func NewUserHandler(userService user.UserServicer, authService auth.AuthService) *UserHandler {
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
// @Success 201 {object} user.RegisterUserResponse
// @Router /users [post]
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var registerUserRequest user.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&registerUserRequest); err != nil {
		commonviews.BadRequest(w, nil, err.Error())
		return
	}
	if err := validator.New().Struct(registerUserRequest); err != nil {
		commonviews.BadRequest(w, nil, err.Error())
		return
	}

	res, err := h.userService.RegisterUser(&registerUserRequest)
	if err != nil {
		commonviews.InternalServerError(w, nil, err.Error())
		return
	}

	commonviews.Created(w, nil, res)
}

// FindUserStatusByEmail godoc
// @Summary 이메일로 유저의 가입 상태를 조회합니다.
// @Description
// @Tags users
// @Accept  json
// @Produce  json
// @Param request body user.UserStatusRequest true "사용자 가입 상태 조회 요청"
// @Success 200 {object} user.UserStatusView
// @Router /users/status [post]
func (h *UserHandler) FindUserStatusByEmail(w http.ResponseWriter, r *http.Request) {
	var providerRequest user.UserStatusRequest
	if err := commonviews.ParseBody(w, r, &providerRequest); err != nil {
		return
	}

	userStatus, err := h.userService.FindUserStatusByEmail(providerRequest.Email)
	if err != nil || userStatus == nil {
		commonviews.OK(w, nil, user.UserStatusView{
			Status: user.UserStatusNotRegistered,
		})
		return
	}

	commonviews.OK(w, nil, user.UserStatusView{
		Status:               user.UserStatusRegistered,
		FirebaseProviderType: userStatus.FirebaseProviderType,
	})
}

// FindMyProfile godoc
// @Summary 내 프로필 정보를 조회합니다.
// @Description
// @Tags users
// @Produce  json
// @Security FirebaseAuth
// @Success 200 {object} user.FindUserResponse
// @Router /users/me [get]
func (h *UserHandler) FindMyProfile(w http.ResponseWriter, r *http.Request) {
	idToken, err := h.authService.VerifyAuth(r.Context(), r.Header.Get("Authorization"))
	if err != nil {
		commonviews.Unauthorized(w, nil, "unauthorized")
		log.Printf("verifyAuth error: %v\n", err)
		return
	}

	uid := idToken.UID

	res, err := h.userService.FindUserByUID(uid)
	if err != nil {
		commonviews.Unauthorized(w, nil, "not registered")
		return
	}

	commonviews.OK(w, nil, res)
}

// UpdateMyProfile godoc
// @Summary 내 프로필 정보를 수정합니다.
// @Description
// @Tags users
// @Accept  json
// @Produce  json
// @Security FirebaseAuth
// @Param request body user.UpdateUserRequest true "사용자 프로필 수정 요청"
// @Success 200 {object} user.UpdateUserResponse
// @Router /users/me [put]
func (h *UserHandler) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	idToken, err := h.authService.VerifyAuth(r.Context(), r.Header.Get("Authorization"))
	if err != nil {
		commonviews.Unauthorized(w, nil, "unauthorized")
		return
	}

	uid := idToken.UID

	var updateUserRequest user.UpdateUserRequest
	if err := commonviews.ParseBody(w, r, &updateUserRequest); err != nil {
		return
	}

	userModel, err := h.userService.UpdateUserByUID(uid, updateUserRequest.Nickname, updateUserRequest.ProfileImageID)
	if err != nil {
		commonviews.InternalServerError(w, nil, err.Error())
		return
	}

	commonviews.OK(w, nil, user.UpdateUserResponse{
		ID:                   userModel.ID,
		Email:                userModel.Email,
		Nickname:             userModel.Nickname,
		Fullname:             userModel.Fullname,
		ProfileImageURL:      userModel.ProfileImageURL,
		FirebaseProviderType: userModel.FirebaseProviderType,
		FirebaseUID:          userModel.FirebaseUID,
	})
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
func (h *UserHandler) AddMyPets(w http.ResponseWriter, r *http.Request) {
	idToken, err := h.authService.VerifyAuth(r.Context(), r.Header.Get("Authorization"))
	if err != nil {
		commonviews.Unauthorized(w, nil, "unauthorized")
		return
	}

	uid := idToken.UID

	var addPetsToOwnerRequest pet.AddPetsToOwnerRequest
	if err := commonviews.ParseBody(w, r, &addPetsToOwnerRequest); err != nil {
		return
	}

	if _, err := h.userService.AddPetsToOwner(uid, addPetsToOwnerRequest); err != nil {
		commonviews.InternalServerError(w, nil, err.Error())
		return
	}

	commonviews.OK(w, nil, nil)
}

// FindMyPets godoc
// @Summary 내 반려동물 목록을 조회합니다.
// @Description
// @Tags users,pets
// @Produce json
// @Security FirebaseAuth
// @Success 200 {object} pet.FindMyPetsView
// @Router /users/me/pets [get]
func (h *UserHandler) FindMyPets(w http.ResponseWriter, r *http.Request) {
	idToken, err := h.authService.VerifyAuth(r.Context(), r.Header.Get("Authorization"))
	if err != nil {
		commonviews.Unauthorized(w, nil, "unauthorized")
		return
	}

	uid := idToken.UID

	res, err := h.userService.FindPetsByOwnerUID(uid)
	if err != nil {
		commonviews.InternalServerError(w, nil, err.Error())
		return
	}

	commonviews.OK(w, nil, res)
}