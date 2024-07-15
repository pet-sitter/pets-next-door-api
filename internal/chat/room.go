package chat

import (
	"context"
	"fmt"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type Room struct {
	ID             int64         `json:"id"`
	Name           string        `json:"name"`
	RoomType       chat.RoomType `json:"roomType"`
	StateManager   StateManager  `json:"-"`
	RegisterChan   chan *Client  `json:"-"`
	UnregisterChan chan *Client  `json:"-"`
	BroadcastChan  chan *Message `json:"-"`
}

const welcomeMessage = "%s 이 참여하셨습니다."

func NewRoom(id int64, name string, roomType chat.RoomType, stateManager StateManager) *Room {
	return &Room{
		ID:             id,
		Name:           name,
		RoomType:       roomType,
		StateManager:   stateManager,
		RegisterChan:   make(chan *Client),
		UnregisterChan: make(chan *Client),
		BroadcastChan:  make(chan *Message),
	}
}

func (room *Room) RunRoom(chatService *service.ChatService) {
	for {
		select {
		case client := <-room.RegisterChan:
			room.registerClientInRoom(client, room.ID, chatService)
		case client := <-room.UnregisterChan:
			room.unregisterClientInRoom(client, room.ID, chatService)
		case message := <-room.BroadcastChan:
			room.broadcastToClientsInRoom(message, chatService)
		}
	}
}

func (room *Room) registerClientInRoom(client *Client, roomID int64, chatService *service.ChatService) *pnd.AppError {
	ctx := context.Background()
	exists, err := chatService.ExistsUserInRoom(ctx, roomID, client.FbUID)
	if exists {
		return nil
	}
	if err != nil {
		return err
	}
	_, err = chatService.JoinRoom(ctx, roomID, client.FbUID)
	if err != nil {
		return err
	}
	room.notifyClientJoined(client, chatService)
	return nil
}

func (room *Room) unregisterClientInRoom(client *Client, roomID int64, chatService *service.ChatService) *pnd.AppError {
	ctx := context.Background()
	err := chatService.LeaveRoom(ctx, roomID, client.FbUID)
	if err != nil {
		return err
	}
	return nil
}

func (room *Room) broadcastToClientsInRoom(message *Message, chatService *service.ChatService) *pnd.AppError {
	ctx := context.Background()
	row, err := chatService.SaveMessage(ctx, room.ID, message.Sender.FbUID, message.Message, message.MessageType)
	if err != nil {
		return err
	}
	clients := room.StateManager.GetRoomClients(room.ID)

	for _, client := range clients {
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
