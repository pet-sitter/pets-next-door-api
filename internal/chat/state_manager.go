package chat

import (
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type StateManager interface {
	RegisterClient(client *Client) *pnd.AppError
	UnregisterClient(client *Client) *pnd.AppError
	FindClientByUID(uid string) (*Client, *pnd.AppError)
	FindRoomByID(roomID int64) *Room
	CreateRoom(name string, roomType chat.RoomType, roomService *service.ChatService) (*Room, *pnd.AppError)
	BroadcastToClients(message []byte) *pnd.AppError
	SetRoom(room *Room) *pnd.AppError
}
