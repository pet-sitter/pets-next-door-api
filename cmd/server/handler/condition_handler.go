package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type ConditionHandler struct {
	conditionService service.ConditionService
}

func NewConditionHandler(conditionService service.ConditionService) *ConditionHandler {
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
func (h *ConditionHandler) FindConditions(c echo.Context) error {
	res, err := h.conditionService.FindConditions(c.Request().Context())
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusOK, res)
}
