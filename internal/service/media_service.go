package service

import (
	"context"
	"fmt"
	"github.com/pet-sitter/pets-next-door-api/internal/datatype"
	"io"

	"github.com/google/uuid"

	bucketinfra "github.com/pet-sitter/pets-next-door-api/internal/infra/bucket"

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

	fmt.Println("created", created)

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
		ID:        datatype.NewUUIDV7(),
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

func (s *MediaService) FindMediaByID(ctx context.Context, id uuid.UUID) (*media.DetailView, *pnd.AppError) {
	mediaData, err := databasegen.New(s.conn).FindSingleMedia(ctx, databasegen.FindSingleMediaParams{
		ID: uuid.NullUUID{UUID: id, Valid: true},
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return media.ToDetailView(mediaData), nil
}
