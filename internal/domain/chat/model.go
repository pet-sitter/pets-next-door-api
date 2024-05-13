package chat

import "time"

type (
	RoomType    string
	MessageType string
)

const (
	RoomTypePersonal  = "personal"
	RoomTypeGathering = "gathering"
)

const (
	MessageTypeNormal  = "normal"
	MessageTypePromise = "promise"
)

type Room struct {
	ID        int64     `field:"id" json:"id"`
	Name      string    `field:"name" json:"name"`
	RoomType  RoomType  `field:"RoomType" json:"RoomType"`
	GroupID   int64     `field:"GroupID" json:"GroupID"`
	CreatedAt time.Time `field:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `field:"updatedAt" json:"updatedAt"`
	DeletedAt time.Time `field:"deletedAt" json:"deletedAt"`
}

type Message struct {
	ID          int64       `field:"id" json:"id"`
	UserID      int64       `field:"userID" json:"userID"`
	RoomID      int64       `field:"roomID" json:"roomID"`
	MessageType MessageType `field:"messageType" json:"messageType"`
	Content     string      `field:"content" json:"content"`
	CreatedAt   time.Time   `field:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time   `field:"updatedAt" json:"updatedAt"`
	DeletedAt   time.Time   `field:"deletedAt" json:"deletedAt"`
}

type UserChatRoom struct {
	ID        int64     `field:"id" json:"id"`
	UserID    int64     `field:"userID" json:"userID"`
	RoomID    int64     `field:"roomID" json:"roomID"`
	CreatedAt time.Time `field:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `field:"updatedAt" json:"updatedAt"`
	DeletedAt time.Time `field:"deletedAt" json:"deletedAt"`
}
