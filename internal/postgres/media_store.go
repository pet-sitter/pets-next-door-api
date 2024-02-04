package postgres

import (
	"context"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

type MediaPostgresStore struct {
	db *database.DB
}

func NewMediaPostgresStore(db *database.DB) *MediaPostgresStore {
	return &MediaPostgresStore{
		db: db,
	}
}

func (s *MediaPostgresStore) CreateMedia(ctx context.Context, media *media.Media) (*media.Media, *pnd.AppError) {
	tx, _ := s.db.BeginTx(ctx)
	err := tx.QueryRow(`
	INSERT INTO
		media
		(
			media_type,
			url,
			created_at,
			updated_at
		)
	VALUES ($1, $2, NOW(), NOW())
	RETURNING id, created_at, updated_at
	`,
		media.MediaType,
		media.URL,
	).Scan(&media.ID, &media.CreatedAt, &media.UpdatedAt)
	tx.Commit()

	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return media, nil
}

func (s *MediaPostgresStore) FindMediaByID(ctx context.Context, id int) (*media.Media, *pnd.AppError) {
	media := &media.Media{}

	tx, _ := s.db.BeginTx(ctx)
	err := tx.QueryRow(`
	SELECT
		id,
		media_type,
		url,
		created_at,
		updated_at
	FROM
		media
	WHERE
		id = $1 AND
		deleted_at IS NULL
	`,
		id,
	).Scan(
		&media.ID,
		&media.MediaType,
		&media.URL,
		&media.CreatedAt,
		&media.UpdatedAt,
	)
	tx.Commit()

	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return media, nil
}
