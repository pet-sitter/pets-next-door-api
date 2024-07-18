package chat

import (
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type StateManager interface {
	RegisterClient(client *Client) *pnd.AppError
	UnregisterClient(client *Client) *pnd.AppError
	FindClientByUID(uid string) *Client
	FindRoomByID(roomID int64) *Room
	CreateRoom(
		name string, roomType chat.RoomType, roomService *service.ChatService,
	) (*Room, *pnd.AppError)
	BroadcastToClients(message []byte) *pnd.AppError
	JoinRoom(roomID int64, clientID string) *pnd.AppError
	LeaveRoom(roomID int64, clientID string) *pnd.AppError
	IsClientInRoom(clientID string, roomID int64) bool
	GetClientRooms(clientID string) map[int64]struct{}
	GetRoomClients(roomID int64) map[string]*Client
	SetRoom(room *Room) *pnd.AppError
}
