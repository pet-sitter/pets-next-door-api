package chat

import (
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
	"strconv"
)

func ToCreateRoom(row databasegen.CreateRoomRow, users *[]JoinUsersSimpleInfo) *RoomSimpleInfo {
	return &RoomSimpleInfo{
		ID:        string(row.ID),
		RoomName:  row.Name,
		RoomType:  row.RoomType,
		JoinUsers: users,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func ToJoinUsers(row []databasegen.FindUsersRow) *[]JoinUsersSimpleInfo {
	if len(row) == 0 {
		return nil
	}

	joinUsers := make([]JoinUsersSimpleInfo, len(row))

	for i, r := range row {
		joinUsers[i] = JoinUsersSimpleInfo{
			ID:               string(r.ID),
			UserNickname:     r.Nickname,
			UserProfileImage: r.ProfileImageUrl,
		}
	}
	return &joinUsers
}

func ToJoinRoom(row databasegen.JoinRoomRow) *JoinRoom {
	return &JoinRoom{
		UserID:   strconv.FormatInt(row.UserID, 10),
		RoomID:   strconv.FormatInt(row.RoomID, 10),
		JoinedAt: row.JoinedAt,
	}
}

func ToUserChatRoomsView(rows []databasegen.FindAllUserChatRoomsRow) *JoinRoomsView {
	if len(rows) == 0 {
		return nil
	}

	// rows를 반복하며 JoinRoomsView로 변환
	roomSimpleInfos := make([]RoomSimpleInfo, len(rows))
	for i, r := range rows {
		roomSimpleInfos[i] = RoomSimpleInfo{
			ID:        strconv.FormatInt(r.UserID, 10),
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

func ToUserChatRoomView(row databasegen.FindRoomByIDRow) *RoomSimpleInfo {
	return &RoomSimpleInfo{
		ID:        strconv.FormatInt(int64(row.ID), 10),
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
			ID:          int64(r.ID),
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
