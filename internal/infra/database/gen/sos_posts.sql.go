// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: sos_posts.sql

package databasegen

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/sqlc-dev/pqtype"
)

const deleteSOSPostConditionBySOSPostID = `-- name: DeleteSOSPostConditionBySOSPostID :exec
UPDATE
    sos_posts_conditions
SET
    deleted_at = NOW()
WHERE
    sos_post_id = $1
`

func (q *Queries) DeleteSOSPostConditionBySOSPostID(ctx context.Context, sosPostID sql.NullInt64) error {
	_, err := q.db.ExecContext(ctx, deleteSOSPostConditionBySOSPostID, sosPostID)
	return err
}

const deleteSOSPostDateBySOSPostID = `-- name: DeleteSOSPostDateBySOSPostID :exec
UPDATE
    sos_posts_dates
SET
    deleted_at = NOW()
WHERE
    sos_post_id = $1
`

func (q *Queries) DeleteSOSPostDateBySOSPostID(ctx context.Context, sosPostID sql.NullInt64) error {
	_, err := q.db.ExecContext(ctx, deleteSOSPostDateBySOSPostID, sosPostID)
	return err
}

const deleteSOSPostPetBySOSPostID = `-- name: DeleteSOSPostPetBySOSPostID :exec
UPDATE
    sos_posts_pets
SET
    deleted_at = NOW()
WHERE
    sos_post_id = $1
`

func (q *Queries) DeleteSOSPostPetBySOSPostID(ctx context.Context, sosPostID sql.NullInt64) error {
	_, err := q.db.ExecContext(ctx, deleteSOSPostPetBySOSPostID, sosPostID)
	return err
}

const findDatesBySOSPostID = `-- name: FindDatesBySOSPostID :many
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
    sos_posts_dates.sos_post_id = $1 AND
    sos_posts_dates.deleted_at IS NULL
`

type FindDatesBySOSPostIDRow struct {
	ID          int32
	DateStartAt sql.NullTime
	DateEndAt   sql.NullTime
	CreatedAt   sql.NullTime
	UpdatedAt   sql.NullTime
}

func (q *Queries) FindDatesBySOSPostID(ctx context.Context, id sql.NullInt64) ([]FindDatesBySOSPostIDRow, error) {
	rows, err := q.db.QueryContext(ctx, findDatesBySOSPostID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindDatesBySOSPostIDRow
	for rows.Next() {
		var i FindDatesBySOSPostIDRow
		if err := rows.Scan(
			&i.ID,
			&i.DateStartAt,
			&i.DateEndAt,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const findSOSPostByID = `-- name: FindSOSPostByID :one
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
    v_sos_posts.id = $1
`

type FindSOSPostByIDRow struct {
	ID             int32
	Title          sql.NullString
	Content        sql.NullString
	Reward         sql.NullString
	RewardType     sql.NullString
	CareType       sql.NullString
	CarerGender    sql.NullString
	ThumbnailID    sql.NullInt64
	AuthorID       sql.NullInt64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Dates          json.RawMessage
	PetsInfo       pqtype.NullRawMessage
	MediaInfo      pqtype.NullRawMessage
	ConditionsInfo pqtype.NullRawMessage
}

func (q *Queries) FindSOSPostByID(ctx context.Context, id sql.NullInt32) (FindSOSPostByIDRow, error) {
	row := q.db.QueryRowContext(ctx, findSOSPostByID, id)
	var i FindSOSPostByIDRow
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Content,
		&i.Reward,
		&i.RewardType,
		&i.CareType,
		&i.CarerGender,
		&i.ThumbnailID,
		&i.AuthorID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Dates,
		&i.PetsInfo,
		&i.MediaInfo,
		&i.ConditionsInfo,
	)
	return i, err
}

const findSOSPosts = `-- name: FindSOSPosts :many
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
    v_sos_posts.earliest_date_start_at >= $1
  AND ($2 = 'all' OR NOT EXISTS
    (SELECT 1
     FROM unnest(pet_type_list) AS pet_type
     WHERE pet_type <> $2))
ORDER BY
    CASE WHEN $3 = 'newest' THEN v_sos_posts.created_at END DESC,
    CASE WHEN $3 = 'deadline' THEN v_sos_posts.earliest_date_start_at END
LIMIT $5
    OFFSET $4
`

type FindSOSPostsParams struct {
	EarliestDateStartAt interface{}
	PetType             interface{}
	SortBy              interface{}
	Offset              sql.NullInt32
	Limit               sql.NullInt32
}

type FindSOSPostsRow struct {
	ID             int32
	Title          sql.NullString
	Content        sql.NullString
	Reward         sql.NullString
	RewardType     sql.NullString
	CareType       sql.NullString
	CarerGender    sql.NullString
	ThumbnailID    sql.NullInt64
	AuthorID       sql.NullInt64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Dates          json.RawMessage
	PetsInfo       pqtype.NullRawMessage
	MediaInfo      pqtype.NullRawMessage
	ConditionsInfo pqtype.NullRawMessage
}

func (q *Queries) FindSOSPosts(ctx context.Context, arg FindSOSPostsParams) ([]FindSOSPostsRow, error) {
	rows, err := q.db.QueryContext(ctx, findSOSPosts,
		arg.EarliestDateStartAt,
		arg.PetType,
		arg.SortBy,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindSOSPostsRow
	for rows.Next() {
		var i FindSOSPostsRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Content,
			&i.Reward,
			&i.RewardType,
			&i.CareType,
			&i.CarerGender,
			&i.ThumbnailID,
			&i.AuthorID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Dates,
			&i.PetsInfo,
			&i.MediaInfo,
			&i.ConditionsInfo,
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

const findSOSPostsByAuthorID = `-- name: FindSOSPostsByAuthorID :many
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
    v_sos_posts.earliest_date_start_at >= $1
  AND v_sos_posts.author_id = $2
  AND ($3 = 'all' OR NOT EXISTS
    (SELECT 1
     FROM unnest(pet_type_list) AS pet_type
     WHERE pet_type <> $3))
ORDER BY
    CASE WHEN $4 = 'newest' THEN v_sos_posts.created_at END DESC,
    CASE WHEN $4 = 'deadline' THEN v_sos_posts.earliest_date_start_at END
LIMIT $6
    OFFSET $5
`

type FindSOSPostsByAuthorIDParams struct {
	EarliestDateStartAt interface{}
	AuthorID            sql.NullInt64
	PetType             interface{}
	SortBy              interface{}
	Offset              sql.NullInt32
	Limit               sql.NullInt32
}

type FindSOSPostsByAuthorIDRow struct {
	ID             int32
	Title          sql.NullString
	Content        sql.NullString
	Reward         sql.NullString
	RewardType     sql.NullString
	CareType       sql.NullString
	CarerGender    sql.NullString
	ThumbnailID    sql.NullInt64
	AuthorID       sql.NullInt64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Dates          json.RawMessage
	PetsInfo       pqtype.NullRawMessage
	MediaInfo      pqtype.NullRawMessage
	ConditionsInfo pqtype.NullRawMessage
}

func (q *Queries) FindSOSPostsByAuthorID(ctx context.Context, arg FindSOSPostsByAuthorIDParams) ([]FindSOSPostsByAuthorIDRow, error) {
	rows, err := q.db.QueryContext(ctx, findSOSPostsByAuthorID,
		arg.EarliestDateStartAt,
		arg.AuthorID,
		arg.PetType,
		arg.SortBy,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindSOSPostsByAuthorIDRow
	for rows.Next() {
		var i FindSOSPostsByAuthorIDRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Content,
			&i.Reward,
			&i.RewardType,
			&i.CareType,
			&i.CarerGender,
			&i.ThumbnailID,
			&i.AuthorID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Dates,
			&i.PetsInfo,
			&i.MediaInfo,
			&i.ConditionsInfo,
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

const insertSOSDate = `-- name: InsertSOSDate :one
INSERT INTO sos_dates
(date_start_at,
 date_end_at,
 created_at,
 updated_at)
VALUES ($1, $2, NOW(), NOW())
RETURNING id, date_start_at, date_end_at, created_at, updated_at
`

type InsertSOSDateParams struct {
	DateStartAt sql.NullTime
	DateEndAt   sql.NullTime
}

type InsertSOSDateRow struct {
	ID          int32
	DateStartAt sql.NullTime
	DateEndAt   sql.NullTime
	CreatedAt   sql.NullTime
	UpdatedAt   sql.NullTime
}

func (q *Queries) InsertSOSDate(ctx context.Context, arg InsertSOSDateParams) (InsertSOSDateRow, error) {
	row := q.db.QueryRowContext(ctx, insertSOSDate, arg.DateStartAt, arg.DateEndAt)
	var i InsertSOSDateRow
	err := row.Scan(
		&i.ID,
		&i.DateStartAt,
		&i.DateEndAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const linkResourceMedia = `-- name: LinkResourceMedia :exec
INSERT INTO resource_media
(media_id,
 resource_id,
 resource_type,
 created_at,
 updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
`

type LinkResourceMediaParams struct {
	MediaID      sql.NullInt64
	ResourceID   sql.NullInt64
	ResourceType sql.NullString
}

func (q *Queries) LinkResourceMedia(ctx context.Context, arg LinkResourceMediaParams) error {
	_, err := q.db.ExecContext(ctx, linkResourceMedia, arg.MediaID, arg.ResourceID, arg.ResourceType)
	return err
}

const linkSOSPostCondition = `-- name: LinkSOSPostCondition :exec
INSERT INTO sos_posts_conditions
(sos_post_id,
 sos_condition_id,
 created_at,
 updated_at)
VALUES ($1, $2, NOW(), NOW())
`

type LinkSOSPostConditionParams struct {
	SosPostID      sql.NullInt64
	SosConditionID sql.NullInt64
}

func (q *Queries) LinkSOSPostCondition(ctx context.Context, arg LinkSOSPostConditionParams) error {
	_, err := q.db.ExecContext(ctx, linkSOSPostCondition, arg.SosPostID, arg.SosConditionID)
	return err
}

const linkSOSPostDate = `-- name: LinkSOSPostDate :exec
INSERT INTO sos_posts_dates
(sos_post_id,
 sos_dates_id,
 created_at,
 updated_at)
VALUES ($1, $2, NOW(), NOW())
`

type LinkSOSPostDateParams struct {
	SosPostID  sql.NullInt64
	SosDatesID sql.NullInt64
}

func (q *Queries) LinkSOSPostDate(ctx context.Context, arg LinkSOSPostDateParams) error {
	_, err := q.db.ExecContext(ctx, linkSOSPostDate, arg.SosPostID, arg.SosDatesID)
	return err
}

const linkSOSPostPet = `-- name: LinkSOSPostPet :exec
INSERT INTO sos_posts_pets
(sos_post_id,
 pet_id,
 created_at,
 updated_at)
VALUES ($1, $2, NOW(), NOW())
`

type LinkSOSPostPetParams struct {
	SosPostID sql.NullInt64
	PetID     sql.NullInt64
}

func (q *Queries) LinkSOSPostPet(ctx context.Context, arg LinkSOSPostPetParams) error {
	_, err := q.db.ExecContext(ctx, linkSOSPostPet, arg.SosPostID, arg.PetID)
	return err
}

const updateSOSPost = `-- name: UpdateSOSPost :one
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
    id, author_id, title, content, reward, care_type, carer_gender, reward_type, thumbnail_id, created_at, updated_at
`

type UpdateSOSPostParams struct {
	Title       sql.NullString
	Content     sql.NullString
	Reward      sql.NullString
	CareType    sql.NullString
	CarerGender sql.NullString
	RewardType  sql.NullString
	ThumbnailID sql.NullInt64
	ID          int32
}

type UpdateSOSPostRow struct {
	ID          int32
	AuthorID    sql.NullInt64
	Title       sql.NullString
	Content     sql.NullString
	Reward      sql.NullString
	CareType    sql.NullString
	CarerGender sql.NullString
	RewardType  sql.NullString
	ThumbnailID sql.NullInt64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (q *Queries) UpdateSOSPost(ctx context.Context, arg UpdateSOSPostParams) (UpdateSOSPostRow, error) {
	row := q.db.QueryRowContext(ctx, updateSOSPost,
		arg.Title,
		arg.Content,
		arg.Reward,
		arg.CareType,
		arg.CarerGender,
		arg.RewardType,
		arg.ThumbnailID,
		arg.ID,
	)
	var i UpdateSOSPostRow
	err := row.Scan(
		&i.ID,
		&i.AuthorID,
		&i.Title,
		&i.Content,
		&i.Reward,
		&i.CareType,
		&i.CarerGender,
		&i.RewardType,
		&i.ThumbnailID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const writeSOSPost = `-- name: WriteSOSPost :one
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
RETURNING id, author_id, title, content, reward, care_type, carer_gender, reward_type, thumbnail_id, created_at, updated_at
`

type WriteSOSPostParams struct {
	AuthorID    sql.NullInt64
	Title       sql.NullString
	Content     sql.NullString
	Reward      sql.NullString
	CareType    sql.NullString
	CarerGender sql.NullString
	RewardType  sql.NullString
	ThumbnailID sql.NullInt64
}

type WriteSOSPostRow struct {
	ID          int32
	AuthorID    sql.NullInt64
	Title       sql.NullString
	Content     sql.NullString
	Reward      sql.NullString
	CareType    sql.NullString
	CarerGender sql.NullString
	RewardType  sql.NullString
	ThumbnailID sql.NullInt64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (q *Queries) WriteSOSPost(ctx context.Context, arg WriteSOSPostParams) (WriteSOSPostRow, error) {
	row := q.db.QueryRowContext(ctx, writeSOSPost,
		arg.AuthorID,
		arg.Title,
		arg.Content,
		arg.Reward,
		arg.CareType,
		arg.CarerGender,
		arg.RewardType,
		arg.ThumbnailID,
	)
	var i WriteSOSPostRow
	err := row.Scan(
		&i.ID,
		&i.AuthorID,
		&i.Title,
		&i.Content,
		&i.Reward,
		&i.CareType,
		&i.CarerGender,
		&i.RewardType,
		&i.ThumbnailID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}