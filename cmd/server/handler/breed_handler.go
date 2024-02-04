package handler

import (
	"net/http"

	"github.com/go-chi/render"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
)

type BreedHandler struct {
	breedService pet.BreedService
}

func NewBreedHandler(breedService pet.BreedService) *BreedHandler {
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
func (h *BreedHandler) FindBreeds(w http.ResponseWriter, r *http.Request) {
	petType := pnd.ParseOptionalStringQuery(r, "pet_type")
	page, size, err := pnd.ParsePaginationQueries(r, 1, 20)
	if err != nil {
		render.Render(w, r, err)
		return
	}

	res, err := h.breedService.FindBreeds(r.Context(), page, size, petType)
	if err != nil {
		render.Render(w, r, err)
		return
	}

	render.JSON(w, r, res)
}
