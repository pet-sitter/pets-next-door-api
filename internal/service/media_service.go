package service

import (
	"context"
	"io"
	"path/filepath"

	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"

	"github.com/aws/aws-sdk-go/private/protocol/rest"
	"github.com/google/uuid"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	s3infra "github.com/pet-sitter/pets-next-door-api/internal/infra/s3"
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

func (s *MediaService) UploadMedia(
	ctx context.Context, file io.ReadSeeker, mediaType media.Type, fileName string,
) (*media.DetailView, *pnd.AppError) {
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

	created, err := s.CreateMedia(ctx, mediaType, req.HTTPRequest.URL.String())
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *MediaService) CreateMedia(
	ctx context.Context, mediaType media.Type, url string,
) (*media.DetailView, *pnd.AppError) {
	tx, err := s.conn.BeginTx(ctx)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	created, err2 := databasegen.New(s.conn).CreateMedia(ctx, databasegen.CreateMediaParams{
		MediaType: mediaType.String(),
		Url:       url,
	})
	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return media.ToDetailViewFromCreated(created), nil
}

func (s *MediaService) FindMediaByID(ctx context.Context, id int64) (*media.DetailView, *pnd.AppError) {
	mediaData, err := databasegen.New(s.conn).FindSingleMedia(ctx, databasegen.FindSingleMediaParams{
		ID: utils.Int64ToNullInt32(id),
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return media.ToDetailView(mediaData), nil
}
