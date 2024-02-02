package postgres

import (
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

func (s *MediaPostgresStore) CreateMedia(media *media.Media) (*media.Media, *pnd.AppError) {
	tx, _ := s.db.Begin()
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
		return nil, pnd.FromPGError(err)
	}

	return media, nil
}

func (s *MediaPostgresStore) FindMediaByID(id int) (*media.Media, *pnd.AppError) {
	media := &media.Media{}

	tx, _ := s.db.Begin()
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
		return nil, pnd.FromPGError(err)
	}

	return media, nil
}
