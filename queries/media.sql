-- name: CreateMedia :one
INSERT INTO media
(id,
 media_type,
 url,
 created_at,
 updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING id, media_type, url, created_at, updated_at;

-- name: FindSingleMedia :one
SELECT id,
       media_type,
       url,
       created_at,
       updated_at
FROM media
WHERE (id = sqlc.narg('id') OR sqlc.narg('id') IS NULL)
  AND (sqlc.arg('include_deleted')::BOOLEAN = TRUE OR
       (sqlc.arg('include_deleted')::BOOLEAN = FALSE AND deleted_at IS NULL));

-- name: FindMediasByIDs :many
SELECT id,
	   media_type,
	   url,
	   created_at,
	   updated_at
FROM media
WHERE id = ANY (sqlc.arg('ids')::uuid[])
  AND (sqlc.arg('include_deleted')::BOOLEAN = TRUE OR
       (sqlc.arg('include_deleted')::BOOLEAN = FALSE AND deleted_at IS NULL));
