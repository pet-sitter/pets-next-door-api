package chat

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 10000
)

var newline = []byte{'\n'}

type Client struct {
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	FbUID    string `json:"id"`
	Name     string `json:"name"`
	rooms    map[*Room]bool
}

// TODO: 이미 조인한 채팅방 조회해서 room에 추가하는 로직 추가
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

// 채팅 읽기
func (client *Client) ReadPump(chatService *service.ChatService) {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	err := client.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Println(err)
		return
	}
	client.conn.SetPongHandler(func(string) error {
		err := client.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return err
		}
		return nil
	})

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
		err := client.conn.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	for {
		select {
		case message, ok := <-client.send:
			err := client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Println(err)
				return
			}
			if !ok {
				if err2 := client.conn.WriteMessage(websocket.CloseMessage, []byte{}); err2 != nil {
					log.Println(err2)
					return
				}
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println(err)
				return
			}

			if _, err := w.Write(message); err != nil {
				log.Println(err)
				w.Close()
				return
			}

			n := len(client.send)
			for i := 0; i < n; i++ {
				if _, err := w.Write(newline); err != nil {
					log.Println(err)
					w.Close()
					return
				}
				if _, err := w.Write(<-client.send); err != nil {
					log.Println(err)
					w.Close()
					return
				}
			}

			if err := w.Close(); err != nil {
				log.Println(err)
				return
			}
		case <-ticker.C:
			err := client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Println(err)
				return
			}
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(err)
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
	client.joinRoom(message, chatService)
}

// 클라이언트가 방을 나갈 때 호출되는 함수
func (client *Client) handleLeaveRoomMessage(roomID int64) {
	room := client.wsServer.findRoomByID(roomID)
	// 클라이언트의 rooms 맵에서 해당 방을 찾아 삭제
	delete(client.rooms, room)
	room.unregister <- client
}

func (client *Client) joinRoom(
	message Message, chatService *service.ChatService,
) {
	room := client.wsServer.findRoomByID(message.Target.GetID())
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
