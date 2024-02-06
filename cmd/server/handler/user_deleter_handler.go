package handler

import (
	"net/http"

	"github.com/go-chi/render"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type UserDeleterHandler struct {
	userDeleter service.UserDeleter
	authService service.AuthService
}

func NewUserDeleterHandler(userDeleter service.UserDeleter, authService service.AuthService) *UserDeleterHandler {
	return &UserDeleterHandler{
		userDeleter: userDeleter,
		authService: authService,
	}
}

// HardDeleteMyAccount godoc
// @Summary 내 계정을 PND 및 파이어베이스에서 완전히 삭제합니다.
// @Description
// @Tags users
// @Produce  json
// @Security FirebaseAuth
// @Success 200
// @Router /users/me/hard [delete]
func (h *UserDeleterHandler) HardDeleteMyAccount(w http.ResponseWriter, r *http.Request) {
	foundUser, err := h.authService.VerifyAuthAndGetUser(r.Context(), r)
	if err != nil {
		render.Render(w, r, err)
		return
	}

	uid := foundUser.FirebaseUID

	ctx := r.Context()
	if err := h.userDeleter.HardDeleteUserByUID(ctx, r, uid); err != nil {
		render.Render(w, r, err)
		return
	}

	pnd.OK(w, nil, nil)
}
