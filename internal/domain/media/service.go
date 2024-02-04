package media

import (
	"context"
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go/private/protocol/rest"
	"github.com/google/uuid"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	s3infra "github.com/pet-sitter/pets-next-door-api/internal/infra/s3"
)

type MediaService struct {
	mediaStore MediaStore
	s3Client   *s3infra.S3Client
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

type UploadFileView struct {
	FileEndpoint string
}

func (s *MediaService) UploadMedia(ctx context.Context, file io.ReadSeeker, mediaType MediaType, fileName string) (*Media, *pnd.AppError) {
	randomFileName := generateRandomFileName(fileName)
	fullPath := "media/" + randomFileName

	if _, err := s.s3Client.UploadFile(file, fullPath, "media"); err != nil {
		return nil, pnd.ErrUnknown(err)
	}

	req, _ := s.s3Client.GetFileRequest(fullPath)
	rest.Build(req)
	if err := req.Send(); err != nil {
		return nil, pnd.ErrUnknown(err)
	}

	created, err := s.CreateMedia(ctx, &Media{
		MediaType: mediaType,
		URL:       req.HTTPRequest.URL.String(),
	})

	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *MediaService) CreateMedia(ctx context.Context, media *Media) (*Media, *pnd.AppError) {
	created, err := s.mediaStore.CreateMedia(ctx, media)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *MediaService) FindMediaByID(ctx context.Context, id int) (*Media, *pnd.AppError) {
	media, err := s.mediaStore.FindMediaByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return media, nil
}
