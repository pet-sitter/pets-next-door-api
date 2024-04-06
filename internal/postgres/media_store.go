package postgres

import (
	"context"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database/pgx"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
)

func CreateMedia(ctx context.Context, tx *pgx.PgxTx, media *media.Media) (*media.Media, *pnd.AppError) {
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

	if err := tx.QueryRow(ctx, sql,
		media.MediaType,
		media.URL,
	).Scan(&media.ID, &media.CreatedAt, &media.UpdatedAt); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return media, nil
}

func FindMediaByID(ctx context.Context, tx *pgx.PgxTx, id int) (*media.Media, *pnd.AppError) {
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
	if err := tx.QueryRow(ctx, sql,
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
