// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: media.sql

package databasegen

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createMedia = `-- name: CreateMedia :one
INSERT INTO media
(id,
 media_type,
 url,
 created_at,
 updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING id, media_type, url, created_at, updated_at
`

type CreateMediaParams struct {
	ID        uuid.UUID
	MediaType string
	Url       string
}

type CreateMediaRow struct {
	ID        uuid.UUID
	MediaType string
	Url       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreateMedia(ctx context.Context, arg CreateMediaParams) (CreateMediaRow, error) {
	row := q.db.QueryRowContext(ctx, createMedia, arg.ID, arg.MediaType, arg.Url)
	var i CreateMediaRow
	err := row.Scan(
		&i.ID,
		&i.MediaType,
		&i.Url,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findSingleMedia = `-- name: FindSingleMedia :one
SELECT id,
       media_type,
       url,
       created_at,
       updated_at
FROM media
WHERE (id = $1 OR $1 IS NULL)
  AND ($2::BOOLEAN = TRUE OR
       ($2::BOOLEAN = FALSE AND deleted_at IS NULL))
`

type FindSingleMediaParams struct {
	ID             uuid.NullUUID
	IncludeDeleted bool
}

type FindSingleMediaRow struct {
	ID        uuid.UUID
	MediaType string
	Url       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) FindSingleMedia(ctx context.Context, arg FindSingleMediaParams) (FindSingleMediaRow, error) {
	row := q.db.QueryRowContext(ctx, findSingleMedia, arg.ID, arg.IncludeDeleted)
	var i FindSingleMediaRow
	err := row.Scan(
		&i.ID,
		&i.MediaType,
		&i.Url,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
