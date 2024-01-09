package handler

import (
	"net/http"
	"strings"

	"github.com/pet-sitter/pets-next-door-api/api/commonviews"
	"github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
)

type mediaHandler struct {
	mediaService media.MediaServicer
}

func NewMediaHandler(mediaService media.MediaServicer) *mediaHandler {
	return &mediaHandler{
		mediaService: mediaService,
	}
}

// findMediaByID godoc
// @Summary 미디어를 ID로 조회합니다.
// @Description
// @Tags media
// @Produce  json
// @Param id path int true "미디어 ID"
// @Success 200 {object} media.MediaView
// @Router /media/{id} [get]
func (h *mediaHandler) FindMediaByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseIdFromPath(r, "id")
	if err != nil || id <= 0 {
		commonviews.NotFound(w, nil, "invalid media ID")
		return
	}

	found, err := h.mediaService.FindMediaByID(id)
	if err != nil {
		commonviews.BadRequest(w, nil, err.Error())
		return
	}

	commonviews.OK(w, nil, media.MediaView{
		ID:        found.ID,
		MediaType: found.MediaType,
		URL:       found.URL,
		CreatedAt: found.CreatedAt,
	})
}

// uploadImage godoc
// @Summary 이미지를 업로드합니다.
// @Description
// @Tags media
// @Accept  multipart/form-data
// @Produce  json
// @Param file formData file true "이미지 파일"
// @Success 201 {object} media.MediaView
// @Router /media/images [post]
func (h *mediaHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		commonviews.BadRequest(w, nil, err.Error())
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		commonviews.BadRequest(w, nil, err.Error())
		return
	}
	defer file.Close()

	if !isValidMimeType(header.Header.Get("Content-Type")) {
		commonviews.BadRequest(w, nil, "invalid MIME type; supported MIME types are: ["+supportedMimeTypeString()+"]")
		return
	}

	res, err := h.mediaService.UploadMedia(file, media.IMAGE_MEDIA_TYPE, header.Filename)
	if err != nil {
		commonviews.BadRequest(w, nil, err.Error())
		return
	}

	commonviews.Created(w,
		nil,
		media.MediaView{
			ID:        res.ID,
			MediaType: res.MediaType,
			URL:       res.URL,
			CreatedAt: res.CreatedAt,
		})
}

var supportedMimeTypes = []string{
	"image/jpeg",
	"image/png",
}

func isValidMimeType(mimeType string) bool {
	for _, supportedMimeType := range supportedMimeTypes {
		if mimeType == supportedMimeType {
			return true
		}
	}
	return false
}

func supportedMimeTypeString() string {
	return strings.Join(supportedMimeTypes, ", ")
}
