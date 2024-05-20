// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: chats.sql

package databasegen

import (
	"context"
	"database/sql"
	"time"
)

const createRoom = `-- name: CreateRoom :one
INSERT INTO chat_rooms
(name,
room_type,
created_at,
updated_at)
VALUES ($1, $2, NOW(), NOW())
RETURNING id, name, room_type, created_at, updated_at
`

type CreateRoomParams struct {
	Name     string
	RoomType string
}

type CreateRoomRow struct {
	ID        int64
	Name      string
	RoomType  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreateRoom(ctx context.Context, arg CreateRoomParams) (CreateRoomRow, error) {
	row := q.db.QueryRowContext(ctx, createRoom, arg.Name, arg.RoomType)
	var i CreateRoomRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.RoomType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findMessageByRoomID = `-- name: FindMessageByRoomID :many
SELECT
    id,
    user_id,
    room_id,
    message_type,
    content
FROM
    chat_messages
WHERE
    room_id = $3
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

type FindMessageByRoomIDParams struct {
	Limit  int32
	Offset int32
	RoomID sql.NullInt64
}

type FindMessageByRoomIDRow struct {
	ID          int64
	UserID      int64
	RoomID      int64
	MessageType string
	Content     string
}

func (q *Queries) FindMessageByRoomID(ctx context.Context, arg FindMessageByRoomIDParams) ([]FindMessageByRoomIDRow, error) {
	rows, err := q.db.QueryContext(ctx, findMessageByRoomID, arg.Limit, arg.Offset, arg.RoomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindMessageByRoomIDRow
	for rows.Next() {
		var i FindMessageByRoomIDRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.RoomID,
			&i.MessageType,
			&i.Content,
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

const findRoomByID = `-- name: FindRoomByID :one
SELECT
    id,
    name,
    room_type,
    created_at,
    updated_at
FROM
    chat_rooms
WHERE
    (chat_rooms.id = $1)
    AND (chat_rooms.deleted_at IS NULL)
`

type FindRoomByIDRow struct {
	ID        int64
	Name      string
	RoomType  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) FindRoomByID(ctx context.Context, id sql.NullInt64) (FindRoomByIDRow, error) {
	row := q.db.QueryRowContext(ctx, findRoomByID, id)
	var i FindRoomByIDRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.RoomType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const joinRoom = `-- name: JoinRoom :one
INSERT INTO user_chat_rooms
(user_id, 
room_id,
created_at,
updated_at) 
VALUES ($1, $2, NOW(), NOW())
RETURNING id, user_id, room_id, created_at, updated_at
`

type JoinRoomParams struct {
	UserID int64
	RoomID int64
}

type JoinRoomRow struct {
	ID        int64
	UserID    int64
	RoomID    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) JoinRoom(ctx context.Context, arg JoinRoomParams) (JoinRoomRow, error) {
	row := q.db.QueryRowContext(ctx, joinRoom, arg.UserID, arg.RoomID)
	var i JoinRoomRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RoomID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const leaveRoom = `-- name: LeaveRoom :exec
UPDATE 
    user_chat_rooms
SET deleted_at = NOW()
WHERE user_id = $1
AND room_id = $2
`

type LeaveRoomParams struct {
	UserID int64
	RoomID int64
}

func (q *Queries) LeaveRoom(ctx context.Context, arg LeaveRoomParams) error {
	_, err := q.db.ExecContext(ctx, leaveRoom, arg.UserID, arg.RoomID)
	return err
}

const writeMessage = `-- name: WriteMessage :one
INSERT INTO chat_messages
(user_id,
room_id,
message_type,
content,
created_at,
updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING id, user_id, room_id, message_type, content, created_at, updated_at
`

type WriteMessageParams struct {
	UserID      int64
	RoomID      int64
	MessageType string
	Content     string
}

type WriteMessageRow struct {
	ID          int64
	UserID      int64
	RoomID      int64
	MessageType string
	Content     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (q *Queries) WriteMessage(ctx context.Context, arg WriteMessageParams) (WriteMessageRow, error) {
	row := q.db.QueryRowContext(ctx, writeMessage,
		arg.UserID,
		arg.RoomID,
		arg.MessageType,
		arg.Content,
	)
	var i WriteMessageRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RoomID,
		&i.MessageType,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}