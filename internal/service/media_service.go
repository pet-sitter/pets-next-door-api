package service

import (
	"context"
	"io"

	bucketinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/bucket"

	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

type MediaService struct {
	conn     *database.DB
	uploader bucketinfra.FileUploader
}

func NewMediaService(conn *database.DB, uploader bucketinfra.FileUploader) *MediaService {
	return &MediaService{
		conn:     conn,
		uploader: uploader,
	}
}

type UploadFileView struct {
	FileEndpoint string
}

func (s *MediaService) UploadMedia(
	ctx context.Context, file io.ReadSeeker, mediaType media.Type, fileName string,
) (*media.DetailView, *pnd.AppError) {
	url, err := s.uploader.UploadFile(file, fileName)
	if err != nil {
		return nil, err
	}

	created, err := s.CreateMedia(ctx, mediaType, url)
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
