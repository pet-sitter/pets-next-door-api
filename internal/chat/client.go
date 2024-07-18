package chat

import (
	"encoding/json"
	"fmt"
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
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 2048
)

var newline = []byte{'\n'}

type Client struct {
	Conn          *websocket.Conn `json:"-"`
	MessageSender chan []byte     `json:"-"`
	FbUID         string          `json:"id"`
	Name          string          `json:"name"`
}

func NewClient(conn *websocket.Conn, name, fbUID string) *Client {
	return &Client{
		FbUID:         fbUID,
		Name:          name,
		Conn:          conn,
		MessageSender: make(chan []byte, 256),
	}
}

func (client *Client) HandleRead(stateManager StateManager, chatService *service.ChatService) *pnd.AppError {
	defer client.disconnect(stateManager)
	client.setupConnection()

	for {
		_, jsonMessage, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown,
					"예상치 못한 연결 종료 오류가 발생했습니다. userID="+client.FbUID)
			}
			break
		}
		if len(jsonMessage) > maxMessageSize {
			errMsg := fmt.Sprintf("메시지 크기가 최대 크기(%d 바이트)를 초과합니다.", maxMessageSize)
			err := client.Conn.WriteMessage(websocket.TextMessage, []byte(errMsg))
			if err != nil {
				return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown,
					"메시지 크기 초과 오류 메시지를 전송하는 데 실패했습니다. userID="+client.FbUID)
			}
			continue
		}
		client.handleNewMessage(jsonMessage, stateManager, chatService)
	}
	return nil
}

func (client *Client) HandleWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.MessageSender:
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
	err := client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown,
			"읽기 제한 시간 설정에 실패했습니다. userID="+client.FbUID)
	}

	client.Conn.SetPongHandler(func(string) error {
		err := client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}

func (client *Client) writeMessage(message []byte) *pnd.AppError {
	if err := client.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown,
			"쓰기 제한 시간 설정에 실패했습니다. userID="+client.FbUID)
	}
	w, err := client.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown,
			"다음 작성자 설정에 실패했습니다. userID="+client.FbUID)
	}
	defer w.Close()

	if _, err := w.Write(message); err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown,
			"메시지 쓰기에 실패했습니다. userID=%s"+client.FbUID)
	}

	n := len(client.MessageSender)
	for i := 0; i < n; i++ {
		if _, err := w.Write(newline); err != nil {
			return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown,
				"개행 문자 쓰기에 실패했습니다. userID=%s"+client.FbUID)
		}

		if _, err := w.Write(<-client.MessageSender); err != nil {
			return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown,
				"채널에서 메시지 읽기에 실패했습니다. userID=%s"+client.FbUID)
		}
	}
	return nil
}

func (client *Client) writeCloseMessage() *pnd.AppError {
	err := client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown,
			"종료 메시지의 쓰기 제한 시간 설정에 실패했습니다. userID=%s"+client.FbUID)
	}
	err = client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
	if err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown,
			"종료 메시지 쓰기에 실패했습니다. userID=%s"+client.FbUID)
	}
	return nil
}

func (client *Client) sendPing() *pnd.AppError {
	if err := client.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown,
			"핑 메시지의 쓰기 제한 시간 설정에 실패했습니다. userID=%s"+client.FbUID)
	}
	if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown,
			"핑 메시지 쓰기에 실패했습니다. userID=%s"+client.FbUID)
	}
	return nil
}

func (client *Client) disconnect(stateManager StateManager) *pnd.AppError {
	stateManager.UnregisterClient(client)
	for roomID := range stateManager.GetClientRooms(client.FbUID) {
		stateManager.LeaveRoom(roomID, client.FbUID)
	}
	close(client.MessageSender)
	if err := client.Conn.Close(); err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown,
			"연결 종료에 실패했습니다. userID=%s"+client.FbUID)
	}
	return nil
}

func (client *Client) handleNewMessage(
	jsonMessage []byte, stateManager StateManager, chatService *service.ChatService,
) *pnd.AppError {
	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		return pnd.NewAppError(err, http.StatusBadRequest, pnd.ErrCodeInvalidBody,
			"JSON 메시지 해독에 실패했습니다. userID=%s"+client.FbUID)
	}
	message.Sender = client
	switch message.Action {
	case SendMessageAction:
		roomID := message.Room.GetID()
		if room := stateManager.FindRoomByID(roomID); room != nil {
			room.BroadcastChan <- &message
		}
	case JoinRoomAction:
		client.handleJoinRoomMessage(message, stateManager, chatService)
	case LeaveRoomAction:
		client.handleLeaveRoomMessage(message.Room.GetID(), stateManager)
	}
	return nil
}

func (client *Client) handleJoinRoomMessage(
	message Message, stateManager StateManager, chatService *service.ChatService,
) *pnd.AppError {
	if message.Room == nil {
		return pnd.NewAppError(nil, http.StatusBadRequest, pnd.ErrCodeInvalidBody,
			"채팅방 정보가 nil입니다. userID=%s"+client.FbUID)
	}

	room, err := client.CreateRoomIfNotExists(message, stateManager, chatService)
	if err != nil {
		return err
	}

	if message.Sender == nil {
		return pnd.NewAppError(nil, http.StatusBadRequest, pnd.ErrCodeInvalidBody,
			"보낸 사람이 nil입니다. userID=%s"+client.FbUID)
	}

	if !stateManager.IsClientInRoom(client.FbUID, message.Room.GetID()) {
		if room.RegisterChan == nil {
			return pnd.NewAppError(nil, http.StatusInternalServerError, pnd.ErrCodeUnknown,
				"방 등록 채널이 nil입니다. userID=%s"+client.FbUID)
		}

		stateManager.JoinRoom(room.ID, client.FbUID)
		room.RegisterChan <- client
		err := client.notifyRoomJoined(room, message.Sender)
		if err != nil {
			return err
		}
	}

	return nil
}

func (client *Client) CreateRoomIfNotExists(
	message Message, stateManager StateManager, chatService *service.ChatService,
) (*Room, *pnd.AppError) {
	if stateManager == nil {
		return nil, pnd.NewAppError(nil, http.StatusInternalServerError, pnd.ErrCodeUnknown,
			"StateManager가 nil입니다. userID=%s"+client.FbUID)
	}
	room := stateManager.FindRoomByID(message.Room.GetID())
	if room == nil {
		log.Info().Msgf("ID %d의 방을 찾을 수 없어 새 방을 생성합니다. userID=%s", message.Room.GetID(), client.FbUID)
		var err *pnd.AppError
		room, err = stateManager.CreateRoom(message.Room.Name, message.Room.RoomType, chatService)
		if err != nil {
			log.Error().Err(err.Err).Msgf("방 생성에 실패했습니다. userID=%s", client.FbUID)
			return nil, err
		}
		return room, nil
	}
	return room, nil
}

func (client *Client) handleLeaveRoomMessage(roomID int64, stateManager StateManager) {
	stateManager.LeaveRoom(roomID, client.FbUID)
	room := stateManager.FindRoomByID(roomID)
	room.UnregisterChan <- client
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
	client.MessageSender <- encodedMessage
	return nil
}

func (client *Client) UpdateConn(conn *websocket.Conn) {
	client.Conn = conn
}

func (client *Client) GetName() string {
	return client.Name
}
