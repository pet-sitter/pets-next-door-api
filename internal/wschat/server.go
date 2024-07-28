package wschat

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
	"github.com/rs/zerolog/log"
)

type WSServer struct {
	// key: UserID, value: WSClient
	clients   map[int64]WSClient
	broadcast chan MessageRequest
	upgrader  websocket.Upgrader

	authService service.AuthService
}

func NewWSServer(
	upgrader websocket.Upgrader,
	authService service.AuthService,
) *WSServer {
	return &WSServer{
		clients:     make(map[int64]WSClient),
		broadcast:   make(chan MessageRequest),
		upgrader:    upgrader,
		authService: authService,
	}
}

func NewDefaultUpgrader() websocket.Upgrader {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(_ *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return upgrader
}

// Server-side WebSocket handler
func (s *WSServer) HandleConnections(
	c echo.Context,
) error {
	log.Info().Msg("Handling connections")

	foundUser, err2 := s.authService.VerifyAuthAndGetUser(c.Request().Context(), c.Request().Header.Get("Authorization"))
	if err2 != nil {
		return c.JSON(err2.StatusCode, err2)
	}
	userID := foundUser.ID

	conn, err := s.upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	defer func() {
		err2 := conn.Close()
		if err2 != nil {
			log.Error().Err(err2).Msg("Failed to close connection")
		}
		delete(s.clients, userID)
	}()

	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade connection")
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	client := NewWSClient(conn, userID)
	s.clients[userID] = client

	for {
		var msgReq MessageRequest
		err := conn.ReadJSON(&msgReq)
		msgReq.Sender = Sender{ID: userID}
		if err != nil {
			log.Error().Err(err).Msg("Failed to read message")
			delete(s.clients, userID)

			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		s.broadcast <- msgReq
	}
}

// Broadcast messages to all clients
func (s *WSServer) LoopOverClientMessages() {
	log.Info().Msg("Looping over client messages")

	for {
		msgReq := <-s.broadcast

		for _, client := range s.clients {
			log.Info().Msg(
				"Message from user: " +
					strconv.Itoa(int(client.userID)) +
					" to user: " + strconv.Itoa(int(msgReq.Sender.ID)))

			// Filter messages from the same user
			if client.userID == msgReq.Sender.ID {
				continue
			}

			// TODO: Check if the message is for the room
			msg := NewPlainMessageResponse(msgReq.Sender, msgReq.Room, msgReq.Message, time.Now())

			if err := client.WriteJSON(msg); err != nil {
				// No way but to close the connection
				log.Error().Err(err).Msg("Failed to write message")
				err := client.Close()
				if err != nil {
					log.Error().Err(err).Msg("Failed to close connection")
					delete(s.clients, client.userID)
					return
				}
				delete(s.clients, client.userID)
				return
			}
		}
	}
}

type WSClient struct {
	conn   *websocket.Conn
	userID int64
}

func NewWSClient(
	conn *websocket.Conn,
	userID int64,
) WSClient {
	return WSClient{conn, userID}
}

func (c *WSClient) WriteJSON(v interface{}) error {
	return c.conn.WriteJSON(v)
}

func (c *WSClient) Close() error {
	return c.conn.Close()
}

type MessageRequest struct {
	Sender      Sender `json:"user"`
	Room        Room   `json:"room"`
	MessageType string `json:"messageType"`
	Media       *Media `json:"media,omitempty"`
	Message     string `json:"message"`
}

type MessageResponse struct {
	Sender      Sender `json:"user"`
	Room        Room   `json:"room"`
	MessageType string `json:"messageType"`
	Media       *Media `json:"media,omitempty"`
	Message     string `json:"message"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type Sender struct {
	ID int64 `json:"id"`
}

type Room struct {
	ID int64 `json:"id"`
}

type Media struct {
	ID        int64  `json:"id"`
	MediaType string `json:"type"`
	URL       string `json:"url"`
}

func NewPlainMessageResponse(sender Sender, room Room, message string, now time.Time) MessageResponse {
	return MessageResponse{
		Sender:      sender,
		Room:        room,
		MessageType: "plain",
		Message:     message,
		CreatedAt:   now.Format(time.RFC3339),
		UpdatedAt:   now.Format(time.RFC3339),
	}
}

func NewMediaMessageResponse(sender Sender, room Room, media *Media, now time.Time) MessageResponse {
	return MessageResponse{
		Sender:      sender,
		Room:        room,
		MessageType: "media",
		Media:       media,
		CreatedAt:   now.Format(time.RFC3339),
		UpdatedAt:   now.Format(time.RFC3339),
	}
}
