package postgres

import (
	"context"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

type MediaPostgresStore struct {
	conn *database.Tx
}

func NewMediaPostgresStore(conn *database.Tx) *MediaPostgresStore {
	return &MediaPostgresStore{
		conn: conn,
	}
}

func (s *MediaPostgresStore) CreateMedia(ctx context.Context, media *media.Media) (*media.Media, *pnd.AppError) {
	return (&mediaQueries{conn: s.conn}).CreateMedia(ctx, media)
}

func (s *MediaPostgresStore) FindMediaByID(ctx context.Context, id int) (*media.Media, *pnd.AppError) {
	return (&mediaQueries{conn: s.conn}).FindMediaByID(ctx, id)
}

type mediaQueries struct {
	conn *database.Tx
}

func (s *mediaQueries) CreateMedia(ctx context.Context, media *media.Media) (*media.Media, *pnd.AppError) {
	const sql = `
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
	`

	if err := s.conn.QueryRowContext(ctx, sql,
		media.MediaType,
		media.URL,
	).Scan(&media.ID, &media.CreatedAt, &media.UpdatedAt); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return media, nil
}

func (s *mediaQueries) FindMediaByID(ctx context.Context, id int) (*media.Media, *pnd.AppError) {
	const sql = `
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
	`

	media := &media.Media{}
	if err := s.conn.QueryRowContext(ctx, sql,
		id,
	).Scan(
		&media.ID,
		&media.MediaType,
		&media.URL,
		&media.CreatedAt,
		&media.UpdatedAt,
	); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return media, nil
}
