package chat

import (
	"context"
	"fmt"
	"net/http"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type Room struct {
	ID       int64         `json:"id"`
	Name     string        `json:"name"`
	RoomType chat.RoomType `json:"roomType"`
	clients  map[string]*Client

	register   chan *Client  // 클라이언트 등록 채널
	unregister chan *Client  // 클라이언트 해제 채널
	broadcast  chan *Message // 메시지 브로드캐스트 채널
}

const welcomeMessage = "%s 이 참여하셨습니다."

func NewRoom(name string, roomType chat.RoomType, roomService *service.ChatService) (*Room, *pnd.AppError) {
	ctx := context.Background()
	row, err := roomService.CreateRoom(ctx, name, roomType)
	if err != nil {
		return nil, pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeRoomCreationFailed, "채팅방 생성에 실패했습니다.")
	}

	return &Room{
		ID:         row.ID,
		Name:       row.Name,
		RoomType:   row.RoomType,
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}, nil
}

func (room *Room) InitRoom(roomID int64, name string, roomType chat.RoomType) *Room {
	return &Room{
		ID:         roomID,
		Name:       name,
		RoomType:   roomType,
		clients:    make(map[string]*Client),
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

func (room *Room) registerClientInRoom(client *Client, roomID int64, chatService *service.ChatService) *pnd.AppError {
	ctx := context.Background()
	_, err := chatService.JoinRoom(ctx, roomID, client.FbUID)
	if err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown, "채팅방에 클라이언트를 등록하는 데 실패했습니다.")
	}
	room.notifyClientJoined(client, chatService)
	room.clients[client.FbUID] = client
	return nil
}

func (room *Room) unregisterClientInRoom(client *Client, roomID int64, chatService *service.ChatService) *pnd.AppError {
	ctx := context.Background()
	err := chatService.LeaveRoom(ctx, roomID, client.FbUID)
	if err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown, "채팅방에서 클라이언트를 해제하는 데 실패했습니다.")
	}

	delete(room.clients, client.FbUID)
	return nil
}

func (room *Room) broadcastToClientsInRoom(message *Message, chatService *service.ChatService) *pnd.AppError {
	ctx := context.Background()
	row, err := chatService.SaveMessage(ctx, room.ID, message.Sender.FbUID, message.Message, message.MessageType)
	if err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown, "메시지를 저장하는 데 실패했습니다.")
	}
	for _, client := range room.clients {
		client.messageSender <- []byte(row.Content)
	}
	return nil
}

func (room *Room) notifyClientJoined(client *Client, chatService *service.ChatService) *pnd.AppError {
	message := &Message{
		Action:      SendMessageAction,
		Room:        room,
		Message:     fmt.Sprintf(welcomeMessage, client.GetName()),
		MessageType: chat.MessageTypeNormal,
		Sender:      client,
	}

	err := room.broadcastToClientsInRoom(message, chatService)
	if err != nil {
		return err
	}
	return nil
}

func (room *Room) GetID() int64 {
	return room.ID
}

func (room *Room) GetName() string {
	return room.Name
}
