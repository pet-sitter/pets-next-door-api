-- name: CreateRoom :one
INSERT INTO chat_rooms
(id,
 name,
 host_id,
 room_type,
 created_at,
 updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING id, name, room_type, created_at, updated_at;

-- name: DeleteRoom :exec
UPDATE
    chat_rooms
SET deleted_at = NOW()
WHERE id = $1;

-- name: ExistsUserInRoom :one
SELECT EXISTS (SELECT 1
               FROM user_chat_rooms
               WHERE room_id = $1
                 AND user_id = $2);

-- name: FindUserChatRoomByIDAndUserID :one
SELECT users.id as host_user_id,
       users.email as host_user_email,
       users.nickname as host_user_nickname,
       media.url as host_user_profile_image_url,
       chat_rooms.id as chat_room_id,
       chat_rooms.host_id as chat_room_host_id,
       chat_rooms.name as chat_room_name,
       chat_rooms.room_type as chat_room_type,
       chat_rooms.created_at as chat_room_created_at,
       chat_rooms.updated_at as chat_room_updated_at
FROM user_chat_rooms
         JOIN users ON users.id = user_chat_rooms.user_id
         LEFT OUTER JOIN media ON users.profile_image_id = media.id
         JOIN chat_rooms ON chat_rooms.id = user_chat_rooms.room_id
WHERE user_chat_rooms.left_at IS NULL
    AND chat_rooms.deleted_at IS NULL
    AND user_chat_rooms.room_id = $1
    AND user_chat_rooms.user_id = $2;

-- name: FindAllUserChatRoomsByUserUID :many
SELECT users.id as host_user_id,
       users.email as host_user_email,
       users.nickname as host_user_nickname,
       media.url as host_user_profile_image_url,
       chat_rooms.id as chat_room_id,
       chat_rooms.host_id as chat_room_host_id,
       chat_rooms.name as chat_room_name,
       chat_rooms.room_type as chat_room_type,
       chat_rooms.created_at as chat_room_created_at,
       chat_rooms.updated_at as chat_room_updated_at
FROM user_chat_rooms
JOIN users ON users.id = user_chat_rooms.user_id
LEFT OUTER JOIN media ON users.profile_image_id = media.id
JOIN chat_rooms ON chat_rooms.id = user_chat_rooms.room_id
WHERE user_chat_rooms.user_id = $1
AND user_chat_rooms.left_at IS NULL;

-- name: FindUserInfoByJoinUserId :many
SELECT users.id            AS join_user_id,
       users.email        AS join_user_email,
       users.nickname    AS join_user_nickname,
       users.fullname   AS join_user_fullname,
       media.url             AS profile_image_url,
       users.fb_provider_type AS join_user_fb_provider_type,
       users.fb_uid      AS join_user_fb_uid,
       users.created_at  AS join_user_created_at,
       users.updated_at AS join_user_updated_at
FROM user_chat_rooms
JOIN users ON users.id = user_chat_rooms.user_id
LEFT OUTER JOIN media
ON users.profile_image_id = media.id
WHERE user_chat_rooms.room_id = $1
AND user_chat_rooms.left_at IS NULL
AND user_chat_rooms.user_id != $2;


-- name: JoinRoom :one
INSERT INTO user_chat_rooms
(id,
 user_id,
 room_id,
 joined_at)
VALUES ($1, $2, $3, NOW())
RETURNING id, user_id, room_id, joined_at;


-- name: JoinRooms :exec
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

-- name: UserExistsInRoom :one
SELECT EXISTS (SELECT 1
               FROM user_chat_rooms
               WHERE room_id = $1
                 AND left_at IS NULL);


-- name: ExistsRoom :one
SELECT EXISTS (SELECT 1
               FROM chat_rooms
               WHERE id = $1
                 AND deleted_at IS NULL);


-- name: FindPrevMessageByRoomID :many
SELECT id,
       user_id,
       room_id,
       message_type,
       content,
       created_at
FROM chat_messages
WHERE chat_messages.deleted_at IS NULL
  AND room_id = $2
  AND (
    id < sqlc.narg('prev')::uuid
    )
ORDER BY chat_messages.created_at DESC
LIMIT $1;

-- name: FindNextMessageByRoomID :many
SELECT id,
       user_id,
       room_id,
       message_type,
       content,
       created_at
FROM chat_messages
WHERE chat_messages.deleted_at IS NULL
  AND room_id = $2
  AND (
    id > sqlc.narg('next')::uuid
    )
ORDER BY chat_messages.created_at ASC
LIMIT $1;

-- name: FindBetweenMessagesByRoomID :many
SELECT id,
       user_id,
       room_id,
       message_type,
       content,
       created_at
FROM chat_messages
WHERE chat_messages.deleted_at IS NULL
  AND room_id = $2
  AND id > sqlc.narg('prev')::uuid
  AND id < sqlc.narg('next')::uuid
ORDER BY chat_messages.created_at ASC
LIMIT $1;

-- name: HasPrevMessages :one
SELECT EXISTS (SELECT 1
               FROM chat_messages
               WHERE chat_messages.deleted_at IS NULL
                 AND room_id = $2
                 AND id < $1 -- 주어진 prev UUID보다 이전 메시지
               LIMIT 1);

-- name: HasNextMessages :one
SELECT EXISTS (SELECT 1
               FROM chat_messages
               WHERE chat_messages.deleted_at IS NULL
                 AND room_id = $2
                 AND id > $1 -- 주어진 next UUID보다 이후 메시지
               LIMIT 1);

-- name: FindMessagesByRoomIDAndSize :many
SELECT id,
       user_id,
       room_id,
       message_type,
       content,
       created_at
FROM chat_messages
WHERE chat_messages.deleted_at IS NULL
  AND room_id = $2
ORDER BY chat_messages.created_at DESC
LIMIT $1;

-- name: SaveChatMessage :one
INSERT INTO chat_messages
(id,
 user_id,
 room_id,
 message_type,
 content,
 created_at)
VALUES ($1, $2, $3, $4, $5, NOW())
RETURNING id, user_id, room_id, message_type, content, created_at;
