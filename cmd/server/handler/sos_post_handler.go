package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/pet-sitter/pets-next-door-api/api/commonviews"
	webutils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/auth"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/sos_post"
)

type SosPostHandler struct {
	sosPostService sos_post.SosPostService
	authService    auth.AuthService
}

func NewSosPostHandler(sosPostService sos_post.SosPostService, authService auth.AuthService) *SosPostHandler {
	return &SosPostHandler{
		sosPostService: sosPostService,
		authService:    authService,
	}
}

// writeSosPost godoc
// @Summary 돌봄급구 게시글을 업로드합니다.
// @Description
// @Tags posts
// @Accept  json
// @Produce  json
// @Param request body sos_post.WriteSosPostRequest true "돌봄급구 게시글 업로드 요청"
// @Security FirebaseAuth
// @Success 201 {object} sos_post.WriteSosPostResponse
// @Router /posts/sos [post]
func (h *SosPostHandler) WriteSosPost(w http.ResponseWriter, r *http.Request) {
	foundUser, err := h.authService.VerifyAuthAndGetUser(r.Context(), r.Header.Get("Authorization"))
	if err != nil {
		commonviews.Unauthorized(w, nil, "unauthorized")
		return
	}

	uid := foundUser.FirebaseUID

	var writeSosPostRequest sos_post.WriteSosPostRequest

	if err := json.NewDecoder(r.Body).Decode(&writeSosPostRequest); err != nil {
		commonviews.BadRequest(w, nil, err.Error())
		return
	}
	if err := validator.New().Struct(writeSosPostRequest); err != nil {
		commonviews.BadRequest(w, nil, err.Error())
		return
	}

	res, err := h.sosPostService.WriteSosPost(uid, &writeSosPostRequest)
	if err != nil {
		commonviews.InternalServerError(w, nil, err.Error())
		return
	}

	commonviews.Created(w, nil, res)
}

// findSosPosts godoc
// @Summary 돌봄급구 게시글을 조회합니다.
// @Description
// @Tags posts
// @Accept  json
// @Produce  json
// @Param author_id query int false "작성자 ID"
// @Param page query int false "페이지 번호" default(1)
// @Param size query int false "페이지 사이즈" default(20)
// @Param sort_by query string false "정렬 기준" Enums(newest, deadline)
// @Success 200 {object} commonviews.PaginatedView[sos_post.FindSosPostResponse]
// @Router /posts/sos [get]
func (h *SosPostHandler) FindSosPosts(w http.ResponseWriter, r *http.Request) {
	authorIDQuery := r.URL.Query().Get("author_id")
	pageQuery := r.URL.Query().Get("page")
	sizeQuery := r.URL.Query().Get("size")
	sortByQuery := r.URL.Query().Get("sort_by")

	page := 1
	size := 20

	var err error
	if pageQuery != "" {
		page, err = strconv.Atoi(pageQuery)
		if err != nil {
			commonviews.BadRequest(w, nil, err.Error())
			return
		}
	}

	if sizeQuery != "" {
		size, err = strconv.Atoi(sizeQuery)
		if err != nil {
			commonviews.BadRequest(w, nil, err.Error())
			return
		}
	}

	var res []sos_post.FindSosPostResponse

	if authorIDQuery != "" {
		var authorID int
		authorID, err = strconv.Atoi(authorIDQuery)
		if err != nil {
			commonviews.BadRequest(w, nil, err.Error())
			return
		}

		res, err = h.sosPostService.FindSosPostsByAuthorID(authorID, page, size)
		if err != nil {
			commonviews.InternalServerError(w, nil, err.Error())
			return
		}
	} else {
		var sortBy string
		if sortByQuery == "" {
			sortBy = "newest"
		} else {
			sortBy = sortByQuery
		}

		res, err = h.sosPostService.FindSosPosts(page, size, sortBy)
		if err != nil {
			commonviews.InternalServerError(w, nil, err.Error())
			return
		}
	}

	commonviews.OK(w, nil, commonviews.NewPaginatedView(page, size, res))
}

// findSosPostsByID godoc
// @Summary 게시글 ID로 돌봄급구 게시글을 조회합니다.
// @Description
// @Tags posts
// @Accept  json
// @Produce  json
// @Param id path string true "게시글 ID"
// @Success 200 {object} sos_post.FindSosPostResponse
// @Router /posts/sos/{id} [get]
func (h *SosPostHandler) FindSosPostByID(w http.ResponseWriter, r *http.Request) {
	SosPostID, err := webutils.ParseIdFromPath(r, "id")
	if err != nil || SosPostID <= 0 {
		commonviews.NotFound(w, nil, "invalid sos_post ID")
		return
	}
	res, err := h.sosPostService.FindSosPostByID(SosPostID)
	if err != nil {
		commonviews.InternalServerError(w, nil, err.Error())
		return
	}

	commonviews.OK(w, nil, res)
}

// updateSosPost godoc
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
	foundUser, err := h.authService.VerifyAuthAndGetUser(r.Context(), r.Header.Get("Authorization"))
	if err != nil {
		commonviews.Unauthorized(w, nil, "unauthorized")
		return
	}

	uid := foundUser.FirebaseUID

	var updateSosPostRequest sos_post.UpdateSosPostRequest
	if err := commonviews.ParseBody(w, r, &updateSosPostRequest); err != nil {
		return
	}

	permission := h.sosPostService.CheckUpdatePermission(uid, updateSosPostRequest.ID)

	if !permission {
		commonviews.Forbidden(w, nil, "forbidden")
		return
	}

	res, err := h.sosPostService.UpdateSosPost(&updateSosPostRequest)
	if err != nil {
		commonviews.InternalServerError(w, nil, err.Error())
	}

	commonviews.OK(w, nil, res)
}
