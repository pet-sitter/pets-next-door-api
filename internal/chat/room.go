package chat

import (
	"context"
	"fmt"
	"log"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type Room struct {
	ID       int64         `json:"id"`       // 방 ID
	Name     string        `json:"name"`     // 방 이름
	RoomType chat.RoomType `json:"roomType"` // 방 유형
	// 키를 *Client로 둬서 유저가 추가되거나 나갈 때 o(1) -> WsServer 의 rooms와 같은 맥락
	// 현재 구현으로는 true 만 올 수 있음, 현재 해당 채팅방을 접속해 보고 있는지 여부를 저장하려고 bool으로 선언
	clients map[*Client]bool // 방에 있는 유저 목록
	// TODO: 채널을 액션 별로 모두 쪼개서 사용하는 것이 맞는지 고민.
	register   chan *Client  // 클라이언트 등록 채널
	unregister chan *Client  // 클라이언트 해제 채널
	broadcast  chan *Message // 메시지 브로드캐스트 채널
}

const welcomeMessage = "%s 이 참여하셨습니다."

func NewRoom(name string, roomType chat.RoomType, roomService *service.ChatService) (*Room, error) {
	ctx := context.Background()
	row, err := roomService.CreateRoom(ctx, name, roomType)
	if err != nil {
		return nil, err
	}

	return &Room{
		ID:         row.ID,
		Name:       row.Name,
		RoomType:   row.RoomType,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}, nil
}

// 채팅방 초기화
func (room *Room) InitRoom(roomID int64, name string, roomType chat.RoomType) *Room {
	return &Room{
		ID:         roomID,
		Name:       name,
		RoomType:   roomType,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}
}

func (room *Room) RunRoom(roomService *service.ChatService) {
	for {
		select {
		case client := <-room.register:
			room.registerClientInRoom(client, room.ID, roomService)
		case client := <-room.unregister:
			room.unregisterClientInRoom(client, room.ID, roomService)
		case message := <-room.broadcast:
			room.broadcastToClientsInRoom(message, roomService)
		}
	}
}

// 클라이언트를 방에 등록하는 함수
func (room *Room) registerClientInRoom(client *Client, roomID int64, chatService *service.ChatService) {
	ctx := context.Background()
	_, err := chatService.JoinRoom(ctx, roomID, client.FbUID)
	if err != nil {
		return
	}
	// 클라이언트가 방에 참여했음을 알림
	room.notifyClientJoined(client, chatService)
	room.clients[client] = true
}

// 클라이언트를 방에서 해제하는 함수
func (room *Room) unregisterClientInRoom(client *Client, roomID int64, chatService *service.ChatService) {
	ctx := context.Background()
	err := chatService.LeaveRoom(ctx, roomID, client.FbUID)
	if err != nil {
		log.Println(err)
	}
	delete(room.clients, client)
}

// 방의 모든 클라이언트에게 메시지를 브로드캐스트하는 함수
func (room *Room) broadcastToClientsInRoom(message *Message, chatService *service.ChatService) {
	ctx := context.Background()
	row, err := chatService.SaveMessage(ctx, room.ID, message.Sender.FbUID, message.Message, message.MessageType)
	if err != nil {
		log.Println(err)
	}
	for client := range room.clients {
		client.send <- []byte(row.Content)
	}
}

// 클라이언트가 방에 참여했음을 알리는 함수
func (room *Room) notifyClientJoined(client *Client, chatService *service.ChatService) {
	message := &Message{
		Action:      SendMessageAction,
		Target:      room,
		Message:     fmt.Sprintf(welcomeMessage, client.GetName()),
		MessageType: chat.MessageTypeNormal,
		Sender:      client,
	}

	// 방의 모든 클라이언트에게 환영 메시지 브로드캐스트
	room.broadcastToClientsInRoom(message, chatService)
}

func (room *Room) GetID() int64 {
	return room.ID
}

func (room *Room) GetName() string {
	return room.Name
}
