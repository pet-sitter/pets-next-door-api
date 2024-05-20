package chat

import (
	"encoding/json"
	"log"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
)

const (
	SendMessageAction = "send-message"
	JoinRoomAction    = "join-room"
	LeaveRoomAction   = "leave-room"
	UserJoinedAction  = "user-join"
	UserLeftAction    = "user-left"
	RoomJoinedAction  = "room-joined"
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
