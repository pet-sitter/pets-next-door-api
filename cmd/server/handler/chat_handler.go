package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/pet-sitter/pets-next-door-api/internal/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type ChatHandler struct {
	wsServer    *chat.WsServer
	upgrader    websocket.Upgrader
	authService service.AuthService
	chatService service.ChatService
}

var upgrader = websocket.Upgrader{
	// TODO: 버퍼 사이즈의 근거
	ReadBufferSize:  4096, // 웹소켓 읽기 버퍼 크기
	WriteBufferSize: 4096, // 웹소켓 쓰기 버퍼 크기
}

func NewChatController(
	wsServer *chat.WsServer, authService service.AuthService, chatService service.ChatService,
) *ChatHandler {
	return &ChatHandler{
		wsServer:    wsServer,
		upgrader:    upgrader,
		authService: authService,
		chatService: chatService,
	}
}

func (h *ChatHandler) ServerWebsocket(
	c echo.Context, w http.ResponseWriter, r *http.Request,
) error {
	foundUser, err := h.authService.VerifyAuthAndGetUser(c.Request().Context(), c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	conn, err2 := upgrader.Upgrade(w, r, nil)
	if err2 != nil {
		log.Println(err2)
		return err2
	}

	client := h.initializeOrUpdateClient(conn, foundUser)

	// 클라이언트의 메시지를 읽고 쓰는 데 사용되는 고루틴을 시작
	go client.HandleWrite()
	go client.HandleRead(&h.chatService)

	return nil
}

// 클라이언트를 초기화하거나 기존 클라이언트를 업데이트하는 함수
func (h *ChatHandler) initializeOrUpdateClient(conn *websocket.Conn, userData *user.InternalView) *chat.Client {
	client := h.wsServer.FindClientByUID(userData.FirebaseUID)
	if client == nil {
		// 클라이언트를 찾지 못한 경우 새로운 클라이언트를 생성
		client = chat.NewClient(conn, h.wsServer, userData.Nickname, userData.FirebaseUID)
		// 새 클라이언트를 웹소켓 서버에 등록
		h.wsServer.RegisterClient(client)
	} else {
		// 기존 클라이언트가 있는 경우 연결을 업데이트
		client.UpdateConn(conn)
	}
	return client
}
