package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"github.com/pet-sitter/pets-next-door-api/internal/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type ChatHandler struct {
	wsServer     *chat.WebSocketServer
	upgrader     websocket.Upgrader
	stateManager *chat.StateManager
	authService  service.AuthService
	chatService  service.ChatService
}

var upgrader = websocket.Upgrader{
	// 사이즈 단위: byte
	// 영어 기준: 2048 글자
	// 한국어(글자당 3바이트) 기준: 약 682 글자
	// 영어 + 한국어 기준(글자당 2바이트로 가정): 약 1024 글자
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}

func NewChatController(
	wsServer *chat.WebSocketServer,
	stateManager chat.StateManager,
	authService service.AuthService,
	chatService service.ChatService,
) *ChatHandler {
	return &ChatHandler{
		wsServer:     wsServer,
		upgrader:     upgrader,
		stateManager: &stateManager,
		authService:  authService,
		chatService:  chatService,
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
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err2.Error(),
		})
	}

	client := h.initializeOrUpdateClient(conn, foundUser)

	// 클라이언트의 메시지를 읽고 쓰는 데 사용되는 고루틴을 시작 (비동기)
	go client.HandleWrite()
	go client.HandleRead(*h.stateManager, &h.chatService)

	return nil
}

// 클라이언트를 초기화하거나 기존 클라이언트를 업데이트하는 함수
func (h *ChatHandler) initializeOrUpdateClient(
	conn *websocket.Conn, userData *user.InternalView,
) *chat.Client {
	client := h.wsServer.StateManager.FindClientByUID(userData.FirebaseUID)
	if client == nil {
		client = chat.NewClient(conn, userData.Nickname, userData.FirebaseUID)
		h.wsServer.StateManager.RegisterClient(client)
	} else {
		// 기존 클라이언트가 있는 경우 연결을 업데이트
		client.UpdateConn(conn)
	}
	return client
}
