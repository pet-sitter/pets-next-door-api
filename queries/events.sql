-- name: CreateEvent :one
INSERT INTO
  events (
    id,
    event_type,
    author_id,
    name,
    description,
    media_id,
    topics,
    max_participants,
    fee,
    start_at,
    created_at,
    updated_at
  )
VALUES
  (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    NOW(),
    NOW()
  )
RETURNING
  *;

-- name: FindEvents :many
SELECT
  events.*
FROM
  events
WHERE
  (
    events.deleted_at IS NULL
    OR sqlc.arg ('include_deleted')::boolean = TRUE
  )
  AND (id > sqlc.narg ('prev')::uuid OR sqlc.narg ('prev') IS NULL)
  AND (id < sqlc.narg ('next')::uuid OR sqlc.narg ('next') IS NULL)
  AND (events.author_id = sqlc.narg ('author_id') OR sqlc.narg ('author_id') IS NULL)
ORDER BY
  events.created_at DESC
LIMIT
  $1;

-- name: FindEvent :one
SELECT
  events.*
FROM
  events
WHERE
  (
    events.deleted_at IS NULL
    OR sqlc.arg ('include_deleted')::boolean = TRUE
  )
  AND (events.id = sqlc.narg('id') OR sqlc.narg('id') IS NULL)
LIMIT
  1;

-- name: UpdateEvent :one
UPDATE events
SET
  name = $1,
  description = $2,
  media_id = $3,
  topics = $4,
  max_participants = $5,
  fee = $6,
  start_at = $7,
  updated_at = NOW()
WHERE
  id = $8
RETURNING
  *;

-- name: DeleteEvent :exec
UPDATE events
SET
  deleted_at = NOW()
WHERE
  id = $1;
