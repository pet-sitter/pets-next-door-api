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

func NewChatController(
	wsServer *chat.WsServer, authService service.AuthService, chatService service.ChatService,
) *ChatHandler {
	return &ChatHandler{
		wsServer: wsServer,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
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

	go client.WritePump()
	go client.ReadPump(&h.chatService)

	return nil
}

func (h *ChatHandler) initializeOrUpdateClient(conn *websocket.Conn, userData *user.InternalView) *chat.Client {
	client := h.wsServer.FindClientByUID(userData.FirebaseUID)
	if client == nil {
		client = chat.NewClient(conn, h.wsServer, userData.Nickname, userData.FirebaseUID)
		h.wsServer.RegisterClient(client)
	} else {
		client.UpdateConn(conn)
	}
	return client
}
