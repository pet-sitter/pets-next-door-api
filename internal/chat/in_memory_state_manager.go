package chat

import (
	"net/http"
	"sync"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type InMemoryStateManager struct {
	clients map[string]*Client
	rooms   map[int64]*Room
	mutex   sync.RWMutex
}

func NewInMemoryStateManager() *InMemoryStateManager {
	return &InMemoryStateManager{
		clients: make(map[string]*Client),
		rooms:   make(map[int64]*Room),
	}
}

func (m *InMemoryStateManager) RegisterClient(client *Client) *pnd.AppError {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.clients[client.FbUID] = client
	return nil
}

func (m *InMemoryStateManager) UnregisterClient(client *Client) *pnd.AppError {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.clients, client.FbUID)
	return nil
}

func (m *InMemoryStateManager) FindClientByUID(uid string) (*Client, *pnd.AppError) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	client, ok := m.clients[uid]
	if !ok {
		return nil, pnd.NewAppError(nil, http.StatusNotFound, pnd.ErrCodeUnknown, "클라이언트를 찾을 수 없습니다.")
	}
	return client, nil
}

func (m *InMemoryStateManager) FindRoomByID(roomID int64) *Room {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	room, ok := m.rooms[roomID]
	if !ok {
		return nil
	}
	return room
}

func (m *InMemoryStateManager) CreateRoom(
	name string, roomType chat.RoomType, roomService *service.ChatService,
) (*Room, *pnd.AppError) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	room, err := NewRoom(name, roomType, roomService)
	if err != nil {
		return nil, err
	}
	go room.RunRoom(roomService)
	m.rooms[room.GetID()] = room
	return room, nil
}

func (m *InMemoryStateManager) BroadcastToClients(message []byte) *pnd.AppError {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	for _, client := range m.clients {
		client.MessageSender <- message
	}
	return nil
}

func (m *InMemoryStateManager) SetRoom(room *Room) *pnd.AppError {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.rooms[room.ID] = room
	return nil
}
