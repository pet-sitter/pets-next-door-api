package chat

import (
	"encoding/json"
	"net/http"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
)

const (
	SendMessageAction = "SEND_MESSAGE" // 클라이언트가 메시지를 보낼 때
	JoinRoomAction    = "JOIN_ROOM"    // 클라이언트가 방에 참여할 때
	LeaveRoomAction   = "LEAVE_ROOM"   // 클라이언트가 방을 떠날 때
	UserJoinedAction  = "USER_JOIN"    // 다른 사용자가 방에 참여했음을 알릴 때
	UserLeftAction    = "USER_LEFT"    // 다른 사용자가 방을 떠났음을 알릴 때
	RoomJoinedAction  = "ROOM_JOINED"  // 클라이언트가 방에 성공적으로 참여했음을 알릴 때
)

type Message struct {
	Action      string           `json:"action"`
	Message     string           `json:"message"`
	MessageType chat.MessageType `json:"messageType"`
	Target      *Room            `json:"target"`
	Sender      *Client          `json:"sender"`
}

func (message *Message) encode() ([]byte, *pnd.AppError) {
	bytes, err := json.Marshal(message)
	if err != nil {
		return nil, pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeMessageEncodingFailed, "메시지 인코딩에 실패했습니다")
	}
	return bytes, nil
}
