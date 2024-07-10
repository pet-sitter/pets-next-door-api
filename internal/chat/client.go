package chat

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

const (
	// WebSocket 연결에서 메시지를 쓰는 데 최대 10초를 기다림
	writeWait = 10 * time.Second
	// WebSocket 연결에서 마지막으로 받은 메시지 이후로 60초를 기다림
	pongWait = 60 * time.Second
	// 서버가 클라이언트에게 ping 메시지를 보내는 주기
	pingPeriod = (pongWait * 9) / 10
	// 메시지의 최대 크기를 10,000 바이트
	maxMessageSize = 10000
)

// 각 메시지의 끝을 나타내기 위해 사용
var newline = []byte{'\n'}

type Client struct {
	conn          *websocket.Conn
	wsServer      *WebSocketServer
	messageSender chan []byte
	FbUID         string `json:"id"`
	Name          string `json:"name"`
	rooms         map[int64]*Room
}

func NewClient(conn *websocket.Conn, wsServer *WebSocketServer, name, fbUID string) *Client {
	return &Client{
		FbUID:         fbUID,
		Name:          name,
		conn:          conn,
		wsServer:      wsServer,
		messageSender: make(chan []byte, 256),
		rooms:         make(map[int64]*Room),
	}
}

func (client *Client) HandleRead(chatService *service.ChatService) *pnd.AppError {
	defer client.disconnect()
	client.setupConnection()

	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown, "예상치 못한 연결 종료 오류가 발생했습니다.")
			}
			break
		}
		client.handleNewMessage(jsonMessage, chatService)
	}
	return nil
}

func (client *Client) HandleWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.messageSender:
			if !ok {
				client.writeCloseMessage()
				return
			}
			client.writeMessage(message)

		case <-ticker.C:
			client.sendPing()
		}
	}
}

func (client *Client) setupConnection() *pnd.AppError {
	client.conn.SetReadLimit(maxMessageSize)
	err := client.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown, "읽기 제한 시간 설정에 실패했습니다.")
	}

	client.conn.SetPongHandler(func(string) error {
		err := client.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}

func (client *Client) writeMessage(message []byte) *pnd.AppError {
	if err := client.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown, "쓰기 제한 시간 설정에 실패했습니다.")
	}
	w, err := client.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown, "다음 작성자 설정에 실패했습니다.")
	}
	defer w.Close()

	if _, err := w.Write(message); err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown, "메시지 쓰기에 실패했습니다.")
	}

	n := len(client.messageSender)
	for i := 0; i < n; i++ {
		if _, err := w.Write(newline); err != nil {
			return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown, "개행 문자 쓰기에 실패했습니다.")
		}
		if _, err := w.Write(<-client.messageSender); err != nil {
			return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown, "채널에서 메시지 읽기에 실패했습니다.")
		}
	}
	return nil
}

func (client *Client) writeCloseMessage() *pnd.AppError {
	err := client.conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown, "종료 메시지의 쓰기 제한 시간 설정에 실패했습니다.")
	}
	err = client.conn.WriteMessage(websocket.CloseMessage, []byte{})
	if err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown, "종료 메시지 쓰기에 실패했습니다.")
	}
	return nil
}

func (client *Client) sendPing() *pnd.AppError {
	if err := client.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown, "핑 메시지의 쓰기 제한 시간 설정에 실패했습니다.")
	}
	if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown, "핑 메시지 쓰기에 실패했습니다.")
	}
	return nil
}

func (client *Client) disconnect() *pnd.AppError {
	client.wsServer.unregister <- client
	for roomID := range client.rooms {
		room := client.rooms[roomID]
		room.unregister <- client
	}
	close(client.messageSender)
	if err := client.conn.Close(); err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown, "연결 종료에 실패했습니다.")
	}
	return nil
}

func (client *Client) handleNewMessage(jsonMessage []byte, chatService *service.ChatService) *pnd.AppError {
	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		return pnd.NewAppError(err, http.StatusBadRequest, pnd.ErrCodeInvalidBody, "JSON 메시지 해독에 실패했습니다.")
	}
	message.Sender = client
	switch message.Action {
	case SendMessageAction:
		roomID := message.Room.GetID()
		if room := client.wsServer.findRoomByID(roomID); room != nil {
			room.broadcast <- &message
		}
	case JoinRoomAction:
		client.handleJoinRoomMessage(message, chatService)
	case LeaveRoomAction:
		client.handleLeaveRoomMessage(message.Room.GetID())
	}
	return nil
}

func (client *Client) handleJoinRoomMessage(message Message, chatService *service.ChatService) *pnd.AppError {
	if client.wsServer == nil {
		return pnd.NewAppError(nil, http.StatusInternalServerError, pnd.ErrCodeUnknown, "WebSocket 서버가 nil입니다.")
	}

	if message.Room == nil {
		return pnd.NewAppError(nil, http.StatusBadRequest, pnd.ErrCodeInvalidBody, "채팅방 정보가 nil입니다.")
	}

	room := client.wsServer.findRoomByID(message.Room.GetID())
	if room == nil {
		log.Info().Msgf("ID %d의 방을 찾을 수 없어 새 방을 생성합니다.", message.Room.GetID())
		var err *pnd.AppError
		room, err = client.wsServer.createRoom(message.Room.Name, message.Room.RoomType, chatService)
		if err != nil {
			log.Error().Err(err.Err).Msg("방 생성에 실패했습니다.")
			return err
		}
	}

	if message.Sender == nil {
		return pnd.NewAppError(nil, http.StatusBadRequest, pnd.ErrCodeInvalidBody, "보낸 사람이 nil입니다.")
	}

	if _, ok := client.rooms[message.Room.GetID()]; !ok {
		if room.register == nil {
			return pnd.NewAppError(nil, http.StatusInternalServerError, pnd.ErrCodeUnknown, "방 등록 채널이 nil입니다.")
		}

		client.rooms[message.Room.GetID()] = room
		room.register <- client
		err := client.notifyRoomJoined(room, message.Sender)
		if err != nil {
			return err
		}
	}

	return nil
}

func (client *Client) handleLeaveRoomMessage(roomID int64) {
	room := client.wsServer.findRoomByID(roomID)
	delete(client.rooms, room.ID)
	room.unregister <- client
}

func (client *Client) isInRoom(room *Room) bool {
	_, ok := client.rooms[room.ID]
	return ok
}

func (client *Client) notifyRoomJoined(room *Room, sender *Client) *pnd.AppError {
	message := Message{
		Action: RoomJoinedAction,
		Room:   room,
		Sender: sender,
	}

	encodedMessage, err := message.encode()
	if err != nil {
		return err
	}
	client.messageSender <- encodedMessage
	return nil
}

func (client *Client) UpdateConn(conn *websocket.Conn) {
	client.conn = conn
}

func (client *Client) GetName() string {
	return client.Name
}
