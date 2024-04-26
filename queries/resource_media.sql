-- name: CreateResourceMedia :one
INSERT INTO resource_media
(resource_id,
 media_id,
 resource_type,
 created_at,
 updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING id, resource_id, media_id, resource_type, created_at, updated_at;

-- name: FindResourceMedia :many
SELECT m.id AS media_id,
       m.media_type,
       m.url,
       m.created_at,
       m.updated_at
FROM resource_media rm
         INNER JOIN
     media m
     ON
         rm.media_id = m.id
WHERE (rm.resource_id = sqlc.narg('resource_id') OR sqlc.narg('resource_id') IS NULL)
  AND (rm.resource_type = sqlc.narg('resource_type') OR sqlc.narg('resource_type') IS NULL)
  AND (sqlc.arg('include_deleted')::BOOLEAN = TRUE OR
       (sqlc.arg('include_deleted')::BOOLEAN = FALSE AND rm.deleted_at IS NULL));
