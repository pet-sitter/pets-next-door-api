package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/chat"
	domain "github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type ChatHandler struct {
	stateManager *chat.StateManager
	authService  service.AuthService
	chatService  service.ChatService
}

func NewChatHandler(
	stateManager chat.StateManager,
	authService service.AuthService,
	chatService service.ChatService,
) *ChatHandler {
	return &ChatHandler{
		stateManager: &stateManager,
		authService:  authService,
		chatService:  chatService,
	}
}

func (h ChatHandler) FindRoomByID(c echo.Context) error {
	roomID, err := pnd.ParseIDFromPath(c, "roomID")
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	res, err := h.chatService.FindRoomByID(c.Request().Context(), int64(*roomID))
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (h ChatHandler) CreateRoom(c echo.Context) error {
	var createRoomRequest domain.CreateRoomRequest
	if err := pnd.ParseBody(c, &createRoomRequest); err != nil {
		return c.JSON(err.StatusCode, err)
	}

	res, err := h.chatService.CreateRoom(c.Request().Context(), createRoomRequest.RoomName, createRoomRequest.RoomType)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusCreated, res)
}

func (h ChatHandler) JoinChatRoom(c echo.Context) error {
	foundUser, err := h.authService.VerifyAuthAndGetUser(c.Request().Context(), c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	roomID, err := pnd.ParseIDFromPath(c, "roomID")
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	res, err := h.chatService.JoinRoom(c.Request().Context(), int64(*roomID), foundUser.FirebaseUID)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (h ChatHandler) LeaveChatRoom(c echo.Context) error {
	foundUser, err := h.authService.VerifyAuthAndGetUser(c.Request().Context(), c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	roomID, err := pnd.ParseIDFromPath(c, "roomID")
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	res := h.chatService.LeaveRoom(c.Request().Context(), int64(*roomID), foundUser.FirebaseUID)
	return c.JSON(http.StatusOK, res)
}

func (h ChatHandler) FindAllRooms(c echo.Context) error {
	rooms, err := h.chatService.MockFindAllChatRooms()
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}
	return c.JSON(http.StatusOK, rooms)
}

func (h ChatHandler) FindMessagesByRoomID(c echo.Context) error {
	roomID, err := pnd.ParseIDFromPath(c, "roomID")
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	res, err := h.chatService.MockFindMessagesByRoomID(int64(*roomID))
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusOK, res)
}
