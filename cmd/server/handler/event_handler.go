package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/event"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/media"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/user"
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

func generateDummyEvent() event.ShortTermView {
	profileImageURL := "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png"
	now := time.Now()
	startAt := now.AddDate(0, 0, 1)
	maxParticipants := 3
	return event.ShortTermView{
		BaseView: event.BaseView{
			ID:        uuid.New(),
			EventType: event.ShortTerm,
			Author: user.WithoutPrivateInfo{
				ID:              uuid.New(),
				Nickname:        "멍냥이",
				ProfileImageURL: &profileImageURL,
			},
			Name:        "name",
			Description: "description",
			Media: media.DetailView{
				ID:        uuid.New(),
				MediaType: media.TypeImage,
				URL: "https://images.unsplash.com/" +
					"photo-1493225457124-a3eb161ffa5f?q=80&w=2970&auto=format" +
					"&fit=crop&ixlib=rb-4.0.3&ixid=" +
					"M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
				CreatedAt: now,
			},
			Topics:          []event.EventTopic{event.ETC},
			MaxParticipants: &maxParticipants,
			Fee:             10000,
			StartAt:         &startAt,
			CreatedAt:       now,
			UpdatedAt:       now,
		},
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
// @Success 200 {object} pnd.CursorPaginatedView[event.View]
// @Router /events [get]
func (h *EventHandler) FindEvents(c echo.Context) error {
	return c.JSON(
		http.StatusOK,
		pnd.CursorPaginatedView[event.ShortTermView]{
			Items: []event.ShortTermView{generateDummyEvent()},
		},
	)
}

// FindEventByID godoc
// @Summary ID로 이벤트를 조회합니다.
// @Description
// @Tags events
// @Produce  json
// @Param id path string true "이벤트 ID"
// @Success 200 {object} event.View
// @Router /events/{id} [get]
func (h *EventHandler) FindEventByID(c echo.Context) error {
	id, err := pnd.ParseIDFromPath(c, "id")
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	res := generateDummyEvent()
	res.ID = id
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
// @Success 201 {object} event.View
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

	created := generateDummyEvent()
	created.ID = uuid.New()

	return c.JSON(http.StatusCreated, created)
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
// @Param id path string true "이벤트 ID"
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
