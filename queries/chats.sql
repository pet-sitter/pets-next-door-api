-- name: CreateRoom :one
INSERT INTO chat_rooms
(name,
room_type,
created_at,
updated_at)
VALUES ($1, $2, NOW(), NOW())
RETURNING id, name, room_type, created_at, updated_at;

-- name: FindRoomByID :one
SELECT
    id,
    name,
    room_type,
    created_at,
    updated_at
FROM
    chat_rooms
WHERE
    (chat_rooms.id = sqlc.narg('id'))
    AND (chat_rooms.deleted_at IS NULL);

-- name: JoinRoom :one
INSERT INTO user_chat_rooms
(user_id, 
room_id,
joined_at)
VALUES ($1, $2, NOW())
RETURNING id, user_id, room_id, joined_at;

-- name: LeaveRoom :exec
UPDATE 
    user_chat_rooms
SET left_at = NOW()
WHERE user_id = $1
AND room_id = $2;

-- name: WriteMessage :one
INSERT INTO chat_messages
(user_id,
room_id,
message_type,
content,
created_at,
updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING id, user_id, room_id, message_type, content, created_at, updated_at;

-- name: FindMessageByRoomID :many
SELECT
    id,
    user_id,
    room_id,
    message_type,
    content
FROM
    chat_messages
WHERE
    room_id = sqlc.narg('room_id')
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;