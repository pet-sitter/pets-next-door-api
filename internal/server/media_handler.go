package server

import (
	"net/http"
	"strings"

	"github.com/pet-sitter/pets-next-door-api/api/views"
	webutils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/media"
	"github.com/pet-sitter/pets-next-door-api/internal/models"
)

type mediaHandler struct {
	mediaService media.MediaServicer
}

func newMediaHandler(mediaService media.MediaServicer) *mediaHandler {
	return &mediaHandler{
		mediaService: mediaService,
	}
}

type mediaView struct {
	ID        int              `json:"id"`
	MediaType models.MediaType `json:"mediaType"`
	URL       string           `json:"url"`
	CreatedAt string           `json:"createdAt"`
}

func (h *mediaHandler) findMediaByID(w http.ResponseWriter, r *http.Request) {
	id, err := webutils.ParseIdFromPath(r, "id")
	if err != nil || id <= 0 {
		views.NotFound(w, nil, "invalid media ID")
		return
	}

	media, err := h.mediaService.FindMediaByID(id)
	if err != nil {
		views.BadRequest(w, nil, err.Error())
		return
	}

	views.OK(w, nil, mediaView{
		ID:        media.ID,
		MediaType: media.MediaType,
		URL:       media.URL,
		CreatedAt: media.CreatedAt,
	})
}

func (h *mediaHandler) uploadImage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		views.BadRequest(w, nil, err.Error())
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		views.BadRequest(w, nil, err.Error())
		return
	}
	defer file.Close()

	if !isValidMimeType(header.Header.Get("Content-Type")) {
		views.BadRequest(w, nil, "invalid MIME type; supported MIME types are: ["+supportedMimeTypeString()+"]")
		return
	}

	res, err := h.mediaService.UploadMedia(file, models.IMAGE_MEDIA_TYPE, header.Filename)
	if err != nil {
		views.BadRequest(w, nil, err.Error())
		return
	}

	views.Created(w,
		nil,
		mediaView{
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
