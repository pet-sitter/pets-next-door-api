package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/pet-sitter/pets-next-door-api/internal/models"
	"github.com/pet-sitter/pets-next-door-api/internal/user"
)

type UserHandler struct {
	userService user.UserServicer
}

func newUserHandler(userService user.UserServicer) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type RegisterUserRequest struct {
	Email                string `json:"email"`
	Nickname             string `json:"nickname"`
	Fullname             string `json:"fullname"`
	FirebaseProviderType string `json:"fbProviderType"`
	FirebaseUID          string `json:"fbUid"`
}

type UserResponse struct {
	ID                   int    `json:"id"`
	Email                string `json:"email"`
	Nickname             string `json:"nickname"`
	Fullname             string `json:"fullname"`
	FirebaseProviderType string `json:"fbProviderType"`
	FirebaseUID          string `json:"fbUid"`
}

// RegisterUser godoc
// @Summary 파이어베이스 가입 이후 정보를 입력 받아 유저를 생성합니다.
// @Description
// @Tags users
// @Accept  json
// @Produce  json
// @Param request body RegisterUserRequest true "사용자 회원가입 요청"
// @Success 201 {object} UserResponse
// @Router /users/register [post]
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var registerUserRequest RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&registerUserRequest); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if err := validator.New().Struct(registerUserRequest); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	userModel, err := h.userService.CreateUser(&models.User{
		Email:                registerUserRequest.Email,
		Nickname:             registerUserRequest.Nickname,
		Fullname:             registerUserRequest.Fullname,
		FirebaseProviderType: registerUserRequest.FirebaseProviderType,
		FirebaseUID:          registerUserRequest.FirebaseUID,
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := UserResponse{
		ID:                   userModel.ID,
		Email:                userModel.Email,
		Nickname:             userModel.Nickname,
		Fullname:             userModel.Fullname,
		FirebaseProviderType: userModel.FirebaseProviderType,
		FirebaseUID:          userModel.FirebaseUID,
	}
	json.NewEncoder(w).Encode(response)
}

// FindMyProfile godoc
// @Summary 내 프로필 정보를 조회합니다.
// @Description
// @Tags users
// @Produce  json
// @Security firebase
// @Success 200 {object} UserResponse
// @Router /users/me [get]
func (h *UserHandler) FindMyProfile(w http.ResponseWriter, r *http.Request) {
	idToken, err := verifyAuth(r.Context(), r.Header.Get("Authorization"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	uid := idToken.UID

	userModel, err := h.userService.FindUserByUID(uid)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	response := UserResponse{
		ID:                   userModel.ID,
		Email:                userModel.Email,
		Nickname:             userModel.Nickname,
		Fullname:             userModel.Fullname,
		FirebaseProviderType: userModel.FirebaseProviderType,
		FirebaseUID:          userModel.FirebaseUID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
}

// UpdateMyProfile godoc
// @Summary 내 프로필 정보를 수정합니다.
// @Description
// @Tags users
// @Accept  json
// @Produce  json
// @Security firebase
// @Param request body UpdateUserRequest true "프로필 정보 수정 요청"
// @Success 200 {object} UserResponse
// @Router /users/me [put]
func (h *UserHandler) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	idToken, err := verifyAuth(r.Context(), r.Header.Get("Authorization"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	uid := idToken.UID

	var updateUserRequest UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&updateUserRequest); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if err := validator.New().Struct(updateUserRequest); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	userModel, err := h.userService.UpdateUserByUID(uid, updateUserRequest.Nickname)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	response := UserResponse{
		ID:                   userModel.ID,
		Email:                userModel.Email,
		Nickname:             userModel.Nickname,
		Fullname:             userModel.Fullname,
		FirebaseProviderType: userModel.FirebaseProviderType,
		FirebaseUID:          userModel.FirebaseUID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
