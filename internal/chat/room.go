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
	ID             int64              `json:"id"`
	Name           string             `json:"name"`
	RoomType       chat.RoomType      `json:"roomType"`
	Clients        map[string]*Client `json:"-"`
	RegisterChan   chan *Client       `json:"-"`
	UnregisterChan chan *Client       `json:"-"`
	BroadcastChan  chan *Message      `json:"-"`
}

const welcomeMessage = "%s 이 참여하셨습니다."

func NewRoom(name string, roomType chat.RoomType, roomService *service.ChatService) (*Room, *pnd.AppError) {
	ctx := context.Background()
	row, err := roomService.CreateRoom(ctx, name, roomType)
	if err != nil {
		return nil, pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeRoomCreationFailed, "채팅방 생성에 실패했습니다.")
	}

	return &Room{
		ID:             row.ID,
		Name:           row.Name,
		RoomType:       row.RoomType,
		Clients:        make(map[string]*Client),
		RegisterChan:   make(chan *Client),
		UnregisterChan: make(chan *Client),
		BroadcastChan:  make(chan *Message),
	}, nil
}

func (room *Room) InitRoom(roomID int64, name string, roomType chat.RoomType) *Room {
	return &Room{
		ID:             roomID,
		Name:           name,
		RoomType:       roomType,
		Clients:        make(map[string]*Client),
		RegisterChan:   make(chan *Client),
		UnregisterChan: make(chan *Client),
		BroadcastChan:  make(chan *Message),
	}
}

func (room *Room) RunRoom(roomService *service.ChatService) {
	for {
		select {
		case client := <-room.RegisterChan:
			room.registerClientInRoom(client, room.ID, roomService)
		case client := <-room.UnregisterChan:
			room.unregisterClientInRoom(client, room.ID, roomService)
		case message := <-room.BroadcastChan:
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
	room.Clients[client.FbUID] = client
	return nil
}

func (room *Room) unregisterClientInRoom(client *Client, roomID int64, chatService *service.ChatService) *pnd.AppError {
	ctx := context.Background()
	err := chatService.LeaveRoom(ctx, roomID, client.FbUID)
	if err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown, "채팅방에서 클라이언트를 해제하는 데 실패했습니다.")
	}

	delete(room.Clients, client.FbUID)
	return nil
}

func (room *Room) broadcastToClientsInRoom(message *Message, chatService *service.ChatService) *pnd.AppError {
	ctx := context.Background()
	row, err := chatService.SaveMessage(ctx, room.ID, message.Sender.FbUID, message.Message, message.MessageType)
	if err != nil {
		return pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeUnknown, "메시지를 저장하는 데 실패했습니다.")
	}
	for _, client := range room.Clients {
		client.MessageSender <- []byte(row.Content)
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
