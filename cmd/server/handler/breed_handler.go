package handler

import (
	"github.com/pet-sitter/pets-next-door-api/api/commonviews"
	"github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/pet"
	"net/http"
)

type BreedHandler struct {
	breedService *pet.BreedService
}

func NewBreedHandler(breedService *pet.BreedService) *BreedHandler {
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
// @Success 200 {object} commonviews.PaginatedView[pet.BreedView]
// @Router /breeds [get]
func (h *BreedHandler) FindBreeds(w http.ResponseWriter, r *http.Request) {
	petTypeQuery := r.URL.Query().Get("pet_type")

	page, size, err := utils.ParsePaginationQueries(r, 1, 20)
	if err != nil {
		commonviews.BadRequest(w, nil, err.Error())
		return
	}

	var petType *string
	if petTypeQuery == "" {
		petType = nil
	} else {
		petType = &petTypeQuery
	}

	res, err := h.breedService.FindBreeds(page, size, petType)
	if err != nil {
		commonviews.InternalServerError(w, nil, err.Error())
		return
	}

	commonviews.OK(w, nil, commonviews.NewPaginatedView(page, size, res))
}
