package handler

import (
	"net/http"

	"github.com/go-chi/render"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sos_post"
)

type ConditionHandler struct {
	conditionService sos_post.ConditionService
}

func NewConditionHandler(conditionService sos_post.ConditionService) *ConditionHandler {
	return &ConditionHandler{conditionService: conditionService}
}

// FindConditions godoc
// @Summary 돌봄 조건을 조회합니다.
// @Description
// @Tags posts
// @Accept  json
// @Produce  json
// @Success 200 {object} []sos_post.ConditionView
// @Router /posts/sos/conditions [get]
func (h *ConditionHandler) FindConditions(w http.ResponseWriter, r *http.Request) {
	res, err := h.conditionService.FindConditions(r.Context())
	if err != nil {
		render.Render(w, r, err)
		return
	}

	pnd.OK(w, nil, res)
}
