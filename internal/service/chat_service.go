package service

import (
	"context"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/infra/database"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
)

type ChatService struct {
	conn *database.DB
}

func NewChatService(conn *database.DB) *ChatService {
	return &ChatService{
		conn: conn,
	}
}

func (s *ChatService) CreateRoom(ctx context.Context, name, roomType string, joinUserIDs *[]int64) (
	*chat.RoomSimpleInfo, *pnd.AppError,
) {
	// 채팅방 생성
	tx, err := s.conn.BeginTx(ctx)
	defer tx.Rollback()

	if err != nil {
		return nil, err
	}
	q := databasegen.New(tx)

	row, err2 := q.CreateRoom(ctx, databasegen.CreateRoomParams{
		Name:     name,
		RoomType: roomType,
	})

	if err2 != nil {
		return nil, pnd.FromPostgresError(err2)
	}

	// 채팅방에 참여하는 인원이 없을 경우 방만 생성
	if joinUserIDs == nil || len(*joinUserIDs) == 0 {
		return chat.ToCreateRoom(row, nil), nil
	}

	// 채팅방에 참여하는 인원이 있을 경우 참여자 추가
	err3 := q.JoinRooms(ctx, databasegen.JoinRoomsParams{
		RoomID:  int64(row.ID),
		UserIDs: *joinUserIDs,
	})
	if err3 != nil {
		return nil, pnd.FromPostgresError(err3)
	}

	tx.Commit()

	joinUsers, err4 := databasegen.New(s.conn).FindUsersByIds(ctx, *joinUserIDs)

	if err4 != nil {
		return nil, pnd.FromPostgresError(err4)
	}

	return chat.ToCreateRoom(row, chat.ToJoinUsers(joinUsers)), nil
}

func (s *ChatService) JoinRoom(ctx context.Context, roomID int64, fbUID string) (*chat.JoinRoom, *pnd.AppError) {
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

	return chat.ToJoinRoom(row), nil
}

func (s *ChatService) LeaveRoom(ctx context.Context, roomID int64, fbUID string) *pnd.AppError {
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

func (s *ChatService) FindAllByUserUID(ctx context.Context, fbUID string) (*chat.JoinRoomsView, *pnd.AppError) {
	userData, err := databasegen.New(s.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	rows, err := databasegen.New(s.conn).FindAllUserChatRoomsByUserUID(ctx, int64(userData.ID))
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	// rows를 반복하며 각 row에 대해 ToJoinRoom을 호출하여 JoinRoom으로 변환
	return chat.ToUserChatRoomsView(rows), nil
}

func (s *ChatService) FindChatRoomByUIDAndRoomID(ctx context.Context, fbUID string, roomID int64) (
	*chat.RoomSimpleInfo, *pnd.AppError,
) {
	userData, err := databasegen.New(s.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	row, err := databasegen.New(s.conn).FindRoomByIDAndUserID(ctx, roomID, int64(userData.ID))
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return chat.ToUserChatRoomView(row), nil
}

func (s *ChatService) FindChatRoomMessagesByRoomID(ctx context.Context, roomID, prev, next, limit int64) (
	*chat.MessageCursorView, *pnd.AppError,
) {
	hasNext, hasPrev, rows, err := databasegen.New(s.conn).FindMessageByRoomID(ctx, databasegen.FindMessageByRoomIDParams{
		Prev:   prev,
		Next:   next,
		Limit:  limit,
		RoomID: roomID,
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return chat.ToUserChatRoomMessageView(rows, hasNext, hasPrev), nil
}
