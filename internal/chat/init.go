package chat

import (
	"context"
	"log"

	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

// 서버가 시작되거나 재시작될 때, 채널 상태 롤백
func InitializeWebSocketServer(ctx context.Context, wsServer *WsServer, chatService *service.ChatService) {
	rows, err := chatService.FindUserChatRoom(ctx)
	if err != nil {
		log.Println("Error finding user chat rooms:", err)
		return
	}

	// 클라이언트를 중복 생성하지 않도록 관리하는 맵
	clientMap := make(map[string]*Client)

	for _, row := range rows {
		// 클라이언트를 생성하거나 기존 클라이언트를 재사용
		client, exists := clientMap[row.UserInfo.FirebaseUID]
		if !exists {
			client = NewClient(nil, wsServer, row.UserInfo.Nickname, row.UserInfo.FirebaseUID)
			wsServer.RegisterClient(client)
			clientMap[row.UserInfo.FirebaseUID] = client
		}

		// 방을 생성하거나 기존 방을 불러옴
		room := wsServer.findRoomByID(row.RoomInfo.ID)
		if room == nil {
			room = room.InitRoom(row.RoomInfo.ID, row.RoomInfo.Name, row.RoomInfo.RoomType)
			wsServer.rooms[room] = true
			go room.RunRoom(chatService)
		}

		// 클라이언트를 방에 등록
		if !client.isInRoom(room) {
			client.rooms[room] = true
			room.register <- client
		}
	}
}
