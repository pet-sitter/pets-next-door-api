package chat

import (
	"context"
	"fmt"
	"log"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type Room struct {
	ID         int64         `json:"id"`
	Name       string        `json:"name"`
	RoomType   chat.RoomType `json:"roomType"`
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
}

const welcomeMessage = "%s 이 참여하셨습니다."

func NewRoom(name string, roomType chat.RoomType, roomService *service.ChatService) (*Room, error) {
	ctx := context.Background()
	row, err := roomService.CreateRoom(ctx, name, roomType)
	if err != nil {
		return nil, err
	}

	return &Room{
		ID:         int64(row.ID),
		Name:       row.Name,
		RoomType:   row.RoomType,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}, nil
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

func (room *Room) registerClientInRoom(client *Client, roomID int64, chatService *service.ChatService) {
	ctx := context.Background()
	_, err := chatService.JoinRoom(ctx, roomID, client.FbUID)
	if err != nil {
		return
	}
	room.notifyClientJoined(client, chatService)
	room.clients[client] = true
}

func (room *Room) unregisterClientInRoom(client *Client, roomID int64, chatService *service.ChatService) {
	ctx := context.Background()
	err := chatService.LeaveRoom(ctx, roomID, client.FbUID)
	if err != nil {
		log.Println(err)
	}
	delete(room.clients, client)
}

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

func (room *Room) notifyClientJoined(client *Client, chatService *service.ChatService) {
	message := &Message{
		Action:      SendMessageAction,
		Target:      room,
		Message:     fmt.Sprintf(welcomeMessage, client.GetName()),
		MessageType: chat.MessageTypeNormal,
		Sender:      client,
	}

	room.broadcastToClientsInRoom(message, chatService)
}

func (room *Room) GetID() int64 {
	return room.ID
}

func (room *Room) GetName() string {
	return room.Name
}
