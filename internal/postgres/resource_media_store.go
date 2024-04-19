package postgres

import (
	"context"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

func CreateResourceMedia(ctx context.Context, tx *database.Tx, resourceID int, mediaID int, resourceType string) (*media.ResourceMedia, *pnd.AppError) {
	const sql = `
	INSERT INTO
		resource_media
		(
			resource_id,
			media_id,
			resource_type,
			created_at,
			updated_at
		)
	VALUES ($1, $2, $3, NOW(), NOW())
	RETURNING id, resource_id, media_id, created_at, updated_at
	`

	resourceMedia := &media.ResourceMedia{}
	err := tx.QueryRowContext(ctx, sql,
		resourceID,
		mediaID,
		resourceType,
	).Scan(
		&resourceMedia.ID,
		&resourceMedia.ResourceID,
		&resourceMedia.MediaID,
		&resourceMedia.CreatedAt,
		&resourceMedia.UpdatedAt,
	)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return resourceMedia, nil
}

func FindResourceMediaByResourceID(
	ctx context.Context, tx *database.Tx, resourceID int, resourceType string,
) (*media.MediaList, *pnd.AppError) {
	const sql = `
	SELECT
		m.id,
		m.media_type,
		m.url,
		m.created_at,
		m.updated_at
	FROM
		resource_media rm
	INNER JOIN
		media m
	ON
		rm.media_id = m.id
	WHERE
		rm.resource_id = $1 AND
		rm.resource_type = $2 AND
		rm.deleted_at IS NULL
	`

	var mediaList media.MediaList
	rows, err := tx.QueryContext(ctx, sql,
		resourceID,
		resourceType,
	)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	defer rows.Close()

	for rows.Next() {
		mediaItem := media.Media{}
		if err := rows.Scan(
			&mediaItem.ID,
			&mediaItem.MediaType,
			&mediaItem.URL,
			&mediaItem.CreatedAt,
			&mediaItem.UpdatedAt,
		); err != nil {
			return nil, pnd.FromPostgresError(err)
		}
		mediaList = append(mediaList, &mediaItem)
	}
	if err := rows.Err(); err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return &mediaList, nil
}
