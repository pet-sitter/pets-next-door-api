package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
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
func (h *SosPostHandler) WriteSosPost(w http.ResponseWriter, r *http.Request) {
	foundUser, err := h.authService.VerifyAuthAndGetUser(r.Context(), r)
	if err != nil {
		render.Render(w, r, err)
		return
	}

	var writeSosPostRequest sos_post.WriteSosPostRequest
	if err := pnd.ParseBody(r, &writeSosPostRequest); err != nil {
		render.Render(w, r, err)
		return
	}

	res, err := h.sosPostService.WriteSosPost(r.Context(), foundUser.FirebaseUID, &writeSosPostRequest)
	if err != nil {
		render.Render(w, r, err)
		return
	}

	pnd.Created(w, nil, res)
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
func (h *SosPostHandler) FindSosPosts(w http.ResponseWriter, r *http.Request) {
	authorID, err := pnd.ParseOptionalIntQuery(r, "author_id")
	if err != nil {
		render.Render(w, r, err)
		return
	}

	sortBy := "newest"
	if sortByQuery := pnd.ParseOptionalStringQuery(r, "sort_by"); sortByQuery != nil {
		sortBy = *sortByQuery
	}

	page, size, err := pnd.ParsePaginationQueries(r, 1, 20)
	if err != nil {
		render.Render(w, r, err)
		return
	}

	var res *sos_post.FindSosPostListView
	if authorID != nil {
		res, err = h.sosPostService.FindSosPostsByAuthorID(r.Context(), *authorID, page, size, sortBy)
		if err != nil {
			render.Render(w, r, err)
			return
		}
	} else {
		res, err = h.sosPostService.FindSosPosts(r.Context(), page, size, sortBy)
		if err != nil {
			render.Render(w, r, err)
			return
		}
	}

	render.JSON(w, r, res)
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
func (h *SosPostHandler) FindSosPostByID(w http.ResponseWriter, r *http.Request) {
	SosPostID, err := pnd.ParseIdFromPath(r, "id")
	if err != nil {
		render.Render(w, r, err)
		return
	}
	res, err := h.sosPostService.FindSosPostByID(r.Context(), *SosPostID)
	if err != nil {
		render.Render(w, r, err)
		return
	}

	pnd.OK(w, nil, res)
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
func (h *SosPostHandler) UpdateSosPost(w http.ResponseWriter, r *http.Request) {
	foundUser, err := h.authService.VerifyAuthAndGetUser(r.Context(), r)
	if err != nil {
		render.Render(w, r, err)
		return
	}

	var updateSosPostRequest sos_post.UpdateSosPostRequest
	if err := pnd.ParseBody(r, &updateSosPostRequest); err != nil {
		render.Render(w, r, err)
		return
	}

	permission, err := h.sosPostService.CheckUpdatePermission(r.Context(), foundUser.FirebaseUID, updateSosPostRequest.ID)
	if err != nil {
		render.Render(w, r, err)
		return
	}
	if !permission {
		render.Render(w, r, pnd.ErrForbidden(fmt.Errorf("해당 게시글에 대한 수정 권한이 없습니다")))
		return
	}

	res, err := h.sosPostService.UpdateSosPost(r.Context(), &updateSosPostRequest)
	if err != nil {
		render.Render(w, r, err)
		return
	}

	pnd.OK(w, nil, res)
}
