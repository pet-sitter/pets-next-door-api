package media

import (
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go/private/protocol/rest"
	"github.com/google/uuid"
	"github.com/pet-sitter/pets-next-door-api/internal/database"
	s3infra "github.com/pet-sitter/pets-next-door-api/internal/infra/s3"
	"github.com/pet-sitter/pets-next-door-api/internal/models"
)

type MediaService struct {
	db       *database.DB
	s3Client *s3infra.S3Client
}

type MediaServicer interface {
	UploadMedia(file io.ReadSeeker, mediaType models.MediaType, fileName string) (*models.Media, error)
	CreateMedia(media *models.Media) (*models.Media, error)
	FindMediaByID(id int) (*models.Media, error)
}

func NewMediaService(db *database.DB, s3Client *s3infra.S3Client) *MediaService {
	return &MediaService{
		db:       db,
		s3Client: s3Client,
	}
}

func generateRandomFileName(originalFileName string) string {
	extension := filepath.Ext(originalFileName)
	return uuid.New().String() + extension
}

type UploadFileResponse struct {
	FileEndpoint string
}

func (s *MediaService) UploadMedia(file io.ReadSeeker, mediaType models.MediaType, fileName string) (*models.Media, error) {
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

	created, err := s.CreateMedia(&models.Media{
		MediaType: mediaType,
		URL:       req.HTTPRequest.URL.String(),
	})

	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *MediaService) CreateMedia(media *models.Media) (*models.Media, error) {
	tx, _ := s.db.Begin()

	created, err := tx.CreateMedia(media)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return created, nil
}

func (s *MediaService) FindMediaByID(id int) (*models.Media, error) {
	tx, _ := s.db.Begin()

	media, err := tx.FindMediaByID(id)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return media, nil
}
