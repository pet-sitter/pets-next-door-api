package chat

import (
	"time"

	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
)

type JoinRoomView struct {
	UserID   int64
	RoomID   int64
	JoinedAt time.Time
}

type UserChatRoomView struct {
	ID       int64
	UserID   int64
	RoomID   int64
	JoinedAt time.Time
	UserInfo user.WithProfileImage
	RoomInfo *Room
}

type UserChatRoomViewList []*UserChatRoomView

func ToJoinRoomView(row databasegen.JoinRoomRow) *JoinRoomView {
	return &JoinRoomView{
		UserID:   row.UserID,
		RoomID:   row.RoomID,
		JoinedAt: row.JoinedAt,
	}
}

func ToCreateRoom(row databasegen.CreateRoomRow) *Room {
	return &Room{
		ID:        int64(row.ID),
		Name:      row.Name,
		RoomType:  RoomType(row.RoomType),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func ToMessage(row databasegen.WriteMessageRow) *Message {
	return &Message{
		ID:        int64(row.ID),
		RoomID:    row.RoomID,
		UserID:    row.UserID,
		Content:   row.Content,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func ToRoom(row databasegen.FindRoomByIDRow) *Room {
	return &Room{
		ID:        int64(row.ID),
		Name:      row.Name,
		RoomType:  RoomType(row.RoomType),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func ToUserChatRoom(row databasegen.FindUserChatRoomsRow) *UserChatRoomView {
	return &UserChatRoomView{
		ID:       int64(row.ID),
		UserID:   row.UserID,
		RoomID:   row.RoomID,
		JoinedAt: row.JoinedAt,
		UserInfo: user.WithProfileImage{
			ID:                   int64(row.ID),
			Email:                row.Email,
			Nickname:             row.Nickname,
			Fullname:             row.Fullname,
			ProfileImageURL:      utils.NullStrToStrPtr(row.ProfileImageUrl),
			FirebaseProviderType: user.FirebaseProviderType(row.FbProviderType.String),
			FirebaseUID:          row.FbUid.String,
			CreatedAt:            row.CreatedAt,
			UpdatedAt:            row.UpdatedAt,
		},
		RoomInfo: &Room{
			ID:        row.RoomID,
			Name:      row.ChatRoomName,
			RoomType:  RoomType(row.ChatRoomType),
			CreatedAt: row.ChatRoomCreatedAt,
			UpdatedAt: row.ChatRoomUpdatedAt,
		},
	}
}

func ToUserChatRoomFromRows(rows []databasegen.FindUserChatRoomsRow) UserChatRoomViewList {
	userChatRooms := make([]*UserChatRoomView, len(rows))
	for i, row := range rows {
		userChatRooms[i] = ToUserChatRoom(row)
	}
	return userChatRooms
}
