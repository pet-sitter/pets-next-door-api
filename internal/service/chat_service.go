package service

import (
	"context"

	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
)

type ChatService struct {
	conn *database.DB
}

func NewChatService(conn *database.DB) *ChatService {
	return &ChatService{
		conn: conn,
	}
}

// 채팅방 생성
func (s *ChatService) CreateRoom(
	ctx context.Context, name string, roomType chat.RoomType,
) (*chat.Room, error) {
	row, err := databasegen.New(s.conn).CreateRoom(ctx, databasegen.CreateRoomParams{
		Name:     name,
		RoomType: string(roomType),
	})
	if err != nil {
		return nil, err
	}
	return chat.ToCreateRoom(row), nil
}

// 채팅방 입장
func (s *ChatService) JoinRoom(
	ctx context.Context, roomID int64, fbUID string,
) (*chat.JoinRoomView, error) {
	userData, err := databasegen.New(s.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err != nil {
		return nil, err
	}

	row, err2 := databasegen.New(s.conn).JoinRoom(ctx, databasegen.JoinRoomParams{
		RoomID: roomID,
		UserID: int64(userData.ID),
	})

	if err2 != nil {
		return nil, err2
	}

	return chat.ToJoinRoomView(row), nil
}

// 채팅방 떠나기
func (s *ChatService) LeaveRoom(
	ctx context.Context, roomID int64, fbUID string,
) error {
	userData, err := databasegen.New(s.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err != nil {
		return err
	}
	err2 := databasegen.New(s.conn).LeaveRoom(ctx, databasegen.LeaveRoomParams{
		RoomID: roomID,
		UserID: int64(userData.ID),
	})
	if err2 != nil {
		return err2
	}
	return nil
}

// 채팅 메시지 저장
func (s *ChatService) SaveMessage(
	ctx context.Context, roomID int64, fbUID, message string, messageType chat.MessageType,
) (*chat.Message, error) {
	userData, err := databasegen.New(s.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err != nil {
		return nil, err
	}
	row, err := databasegen.New(s.conn).WriteMessage(ctx, databasegen.WriteMessageParams{
		RoomID:      roomID,
		UserID:      int64(userData.ID),
		MessageType: string(messageType),
		Content:     message,
	})
	if err != nil {
		return nil, err
	}
	return chat.ToMessage(row), nil
}

// 채팅방 목록 조회
func (s *ChatService) FindRoomByID(ctx context.Context, roomID *int64) (*chat.Room, error) {
	row, err := databasegen.New(s.conn).FindRoomByID(ctx, utils.Int64PtrToNullInt64(roomID))
	if err != nil {
		return nil, err
	}
	return chat.ToRoom(row), nil
}
