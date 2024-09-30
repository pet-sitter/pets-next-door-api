-- name: CreateBreed :one
INSERT INTO breeds (id,
                    name,
                    pet_type,
                    created_at,
                    updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING id, pet_type, name, created_at, updated_at;

-- name: FindBreeds :many
SELECT id,
       name,
       pet_type,
       created_at,
       updated_at
FROM breeds
WHERE (pet_type = sqlc.narg('pet_type') OR sqlc.narg('pet_type') IS NULL)
  AND (name = sqlc.narg('name') OR sqlc.narg('name') IS NULL)
  AND (deleted_at IS NULL OR sqlc.arg('include_deleted')::boolean = TRUE)
ORDER BY id
LIMIT $1 OFFSET $2;
