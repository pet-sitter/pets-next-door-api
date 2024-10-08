// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: pets.sql

package databasegen

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createPet = `-- name: CreatePet :one
INSERT INTO pets
(id,
 owner_id,
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
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
RETURNING id, created_at, updated_at
`

type CreatePetParams struct {
	ID             uuid.UUID
	OwnerID        uuid.UUID
	Name           string
	PetType        string
	Sex            string
	Neutered       bool
	Breed          string
	BirthDate      time.Time
	WeightInKg     string
	Remarks        string
	ProfileImageID uuid.NullUUID
}

type CreatePetRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreatePet(ctx context.Context, arg CreatePetParams) (CreatePetRow, error) {
	row := q.db.QueryRowContext(ctx, createPet,
		arg.ID,
		arg.OwnerID,
		arg.Name,
		arg.PetType,
		arg.Sex,
		arg.Neutered,
		arg.Breed,
		arg.BirthDate,
		arg.WeightInKg,
		arg.Remarks,
		arg.ProfileImageID,
	)
	var i CreatePetRow
	err := row.Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const deletePet = `-- name: DeletePet :exec
UPDATE
    pets
SET deleted_at = NOW()
WHERE id = $1
`

func (q *Queries) DeletePet(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deletePet, id)
	return err
}

const findPet = `-- name: FindPet :one
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
WHERE (pets.id = $1 OR $1 IS NULL)
  AND (pets.owner_id = $2 OR $2 IS NULL)
  AND ($3::boolean = TRUE OR
       ($3::boolean = FALSE AND pets.deleted_at IS NULL))
LIMIT 1
`

type FindPetParams struct {
	ID             uuid.NullUUID
	OwnerID        uuid.NullUUID
	IncludeDeleted bool
}

type FindPetRow struct {
	ID              uuid.UUID
	OwnerID         uuid.UUID
	Name            string
	PetType         string
	Sex             string
	Neutered        bool
	Breed           string
	BirthDate       time.Time
	WeightInKg      string
	Remarks         string
	ProfileImageUrl sql.NullString
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       sql.NullTime
}

func (q *Queries) FindPet(ctx context.Context, arg FindPetParams) (FindPetRow, error) {
	row := q.db.QueryRowContext(ctx, findPet, arg.ID, arg.OwnerID, arg.IncludeDeleted)
	var i FindPetRow
	err := row.Scan(
		&i.ID,
		&i.OwnerID,
		&i.Name,
		&i.PetType,
		&i.Sex,
		&i.Neutered,
		&i.Breed,
		&i.BirthDate,
		&i.WeightInKg,
		&i.Remarks,
		&i.ProfileImageUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const findPets = `-- name: FindPets :many
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
WHERE (pets.id = $3 OR $3 IS NULL)
  AND (pets.owner_id = $4 OR $4 IS NULL)
  AND ($5::boolean = TRUE OR
       ($5::boolean = FALSE AND pets.deleted_at IS NULL))
ORDER BY pets.created_at DESC
LIMIT $1 OFFSET $2
`

type FindPetsParams struct {
	Limit          int32
	Offset         int32
	ID             uuid.NullUUID
	OwnerID        uuid.NullUUID
	IncludeDeleted bool
}

type FindPetsRow struct {
	ID              uuid.UUID
	OwnerID         uuid.UUID
	Name            string
	PetType         string
	Sex             string
	Neutered        bool
	Breed           string
	BirthDate       time.Time
	WeightInKg      string
	Remarks         string
	ProfileImageUrl sql.NullString
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       sql.NullTime
}

func (q *Queries) FindPets(ctx context.Context, arg FindPetsParams) ([]FindPetsRow, error) {
	rows, err := q.db.QueryContext(ctx, findPets,
		arg.Limit,
		arg.Offset,
		arg.ID,
		arg.OwnerID,
		arg.IncludeDeleted,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindPetsRow
	for rows.Next() {
		var i FindPetsRow
		if err := rows.Scan(
			&i.ID,
			&i.OwnerID,
			&i.Name,
			&i.PetType,
			&i.Sex,
			&i.Neutered,
			&i.Breed,
			&i.BirthDate,
			&i.WeightInKg,
			&i.Remarks,
			&i.ProfileImageUrl,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findPetsByIDs = `-- name: FindPetsByIDs :many
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
WHERE pets.id = ANY ($1::uuid[])
  AND ($2::boolean = TRUE OR
       ($2::boolean = FALSE AND pets.deleted_at IS NULL))
ORDER BY pets.created_at DESC
`

type FindPetsByIDsParams struct {
	Ids            []uuid.UUID
	IncludeDeleted bool
}

type FindPetsByIDsRow struct {
	ID              uuid.UUID
	OwnerID         uuid.UUID
	Name            string
	PetType         string
	Sex             string
	Neutered        bool
	Breed           string
	BirthDate       time.Time
	WeightInKg      string
	Remarks         string
	ProfileImageUrl sql.NullString
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       sql.NullTime
}

func (q *Queries) FindPetsByIDs(ctx context.Context, arg FindPetsByIDsParams) ([]FindPetsByIDsRow, error) {
	rows, err := q.db.QueryContext(ctx, findPetsByIDs, pq.Array(arg.Ids), arg.IncludeDeleted)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindPetsByIDsRow
	for rows.Next() {
		var i FindPetsByIDsRow
		if err := rows.Scan(
			&i.ID,
			&i.OwnerID,
			&i.Name,
			&i.PetType,
			&i.Sex,
			&i.Neutered,
			&i.Breed,
			&i.BirthDate,
			&i.WeightInKg,
			&i.Remarks,
			&i.ProfileImageUrl,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findPetsBySOSPostID = `-- name: FindPetsBySOSPostID :many
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
         INNER JOIN
     sos_posts_pets
     ON
         pets.id = sos_posts_pets.pet_id
         LEFT JOIN
     media
     ON
         pets.profile_image_id = media.id
WHERE sos_posts_pets.sos_post_id = $1
  AND sos_posts_pets.deleted_at IS NULL
`

type FindPetsBySOSPostIDRow struct {
	ID              uuid.UUID
	OwnerID         uuid.UUID
	Name            string
	PetType         string
	Sex             string
	Neutered        bool
	Breed           string
	BirthDate       time.Time
	WeightInKg      string
	Remarks         string
	ProfileImageUrl sql.NullString
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       sql.NullTime
}

func (q *Queries) FindPetsBySOSPostID(ctx context.Context, sosPostID uuid.UUID) ([]FindPetsBySOSPostIDRow, error) {
	rows, err := q.db.QueryContext(ctx, findPetsBySOSPostID, sosPostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindPetsBySOSPostIDRow
	for rows.Next() {
		var i FindPetsBySOSPostIDRow
		if err := rows.Scan(
			&i.ID,
			&i.OwnerID,
			&i.Name,
			&i.PetType,
			&i.Sex,
			&i.Neutered,
			&i.Breed,
			&i.BirthDate,
			&i.WeightInKg,
			&i.Remarks,
			&i.ProfileImageUrl,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePet = `-- name: UpdatePet :exec
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
WHERE id = $1
`

type UpdatePetParams struct {
	ID             uuid.UUID
	Name           string
	Neutered       bool
	Breed          string
	BirthDate      time.Time
	WeightInKg     string
	Remarks        string
	ProfileImageID uuid.NullUUID
}

func (q *Queries) UpdatePet(ctx context.Context, arg UpdatePetParams) error {
	_, err := q.db.ExecContext(ctx, updatePet,
		arg.ID,
		arg.Name,
		arg.Neutered,
		arg.Breed,
		arg.BirthDate,
		arg.WeightInKg,
		arg.Remarks,
		arg.ProfileImageID,
	)
	return err
}
