package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sos_post"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type SosPostHandler struct {
	sosPostService service.SosPostService
	authService    service.AuthService
}

func NewSosPostHandler(sosPostService service.SosPostService, authService service.AuthService) *SosPostHandler {
	return &SosPostHandler{
		sosPostService: sosPostService,
		authService:    authService,
	}
}

// WriteSosPost godoc
// @Summary 돌봄급구 게시글을 업로드합니다.
// @Description
// @Tags posts
// @Accept  json
// @Produce  json
// @Param request body sos_post.WriteSosPostRequest true "돌봄급구 게시글 업로드 요청"
// @Security FirebaseAuth
// @Success 201 {object} sos_post.WriteSosPostView
// @Router /posts/sos [post]
func (h *SosPostHandler) WriteSosPost(c echo.Context) error {
	foundUser, err := h.authService.VerifyAuthAndGetUser(c.Request().Context(), c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	var writeSosPostRequest sos_post.WriteSosPostRequest
	if err := pnd.ParseBody(c, &writeSosPostRequest); err != nil {
		return c.JSON(err.StatusCode, err)
	}

	res, err := h.sosPostService.WriteSosPost(c.Request().Context(), foundUser.FirebaseUID, &writeSosPostRequest)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusCreated, res)
}

// FindSosPosts godoc
// @Summary 돌봄급구 게시글을 조회합니다.
// @Description
// @Tags posts
// @Accept  json
// @Produce  json
// @Param author_id query int false "작성자 ID"
// @Param page query int false "페이지 번호" default(1)
// @Param size query int false "페이지 사이즈" default(20)
// @Param sort_by query string false "정렬 기준" Enums(newest, deadline)
// @Success 200 {object} sos_post.FindSosPostListView
// @Router /posts/sos [get]
func (h *SosPostHandler) FindSosPosts(c echo.Context) error {
	authorID, err := pnd.ParseOptionalIntQuery(c, "author_id")
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	sortBy := "newest"
	if sortByQuery := pnd.ParseOptionalStringQuery(c, "sort_by"); sortByQuery != nil {
		sortBy = *sortByQuery
	}

	page, size, err := pnd.ParsePaginationQueries(c, 1, 20)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	var res *sos_post.FindSosPostListView
	if authorID != nil {
		res, err = h.sosPostService.FindSosPostsByAuthorID(c.Request().Context(), *authorID, page, size, sortBy)
		if err != nil {
			return c.JSON(err.StatusCode, err)
		}
	} else {
		res, err = h.sosPostService.FindSosPosts(c.Request().Context(), page, size, sortBy)
		if err != nil {
			return c.JSON(err.StatusCode, err)
		}
	}

	return c.JSON(http.StatusOK, res)
}

// FindSosPostByID godoc
// @Summary 게시글 ID로 돌봄급구 게시글을 조회합니다.
// @Description
// @Tags posts
// @Accept  json
// @Produce  json
// @Param id path string true "게시글 ID"
// @Success 200 {object} sos_post.FindSosPostView
// @Router /posts/sos/{id} [get]
func (h *SosPostHandler) FindSosPostByID(c echo.Context) error {
	SosPostID, err := pnd.ParseIDFromPath(c, "id")
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}
	res, err := h.sosPostService.FindSosPostByID(c.Request().Context(), *SosPostID)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusOK, res)
}

// UpdateSosPost godoc
// @Summary 돌봄급구 게시글을 수정합니다.
// @Description
// @Tags posts
// @Accept  json
// @Produce  json
// @Security FirebaseAuth
// @Param request body sos_post.UpdateSosPostRequest true "돌봄급구 수정 요청"
// @Success 200
// @Router /posts/sos [put]
func (h *SosPostHandler) UpdateSosPost(c echo.Context) error {
	foundUser, err := h.authService.VerifyAuthAndGetUser(c.Request().Context(), c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	var updateSosPostRequest sos_post.UpdateSosPostRequest
	if err := pnd.ParseBody(c, &updateSosPostRequest); err != nil {
		return c.JSON(err.StatusCode, err)
	}

	permission, err := h.sosPostService.CheckUpdatePermission(c.Request().Context(), foundUser.FirebaseUID, updateSosPostRequest.ID)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}
	if !permission {
		pndErr := pnd.ErrForbidden(fmt.Errorf("해당 게시글에 대한 수정 권한이 없습니다"))
		return c.JSON(pndErr.StatusCode, pndErr)
	}

	res, err := h.sosPostService.UpdateSosPost(c.Request().Context(), &updateSosPostRequest)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusOK, res)
}
