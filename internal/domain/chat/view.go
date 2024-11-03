package chat

import (
	"github.com/google/uuid"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
)

func ToCreateRoom(row databasegen.CreateRoomRow, users *JoinUsersSimpleInfo) *RoomSimpleInfo {
	return &RoomSimpleInfo{
		ID:        row.ID,
		RoomName:  row.Name,
		RoomType:  row.RoomType,
		JoinUser:  users,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func ToJoinUsers(row databasegen.FindUserRow) *JoinUsersSimpleInfo {
	return &JoinUsersSimpleInfo{
		ID:               row.ID,
		UserNickname:     row.Nickname,
		UserProfileImage: row.ProfileImageUrl.String,
	}
}

func ToJoinRoom(row databasegen.JoinRoomRow) *JoinRoom {
	return &JoinRoom{
		UserID:   row.UserID,
		RoomID:   row.RoomID,
		JoinedAt: row.JoinedAt,
	}
}

func ToUserChatRoomsView(
	rows []databasegen.FindAllUserChatRoomsByUserUIDRow,
) *JoinRoomsView {
	if len(rows) == 0 {
		// row가 없으면 빈 배열 반환
		return &JoinRoomsView{
			Items: []RoomSimpleInfo{},
		}
	}

	// rows를 반복하며 JoinRoomsView로 변환
	roomSimpleInfos := make([]RoomSimpleInfo, len(rows))
	for i, r := range rows {
		roomSimpleInfos[i] = RoomSimpleInfo{
			ID:        r.ChatRoomID,
			RoomName:  r.ChatRoomName,
			RoomType:  r.ChatRoomType,
			CreatedAt: r.ChatRoomCreatedAt,
			UpdatedAt: r.ChatRoomUpdatedAt,
		}
	}

	return &JoinRoomsView{
		Items: roomSimpleInfos,
	}
}

func ToUserChatRoomView(row databasegen.FindRoomByIDAndUserIDRow) *RoomSimpleInfo {
	return &RoomSimpleInfo{
		ID:        row.ID,
		RoomName:  row.Name,
		RoomType:  row.RoomType,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func createMessageCursorView(
	row interface{},
	hasNext, hasPrev bool,
	nextMessageID, prevMessageID *uuid.UUID,
) *MessageCursorView {
	var messages []Message

	switch v := row.(type) {
	case []databasegen.FindBetweenMessagesByRoomIDRow:
		messages = make([]Message, len(v))
		for i, r := range v {
			messages[i] = Message{
				ID:          r.ID,
				UserID:      r.UserID,
				RoomID:      r.RoomID,
				MessageType: r.MessageType,
				Content:     r.Content,
				CreatedAt:   r.CreatedAt,
			}
		}
	case []databasegen.FindPrevMessageByRoomIDRow:
		messages = make([]Message, len(v))
		for i, r := range v {
			messages[i] = Message{
				ID:          r.ID,
				UserID:      r.UserID,
				RoomID:      r.RoomID,
				MessageType: r.MessageType,
				Content:     r.Content,
				CreatedAt:   r.CreatedAt,
			}
		}
	case []databasegen.FindNextMessageByRoomIDRow:
		messages = make([]Message, len(v))
		for i, r := range v {
			messages[i] = Message{
				ID:          r.ID,
				UserID:      r.UserID,
				RoomID:      r.RoomID,
				MessageType: r.MessageType,
				Content:     r.Content,
				CreatedAt:   r.CreatedAt,
			}
		}
	case []databasegen.FindMessagesByRoomIDAndSizeRow:
		messages = make([]Message, len(v))
		for i, r := range v {
			messages[i] = Message{
				ID:          r.ID,
				UserID:      r.UserID,
				RoomID:      r.RoomID,
				MessageType: r.MessageType,
				Content:     r.Content,
				CreatedAt:   r.CreatedAt,
			}
		}
	default:
		return &MessageCursorView{
			HasNext: false,
			HasPrev: false,
			Items:   &[]Message{},
		}
	}

	nextID := nextMessageID
	prevID := prevMessageID
	if !hasNext {
		nextID = nil
	}
	if !hasPrev {
		prevID = nil
	}

	return &MessageCursorView{
		Items:   &messages,
		HasNext: hasNext,
		HasPrev: hasPrev,
		NextID:  nextID,
		PrevID:  prevID,
	}
}

func ToUserChatRoomMessageBetweenView(
	row []databasegen.FindBetweenMessagesByRoomIDRow,
	hasNext, hasPrev bool,
	nextMessageID, prevMessageID *uuid.UUID,
) *MessageCursorView {
	return createMessageCursorView(row, hasNext, hasPrev, nextMessageID, prevMessageID)
}

func ToUserChatRoomMessagePrevView(
	row []databasegen.FindPrevMessageByRoomIDRow, hasNext, hasPrev bool,
	nextMessageID, prevMessageID *uuid.UUID,
) *MessageCursorView {
	return createMessageCursorView(row, hasNext, hasPrev, nextMessageID, prevMessageID)
}

func ToUserChatRoomMessageNextView(
	row []databasegen.FindNextMessageByRoomIDRow, hasNext, hasPrev bool,
	nextMessageID, prevMessageID *uuid.UUID,
) *MessageCursorView {
	return createMessageCursorView(row, hasNext, hasPrev, nextMessageID, prevMessageID)
}

func ToUserChatRoomMessageView(
	row []databasegen.FindMessagesByRoomIDAndSizeRow, hasNext, hasPrev bool,
	nextMessageID, prevMessageID *uuid.UUID,
) *MessageCursorView {
	return createMessageCursorView(row, hasNext, hasPrev, nextMessageID, prevMessageID)
}
