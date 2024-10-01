package chat

import (
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

func ToUserChatRoomsView(rows []databasegen.FindAllUserChatRoomsByUserUIDRow) *JoinRoomsView {
	if len(rows) == 0 {
		return nil
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

func ToUserChatRoomMessageView(row []databasegen.FindMessageByRoomIDRow, hasNext, hasPrev *bool) *MessageCursorView {
	if len(row) == 0 {
		return nil
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

	return &MessageCursorView{
		Items:   messages,
		HasNext: hasNext,
		HasPrev: hasPrev,
	}
}
