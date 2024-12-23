package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	domain "github.com/pet-sitter/pets-next-door-api/internal/domain/chat"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
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
// @Param roomID path string true "채팅방 ID"
// @Security FirebaseAuth
// @Success 200 {object} domain.RoomSimpleInfo
// @Router /chat/rooms/{roomID} [get]
func (h ChatHandler) FindRoomByID(c echo.Context) error {
	foundUser, err := h.authService.VerifyAuthAndGetUser(
		c.Request().Context(),
		c.Request().Header.Get("Authorization"),
	)
	if err != nil {
		return err
	}

	roomID, err := pnd.ParseIDFromPath(c, "roomID")
	if err != nil {
		return err
	}

	res, err := h.chatService.FindChatRoomByUIDAndRoomID(
		c.Request().Context(),
		foundUser.FirebaseUID,
		roomID,
	)
	if err != nil {
		return err
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
// @Success 201 {object} domain.RoomSimpleInfo
// @Router /chat/rooms [post]
func (h ChatHandler) CreateRoom(c echo.Context) error {
	user, err := h.authService.VerifyAuthAndGetUser(
		c.Request().Context(),
		c.Request().Header.Get("Authorization"),
	)
	if err != nil {
		return err
	}
	var createRoomRequest domain.CreateRoomRequest

	if bodyError := pnd.ParseBody(c, &createRoomRequest); bodyError != nil {
		return err
	}

	res, err := h.chatService.CreateRoom(
		c.Request().Context(),
		createRoomRequest.RoomName,
		createRoomRequest.RoomType,
		user.FirebaseUID,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, res)
}

// JoinChatRoom godoc
// @Summary 채팅방에 참가합니다.
// @Description 채팅방에 참가합니다.
// @Tags chat
// @Accept  json
// @Produce  json
// @Param roomID path string true "채팅방 ID"
// @Security FirebaseAuth
// @Success 200 {object} domain.JoinRoomsView
// @Router /chat/rooms/{roomID}/join [post]
func (h ChatHandler) JoinChatRoom(c echo.Context) error {
	foundUser, err := h.authService.VerifyAuthAndGetUser(
		c.Request().Context(),
		c.Request().Header.Get("Authorization"),
	)
	if err != nil {
		return err
	}

	roomID, err := pnd.ParseIDFromPath(c, "roomID")
	if err != nil {
		return err
	}

	res, err := h.chatService.JoinRoom(c.Request().Context(), roomID, foundUser.FirebaseUID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

// LeaveChatRoom godoc
// @Summary 채팅방을 나갑니다.
// @Description 채팅방을 나갑니다.
// @Tags chat
// @Accept  json
// @Produce  json
// @Param roomID path string true "채팅방 ID"
// @Security FirebaseAuth
// @Success 200
// @Router /chat/rooms/{roomID}/leave [post]
func (h ChatHandler) LeaveChatRoom(c echo.Context) error {
	foundUser, err := h.authService.VerifyAuthAndGetUser(
		c.Request().Context(),
		c.Request().Header.Get("Authorization"),
	)
	if err != nil {
		return err
	}

	roomID, err := pnd.ParseIDFromPath(c, "roomID")
	if err != nil {
		return err
	}

	res := h.chatService.LeaveRoom(c.Request().Context(), roomID, foundUser.FirebaseUID)
	return c.JSON(http.StatusOK, res)
}

// FindAllRooms godoc
// @Summary 사용자의 채팅방 목록을 조회합니다.
// @Description 사용자의 채팅방 목록을 조회합니다.
// @Tags chat
// @Accept  json
// @Produce  json
// @Security FirebaseAuth
// @Success 200 {object} domain.JoinRoomsView
// @Router /chat/rooms [get]
func (h ChatHandler) FindAllRooms(c echo.Context) error {
	foundUser, err := h.authService.VerifyAuthAndGetUser(
		c.Request().Context(),
		c.Request().Header.Get("Authorization"),
	)
	if err != nil {
		return err
	}

	rooms, err := h.chatService.FindAllByUserUID(c.Request().Context(), foundUser.FirebaseUID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, rooms)
}

// FindMessagesByRoomID godoc
// @Summary 채팅방의 메시지 목록을 조회합니다.
// @Description 채팅방의 메시지 목록을 조회합니다.
// @Tags chat
// @Accept  json
// @Produce  json
// @Param roomID path string true "채팅방 ID"
// @Param prev query int false "이전 페이지"
// @Param next query int false "다음 페이지"
// @Param size query int false "페이지 사이즈" default(30)
// @Security FirebaseAuth
// @Success 200 {object} domain.MessageCursorView
// @Router /chat/rooms/{roomID}/messages [get]
func (h ChatHandler) FindMessagesByRoomID(c echo.Context) error {
	roomID, err := pnd.ParseIDFromPath(c, "roomID")
	if err != nil {
		return err
	}

	prev, next, limit, err := pnd.ParseCursorPaginationQueries(c, 30)
	if err != nil {
		return err
	}

	res, err := h.chatService.FindChatRoomMessagesByRoomID(
		c.Request().Context(),
		roomID,
		prev,
		next,
		int64(limit),
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
