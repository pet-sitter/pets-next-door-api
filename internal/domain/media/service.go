package media

import (
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go/private/protocol/rest"
	"github.com/google/uuid"
	s3infra "github.com/pet-sitter/pets-next-door-api/internal/infra/s3"
)

type MediaService struct {
	mediaStore MediaStore
	s3Client   *s3infra.S3Client
}

type MediaServicer interface {
	UploadMedia(file io.ReadSeeker, mediaType MediaType, fileName string) (*Media, error)
	CreateMedia(media *Media) (*Media, error)
	FindMediaByID(id int) (*Media, error)
}

func NewMediaService(mediaStore MediaStore, s3Client *s3infra.S3Client) *MediaService {
	return &MediaService{
		mediaStore: mediaStore,
		s3Client:   s3Client,
	}
}

func generateRandomFileName(originalFileName string) string {
	extension := filepath.Ext(originalFileName)
	return uuid.New().String() + extension
}

type UploadFileResponse struct {
	FileEndpoint string
}

func (s *MediaService) UploadMedia(file io.ReadSeeker, mediaType MediaType, fileName string) (*Media, error) {
	randomFileName := generateRandomFileName(fileName)
	fullPath := "media/" + randomFileName

	if _, err := s.s3Client.UploadFile(file, fullPath, "media"); err != nil {
		return nil, err
	}

	req, _ := s.s3Client.GetFileRequest(fullPath)
	rest.Build(req)
	if err := req.Send(); err != nil {
		return nil, err
	}

	created, err := s.CreateMedia(&Media{
		MediaType: mediaType,
		URL:       req.HTTPRequest.URL.String(),
	})

	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *MediaService) CreateMedia(media *Media) (*Media, error) {
	created, err := s.mediaStore.CreateMedia(media)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *MediaService) FindMediaByID(id int) (*Media, error) {
	media, err := s.mediaStore.FindMediaByID(id)
	if err != nil {
		return nil, err
	}

	return media, nil
}
