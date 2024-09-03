package service

import (
	"context"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
	"time"
)

type ChatService struct {
	conn *database.DB
}

func NewChatService(conn *database.DB) *ChatService {
	return &ChatService{
		conn: conn,
	}
}

func (s *ChatService) CreateRoom(
	ctx context.Context, name string, roomType chat.RoomType,
) (*chat.Room, *pnd.AppError) {
	row, err := databasegen.New(s.conn).CreateRoom(ctx, databasegen.CreateRoomParams{
		Name:     name,
		RoomType: string(roomType),
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	return chat.ToCreateRoom(row), nil
}

func (s *ChatService) JoinRoom(
	ctx context.Context, roomID int64, fbUID string,
) (*chat.JoinRoomView, *pnd.AppError) {
	userData, err := databasegen.New(s.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	row, err := databasegen.New(s.conn).JoinRoom(ctx, databasegen.JoinRoomParams{
		RoomID: roomID,
		UserID: int64(userData.ID),
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return chat.ToJoinRoomView(row), nil
}

func (s *ChatService) LeaveRoom(
	ctx context.Context, roomID int64, fbUID string,
) *pnd.AppError {
	userData, err := databasegen.New(s.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err != nil {
		return pnd.FromPostgresError(err)
	}
	err = databasegen.New(s.conn).LeaveRoom(ctx, databasegen.LeaveRoomParams{
		RoomID: roomID,
		UserID: int64(userData.ID),
	})
	if err != nil {
		return pnd.FromPostgresError(err)
	}
	exists, err := databasegen.New(s.conn).UserExistsInRoom(ctx, roomID)
	if err != nil {
		return pnd.FromPostgresError(err)
	}

	if !exists {
		err = databasegen.New(s.conn).DeleteRoom(ctx, int32(roomID))
		if err != nil {
			return pnd.FromPostgresError(err)
		}
	}
	return nil
}

func (s *ChatService) SaveMessage(
	ctx context.Context, roomID int64, fbUID, message string, messageType chat.MessageType,
) (*chat.Message, *pnd.AppError) {
	userData, err := databasegen.New(s.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	row, err := databasegen.New(s.conn).WriteMessage(ctx, databasegen.WriteMessageParams{
		RoomID:      roomID,
		UserID:      int64(userData.ID),
		MessageType: string(messageType),
		Content:     message,
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	return chat.ToMessage(row), nil
}

func (s *ChatService) FindRoomByID(ctx context.Context, roomID int64) (*chat.Room, *pnd.AppError) {
	row, err := databasegen.New(s.conn).FindRoomByID(ctx, utils.Int64ToNullInt32(roomID))
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	return chat.ToRoom(row), nil
}

func (s *ChatService) FindUserChatRoom(ctx context.Context) (chat.UserChatRoomViewList, *pnd.AppError) {
	rows, err := databasegen.New(s.conn).FindUserChatRooms(ctx)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	return chat.ToUserChatRoomFromRows(rows), nil
}

func (s *ChatService) ExistsUserInRoom(ctx context.Context, roomID int64, fbUID string) (bool, *pnd.AppError) {
	userData, err := databasegen.New(s.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err != nil {
		return false, pnd.FromPostgresError(err)
	}
	exists, err := databasegen.New(s.conn).ExistsUserInRoom(ctx, databasegen.ExistsUserInRoomParams{
		RoomID: roomID,
		UserID: int64(userData.ID),
	})
	if err != nil {
		return false, pnd.FromPostgresError(err)
	}
	return exists, nil
}

func (s *ChatService) MockFindAllChatRooms(ctx context.Context) (*[]chat.Room, *pnd.AppError) {

	rooms := []chat.Room{
		{
			ID:        1,
			Name:      "Room 1",
			RoomType:  chat.RoomTypePersonal,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Now(),
		},
		{
			ID:        2,
			Name:      "Room 2",
			RoomType:  chat.RoomTypeGathering,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Now(),
		},
	}

	return &rooms, nil
}

func (s *ChatService) MockFindMessagesByRoomID(ctx context.Context, roomID int64) (*[]chat.Message, *pnd.AppError) {
	messages := []chat.Message{
		{
			ID:          1,
			RoomID:      roomID,
			UserID:      1,
			Content:     "Hello",
			MessageType: "normal",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			RoomID:      roomID,
			UserID:      2,
			Content:     "Hi",
			MessageType: "normal",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	return &messages, nil
}
