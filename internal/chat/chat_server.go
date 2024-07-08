package chat

import (
	"log"

	"github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type WsServer struct {
	// TODO: 전반적인 타입에 대한 고민. 특히, clients의 대한 타입에 대한 고민.
	clients          map[*Client]bool   // 접속한 클라이언트들을 관리하는 맵
	clientUIDMapping map[string]*Client // UID로 클라이언트를 찾기 위한 맵
	register         chan *Client       // 클라이언트 등록 채널
	unregister       chan *Client       // 클라이언트 해제 채널
	broadcast        chan []byte        // 모든 클라이언트에게 브로드캐스트 메시지를 보내기 위한 채널
	// 생성된 방을 관리하는 맵, 키를 *Room로 둬서 채팅방을 추가하거나 삭제할 때 o(1), true/false로 방 활성상태를 알 수 있음
	rooms map[*Room]bool
}

func NewWebsocketServer() *WsServer {
	return &WsServer{
		clients:          make(map[*Client]bool),
		clientUIDMapping: make(map[string]*Client),
		register:         make(chan *Client),
		unregister:       make(chan *Client),
		broadcast:        make(chan []byte),
		rooms:            make(map[*Room]bool),
	}
}

func (server *WsServer) Run() {
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

// 새로운 클라이언트 등록
func (server *WsServer) RegisterClient(client *Client) {
	server.notifyClientJoined(client)
	server.listOnlineClients(client)
	server.clients[client] = true
	server.clientUIDMapping[client.FbUID] = client
}

// 클라이언트 해제
func (server *WsServer) unregisterClient(client *Client) {
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
		server.notifyClientLeft(client)
	}
}

// UID로 클라이언트 찾음
func (server *WsServer) FindClientByUID(uid string) *Client {
	return server.clientUIDMapping[uid]
}

// 클라이언트가 접속했음을 모든 클라이언트에게 알림
func (server *WsServer) notifyClientJoined(client *Client) {
	message := &Message{
		Action: UserJoinedAction,
		Sender: client,
	}

	server.broadcastToClients(message.encode())
}

// 클라이언트가 떠났음을 모든 클라이언트에게 알림
func (server *WsServer) notifyClientLeft(client *Client) {
	message := &Message{
		Action: UserLeftAction,
		Sender: client,
	}

	server.broadcastToClients(message.encode())
}

// 현재 온라인 상태인 클라이언트들을 새 클라이언트에게 전송
func (server *WsServer) listOnlineClients(client *Client) {
	for existingClient := range server.clients {
		message := &Message{
			Action: UserJoinedAction,
			Sender: existingClient,
		}
		client.send <- message.encode()
	}
}

// 주어진 메시지를 모든 클라이언트에게 전송하는 함수
func (server *WsServer) broadcastToClients(message []byte) {
	for client := range server.clients {
		client.send <- message
	}
}

// 방 ID로 방을 찾음
func (server *WsServer) findRoomByID(roomID int64) *Room {
	var foundRoom *Room
	for room := range server.rooms {
		if room.GetID() == roomID {
			foundRoom = room
			break
		}
	}
	return foundRoom
}

// 새로운 방을 생성
func (server *WsServer) createRoom(name string, roomType chat.RoomType, roomService *service.ChatService) *Room {
	room, err := NewRoom(name, roomType, roomService)
	if err != nil {
		log.Println(err)
	}
	// 고루틴으로 실행하면 방이 독립적으로 동작
	go room.RunRoom(roomService)
	server.rooms[room] = true

	return room
}
