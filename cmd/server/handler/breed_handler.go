package handler

import (
	"github.com/labstack/echo/v4"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
	"net/http"
)

type BreedHandler struct {
	breedService service.BreedService
}

func NewBreedHandler(breedService service.BreedService) *BreedHandler {
	return &BreedHandler{breedService: breedService}
}

// FindBreeds godoc
// @Summary 견/묘종을 조회합니다.
// @Description
// @Tags pets
// @Accept  json
// @Produce  json
// @Param page query int false "페이지 번호" default(1)
// @Param size query int false "페이지 사이즈" default(20)
// @Param pet_type query string false "펫 종류" Enums(dog, cat)
// @Success 200 {object} pet.BreedListView
// @Router /breeds [get]
func (h *BreedHandler) FindBreeds(c echo.Context) error {
	petType := pnd.ParseOptionalStringQuery(c, "pet_type")
	page, size, err := pnd.ParsePaginationQueries(c, 1, 20)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	res, err := h.breedService.FindBreeds(c.Request().Context(), page, size, petType)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusOK, res)
}
