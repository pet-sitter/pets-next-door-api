package handler

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/event"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type EventHandler struct {
	authService service.AuthService
}

func NewEventHandler(authService service.AuthService) *EventHandler {
	return &EventHandler{
		authService: authService,
	}
}

func generateDummyEvent() event.ListView {
	return event.ListView{
		ID: uuid.New(),
	}
}

// FindEvents godoc
// @Summary 이벤트를 조회합니다.
// @Description
// @Tags events
// @Accept  json
// @Produce  json
// @Param author_id query string false "작성자 ID"
// @Param page query int false "페이지 번호" default(1)
// @Param size query int false "페이지 사이즈" default(20)
// @Success 200 {object} pnd.CursorPaginatedView[event.ListView]
// @Router /events [get]
func (h *EventHandler) FindEvents(c echo.Context) error {
	return c.JSON(
		http.StatusOK,
		pnd.CursorPaginatedView[event.ListView]{Items: []event.ListView{generateDummyEvent()}},
	)
}

// FindEventByID godoc
// @Summary ID로 이벤트를 조회합니다.
// @Description
// @Tags events
// @Produce  json
// @Param id path int true "이벤트 ID"
// @Success 200 {object} event.DetailView
// @Router /events [get]
func (h *EventHandler) FindEventByID(c echo.Context) error {
	id, err := pnd.ParseIDFromPath(c, "id")
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	res := event.DetailView{
		ID: id,
	}

	return c.JSON(http.StatusOK, res)
}

// CreateEvent  godoc
// @Summary 이벤트를 생성합니다.
// @Description
// @Tags events
// @Accept  json
// @Produce  json
// @Param request body event.CreateRequest true "이벤트 생성 요청"
// @Security FirebaseAuth
// @Success 201 {object} event.DetailView
// @Router /events [post]
func (h *EventHandler) CreateEvent(c echo.Context) error {
	foundUser, err := h.authService.VerifyAuthAndGetUser(
		c.Request().Context(),
		c.Request().Header.Get("Authorization"),
	)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}
	uid := foundUser.FirebaseUID

	var reqBody event.CreateRequest
	if err := pnd.ParseBody(c, &reqBody); err != nil {
		return c.JSON(err.StatusCode, err)
	}

	log.Printf("uid: %s, reqBody: %+v", uid, reqBody)
	// TODO: Implement create event logic

	return c.JSON(http.StatusCreated, event.DetailView{ID: uuid.New()})
}

// UpdateEvent godoc
// @Summary 이벤트를 수정합니다.
// @Description
// @Tags events
// @Accept  json
// @Produce  json
// @Security FirebaseAuth
// @Param request body event.UpdateRequest true "이벤트 수정 요청"
// @Success 200
// @Router /events [put]
func (h *EventHandler) UpdateEvent(c echo.Context) error {
	foundUser, err := h.authService.VerifyAuthAndGetUser(
		c.Request().Context(),
		c.Request().Header.Get("Authorization"),
	)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}
	uid := foundUser.FirebaseUID

	var reqBody event.UpdateRequest
	if err := pnd.ParseBody(c, &reqBody); err != nil {
		return c.JSON(err.StatusCode, err)
	}

	log.Printf("uid: %s, reqBody: %+v", uid, reqBody)
	// TODO: Implement update event logic

	return c.JSON(http.StatusOK, nil)
}

// DeleteEvent godoc
// @Summary 이벤트를 삭제합니다.
// @Description
// @Tags events
// @Security FirebaseAuth
// @Param id path int true "이벤트 ID"
// @Success 200
// @Router /events/{id} [delete]
func (h *EventHandler) DeleteEvent(c echo.Context) error {
	foundUser, err := h.authService.VerifyAuthAndGetUser(
		c.Request().Context(),
		c.Request().Header.Get("Authorization"),
	)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}
	uid := foundUser.FirebaseUID

	id, err := pnd.ParseIDFromPath(c, "id")
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	log.Printf("uid: %s, id: %s", uid, id)
	// TODO: Implement delete event logic

	return c.JSON(http.StatusOK, nil)
}
