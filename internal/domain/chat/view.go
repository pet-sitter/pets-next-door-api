package chat

import (
	"time"

	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
)

type JoinRoomView struct {
	UserID    int64
	RoomID    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ToJoinRoomView(row databasegen.JoinRoomRow) *JoinRoomView {
	return &JoinRoomView{
		UserID:    row.UserID,
		RoomID:    row.RoomID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func ToCreateRoom(row databasegen.CreateRoomRow) *Room {
	return &Room{
		ID:        row.ID,
		Name:      row.Name,
		RoomType:  RoomType(row.RoomType),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func ToMessage(row databasegen.WriteMessageRow) *Message {
	return &Message{
		ID:        row.ID,
		RoomID:    row.RoomID,
		UserID:    row.UserID,
		Content:   row.Content,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func ToRoom(row databasegen.FindRoomByIDRow) *Room {
	return &Room{
		ID:        row.ID,
		Name:      row.Name,
		RoomType:  RoomType(row.RoomType),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}
