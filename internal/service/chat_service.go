package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

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

func (s *ChatService) CreateRoom(
	ctx context.Context, name, roomType, userFirebaseUID string,
) (
	*chat.RoomSimpleInfo, *pnd.AppError,
) {
	userData, err := databasegen.New(s.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(userFirebaseUID),
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	// 채팅방 생성
	tx, transactionError := s.conn.BeginTx(ctx)
	defer tx.Rollback()

	if transactionError != nil {
		return nil, transactionError
	}

	q := databasegen.New(tx)
	row, databaseGenError := q.CreateRoom(ctx, databasegen.CreateRoomParams{
		Name:     name,
		RoomType: roomType,
	})

	if databaseGenError != nil {
		return nil, pnd.FromPostgresError(databaseGenError)
	}

	_, err3 := q.JoinRoom(ctx, databasegen.JoinRoomParams{
		UserID: userData.ID,
		RoomID: row.ID,
	})

	if err3 != nil {
		return nil, pnd.FromPostgresError(err3)
	}

	tx.Commit()

	return chat.ToCreateRoom(row, chat.ToJoinUsers(userData)), nil
}

func (s *ChatService) JoinRoom(ctx context.Context, roomID uuid.UUID, fbUID string) (*chat.JoinRoom, *pnd.AppError) {
	userData, err := databasegen.New(s.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}
	// 채팅방에 이미 참여중인지 확인
	exists, err := databasegen.New(s.conn).UserExistsInRoom(ctx, roomID)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	if !exists {
		row, err := databasegen.New(s.conn).JoinRoom(ctx, databasegen.JoinRoomParams{
			RoomID: roomID,
			UserID: userData.ID,
		})
		if err != nil {
			return nil, pnd.FromPostgresError(err)
		}

		return chat.ToJoinRoom(row), nil
	}

	return nil, pnd.ErrBadRequest(errors.New("이미 참여중인 채팅방입니다"))
}

func (s *ChatService) LeaveRoom(ctx context.Context, roomID uuid.UUID, fbUID string) *pnd.AppError {
	userData, err := databasegen.New(s.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err != nil {
		return pnd.FromPostgresError(err)
	}

	err = databasegen.New(s.conn).LeaveRoom(ctx, databasegen.LeaveRoomParams{
		RoomID: roomID,
		UserID: userData.ID,
	})
	if err != nil {
		return pnd.FromPostgresError(err)
	}
	exists, err := databasegen.New(s.conn).UserExistsInRoom(ctx, roomID)
	if err != nil {
		return pnd.FromPostgresError(err)
	}

	if !exists {
		err = databasegen.New(s.conn).DeleteRoom(ctx, roomID)
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
	rows, err := databasegen.New(s.conn).FindAllUserChatRoomsByUserUID(ctx, userData.ID)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	// rows를 반복하며 각 row에 대해 ToJoinRoom을 호출하여 JoinRoom으로 변환
	return chat.ToUserChatRoomsView(rows), nil
}

func (s *ChatService) FindChatRoomByUIDAndRoomID(ctx context.Context, fbUID string, roomID uuid.UUID) (
	*chat.RoomSimpleInfo, *pnd.AppError,
) {
	userData, err := databasegen.New(s.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	row, err := databasegen.New(s.conn).FindRoomByIDAndUserID(
		ctx,
		databasegen.FindRoomByIDAndUserIDParams{ID: roomID, UserID: userData.ID},
	)
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	return chat.ToUserChatRoomView(row), nil
}

func (s *ChatService) FindChatRoomMessagesByRoomID(
	ctx context.Context, roomID uuid.UUID, prev, next uuid.NullUUID, limit int64,
) (*chat.MessageCursorView, *pnd.AppError) {
	rows, err := databasegen.New(s.conn).FindMessageByRoomID(ctx, databasegen.FindMessageByRoomIDParams{
		Prev:   prev,
		Next:   next,
		Limit:  int32(limit),
		RoomID: roomID,
	})
	if err != nil {
		return nil, pnd.FromPostgresError(err)
	}

	// TODO: Implement hasNext, hasPrev
	// return chat.ToUserChatRoomMessageView(rows, hasNext, hasPrev), nil
	return chat.ToUserChatRoomMessageView(rows, nil, nil), nil
}
