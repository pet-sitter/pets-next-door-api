-- name: WriteSOSPost :one
INSERT INTO sos_posts
(author_id,
 title,
 content,
 reward,
 care_type,
 carer_gender,
 reward_type,
 thumbnail_id,
 created_at,
 updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
RETURNING id, author_id, title, content, reward, care_type, carer_gender, reward_type, thumbnail_id, created_at, updated_at;

-- name: InsertSOSDate :one
INSERT INTO sos_dates
(date_start_at,
 date_end_at,
 created_at,
 updated_at)
VALUES ($1, $2, NOW(), NOW())
RETURNING id, date_start_at, date_end_at, created_at, updated_at;

-- name: LinkSOSPostDate :exec
INSERT INTO sos_posts_dates
(sos_post_id,
 sos_dates_id,
 created_at,
 updated_at)
VALUES ($1, $2, NOW(), NOW());

-- name: LinkResourceMedia :exec
INSERT INTO resource_media
(media_id,
 resource_id,
 resource_type,
 created_at,
 updated_at)
VALUES ($1, $2, $3, NOW(), NOW());

-- name: LinkSOSPostCondition :exec
INSERT INTO sos_posts_conditions
(sos_post_id,
 sos_condition_id,
 created_at,
 updated_at)
VALUES ($1, $2, NOW(), NOW());

-- name: LinkSOSPostPet :exec
INSERT INTO sos_posts_pets
(sos_post_id,
 pet_id,
 created_at,
 updated_at)
VALUES ($1, $2, NOW(), NOW());

-- name: FindSOSPosts :many
SELECT
    v_sos_posts.id,
    v_sos_posts.title,
    v_sos_posts.content,
    v_sos_posts.reward,
    v_sos_posts.reward_type,
    v_sos_posts.care_type,
    v_sos_posts.carer_gender,
    v_sos_posts.thumbnail_id,
    v_sos_posts.author_id,
    v_sos_posts.created_at,
    v_sos_posts.updated_at,
    v_sos_posts.dates,
    v_pets_for_sos_posts.pets_info,
    v_media_for_sos_posts.media_info,
    v_conditions.conditions_info
FROM
    v_sos_posts
        LEFT JOIN v_pets_for_sos_posts ON v_sos_posts.id = v_pets_for_sos_posts.sos_post_id
        LEFT JOIN v_media_for_sos_posts ON v_sos_posts.id = v_media_for_sos_posts.sos_post_id
        LEFT JOIN v_conditions ON v_sos_posts.id = v_conditions.sos_post_id
WHERE
    v_sos_posts.earliest_date_start_at >= sqlc.narg('earliest_date_start_at')
  AND (sqlc.narg('pet_type') = 'all' OR NOT EXISTS
    (SELECT 1
     FROM unnest(pet_type_list) AS pet_type
     WHERE pet_type <> sqlc.narg('pet_type')))
ORDER BY
    CASE WHEN sqlc.narg('sort_by') = 'newest' THEN v_sos_posts.created_at END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'deadline' THEN v_sos_posts.earliest_date_start_at END
LIMIT sqlc.narg('limit')
    OFFSET sqlc.narg('offset');


-- name: FindSOSPostsByAuthorID :many
SELECT
    v_sos_posts.id,
    v_sos_posts.title,
    v_sos_posts.content,
    v_sos_posts.reward,
    v_sos_posts.reward_type,
    v_sos_posts.care_type,
    v_sos_posts.carer_gender,
    v_sos_posts.thumbnail_id,
    v_sos_posts.author_id,
    v_sos_posts.created_at,
    v_sos_posts.updated_at,
    v_sos_posts.dates,
    v_pets_for_sos_posts.pets_info,
    v_media_for_sos_posts.media_info,
    v_conditions.conditions_info
FROM
    v_sos_posts
        LEFT JOIN v_pets_for_sos_posts ON v_sos_posts.id = v_pets_for_sos_posts.sos_post_id
        LEFT JOIN v_media_for_sos_posts ON v_sos_posts.id = v_media_for_sos_posts.sos_post_id
        LEFT JOIN v_conditions ON v_sos_posts.id = v_conditions.sos_post_id
WHERE
    v_sos_posts.earliest_date_start_at >= sqlc.narg('earliest_date_start_at')
  AND v_sos_posts.author_id = sqlc.narg('author_id')
  AND (sqlc.narg('pet_type') = 'all' OR NOT EXISTS
    (SELECT 1
     FROM unnest(pet_type_list) AS pet_type
     WHERE pet_type <> sqlc.narg('pet_type')))
ORDER BY
    CASE WHEN sqlc.narg('sort_by') = 'newest' THEN v_sos_posts.created_at END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'deadline' THEN v_sos_posts.earliest_date_start_at END
LIMIT sqlc.narg('limit')
    OFFSET sqlc.narg('offset');

-- name: FindSOSPostByID :one
SELECT
    v_sos_posts.id,
    v_sos_posts.title,
    v_sos_posts.content,
    v_sos_posts.reward,
    v_sos_posts.reward_type,
    v_sos_posts.care_type,
    v_sos_posts.carer_gender,
    v_sos_posts.thumbnail_id,
    v_sos_posts.author_id,
    v_sos_posts.created_at,
    v_sos_posts.updated_at,
    v_sos_posts.dates,
    v_pets_for_sos_posts.pets_info,
    v_media_for_sos_posts.media_info,
    v_conditions.conditions_info
FROM
    v_sos_posts
        LEFT JOIN v_pets_for_sos_posts ON v_sos_posts.id = v_pets_for_sos_posts.sos_post_id
        LEFT JOIN v_media_for_sos_posts ON v_sos_posts.id = v_media_for_sos_posts.sos_post_id
        LEFT JOIN v_conditions ON v_sos_posts.id = v_conditions.sos_post_id
WHERE
    v_sos_posts.id = sqlc.narg('id');

-- name: FindDatesBySOSPostID :many
SELECT
    sos_dates.id,
    sos_dates.date_start_at,
    sos_dates.date_end_at,
    sos_dates.created_at,
    sos_dates.updated_at
FROM
    sos_dates
        INNER JOIN
    sos_posts_dates
    ON sos_dates.id = sos_posts_dates.sos_dates_id
WHERE
    sos_posts_dates.sos_post_id = sqlc.narg('id') AND
    sos_posts_dates.deleted_at IS NULL;

-- name: UpdateSOSPost :one
UPDATE
    sos_posts
SET
    title = $1,
    content = $2,
    reward = $3,
    care_type = $4,
    carer_gender = $5,
    reward_type = $6,
    thumbnail_id = $7,
    updated_at = NOW()
WHERE
    id = $8
RETURNING
    id, author_id, title, content, reward, care_type, carer_gender, reward_type, thumbnail_id, created_at, updated_at;

-- name: DeleteSOSPostDateBySOSPostID :exec
UPDATE
    sos_posts_dates
SET
    deleted_at = NOW()
WHERE
    sos_post_id = $1;

-- name: DeleteSOSPostConditionBySOSPostID :exec
UPDATE
    sos_posts_conditions
SET
    deleted_at = NOW()
WHERE
    sos_post_id = $1;

-- name: DeleteSOSPostPetBySOSPostID :exec
UPDATE
    sos_posts_pets
SET
    deleted_at = NOW()
WHERE
    sos_post_id = $1;