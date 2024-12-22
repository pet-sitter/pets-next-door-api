package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type MediaHandler struct {
	mediaService service.MediaService
}

func NewMediaHandler(mediaService service.MediaService) *MediaHandler {
	return &MediaHandler{
		mediaService: mediaService,
	}
}

// FindMediaByID godoc
// @Summary 미디어를 ID로 조회합니다.
// @Description
// @Tags media
// @Produce  json
// @Param id path int true "미디어 ID"
// @Success 200 {object} media.DetailView
// @Router /media/{id} [get]
func (h *MediaHandler) FindMediaByID(c echo.Context) error {
	id, err := pnd.ParseIDFromPath(c, "id")
	if err != nil {
		return err
	}

	found, err := h.mediaService.FindMediaByID(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, found)
}

// UploadImage godoc
// @Summary 이미지를 업로드합니다.
// @Description
// @Tags media
// @Accept  multipart/form-data
// @Produce  json
// @Param file formData file true "이미지 파일"
// @Success 201 {object} media.DetailView
// @Router /media/images [post]
func (h *MediaHandler) UploadImage(c echo.Context) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return pnd.ErrMultipartFormError(errors.New("file must be provided"))
	}

	if fileHeader.Size > 10<<20 {
		return pnd.ErrMultipartFormError(errors.New("file size must be less than 10MB"))
	}

	file, err := fileHeader.Open()
	if err != nil {
		return pnd.ErrMultipartFormError(errors.New("failed to open file"))
	}
	defer file.Close()

	if !isValidMimeType(fileHeader.Header.Get("Content-Type")) {
		return pnd.ErrMultipartFormError(
			errors.New(
				"invalid MIME type; supported MIME types are: [" + supportedMimeTypeString() + "]",
			),
		)
	}

	res, err := h.mediaService.UploadMedia(
		c.Request().Context(),
		file,
		media.TypeImage,
		fileHeader.Filename,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, res)
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
