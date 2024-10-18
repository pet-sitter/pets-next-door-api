package chat

import (
	"time"

	"github.com/google/uuid"
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
	ID        uuid.UUID            `field:"id" json:"id"`
	RoomName  string               `field:"roomName" json:"roomName"`
	RoomType  string               `field:"roomType" json:"roomType"`
	JoinUser  *JoinUsersSimpleInfo `field:"joinUser" json:"joinUser"`
	CreatedAt time.Time            `field:"createdAt" json:"createdAt"`
	UpdatedAt time.Time            `field:"updatedAt" json:"updatedAt"`
}

type JoinUsersSimpleInfo struct {
	ID               uuid.UUID `field:"id" json:"userId"`
	UserNickname     string    `field:"nickname" json:"userNickname"`
	UserProfileImage string    `field:"profileImage" json:"profileImageUrl"`
}

type JoinRoom struct {
	UserID   uuid.UUID
	RoomID   uuid.UUID
	JoinedAt time.Time
}

// 조회 시 Room 정보를 반환하는 View
type JoinRoomsView struct {
	Items []RoomSimpleInfo `field:"items" json:"items"`
}

type UserChatRoomMessageView struct {
	ID          uuid.UUID `field:"id" json:"id"`
	MessageType string    `field:"messageType" json:"messageType"`
}

type Message struct {
	ID          uuid.UUID `field:"id" json:"id"`
	UserID      uuid.UUID `field:"userID" json:"userID"`
	RoomID      uuid.UUID `field:"roomID" json:"roomID"`
	MessageType string    `field:"messageType" json:"messageType"`
	Content     string    `field:"content" json:"content"`
	CreatedAt   time.Time `field:"createdAt" json:"createdAt"`
}

type MessageCursorView struct {
	HasNext bool      `field:"hasNext" json:"hasNext"`
	HasPrev bool      `field:"hasPrev" json:"hasPrev"`
	Items   []Message `field:"items" json:"items,omitempty"`
}
