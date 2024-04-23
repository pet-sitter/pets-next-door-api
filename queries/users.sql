-- name: CreateUser :one
INSERT INTO users
(email,
 nickname,
 fullname,
 password,
 profile_image_id,
 fb_provider_type,
 fb_uid,
 created_at,
 updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
RETURNING id, email, nickname, fullname, profile_image_id, fb_provider_type, fb_uid, created_at, updated_at;

-- name: FindUsers :many
SELECT users.id,
       users.nickname,
       media.url AS profile_image_url
FROM users
         LEFT OUTER JOIN
     media
     ON
         users.profile_image_id = media.id
WHERE (users.nickname = $1 OR $1 IS NULL)
  AND users.deleted_at IS NULL
ORDER BY users.created_at DESC
LIMIT $2 OFFSET $3;

-- name: FindUsersByID :one
SELECT users.id,
       users.email,
       users.nickname,
       users.fullname,
       media.url AS profile_image_url,
       users.fb_provider_type,
       users.fb_uid,
       users.created_at,
       users.updated_at,
       users.deleted_at
FROM users
         LEFT OUTER JOIN
     media
     ON
         users.profile_image_id = media.id
WHERE users.id = $1
  AND (users.deleted_at IS NULL OR $2);

-- name: FindUserByEmail :one
SELECT users.id,
       users.email,
       users.nickname,
       users.fullname,
       media.url AS profile_image_url,
       users.fb_provider_type,
       users.fb_uid,
       users.created_at,
       users.updated_at
FROM users
         LEFT OUTER JOIN
     media
     ON
         users.profile_image_id = media.id
WHERE users.email = $1
  AND users.deleted_at IS NULL;

-- name: FindUserByUID :one
SELECT users.id,
       users.email,
       users.nickname,
       users.fullname,
       media.url AS profile_image_url,
       users.fb_provider_type,
       users.fb_uid,
       users.created_at,
       users.updated_at
FROM users
         LEFT JOIN
     media
     ON
         users.profile_image_id = media.id
WHERE users.fb_uid = $1
  AND users.deleted_at IS NULL;

-- name: FindUserIDByFbUID :one
SELECT id
FROM users
WHERE fb_uid = $1
  AND deleted_at IS NULL;

-- name: ExistsUserByNickname :one
SELECT CASE
           WHEN
               EXISTS (SELECT 1
                       FROM users
                       WHERE nickname = $1
                         AND deleted_at IS NULL)
               THEN TRUE
           ELSE FALSE
           END;

-- name: FindUserStatusByEmail :one
SELECT fb_provider_type
FROM users
WHERE email = $1
  AND deleted_at IS NULL;

-- name: UpdateUserByUID :one
UPDATE
    users
SET nickname         = $1,
    profile_image_id = $2,
    updated_at       = NOW()
WHERE fb_uid = $3
  AND deleted_at IS NULL
RETURNING
    id,
    email,
    nickname,
    fullname,
    profile_image_id,
    fb_provider_type,
    fb_uid,
    created_at,
    updated_at;

-- name: DeleteUserByUID :exec
UPDATE
    users
SET deleted_at = NOW()
WHERE fb_uid = $1;
