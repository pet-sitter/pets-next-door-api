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
			ID:        r.UserID,
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

func ToUserChatRoomMessageBetweenView(
	row []databasegen.FindBetweenMessagesByRoomIDRow,
	hasNext, hasPrev bool,
	nextMessageID, prevMessageID *uuid.UUID,
) *MessageCursorView {
	if len(row) == 0 {
		return &MessageCursorView{
			HasNext: false,
			HasPrev: false,
			Items:   &[]Message{},
		}
	}

	messages := make([]Message, len(row))
	for i, r := range row {
		messages[i] = Message{
			ID:          r.ID,
			UserID:      r.UserID,
			RoomID:      r.RoomID,
			MessageType: r.MessageType,
			Content:     r.Content,
			CreatedAt:   r.CreatedAt,
		}
	}

	if hasNext && hasPrev {
		return &MessageCursorView{
			Items:   &messages,
			HasNext: hasNext,
			NextID:  nextMessageID,
			HasPrev: hasPrev,
			PrevID:  prevMessageID,
		}
	}

	if hasNext && !hasPrev {
		return &MessageCursorView{
			Items:   &messages,
			HasNext: hasNext,
			HasPrev: hasPrev,
			NextID:  nextMessageID,
			PrevID:  nil,
		}
	}

	if hasPrev && !hasNext {
		return &MessageCursorView{
			Items:   &messages,
			HasNext: hasNext,
			NextID:  nil,
			HasPrev: hasPrev,
			PrevID:  prevMessageID,
		}
	}

	return &MessageCursorView{
		Items:   &messages,
		HasNext: hasNext,
		HasPrev: hasPrev,
		NextID:  nil,
		PrevID:  nil,
	}
}

func ToUserChatRoomMessagePrevView(
	row []databasegen.FindPrevMessageByRoomIDRow, hasNext, hasPrev bool,
	nextMessageID, prevMessageID *uuid.UUID,
) *MessageCursorView {
	if len(row) == 0 {
		return &MessageCursorView{
			HasNext: false,
			HasPrev: false,
			Items:   &[]Message{},
		}
	}

	messages := make([]Message, len(row))
	for i, r := range row {
		messages[i] = Message{
			ID:          r.ID,
			UserID:      r.UserID,
			RoomID:      r.RoomID,
			MessageType: r.MessageType,
			Content:     r.Content,
			CreatedAt:   r.CreatedAt,
		}
	}

	if hasNext && hasPrev {
		return &MessageCursorView{
			Items:   &messages,
			HasNext: hasNext,
			NextID:  nextMessageID,
			HasPrev: hasPrev,
			PrevID:  prevMessageID,
		}
	}

	if hasNext && !hasPrev {
		return &MessageCursorView{
			Items:   &messages,
			HasNext: hasNext,
			HasPrev: hasPrev,
			NextID:  nextMessageID,
			PrevID:  nil,
		}
	}

	if hasPrev && !hasNext {
		return &MessageCursorView{
			Items:   &messages,
			HasNext: hasNext,
			NextID:  nil,
			HasPrev: hasPrev,
			PrevID:  prevMessageID,
		}
	}

	return &MessageCursorView{
		Items:   &messages,
		HasNext: hasNext,
		HasPrev: hasPrev,
		NextID:  nil,
		PrevID:  nil,
	}
}

func ToUserChatRoomMessageNextView(
	row []databasegen.FindNextMessageByRoomIDRow, hasNext, hasPrev bool,
	nextMessageID, prevMessageID *uuid.UUID,
) *MessageCursorView {
	if len(row) == 0 {
		return &MessageCursorView{
			HasNext: false,
			HasPrev: false,
			Items:   &[]Message{},
		}
	}

	messages := make([]Message, len(row))
	for i, r := range row {
		messages[i] = Message{
			ID:          r.ID,
			UserID:      r.UserID,
			RoomID:      r.RoomID,
			MessageType: r.MessageType,
			Content:     r.Content,
			CreatedAt:   r.CreatedAt,
		}
	}

	if hasNext && hasPrev {
		return &MessageCursorView{
			Items:   &messages,
			HasNext: hasNext,
			NextID:  nextMessageID,
			HasPrev: hasPrev,
			PrevID:  prevMessageID,
		}
	}

	if hasNext && !hasPrev {
		return &MessageCursorView{
			Items:   &messages,
			HasNext: hasNext,
			HasPrev: hasPrev,
			NextID:  nextMessageID,
			PrevID:  nil,
		}
	}

	if hasPrev && !hasNext {
		return &MessageCursorView{
			Items:   &messages,
			HasNext: hasNext,
			NextID:  nil,
			HasPrev: hasPrev,
			PrevID:  prevMessageID,
		}
	}

	return &MessageCursorView{
		Items:   &messages,
		HasNext: hasNext,
		HasPrev: hasPrev,
		NextID:  nil,
		PrevID:  nil,
	}
}

func ToUserChatRoomMessageView(
	row []databasegen.FindMessagesByRoomIDAndSizeRow, hasNext, hasPrev bool,
	nextMessageID, prevMessageID *uuid.UUID,
) *MessageCursorView {
	if len(row) == 0 {
		return &MessageCursorView{
			HasNext: false,
			HasPrev: false,
			Items:   &[]Message{},
		}
	}

	messages := make([]Message, len(row))
	for i, r := range row {
		messages[i] = Message{
			ID:          r.ID,
			UserID:      r.UserID,
			RoomID:      r.RoomID,
			MessageType: r.MessageType,
			Content:     r.Content,
			CreatedAt:   r.CreatedAt,
		}
	}

	if hasNext && hasPrev {
		return &MessageCursorView{
			Items:   &messages,
			HasNext: hasNext,
			NextID:  nextMessageID,
			HasPrev: hasPrev,
			PrevID:  prevMessageID,
		}
	}

	if hasNext && !hasPrev {
		return &MessageCursorView{
			Items:   &messages,
			HasNext: hasNext,
			HasPrev: hasPrev,
			NextID:  nextMessageID,
			PrevID:  nil,
		}
	}

	if hasPrev && !hasNext {
		return &MessageCursorView{
			Items:   &messages,
			HasNext: hasNext,
			HasPrev: hasPrev,
			PrevID:  prevMessageID,
			NextID:  nil,
		}
	}

	return &MessageCursorView{
		Items:   &messages,
		HasNext: hasNext,
		HasPrev: hasPrev,
		PrevID:  nil,
		NextID:  nil,
	}
}
