package wschat

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
	"github.com/rs/zerolog/log"
)

type WSServer struct {
	// key: UserID, value: WSClient
	clients   map[uuid.UUID]WSClient
	broadcast chan MessageRequest
	upgrader  websocket.Upgrader

	authService  service.AuthService
	mediaService service.MediaService
	chatService  service.ChatService
}

func NewWSServer(
	upgrader websocket.Upgrader,
	authService service.AuthService,
	mediaService service.MediaService,
	chatService service.ChatService,
) *WSServer {
	return &WSServer{
		clients:      make(map[uuid.UUID]WSClient),
		broadcast:    make(chan MessageRequest),
		upgrader:     upgrader,
		authService:  authService,
		mediaService: mediaService,
		chatService:  chatService,
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

	foundUser, err := s.authService.VerifyAuthAndGetUser(
		c.Request().Context(),
		c.Request().Header.Get("Authorization"),
	)
	if err != nil {
		return err
	}
	userID := foundUser.ID

	conn, err := s.upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close connection")
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

			return c.JSON(
				http.StatusInternalServerError,
				map[string]interface{}{"error": err.Error()},
			)
		}

		s.broadcast <- msgReq
	}
}

// Broadcast messages to all clients
// TODO : 각 메시지 별 타입에 따라, 메시지 처리하고 있음 (plain, media) 저장로직 추가해야 함
func (s *WSServer) LoopOverClientMessages() {
	log.Info().Msg("Looping over client messages")
	ctx := context.Background()

	for {
		msgReq := <-s.broadcast

		for _, client := range s.clients {
			log.Info().Msg(
				"Message from user: " +
					client.userID.String() +
					" to user: " + msgReq.Sender.ID.String())

			// Message print
			log.Info().Msg("Message: " + msgReq.String())

			// TODO: Check if the message is for the room
			var msg MessageResponse
			switch msgReq.MessageType {
			case "plain":
				_, err := s.chatService.SaveChatMessage(
					ctx,
					msgReq.Sender.ID,
					msgReq.Room.ID,
					msgReq.MessageType,
					msgReq.Message,
				)
				if err != nil {
					log.Error().Err(err).Msg("Failed to save chat message")
					msg = NewErrorMessageResponse(
						msgReq.MessageID,
						msgReq.Sender,
						msgReq.Room,
						"Failed to save chat message",
						time.Now(),
					)
				} else {
					msg = NewPlainMessageResponse(
						msgReq.MessageID,
						msgReq.Sender,
						msgReq.Room,
						msgReq.Message,
						time.Now(),
					)
				}

			case "media":
				if len(msgReq.Medias) == 0 {
					log.Error().Msg("No media found")
					msg = NewErrorMessageResponse(
						msgReq.MessageID,
						msgReq.Sender,
						msgReq.Room,
						"No media found",
						time.Now(),
					)
				} else {
					ids := make([]uuid.UUID, 0)
					for _, mediaReq := range msgReq.Medias {
						ids = append(ids, mediaReq.ID)
					}
					medias, err := s.mediaService.FindMediasByIDs(ctx, ids)
					if err != nil {
						log.Error().Err(err).Msg("Failed to find media")
						msg = NewErrorMessageResponse(msgReq.MessageID, msgReq.Sender, msgReq.Room, "Failed to find media", time.Now())
					} else {
						msg = NewMediaMessageResponse(msgReq.MessageID, msgReq.Sender, msgReq.Room, medias, time.Now())
					}
				}
			default:
				log.Error().Msg("Unknown message type")
				return
			}

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
	userID uuid.UUID
}

func NewWSClient(
	conn *websocket.Conn,
	userID uuid.UUID,
) WSClient {
	return WSClient{conn, userID}
}

func (c *WSClient) WriteJSON(v interface{}) error {
	return c.conn.WriteJSON(v)
}

func (c *WSClient) Close() error {
	return c.conn.Close()
}

type MediaRequest struct {
	ID uuid.UUID `json:"id"`
}

type MessageRequest struct {
	Sender      Sender         `json:"sender"`
	Room        Room           `json:"room"`
	MessageID   string         `json:"messageId"`
	MessageType string         `json:"messageType"`
	Medias      []MediaRequest `json:"medias,omitempty"`
	Message     string         `json:"message"`
}

func (m MessageRequest) String() string {
	return "Sender: " + m.Sender.ID.String() +
		" Room: " + m.Room.ID.String() +
		" MessageID: " + m.MessageID +
		" MessageType: " + m.MessageType +
		" Message: " + m.Message +
		" Medias: " + strconv.Itoa(len(m.Medias))
}

type MessageResponse struct {
	Sender      Sender             `json:"sender"`
	Room        Room               `json:"room"`
	MessageID   string             `json:"messageId"`
	MessageType string             `json:"messageType"`
	Medias      []media.DetailView `json:"medias,omitempty"`
	Message     string             `json:"message"`
	CreatedAt   string             `json:"createdAt"`
	UpdatedAt   string             `json:"updatedAt"`
}

type Sender struct {
	ID uuid.UUID `json:"id"`
}

type Room struct {
	ID uuid.UUID `json:"id"`
}

type Media struct {
	ID        uuid.UUID `json:"id"`
	MediaType string    `json:"type"`
	URL       string    `json:"url"`
}

func NewPlainMessageResponse(
	messageID string,
	sender Sender,
	room Room,
	message string,
	now time.Time,
) MessageResponse {
	return MessageResponse{
		MessageID:   messageID,
		Sender:      sender,
		Room:        room,
		MessageType: "plain",
		Message:     message,
		CreatedAt:   now.Format(time.RFC3339),
		UpdatedAt:   now.Format(time.RFC3339),
	}
}

func NewErrorMessageResponse(
	messageID string,
	sender Sender,
	room Room,
	message string,
	now time.Time,
) MessageResponse {
	return MessageResponse{
		MessageID:   messageID,
		Sender:      sender,
		Room:        room,
		MessageType: "error",
		Message:     message,
		CreatedAt:   now.Format(time.RFC3339),
		UpdatedAt:   now.Format(time.RFC3339),
	}
}

func NewMediaMessageResponse(
	messageID string,
	sender Sender,
	room Room,
	medias []media.DetailView,
	now time.Time,
) MessageResponse {
	return MessageResponse{
		MessageID:   messageID,
		Sender:      sender,
		Room:        room,
		MessageType: "media",
		Medias:      medias,
		CreatedAt:   now.Format(time.RFC3339),
		UpdatedAt:   now.Format(time.RFC3339),
	}
}
