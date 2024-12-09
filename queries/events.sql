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
  id,
  author_id,
  name,
  description,
  media_id,
  topics,
  max_participants,
  fee,
  start_at,
  created_at,
  updated_at;

-- name: FindEvents :many
SELECT
  (
    events.id,
    events.event_type,
    events.author_id,
    events.name,
    events.description,
    events.media_id,
    events.topics,
    events.max_participants,
    events.fee,
    events.start_at,
    events.created_at,
    events.updated_at
  )
FROM
  events
  LEFT OUTER JOIN media ON events.media_id = media.id
WHERE
  (
    events.deleted_at IS NULL
    OR sqlc.arg ('include_deleted')::boolean = TRUE
  )
  AND id > sqlc.narg ('prev')::uuid
  AND id < sqlc.narg ('next')::uuid
  AND events.author_id = sqlc.narg ('author_id')
ORDER BY
  events.created_at DESC
LIMIT
  $1;

-- name: FindEventByID :one
SELECT
  (
    events.id,
    events.event_type,
    events.author_id,
    events.name,
    events.description,
    events.media_id,
    events.topics,
    events.max_participants,
    events.fee,
    events.start_at,
    events.created_at,
    events.updated_at
  )
FROM
  events
WHERE
  events.id = sqlc.narg ('id')
  AND (
    events.deleted_at IS NULL
    OR sqlc.arg ('include_deleted')::boolean = TRUE
  )
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
  id,
  author_id,
  name,
  description,
  media_id,
  topics,
  max_participants,
  fee,
  start_at,
  created_at,
  updated_at;

-- name: DeleteEvent :exec
UPDATE events
SET
  deleted_at = NOW()
WHERE
  id = $1;
