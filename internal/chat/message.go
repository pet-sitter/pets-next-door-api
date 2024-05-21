package chat

import (
	"encoding/json"
	"log"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
)

const (
	SendMessageAction = "SEND_MESSAGE"
	JoinRoomAction    = "JOIN_ROOM"
	LeaveRoomAction   = "LEAVE_ROOM"
	UserJoinedAction  = "USER_JOIN"
	UserLeftAction    = "USER_LEFT"
	RoomJoinedAction  = "ROOM_JOINED"
)

type Message struct {
	Action      string           `json:"action"`
	Message     string           `json:"message"`
	MessageType chat.MessageType `json:"messageType"`
	Target      *Room            `json:"target"`
	Sender      *Client          `json:"sender"`
}

func (message *Message) encode() []byte {
	bytes, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return bytes
}
