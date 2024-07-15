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

-- name: FindUserChatRooms :many
SELECT
    user_chat_rooms.id,
    user_chat_rooms.user_id,
    user_chat_rooms.room_id,
    user_chat_rooms.joined_at,
    users.email,
    users.nickname,
    users.fullname,
    media.url AS profile_image_url,
    users.fb_provider_type,
    users.fb_uid,
    users.created_at,
    users.updated_at,
    chat_rooms.id AS chat_room_id,
    chat_rooms.name AS chat_room_name,
    chat_rooms.room_type AS chat_room_type,
    chat_rooms.created_at AS chat_room_created_at,
    chat_rooms.updated_at AS chat_room_updated_at
FROM
    user_chat_rooms
    JOIN users
        ON users.id = user_chat_rooms.user_id
    JOIN chat_rooms
        ON chat_rooms.id = user_chat_rooms.room_id
    LEFT OUTER JOIN media
        ON users.profile_image_id = media.id
WHERE
    user_chat_rooms.left_at IS NULL;

-- name: ExistsUserInRoom :one
SELECT EXISTS (
    SELECT 1
    FROM user_chat_rooms
    WHERE room_id = $1 AND user_id = $2
) AS is_in_room;
