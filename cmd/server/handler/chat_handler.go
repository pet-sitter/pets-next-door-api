package handler

import (
	"github.com/labstack/echo/v4"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	domain "github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
	"net/http"
)

type ChatHandler struct {
	authService service.AuthService
	chatService service.ChatService
}

func NewChatHandler(
	authService service.AuthService,
	chatService service.ChatService,
) *ChatHandler {
	return &ChatHandler{
		authService: authService,
		chatService: chatService,
	}
}

// FindRoomByID godoc
// @Summary 채팅방을 조회합니다.
// @Description
// @Tags chat
// @Accept  json
// @Produce  json
// @Param roomID path int true "채팅방 ID"
// @Security FirebaseAuth
// @Success 200 {object} domain.Room
// @Router /chat/rooms/{roomID} [get]
func (h ChatHandler) FindRoomByID(c echo.Context) error {
	foundUser, err := h.authService.VerifyAuthAndGetUser(c.Request().Context(), c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	roomID, err := pnd.ParseIDFromPath(c, "roomID")
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	res, err := h.chatService.FindChatRoomByUIDAndRoomID(c.Request().Context(), foundUser.FirebaseUID, int64(*roomID))
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusOK, res)
}

// CreateRoom godoc
// @Summary 채팅방을 생성합니다.
// @Description
// @Tags chat
// @Accept  json
// @Produce  json
// @Param request body domain.CreateRoomRequest true "채팅방 생성 요청"
// @Security FirebaseAuth
// @Success 201 {object} domain.Room
// @Router /chat/rooms [post]
func (h ChatHandler) CreateRoom(c echo.Context) error {
	var createRoomRequest domain.CreateRoomRequest

	if err := pnd.ParseBody(c, &createRoomRequest); err != nil {
		return c.JSON(err.StatusCode, err)
	}

	res, err := h.chatService.CreateRoom(
		c.Request().Context(),
		createRoomRequest.RoomName,
		createRoomRequest.RoomType,
		createRoomRequest.JoinUserIds,
	)

	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusCreated, res)
}

// JoinChatRoom godoc
// @Summary 채팅방에 참가합니다.
// @Description 채팅방에 참가합니다.
// @Tags chat
// @Accept  json
// @Produce  json
// @Param roomID path int true "채팅방 ID"
// @Security FirebaseAuth
// @Success 200 {object} domain.JoinRoomsView
// @Router /chat/rooms/{roomID}/join [post]
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

// LeaveChatRoom godoc
// @Summary 채팅방을 나갑니다.
// @Description 채팅방을 나갑니다.
// @Tags chat
// @Accept  json
// @Produce  json
// @Param roomID path int true "채팅방 ID"
// @Security FirebaseAuth
// @Success 200
// @Router /chat/rooms/{roomID}/leave [post]
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

// FindAllRooms godoc
// @Summary 유저가 소속되어 있는 모든 채팅방을 조회합니다.
// @Description 유저가 소속되어 있는 채팅방 전체 목록을 조회합니다.
// @Tags chat
// @Accept  json
// @Produce  json
// @Security FirebaseAuth
// @Success 200 {object} []domain.Room
// @Router /chat/rooms [get]
func (h ChatHandler) FindAllRooms(c echo.Context) error {
	foundUser, err := h.authService.VerifyAuthAndGetUser(c.Request().Context(), c.Request().Header.Get("Authorization"))

	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	rooms, err := h.chatService.FindAllByUserUID(c.Request().Context(), foundUser.FirebaseUID)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}
	return c.JSON(http.StatusOK, rooms)
}

// FindMessagesByRoomID godoc
// @Summary 채팅방의 메시지를 조회합니다.
// @Description 채팅방의 메시지를 조회합니다.
// @Tags chat
// @Accept  json
// @Produce  json
// @Param roomID path int true "채팅방 ID"
// @Security FirebaseAuth
// @Success 200 {object} []domain.Message
// @Router /chat/rooms/{roomID}/messages [get]
func (h ChatHandler) FindMessagesByRoomID(c echo.Context) error {
	roomID, err := pnd.ParseIDFromPath(c, "roomID")
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	prev, next, limit, appError := pnd.ParseCursorPaginationQueries(c, 30)
	if appError != nil {
		return c.JSON(appError.StatusCode, appError)
	}

	res, err := h.chatService.FindChatRoomMessagesByRoomID(
		c.Request().Context(),
		int64(*roomID),
		int64(prev),
		int64(next),
		int64(limit),
	)

	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	return c.JSON(http.StatusOK, res)
}
