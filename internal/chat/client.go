package chat

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
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
	conn     *websocket.Conn // 웹소켓 연결
	wsServer *WsServer       // 웹소켓 서버
	send     chan []byte     // 클라이언트로 보낼 메시지를 저장하는 채널
	FbUID    string          `json:"id"`   // 유저의 UID
	Name     string          `json:"name"` // 유저의 이름
	rooms    map[*Room]bool  // 클라이언트가 속한 방들
}

func NewClient(conn *websocket.Conn, wsServer *WsServer, name, fbUID string) *Client {
	return &Client{
		FbUID:    fbUID,
		Name:     name,
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		rooms:    make(map[*Room]bool),
	}
}

func (client *Client) setupConnection() {
	client.conn.SetReadLimit(maxMessageSize)
	err := client.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Println(err)
	}

	client.conn.SetPongHandler(func(string) error {
		err := client.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return err
		}
		return nil
	})
}

func (client *Client) writeMessage(message []byte) {
	err := client.conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err != nil {
		log.Println(err)
		return
	}
	w, err := client.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		log.Println(err)
		return
	}
	defer w.Close()

	if _, err := w.Write(message); err != nil {
		log.Println(err)
		return
	}

	n := len(client.send)
	for i := 0; i < n; i++ {
		if _, err := w.Write(newline); err != nil {
			log.Println(err)
			return
		}
		if _, err := w.Write(<-client.send); err != nil {
			log.Println(err)
			return
		}
	}
}

func (client *Client) HandleRead(chatService *service.ChatService) {
	defer client.disconnect()
	client.setupConnection()

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

func (client *Client) HandleWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
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

func (client *Client) writeCloseMessage() {
	err := client.conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err != nil {
		log.Println(err)
	}
	err = client.conn.WriteMessage(websocket.CloseMessage, []byte{})
	if err != nil {
		log.Println(err)
	}
}

func (client *Client) sendPing() {
	err := client.conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err != nil {
		log.Println(err)
	}
	err = client.conn.WriteMessage(websocket.PingMessage, nil)
	if err != nil {
		log.Println(err)
	}
}

// 클라이언트 연결을 종료하는 함수
func (client *Client) disconnect() {
	// 서버에서 클라이언트 등록 해제
	client.wsServer.unregister <- client
	// 클라이언트가 속한 모든 방에서 클라이언트 등록 해제
	for room := range client.rooms {
		room.unregister <- client
	}
	// 채널 닫음
	close(client.send)
	client.conn.Close()
}

// 클라이언트가 새로운 메시지를 보낼 때 호출되는 함수
func (client *Client) handleNewMessage(jsonMessage []byte, chatService *service.ChatService) {
	var message Message
	// []byte -> JSON
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return
	}
	message.Sender = client
	switch message.Action {
	case SendMessageAction:
		roomID := message.Target.GetID()
		if room := client.wsServer.findRoomByID(roomID); room != nil {
			room.broadcast <- &message
		}
	case JoinRoomAction:
		client.handleJoinRoomMessage(message, chatService)
	case LeaveRoomAction:
		client.handleLeaveRoomMessage(message.Target.GetID())
	}
}

// 클라이언트가 방에 들어가거나 나갈 때 호출되는 함수
func (client *Client) handleJoinRoomMessage(message Message, chatService *service.ChatService) {
	// 방이 없으면 새로 생성
	room := client.wsServer.findRoomByID(message.Target.GetID())
	if room == nil {
		room = client.wsServer.createRoom(message.Target.Name, message.Target.RoomType, chatService)
	}

	if message.Sender == nil {
		log.Println("Sender is nil")
		return
	}

	// 클라이언트가 방에 참여하지 않았으면 방에 참여시킴
	if !client.isInRoom(room) {
		client.rooms[room] = true
		room.register <- client

		client.notifyRoomJoined(room, message.Sender)
	}
}

// 클라이언트가 방을 나갈 때 호출되는 함수
func (client *Client) handleLeaveRoomMessage(roomID int64) {
	room := client.wsServer.findRoomByID(roomID)
	// 클라이언트의 rooms 맵에서 해당 방을 찾아 삭제
	delete(client.rooms, room)
	room.unregister <- client
}

// 클라이언트가 특정 방에 있는지 확인하는 함수
func (client *Client) isInRoom(room *Room) bool {
	if _, ok := client.rooms[room]; ok {
		return true
	}

	return false
}

// 클라이언트가 방에 참여했음을 알리는 함수
func (client *Client) notifyRoomJoined(room *Room, sender *Client) {
	message := Message{
		Action: RoomJoinedAction,
		Target: room,
		Sender: sender,
	}

	client.send <- message.encode()
}

func (client *Client) UpdateConn(conn *websocket.Conn) {
	client.conn = conn
}

func (client *Client) GetName() string {
	return client.Name
}
