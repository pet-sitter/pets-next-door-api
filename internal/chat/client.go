package chat

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
	"log"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 10000
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	FbUID    string `json:"id"`
	Name     string `json:"name"`
	rooms    map[*Room]bool
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

// TODO: 이미 조인한 채팅방 조회해서 room에 추가하는 로직 추가
func NewClient(conn *websocket.Conn, wsServer *WsServer, name string, fbUID string) *Client {
	return &Client{
		FbUID:    fbUID,
		Name:     name,
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		rooms:    make(map[*Room]bool),
	}
}

// 채팅 읽기
func (client *Client) ReadPump(chatService *service.ChatService) {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		client.handleNewMessage(jsonMessage, chatService)
	}

}

// 채팅 쓰기
func (client *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// 연결 끊기
func (client *Client) disconnect() {
	client.wsServer.unregister <- client
	for room := range client.rooms {
		room.unregister <- client
	}
	close(client.send)
	client.conn.Close()
}

// 클라이언트가 새로운 메시지를 보낼 때 호출되는 함수
func (client *Client) handleNewMessage(jsonMessage []byte, chatService *service.ChatService) {

	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return
	}
	message.Sender = client
	switch message.Action {
	case SendMessageAction:
		roomID := message.Target.GetId()
		if room := client.wsServer.findRoomByID(roomID); room != nil {
			room.broadcast <- &message
		}

	case JoinRoomAction:
		client.handleJoinRoomMessage(message, chatService)

	case LeaveRoomAction:
		client.handleLeaveRoomMessage(message.Target.GetId())
	}

}

// 클라이언트가 방에 들어가거나 나갈 때 호출되는 함수
func (client *Client) handleJoinRoomMessage(message Message, chatService *service.ChatService) {
	client.joinRoom(message, chatService)
}

// 클라이언트가 방을 나갈 때 호출되는 함수
func (client *Client) handleLeaveRoomMessage(roomID int64) {
	room := client.wsServer.findRoomByID(roomID)
	if _, ok := client.rooms[room]; ok {
		// 클라이언트의 rooms 맵에서 해당 방을 찾아 삭제
		delete(client.rooms, room)
	}

	room.unregister <- client
}

func (client *Client) joinRoom(
	message Message, chatService *service.ChatService,
) {

	room := client.wsServer.findRoomByID(message.Target.GetId())
	if room == nil {
		room = client.wsServer.createRoom(message.Target.Name, message.Target.RoomType, chatService)
	}

	if message.Sender == nil {
		log.Println("Sender is nil")
		return
	}

	if !client.isInRoom(room) {

		client.rooms[room] = true
		room.register <- client

		client.notifyRoomJoined(room, message.Sender)
	}

}

func (client *Client) isInRoom(room *Room) bool {
	if _, ok := client.rooms[room]; ok {
		return true
	}

	return false
}

func (client *Client) notifyRoomJoined(room *Room, sender *Client) {
	message := Message{
		Action: RoomJoinedAction,
		Target: room,
		Sender: sender,
	}

	client.send <- message.encode()
}

func (client *Client) GetName() string {
	return client.Name
}
