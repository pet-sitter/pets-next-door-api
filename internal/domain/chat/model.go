package chat

import (
	"database/sql"
	"time"
)

type (
	RoomType    string
	MessageType string
)

func (t RoomType) IsValid() bool {
	switch t {
	case EventRoomType:
		return true
	default:
		return false
	}
}

const (
	EventRoomType = "event"
)

const (
	EventMessage = "event"
)

type RoomSimpleInfo struct {
	ID        string                 `field:"id" json:"id"`
	RoomName  string                 `field:"roomName" json:"roomName"`
	RoomType  string                 `field:"roomType" json:"roomType"`
	JoinUsers *[]JoinUsersSimpleInfo `field:"joinUsers" json:"joinUsers"`
	CreatedAt time.Time              `field:"createdAt" json:"createdAt"`
	UpdatedAt time.Time              `field:"updatedAt" json:"updatedAt"`
}

type JoinUsersSimpleInfo struct {
	ID               string         `field:"id" json:"userId"`
	UserNickname     string         `field:"nickname" json:"userNickname"`
	UserProfileImage sql.NullString `field:"profileImage" json:"profileImageId"`
}

type JoinRoom struct {
	UserID   string
	RoomID   string
	JoinedAt time.Time
}

// 조회 시 Room 정보를 반환하는 View
type JoinRoomsView struct {
	Items []RoomSimpleInfo `field:"items" json:"items"`
}

type UserChatRoomMessageView struct {
	ID          string `field:"id" json:"id"`
	MessageType string `field:"messageType" json:"messageType"`
}

type Message struct {
	ID          int64     `field:"id" json:"id"`
	UserID      int64     `field:"userID" json:"userID"`
	RoomID      int64     `field:"roomID" json:"roomID"`
	MessageType string    `field:"messageType" json:"messageType"`
	Content     string    `field:"content" json:"content"`
	CreatedAt   time.Time `field:"createdAt" json:"createdAt"`
}

type MessageCursorView struct {
	HasNext *bool     `field:"hasNext" json:"hasNext"`
	HasPrev *bool     `field:"hasPrev" json:"hasPrev"`
	Items   []Message `field:"items" json:"items,omitempty"`
}
