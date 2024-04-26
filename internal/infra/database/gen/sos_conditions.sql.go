// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: sos_conditions.sql

package databasegen

import (
	"context"
	"database/sql"
	"time"
)

const createSOSCondition = `-- name: CreateSOSCondition :one
INSERT INTO sos_conditions
(id,
 name,
 created_at,
 updated_at)
SELECT $1, $2, now(), now()
WHERE NOT EXISTS (SELECT 1
                  FROM sos_conditions
                  WHERE name = $2::VARCHAR(50))
RETURNING id, name, created_at, updated_at, deleted_at
`

type CreateSOSConditionParams struct {
	ID   int32
	Name sql.NullString
}

func (q *Queries) CreateSOSCondition(ctx context.Context, arg CreateSOSConditionParams) (SosCondition, error) {
	row := q.db.QueryRowContext(ctx, createSOSCondition, arg.ID, arg.Name)
	var i SosCondition
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const findConditions = `-- name: FindConditions :many
SELECT id,
       name,
       created_at,
       updated_at,
       deleted_at
FROM sos_conditions
WHERE ($1::BOOLEAN = TRUE OR
       ($1::BOOLEAN = FALSE AND deleted_at IS NULL))
`

func (q *Queries) FindConditions(ctx context.Context, includeDeleted bool) ([]SosCondition, error) {
	rows, err := q.db.QueryContext(ctx, findConditions, includeDeleted)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SosCondition
	for rows.Next() {
		var i SosCondition
		if err := rows.Scan(
			&i.ID,
			&i.Name,
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

const findSOSPostConditions = `-- name: FindSOSPostConditions :many
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
  AND ($2::BOOLEAN = TRUE OR
       ($2::BOOLEAN = FALSE AND sos_posts_conditions.deleted_at IS NULL))
`

type FindSOSPostConditionsParams struct {
	SosPostID      sql.NullInt64
	IncludeDeleted bool
}

type FindSOSPostConditionsRow struct {
	ID        int32
	Name      sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) FindSOSPostConditions(ctx context.Context, arg FindSOSPostConditionsParams) ([]FindSOSPostConditionsRow, error) {
	rows, err := q.db.QueryContext(ctx, findSOSPostConditions, arg.SosPostID, arg.IncludeDeleted)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindSOSPostConditionsRow
	for rows.Next() {
		var i FindSOSPostConditionsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
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
