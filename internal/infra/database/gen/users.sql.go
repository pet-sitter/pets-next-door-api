// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users.sql

package databasegen

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users
(id,
 email,
 nickname,
 fullname,
 password,
 profile_image_id,
 fb_provider_type,
 fb_uid,
 created_at,
 updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
RETURNING id, email, nickname, fullname, profile_image_id, fb_provider_type, fb_uid, created_at, updated_at
`

type CreateUserParams struct {
	ID             uuid.UUID
	Email          string
	Nickname       string
	Fullname       string
	Password       string
	ProfileImageID uuid.NullUUID
	FbProviderType sql.NullString
	FbUid          sql.NullString
}

type CreateUserRow struct {
	ID             uuid.UUID
	Email          string
	Nickname       string
	Fullname       string
	ProfileImageID uuid.NullUUID
	FbProviderType sql.NullString
	FbUid          sql.NullString
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Email,
		arg.Nickname,
		arg.Fullname,
		arg.Password,
		arg.ProfileImageID,
		arg.FbProviderType,
		arg.FbUid,
	)
	var i CreateUserRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Nickname,
		&i.Fullname,
		&i.ProfileImageID,
		&i.FbProviderType,
		&i.FbUid,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUserByFbUID = `-- name: DeleteUserByFbUID :exec
UPDATE
    users
SET deleted_at = NOW()
WHERE fb_uid = $1
`

func (q *Queries) DeleteUserByFbUID(ctx context.Context, fbUid sql.NullString) error {
	_, err := q.db.ExecContext(ctx, deleteUserByFbUID, fbUid)
	return err
}

const existsUserByNickname = `-- name: ExistsUserByNickname :one
SELECT CASE
           WHEN
               EXISTS (SELECT 1
                       FROM users
                       WHERE nickname = $1
                         AND deleted_at IS NULL)
               THEN TRUE
           ELSE FALSE
           END
`

func (q *Queries) ExistsUserByNickname(ctx context.Context, nickname string) (bool, error) {
	row := q.db.QueryRowContext(ctx, existsUserByNickname, nickname)
	var column_1 bool
	err := row.Scan(&column_1)
	return column_1, err
}

const findUser = `-- name: FindUser :one
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
WHERE (users.id = $1 OR $1 IS NULL)
  AND (users.nickname = $2 OR $2 IS NULL)
  AND (users.email = $3 OR $3 IS NULL)
  AND (users.fb_uid = $4 OR $4 IS NULL)
  AND (users.deleted_at IS NULL OR $5::boolean = TRUE)
LIMIT 1
`

type FindUserParams struct {
	ID             uuid.NullUUID
	Nickname       sql.NullString
	Email          sql.NullString
	FbUid          sql.NullString
	IncludeDeleted bool
}

type FindUserRow struct {
	ID              uuid.UUID
	Email           string
	Nickname        string
	Fullname        string
	ProfileImageUrl sql.NullString
	FbProviderType  sql.NullString
	FbUid           sql.NullString
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       sql.NullTime
}

func (q *Queries) FindUser(ctx context.Context, arg FindUserParams) (FindUserRow, error) {
	row := q.db.QueryRowContext(ctx, findUser,
		arg.ID,
		arg.Nickname,
		arg.Email,
		arg.FbUid,
		arg.IncludeDeleted,
	)
	var i FindUserRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Nickname,
		&i.Fullname,
		&i.ProfileImageUrl,
		&i.FbProviderType,
		&i.FbUid,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const findUsers = `-- name: FindUsers :many
SELECT users.id,
       users.nickname,
       media.url AS profile_image_url
FROM users
         LEFT OUTER JOIN
     media
     ON
         users.profile_image_id = media.id
WHERE (users.id = $3 OR $3 IS NULL)
  AND (users.nickname = $4 OR $4 IS NULL)
  AND (users.email = $5 OR $5 IS NULL)
  AND (users.fb_uid = $6 OR $6 IS NULL)
  AND (users.deleted_at IS NULL OR $7::boolean = TRUE)
ORDER BY users.created_at DESC
LIMIT $1 OFFSET $2
`

type FindUsersParams struct {
	Limit          int32
	Offset         int32
	ID             uuid.NullUUID
	Nickname       sql.NullString
	Email          sql.NullString
	FbUid          sql.NullString
	IncludeDeleted bool
}

type FindUsersRow struct {
	ID              uuid.UUID
	Nickname        string
	ProfileImageUrl sql.NullString
}

func (q *Queries) FindUsers(ctx context.Context, arg FindUsersParams) ([]FindUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, findUsers,
		arg.Limit,
		arg.Offset,
		arg.ID,
		arg.Nickname,
		arg.Email,
		arg.FbUid,
		arg.IncludeDeleted,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindUsersRow
	for rows.Next() {
		var i FindUsersRow
		if err := rows.Scan(&i.ID, &i.Nickname, &i.ProfileImageUrl); err != nil {
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

const updateUserByFbUID = `-- name: UpdateUserByFbUID :one
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
    updated_at
`

type UpdateUserByFbUIDParams struct {
	Nickname       string
	ProfileImageID uuid.NullUUID
	FbUid          sql.NullString
}

type UpdateUserByFbUIDRow struct {
	ID             uuid.UUID
	Email          string
	Nickname       string
	Fullname       string
	ProfileImageID uuid.NullUUID
	FbProviderType sql.NullString
	FbUid          sql.NullString
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (q *Queries) UpdateUserByFbUID(ctx context.Context, arg UpdateUserByFbUIDParams) (UpdateUserByFbUIDRow, error) {
	row := q.db.QueryRowContext(ctx, updateUserByFbUID, arg.Nickname, arg.ProfileImageID, arg.FbUid)
	var i UpdateUserByFbUIDRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Nickname,
		&i.Fullname,
		&i.ProfileImageID,
		&i.FbProviderType,
		&i.FbUid,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
