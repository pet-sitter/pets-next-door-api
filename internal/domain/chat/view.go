package chat

import (
	"github.com/google/uuid"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
)

func ToCreateRoom(row databasegen.CreateRoomRow, users *JoinUsersSimpleInfo) *RoomSimpleInfo {
	return &RoomSimpleInfo{
		ID:            row.ID,
		HostInfo:      users,
		RoomName:      row.Name,
		RoomType:      row.RoomType,
		JoinUsersInfo: nil,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}

func ToJoinUsersByFindUserRow(row databasegen.FindUserRow) *JoinUsersSimpleInfo {
	return &JoinUsersSimpleInfo{
		ID:               row.ID,
		UserNickname:     row.Nickname,
		UserProfileImage: row.ProfileImageUrl.String,
	}
}

func ToJoinUsersByFindUserInfoByJoinUserIDRow(row []databasegen.FindUserInfoByJoinUserIdRow) *[]JoinUsersSimpleInfo {
	if len(row) == 0 {
		return nil
	}

	users := make([]JoinUsersSimpleInfo, len(row))
	for i, row := range row {
		users[i] = JoinUsersSimpleInfo{
			ID:               row.JoinUserID,
			UserNickname:     row.JoinUserNickname,
			UserProfileImage: row.ProfileImageUrl.String,
		}
	}
	return &users
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
	joinUserMap map[uuid.UUID]*[]JoinUsersSimpleInfo,
) *JoinRoomsView {
	if len(rows) == 0 {
		// row가 없으면 빈 배열 반환
		return &JoinRoomsView{
			Items: []RoomSimpleInfo{},
		}
	}

	// rows를 반복하며 JoinRoomsView로 변환
	roomSimpleInfos := make([]RoomSimpleInfo, len(rows))
	for i, row := range rows {
		joinUsers := joinUserMap[row.ChatRoomID]

		roomSimpleInfos[i] = RoomSimpleInfo{
			ID: row.ChatRoomID,
			HostInfo: &JoinUsersSimpleInfo{
				ID:               row.HostUserID,
				UserNickname:     row.HostUserNickname,
				UserProfileImage: row.HostUserProfileImageUrl.String,
			},
			RoomName:      row.ChatRoomName,
			RoomType:      row.ChatRoomType,
			JoinUsersInfo: joinUsers,
			CreatedAt:     row.ChatRoomCreatedAt,
			UpdatedAt:     row.ChatRoomUpdatedAt,
		}
	}

	return &JoinRoomsView{
		Items: roomSimpleInfos,
	}
}

func ToUserChatRoomView(
	row databasegen.FindUserChatRoomByIDAndUserIDRow,
	joinUsers []databasegen.FindUserInfoByJoinUserIdRow,
) *RoomSimpleInfo {
	return &RoomSimpleInfo{
		ID: row.ChatRoomID,
		HostInfo: &JoinUsersSimpleInfo{
			ID:               row.HostUserID,
			UserNickname:     row.HostUserNickname,
			UserProfileImage: row.HostUserProfileImageUrl.String,
		},
		RoomName:      row.ChatRoomName,
		RoomType:      row.ChatRoomType,
		JoinUsersInfo: ToJoinUsersByFindUserInfoByJoinUserIDRow(joinUsers),
		CreatedAt:     row.ChatRoomCreatedAt,
		UpdatedAt:     row.ChatRoomUpdatedAt,
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

func ToChatRoomMessage(
	row databasegen.SaveChatMessageRow,
) *Message {
	return &Message{
		ID:          row.ID,
		UserID:      row.UserID,
		RoomID:      row.RoomID,
		MessageType: row.MessageType,
		Content:     row.Content,
		CreatedAt:   row.CreatedAt,
	}
}
