package chat

import (
	"context"
	"net/http"
	"sync"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type InMemoryStateManager struct {
	clients     map[string]*Client
	rooms       map[int64]*Room
	clientRooms map[string]map[int64]struct{}
	roomClients map[int64]map[string]*Client
	mutex       sync.RWMutex
}

func NewInMemoryStateManager() *InMemoryStateManager {
	return &InMemoryStateManager{
		clients:     make(map[string]*Client),
		rooms:       make(map[int64]*Room),
		clientRooms: make(map[string]map[int64]struct{}),
		roomClients: make(map[int64]map[string]*Client),
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
	name string, roomType chat.RoomType, roomService *service.ChatService, stateManager StateManager,
) (*Room, *pnd.AppError) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	ctx := context.Background()
	row, err := roomService.CreateRoom(ctx, name, roomType)
	if err != nil {
		return nil, pnd.NewAppError(err, http.StatusInternalServerError, pnd.ErrCodeRoomCreationFailed, "채팅방 생성에 실패했습니다.")
	}
	room := NewRoom(row.ID, row.Name, row.RoomType, stateManager)
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

func (m *InMemoryStateManager) JoinRoom(roomID int64, clientID string) *pnd.AppError {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, ok := m.clientRooms[clientID]; !ok {
		m.clientRooms[clientID] = make(map[int64]struct{})
	}
	m.clientRooms[clientID][roomID] = struct{}{}
	if _, ok := m.roomClients[roomID]; !ok {
		m.roomClients[roomID] = make(map[string]*Client)
	}
	m.roomClients[roomID][clientID] = m.clients[clientID]
	return nil
}

func (m *InMemoryStateManager) LeaveRoom(roomID int64, clientID string) *pnd.AppError {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.clientRooms[clientID], roomID)
	delete(m.roomClients[roomID], clientID)
	return nil
}

func (m *InMemoryStateManager) IsClientInRoom(clientID string, roomID int64) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	_, ok := m.clientRooms[clientID][roomID]
	return ok
}

func (m *InMemoryStateManager) GetClientRooms(clientID string) map[int64]struct{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.clientRooms[clientID]
}

func (m *InMemoryStateManager) GetRoomClients(roomID int64) map[string]*Client {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.roomClients[roomID]
}

func (m *InMemoryStateManager) SetRoom(room *Room) *pnd.AppError {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.rooms[room.ID] = room
	return nil
}
