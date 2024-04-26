-- name: CreateSOSCondition :one
INSERT INTO sos_conditions
(id,
 name,
 created_at,
 updated_at)
SELECT $1, $2, now(), now()
WHERE NOT EXISTS (SELECT 1
                  FROM sos_conditions
                  WHERE name = $2::VARCHAR(50))
RETURNING *;

-- name: FindConditions :many
SELECT id,
       name,
       created_at,
       updated_at,
       deleted_at
FROM sos_conditions
WHERE (sqlc.arg('include_deleted')::BOOLEAN = TRUE OR
       (sqlc.arg('include_deleted')::BOOLEAN = FALSE AND deleted_at IS NULL));

-- name: FindSOSPostConditions :many
SELECT sos_conditions.id,
       sos_conditions.name,
       sos_conditions.created_at,
       sos_conditions.updated_at
FROM sos_conditions
         INNER JOIN
     sos_posts_conditions
     ON
         sos_conditions.id = sos_posts_conditions.sos_condition_id
WHERE sos_posts_conditions.sos_post_id = $1
  AND (sqlc.arg('include_deleted')::BOOLEAN = TRUE OR
       (sqlc.arg('include_deleted')::BOOLEAN = FALSE AND sos_posts_conditions.deleted_at IS NULL));
