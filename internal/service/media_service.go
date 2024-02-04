package service

import (
	"context"
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go/private/protocol/rest"
	"github.com/google/uuid"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	s3infra "github.com/pet-sitter/pets-next-door-api/internal/infra/s3"
	"github.com/pet-sitter/pets-next-door-api/internal/postgres"
)

type MediaService struct {
	conn     *database.DB
	s3Client *s3infra.S3Client
}

func NewMediaService(conn *database.DB, s3Client *s3infra.S3Client) *MediaService {
	return &MediaService{
		conn:     conn,
		s3Client: s3Client,
	}
}

func generateRandomFileName(originalFileName string) string {
	extension := filepath.Ext(originalFileName)
	return uuid.New().String() + extension
}

type UploadFileView struct {
	FileEndpoint string
}

func (s *MediaService) UploadMedia(ctx context.Context, file io.ReadSeeker, mediaType media.MediaType, fileName string) (*media.Media, *pnd.AppError) {
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

	created, err := s.CreateMedia(ctx, &media.Media{
		MediaType: mediaType,
		URL:       req.HTTPRequest.URL.String(),
	})

	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *MediaService) CreateMedia(ctx context.Context, mediaData *media.Media) (*media.Media, *pnd.AppError) {
	var created *media.Media
	var err *pnd.AppError

	err = database.WithTransaction(ctx, s.conn, func(tx *database.Tx) *pnd.AppError {
		mediaStore := postgres.NewMediaPostgresStore(tx)

		created, err = mediaStore.CreateMedia(ctx, mediaData)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *MediaService) FindMediaByID(ctx context.Context, id int) (*media.Media, *pnd.AppError) {
	var media *media.Media
	var err *pnd.AppError

	err = database.WithTransaction(ctx, s.conn, func(tx *database.Tx) *pnd.AppError {
		mediaStore := postgres.NewMediaPostgresStore(tx)

		media, err = mediaStore.FindMediaByID(ctx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return media, nil
}
