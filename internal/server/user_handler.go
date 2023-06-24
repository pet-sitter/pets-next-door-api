package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
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

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var registerUserRequest struct {
		Email     string   `json:"email"`
		UID       string   `json:"uid"`
		Nickname  string   `json:"nickname"`
		Providers []string `json:"providers"`
	}
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

	userModel, err := h.userService.CreateUser(&user.UserModel{
		Email:     registerUserRequest.Email,
		UID:       registerUserRequest.UID,
		Nickname:  registerUserRequest.Nickname,
		Providers: strings.Join(registerUserRequest.Providers, ","),
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userModel)
}

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userModel)
}

func (h *UserHandler) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	idToken, err := verifyAuth(r.Context(), r.Header.Get("Authorization"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	uid := idToken.UID

	var updateUserRequest struct {
		Nickname string `json:"nickname"`
	}
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userModel)
}
