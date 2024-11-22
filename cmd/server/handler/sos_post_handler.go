package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sospost"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type SOSPostHandler struct {
	sosPostService service.SOSPostService
	authService    service.AuthService
}

func NewSOSPostHandler(
	sosPostService service.SOSPostService,
	authService service.AuthService,
) *SOSPostHandler {
	return &SOSPostHandler{
		sosPostService: sosPostService,
		authService:    authService,
	}
}

// WriteSOSPost godoc
// @Summary 돌봄급구 게시글을 업로드합니다.
// @Description
// @Tags posts
// @Accept  json
// @Produce  json
// @Param request body sospost.WriteSOSPostRequest true "돌봄급구 게시글 업로드 요청"
// @Security FirebaseAuth
// @Success 201 {object} sospost.DetailView
// @Router /posts/sos [post]
func (h *SOSPostHandler) WriteSOSPost(c echo.Context) error {
	foundUser, err := h.authService.VerifyAuthAndGetUser(
		c.Request().Context(),
		c.Request().Header.Get("Authorization"),
	)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	var writeSOSPostRequest sospost.WriteSOSPostRequest
	if err = pnd.ParseBody(c, &writeSOSPostRequest); err != nil {
		return c.JSON(err.StatusCode, err)
	}

	res, err := h.sosPostService.WriteSOSPost(
		c.Request().Context(),
		foundUser.FirebaseUID,
		&writeSOSPostRequest,
	)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusCreated, res)
}

// FindSOSPosts godoc
// @Summary 돌봄급구 게시글을 조회합니다.
// @Description
// @Tags posts
// @Accept  json
// @Produce  json
// @Param author_id query string false "작성자 ID"
// @Param page query int false "페이지 번호" default(1)
// @Param size query int false "페이지 사이즈" default(20)
// @Param sort_by query string false "정렬 기준" Enums(newest, deadline)
// @Param filter_type query string false "필터링 기준" Enums(dog, cat, all)
// @Success 200 {object} sospost.FindSOSPostListView
// @Router /posts/sos [get]
func (h *SOSPostHandler) FindSOSPosts(c echo.Context) error {
	authorID, err := pnd.ParseOptionalUUIDQuery(c, "author_id")
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	sortBy := "newest"
	if sortByQuery := pnd.ParseOptionalStringQuery(c, "sort_by"); sortByQuery != nil {
		sortBy = *sortByQuery
	}
	filterType := "all"
	if filterTypeQuery := pnd.ParseOptionalStringQuery(c, "filter_type"); filterTypeQuery != nil {
		filterType = *filterTypeQuery
	}

	page, size, err := pnd.ParsePaginationQueries(c, 1, 20)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	var res *sospost.FindSOSPostListView
	if authorID.Valid {
		res, err = h.sosPostService.FindSOSPostsByAuthorID(
			c.Request().Context(), authorID.UUID, page, size, sortBy, filterType)
		if err != nil {
			return c.JSON(err.StatusCode, err)
		}
	} else {
		res, err = h.sosPostService.FindSOSPosts(c.Request().Context(), page, size, sortBy, filterType)
		if err != nil {
			return c.JSON(err.StatusCode, err)
		}
	}

	return c.JSON(http.StatusOK, res)
}

// FindSOSPostByID godoc
// @Summary 게시글 ID로 돌봄급구 게시글을 조회합니다.
// @Description
// @Tags posts
// @Produce  json
// @Param id path int true "게시글 ID"
// @Success 200 {object} sospost.FindSOSPostView
// @Router /posts/sos/{id} [get]
func (h *SOSPostHandler) FindSOSPostByID(c echo.Context) error {
	id, err := pnd.ParseIDFromPath(c, "id")
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}
	res, err := h.sosPostService.FindSOSPostByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusOK, res)
}

// UpdateSOSPost godoc
// @Summary 돌봄급구 게시글을 수정합니다.
// @Description
// @Tags posts
// @Accept  json
// @Produce  json
// @Security FirebaseAuth
// @Param request body sospost.UpdateSOSPostRequest true "돌봄급구 수정 요청"
// @Success 200
// @Router /posts/sos [put]
func (h *SOSPostHandler) UpdateSOSPost(c echo.Context) error {
	foundUser, err := h.authService.VerifyAuthAndGetUser(
		c.Request().Context(),
		c.Request().Header.Get("Authorization"),
	)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	var updateSOSPostRequest sospost.UpdateSOSPostRequest
	if err = pnd.ParseBody(c, &updateSOSPostRequest); err != nil {
		return c.JSON(err.StatusCode, err)
	}

	permission, err := h.sosPostService.CheckUpdatePermission(
		c.Request().Context(),
		foundUser.FirebaseUID,
		updateSOSPostRequest.ID,
	)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}
	if !permission {
		pndErr := pnd.ErrForbidden(errors.New("해당 게시글에 대한 수정 권한이 없습니다"))
		return c.JSON(pndErr.StatusCode, pndErr)
	}

	res, err := h.sosPostService.UpdateSOSPost(c.Request().Context(), &updateSOSPostRequest)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusOK, res)
}
