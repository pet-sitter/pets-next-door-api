package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/pet-sitter/pets-next-door-api/internal/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type ChatHandler struct {
	wsServer    *chat.WsServer
	upgrader    websocket.Upgrader
	authService service.AuthService
	chatService service.ChatService
}

func NewChatController(
	wsServer *chat.WsServer, authService service.AuthService, chatService service.ChatService,
) *ChatHandler {
	return &ChatHandler{
		wsServer: wsServer,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
			CheckOrigin: func(r *http.Request) bool {
				// TODO: 검사로직 추가
				return true
			},
		},
		authService: authService,
		chatService: chatService,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

func (h *ChatHandler) ServerWebsocket(
	c echo.Context, wsServer *chat.WsServer, w http.ResponseWriter, r *http.Request,
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
	client := chat.NewClient(conn, wsServer, foundUser.Nickname, foundUser.FirebaseUID)

	go client.WritePump()
	go client.ReadPump(&h.chatService)

	return nil
}
