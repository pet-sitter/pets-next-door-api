package chat

import (
	"context"
	"sync"

	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type InMemoryStateManager struct {
	clients map[string]*Client
	rooms   map[int64]*Room
	// 클라이언트의 고유 식별자(FbUID)를 키로 사용하고, 클라이언트가 참여한 방의 ID를 값으로 갖는 또 다른 맵을 값으로 가집니다.
	// 값이 struct{}인 이유는 이중 맵에서 값의 실제 데이터가 필요 없기 때문입니다. struct{}는 메모리를 거의 차지하지 않으므로 효율적입니다.
	clientRooms    map[string]map[int64]struct{}
	roomClientUIDs map[int64][]string
	mutex          sync.RWMutex
}

func NewInMemoryStateManager() *InMemoryStateManager {
	return &InMemoryStateManager{
		clients:        make(map[string]*Client),
		rooms:          make(map[int64]*Room),
		clientRooms:    make(map[string]map[int64]struct{}),
		roomClientUIDs: make(map[int64][]string),
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

func (m *InMemoryStateManager) FindClientByUID(uid string) *Client {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	client, ok := m.clients[uid]
	if !ok {
		return nil
	}
	return client
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
	name string, roomType chat.RoomType, chatService *service.ChatService,
) (*Room, *pnd.AppError) {
	ctx := context.Background()

	row, err := chatService.CreateRoom(ctx, name, roomType)
	if err != nil {
		return nil, err
	}
	room := NewRoom(row.ID, row.Name, row.RoomType, m)
	go room.RunRoom(chatService)
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
	m.roomClientUIDs[roomID] = append(m.roomClientUIDs[roomID], clientID)
	return nil
}

func (m *InMemoryStateManager) LeaveRoom(roomID int64, clientID string) *pnd.AppError {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.clientRooms[clientID], roomID)
	for i, id := range m.roomClientUIDs[roomID] {
		if id == clientID {
			m.roomClientUIDs[roomID] = append(m.roomClientUIDs[roomID][:i], m.roomClientUIDs[roomID][i+1:]...)
			break
		}
	}
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
	clientMap := make(map[string]*Client)
	for _, clientID := range m.roomClientUIDs[roomID] {
		clientMap[clientID] = m.clients[clientID]
	}
	return clientMap
}

func (m *InMemoryStateManager) SetRoom(room *Room) *pnd.AppError {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.rooms[room.ID] = room
	return nil
}
