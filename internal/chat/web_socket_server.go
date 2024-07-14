package chat

import (
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type WebSocketServer struct {
	StateManager   StateManager
	RegisterChan   chan *Client
	UnregisterChan chan *Client
	BroadcastChan  chan []byte
}

func NewWebSocketServer(stateManager StateManager) *WebSocketServer {
	return &WebSocketServer{
		StateManager:   stateManager,
		RegisterChan:   make(chan *Client),
		UnregisterChan: make(chan *Client),
		BroadcastChan:  make(chan []byte),
	}
}

func (server *WebSocketServer) Run() {
	for {
		select {
		case client := <-server.RegisterChan:
			server.StateManager.RegisterClient(client)
		case client := <-server.UnregisterChan:
			server.StateManager.UnregisterClient(client)
		case message := <-server.BroadcastChan:
			server.StateManager.BroadcastToClients(message)
		}
	}
}

func (server *WebSocketServer) RegisterClient(client *Client) {
	server.StateManager.RegisterClient(client)
}

func (server *WebSocketServer) UnregisterClient(client *Client) {
	server.StateManager.UnregisterClient(client)
}

func (server *WebSocketServer) FindClientByUID(uid string) *Client {
	client, _ := server.StateManager.FindClientByUID(uid)
	return client
}

func (server *WebSocketServer) CreateRoom(
	name string, roomType chat.RoomType, roomService *service.ChatService,
) (*Room, *pnd.AppError) {
	return server.StateManager.CreateRoom(name, roomType, roomService)
}

func (server *WebSocketServer) BroadcastToClients(message []byte) {
	server.StateManager.BroadcastToClients(message)
}

func (server *WebSocketServer) FindRoomByID(roomID int64) *Room {
	room := server.StateManager.FindRoomByID(roomID)
	return room
}
