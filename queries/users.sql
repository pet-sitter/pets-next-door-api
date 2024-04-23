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
WHERE (users.id = sqlc.narg('id') OR sqlc.narg('id') IS NULL)
  AND (users.nickname = sqlc.narg('nickname') OR sqlc.narg('nickname') IS NULL)
  AND (users.email = sqlc.narg('email') OR sqlc.narg('email') IS NULL)
  AND (users.fb_uid = sqlc.narg('fb_uid') OR sqlc.narg('fb_uid') IS NULL)
  AND (users.deleted_at IS NULL OR sqlc.arg('include_deleted')::boolean = TRUE)
ORDER BY users.created_at DESC
LIMIT $1 OFFSET $2;

-- name: FindUser :one
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
WHERE (users.id = sqlc.narg('id') OR sqlc.narg('id') IS NULL)
  AND (users.nickname = sqlc.narg('nickname') OR sqlc.narg('nickname') IS NULL)
  AND (users.email = sqlc.narg('email') OR sqlc.narg('email') IS NULL)
  AND (users.fb_uid = sqlc.narg('fb_uid') OR sqlc.narg('fb_uid') IS NULL)
  AND (users.deleted_at IS NULL OR sqlc.arg('include_deleted')::boolean = TRUE);

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

-- name: UpdateUserByFbUID :one
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

-- name: DeleteUserByFbUID :exec
UPDATE
    users
SET deleted_at = NOW()
WHERE fb_uid = $1;
