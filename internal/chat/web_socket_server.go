package chat

import (
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type WebSocketServer struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	rooms      map[int64]*Room
}

func NewWebsocketServer() *WebSocketServer {
	return &WebSocketServer{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		rooms:      make(map[int64]*Room),
	}
}

func (server *WebSocketServer) Run() {
	for {
		// 해당하는 채널에 메시지가 들어올 때 작동
		select {
		case client := <-server.register:
			server.RegisterClient(client)
		case client := <-server.unregister:
			server.unregisterClient(client)
		case message := <-server.broadcast:
			server.broadcastToClients(message)
		}
	}
}

func (server *WebSocketServer) RegisterClient(client *Client) {
	server.notifyClientJoined(client)
	server.listOnlineClients(client)
	server.clients[client.FbUID] = client
}

func (server *WebSocketServer) unregisterClient(client *Client) {
	if _, ok := server.clients[client.FbUID]; ok {
		delete(server.clients, client.FbUID)
		server.notifyClientLeft(client)
	}
}

func (server *WebSocketServer) FindClientByUID(uid string) *Client {
	return server.clients[uid]
}

func (server *WebSocketServer) notifyClientJoined(client *Client) *pnd.AppError {
	message := &Message{
		Action: UserJoinedAction,
		Sender: client,
	}
	messageText, err := message.encode()
	if err != nil {
		return err
	}
	server.broadcastToClients(messageText)
	return nil
}

func (server *WebSocketServer) notifyClientLeft(client *Client) *pnd.AppError {
	message := &Message{
		Action: UserLeftAction,
		Sender: client,
	}
	messageText, err := message.encode()
	if err != nil {
		return err
	}
	server.broadcastToClients(messageText)
	return nil
}

func (server *WebSocketServer) listOnlineClients(newClient *Client) *pnd.AppError {
	for _, client := range server.clients {
		message := &Message{
			Action: UserJoinedAction,
			Sender: client,
		}
		encodedMessage, err := message.encode()
		if err != nil {
			return err
		}
		newClient.messageSender <- encodedMessage
	}
	return nil
}

func (server *WebSocketServer) broadcastToClients(message []byte) {
	for _, client := range server.clients {
		client.messageSender <- message
	}
}

func (server *WebSocketServer) findRoomByID(roomID int64) *Room {
	if room, ok := server.rooms[roomID]; ok {
		return room
	}
	return nil
}

func (server *WebSocketServer) createRoom(
	name string, roomType chat.RoomType, roomService *service.ChatService,
) (*Room, *pnd.AppError) {
	room, err := NewRoom(name, roomType, roomService)
	if err != nil {
		return nil, err
	}
	go room.RunRoom(roomService)
	server.rooms[room.GetID()] = room

	return room, nil
}
