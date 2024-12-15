package service

import (
	"context"
	"errors"
	"fmt"

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
	*chat.RoomSimpleInfo, error,
) {
	userData, err := databasegen.New(s.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(userFirebaseUID),
	})
	if err != nil {
		return nil, err
	}

	// 채팅방 생성
	tx, transactionError := s.conn.BeginTx(ctx)
	defer tx.Rollback()

	if transactionError != nil {
		return nil, transactionError
	}

	chatRoomUUID, uuidError := uuid.NewV7()
	if uuidError != nil {
		return nil, pnd.ErrUnknown(fmt.Errorf("failed to generate UUID: %w", uuidError))
	}

	q := databasegen.New(tx)
	row, err := q.CreateRoom(ctx, databasegen.CreateRoomParams{
		ID:       chatRoomUUID,
		Name:     name,
		RoomType: roomType,
	})
	if err != nil {
		return nil, err
	}

	joinRoomUUID, joinRoomUUIDError := uuid.NewV7()
	if joinRoomUUIDError != nil {
		return nil, pnd.ErrUnknown(fmt.Errorf("failed to generate UUID: %w", joinRoomUUIDError))
	}

	_, err = q.JoinRoom(ctx, databasegen.JoinRoomParams{
		ID:     joinRoomUUID,
		UserID: userData.ID,
		RoomID: row.ID,
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return chat.ToCreateRoom(row, chat.ToJoinUsers(userData)), nil
}

func (s *ChatService) JoinRoom(
	ctx context.Context,
	roomID uuid.UUID,
	fbUID string,
) (*chat.JoinRoom, error) {
	userData, err := databasegen.New(s.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err != nil {
		return nil, err
	}
	// 채팅방이 현재 존재하는지 확인
	existRoom, err := databasegen.New(s.conn).ExistsRoom(ctx, roomID)
	if err != nil {
		return nil, err
	}

	if !existRoom {
		return nil, pnd.ErrBadRequest(errors.New("chat room does not exist"))
	}

	// 채팅방에 이미 참여중인지 확인
	existsUser, existsUserError := databasegen.New(s.conn).UserExistsInRoom(ctx, roomID)
	if existsUserError != nil {
		return nil, err
	}

	// 채팅방에 참여하지 않은 경우 참여
	if !existsUser {
		joinRoomUUID, joinRoomUUIDError := uuid.NewV7()

		if joinRoomUUIDError != nil {
			return nil, pnd.ErrUnknown(fmt.Errorf("failed to generate UUID: %w", joinRoomUUIDError))
		}

		row, err := databasegen.New(s.conn).JoinRoom(ctx, databasegen.JoinRoomParams{
			ID:     joinRoomUUID,
			RoomID: roomID,
			UserID: userData.ID,
		})
		if err != nil {
			return nil, err
		}

		return chat.ToJoinRoom(row), nil
	}

	return nil, pnd.ErrBadRequest(errors.New("user already joined in the room"))
}

func (s *ChatService) LeaveRoom(ctx context.Context, roomID uuid.UUID, fbUID string) error {
	userData, err := databasegen.New(s.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err != nil {
		return err
	}

	err = databasegen.New(s.conn).LeaveRoom(ctx, databasegen.LeaveRoomParams{
		RoomID: roomID,
		UserID: userData.ID,
	})
	if err != nil {
		return err
	}
	exists, err := databasegen.New(s.conn).UserExistsInRoom(ctx, roomID)
	if err != nil {
		return err
	}

	if !exists {
		err = databasegen.New(s.conn).DeleteRoom(ctx, roomID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ChatService) FindAllByUserUID(
	ctx context.Context,
	fbUID string,
) (*chat.JoinRoomsView, error) {
	userData, err := databasegen.New(s.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err != nil {
		return nil, err
	}
	rows, err := databasegen.New(s.conn).FindAllUserChatRoomsByUserUID(ctx, userData.ID)
	if err != nil {
		return nil, err
	}

	// rows를 반복하며 각 row에 대해 ToJoinRoom을 호출하여 JoinRoom으로 변환
	return chat.ToUserChatRoomsView(rows), nil
}

func (s *ChatService) FindChatRoomByUIDAndRoomID(
	ctx context.Context,
	fbUID string,
	roomID uuid.UUID,
) (
	*chat.RoomSimpleInfo, error,
) {
	userData, err := databasegen.New(s.conn).FindUser(ctx, databasegen.FindUserParams{
		FbUid: utils.StrToNullStr(fbUID),
	})
	if err != nil {
		return nil, err
	}

	row, err := databasegen.New(s.conn).FindRoomByIDAndUserID(
		ctx,
		databasegen.FindRoomByIDAndUserIDParams{ID: roomID, UserID: userData.ID},
	)
	if err != nil {
		return nil, err
	}

	return chat.ToUserChatRoomView(row), nil
}

/**
 * 채팅방의 메시지를 조회한다. 채팅메시지는 최신순으로 DESC 정렬을 진행한다.
 * prev - 이전 메시지의 ID
 * next - 다음 메시지의 ID
 */
func (s *ChatService) FindChatRoomMessagesByRoomID(
	ctx context.Context, roomID uuid.UUID, prev, next uuid.NullUUID, limit int64,
) (*chat.MessageCursorView, error) {
	// prev와 next에 따라 다른 쿼리를 실행
	if prev.Valid && next.Valid {
		// prev와 next 모두 존재하는 경우
		rows, err := databasegen.New(s.conn).
			FindBetweenMessagesByRoomID(ctx, databasegen.FindBetweenMessagesByRoomIDParams{
				Prev:   prev,
				Next:   next,
				Limit:  int32(limit),
				RoomID: roomID,
			})
		if err != nil {
			return nil, err
		}

		// rows의 맨 앞의 ID값을 가져온다.
		// 이 값이 prev가 존재하는지 여부를 판단하는데 사용된다.
		if len(rows) == 0 {
			return chat.ToUserChatRoomMessageBetweenView(rows, false, false, nil, nil), nil
		}

		// 가장 최신 메시지의 ID를 가져온다.
		firstID := rows[0].ID

		// 가장 오래된 메시지의 ID를 가져온다.
		lastID := rows[len(rows)-1].ID

		hasPrev, err := databasegen.New(s.conn).
			HasPrevMessages(ctx, databasegen.HasPrevMessagesParams{
				ID:     lastID,
				RoomID: roomID,
			})
		if err != nil {
			return nil, err
		}

		hasNext, err := databasegen.New(s.conn).
			HasNextMessages(ctx, databasegen.HasNextMessagesParams{
				ID:     firstID,
				RoomID: roomID,
			})
		if err != nil {
			return nil, err
		}

		return chat.ToUserChatRoomMessageBetweenView(rows, hasNext, hasPrev, &firstID, &lastID), nil
	}

	if prev.Valid {
		// prev만 존재하는 경우
		rows, err := databasegen.New(s.conn).
			FindPrevMessageByRoomID(ctx, databasegen.FindPrevMessageByRoomIDParams{
				Prev:   prev,
				Limit:  int32(limit),
				RoomID: roomID,
			})
		if err != nil {
			return nil, err
		}

		if len(rows) == 0 {
			return chat.ToUserChatRoomMessagePrevView(rows, false, false, nil, nil), nil
		}

		firstID := rows[0].ID
		lastID := rows[len(rows)-1].ID

		hasPrev, hasPrevError := s.HasPrevMessages(ctx, roomID, lastID)

		if hasPrevError != nil {
			return nil, hasPrevError
		}

		hasNext, hasNextError := s.HasNextMessages(ctx, roomID, firstID)
		if hasNextError != nil {
			return nil, hasNextError
		}

		return chat.ToUserChatRoomMessagePrevView(rows, hasNext, hasPrev, &firstID, &lastID), nil
	}

	if next.Valid {
		// next만 존재하는 경우
		rows, err := databasegen.New(s.conn).
			FindNextMessageByRoomID(ctx, databasegen.FindNextMessageByRoomIDParams{
				Next:   next,
				Limit:  int32(limit),
				RoomID: roomID,
			})
		if err != nil {
			return nil, err
		}

		if len(rows) == 0 {
			return chat.ToUserChatRoomMessageNextView(rows, false, false, nil, nil), nil
		}

		firstID := rows[0].ID
		lastID := rows[len(rows)-1].ID

		hasPrev, hasPrevError := s.HasPrevMessages(ctx, roomID, lastID)

		if hasPrevError != nil {
			return nil, hasPrevError
		}

		hasNext, hasNextError := s.HasNextMessages(ctx, roomID, firstID)
		if hasNextError != nil {
			return nil, hasNextError
		}

		return chat.ToUserChatRoomMessageNextView(rows, hasNext, hasPrev, &firstID, &lastID), nil
	}

	// prev와 next가 모두 없는 경우 Size만큼 최신 메시지를 가져온다.
	rows, err := databasegen.New(s.conn).
		FindMessagesByRoomIDAndSize(ctx, databasegen.FindMessagesByRoomIDAndSizeParams{
			Limit:  int32(limit),
			RoomID: roomID,
		})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return chat.ToUserChatRoomMessageView(rows, false, false, nil, nil), nil
	}

	firstID := rows[0].ID
	lastID := rows[len(rows)-1].ID

	hasPrev, hasPrevError := s.HasPrevMessages(ctx, roomID, lastID)

	if hasPrevError != nil {
		return nil, hasPrevError
	}

	hasNext, hasNextError := s.HasNextMessages(ctx, roomID, firstID)
	if hasNextError != nil {
		return nil, hasNextError
	}

	return chat.ToUserChatRoomMessageView(rows, hasNext, hasPrev, &firstID, &lastID), nil
}

// hasPrev 메시지가 있는지 확인
func (s *ChatService) HasPrevMessages(
	ctx context.Context, roomID, messageID uuid.UUID,
) (bool, error) {
	hasPrev, err := databasegen.New(s.conn).HasPrevMessages(ctx, databasegen.HasPrevMessagesParams{
		ID:     messageID,
		RoomID: roomID,
	})
	if err != nil {
		return false, err
	}

	return hasPrev, nil
}

// hasNext 메시지가 있는지 확인
func (s *ChatService) HasNextMessages(
	ctx context.Context, roomID, messageID uuid.UUID,
) (bool, error) {
	hasNext, err := databasegen.New(s.conn).HasNextMessages(ctx, databasegen.HasNextMessagesParams{
		ID:     messageID,
		RoomID: roomID,
	})
	if err != nil {
		return false, err
	}

	return hasNext, nil
}

// 채팅 메시지를 저장합니다.
func (s *ChatService) SaveChatMessage(
	ctx context.Context, userID, roomID uuid.UUID, messageType, content string,
) (*chat.Message, *pnd.AppError) {
	chatMessageID, uuidError := uuid.NewV7()
	if uuidError != nil {
		return nil, pnd.ErrUnknown(fmt.Errorf("failed to generate UUID: %w", uuidError))
	}

	tx, transactionError := s.conn.BeginTx(ctx)
	defer tx.Rollback()

	if transactionError != nil {
		return nil, pnd.FromPostgresError(transactionError.Err)
	}

	q := databasegen.New(tx)
	row, databaseGenError := q.SaveChatMessage(ctx, databasegen.SaveChatMessageParams{
		ID:          chatMessageID,
		UserID:      userID,
		RoomID:      roomID,
		MessageType: messageType,
		Content:     content,
	})

	if databaseGenError != nil {
		return nil, pnd.FromPostgresError(databaseGenError)
	}

	tx.Commit()

	return chat.ToChatRoomMessage(row), nil
}
