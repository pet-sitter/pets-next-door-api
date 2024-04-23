package postgres

import (
	"context"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

func CreateMedia(ctx context.Context, tx *database.Tx, mediaData *media.Media) (*media.Media, *pnd.AppError) {
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

	if err := tx.QueryRowContext(ctx, sql,
		mediaData.MediaType,
		mediaData.URL,
	).Scan(&mediaData.ID, &mediaData.CreatedAt, &mediaData.UpdatedAt); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return mediaData, nil
}

func FindMediaByID(ctx context.Context, tx *database.Tx, id int) (*media.Media, *pnd.AppError) {
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

	mediaData := &media.Media{}
	if err := tx.QueryRowContext(ctx, sql,
		id,
	).Scan(
		&mediaData.ID,
		&mediaData.MediaType,
		&mediaData.URL,
		&mediaData.CreatedAt,
		&mediaData.UpdatedAt,
	); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return mediaData, nil
}
