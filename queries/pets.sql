-- name: CreatePet :one
INSERT INTO pets
(owner_id,
 name,
 pet_type,
 sex,
 neutered,
 breed,
 birth_date,
 weight_in_kg,
 remarks,
 profile_image_id,
 created_at,
 updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
RETURNING id, created_at, updated_at;

-- name: FindPet :one
SELECT pets.id,
       pets.owner_id,
       pets.name,
       pets.pet_type,
       pets.sex,
       pets.neutered,
       pets.breed,
       pets.birth_date,
       pets.weight_in_kg,
       pets.remarks,
       media.url AS profile_image_url,
       pets.created_at,
       pets.updated_at,
       pets.deleted_at
FROM pets
         LEFT OUTER JOIN
     media
     ON
         pets.profile_image_id = media.id
WHERE (pets.id = sqlc.narg('id') OR sqlc.narg('id') IS NULL)
  AND (pets.owner_id = sqlc.narg('owner_id') OR sqlc.narg('owner_id') IS NULL)
  AND (sqlc.arg('include_deleted')::boolean = TRUE OR
       (sqlc.arg('include_deleted')::boolean = FALSE AND pets.deleted_at IS NULL))
LIMIT 1;

-- name: FindPets :many
SELECT pets.id,
       pets.owner_id,
       pets.name,
       pets.pet_type,
       pets.sex,
       pets.neutered,
       pets.breed,
       pets.birth_date,
       pets.weight_in_kg,
       pets.remarks,
       media.url AS profile_image_url,
       pets.created_at,
       pets.updated_at,
       pets.deleted_at
FROM pets
         LEFT OUTER JOIN
     media
     ON
         pets.profile_image_id = media.id
WHERE (pets.id = sqlc.narg('id') OR sqlc.narg('id') IS NULL)
  AND (pets.owner_id = sqlc.narg('owner_id') OR sqlc.narg('owner_id') IS NULL)
  AND (sqlc.arg('include_deleted')::boolean = TRUE OR
       (sqlc.arg('include_deleted')::boolean = FALSE AND pets.deleted_at IS NULL))
ORDER BY pets.created_at DESC
LIMIT $1 OFFSET $2;

-- name: FindPetsByIDs :many
SELECT pets.id,
       pets.owner_id,
       pets.name,
       pets.pet_type,
       pets.sex,
       pets.neutered,
       pets.breed,
       pets.birth_date,
       pets.weight_in_kg,
       pets.remarks,
       media.url AS profile_image_url,
       pets.created_at,
       pets.updated_at,
       pets.deleted_at
FROM pets
         LEFT OUTER JOIN
     media
     ON
         pets.profile_image_id = media.id
WHERE pets.id = ANY (sqlc.arg('ids')::int[])
  AND (sqlc.arg('include_deleted')::boolean = TRUE OR
       (sqlc.arg('include_deleted')::boolean = FALSE AND pets.deleted_at IS NULL))
ORDER BY pets.created_at DESC;

-- name: UpdatePet :exec
UPDATE
    pets
SET name             = $2,
    neutered         = $3,
    breed            = $4,
    birth_date       = $5,
    weight_in_kg     = $6,
    remarks          = $7,
    profile_image_id = $8,
    updated_at       = NOW()
WHERE id = $1;

-- name: DeletePet :exec
UPDATE
    pets
SET deleted_at = NOW()
WHERE id = $1;
