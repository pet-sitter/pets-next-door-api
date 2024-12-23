package handler

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	pnd "github.com/pet-sitter/pets-next-door-api/api"
	utils "github.com/pet-sitter/pets-next-door-api/internal/common"
	"github.com/pet-sitter/pets-next-door-api/internal/domain/event"
	databasegen "github.com/pet-sitter/pets-next-door-api/internal/infra/database/gen"
	"github.com/pet-sitter/pets-next-door-api/internal/service"
)

type EventHandler struct {
	authService  service.AuthService
	eventService service.EventService
}

func NewEventHandler(
	authService service.AuthService,
	eventService service.EventService,
) *EventHandler {
	return &EventHandler{
		authService:  authService,
		eventService: eventService,
	}
}

// FindEvents godoc
// @Summary 이벤트를 조회합니다.
// @Description
// @Tags events
// @Accept  json
// @Produce  json
// @Param author_id query string false "작성자 ID"
// @Param prev query int false "이전 페이지"
// @Param next query int false "다음 페이지"
// @Param size query int false "페이지 사이즈" default(20)
// @Success 200 {object} pnd.CursorPaginatedView[event.View]
// @Router /events [get]
func (h *EventHandler) FindEvents(c echo.Context) error {
	prev, next, size, err := pnd.ParseCursorPaginationQueries(c, 20)
	if err != nil {
		return err
	}
	authorID, err := pnd.ParseOptionalUUIDQuery(c, "author_id")
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	events, err := h.eventService.FindEvents(ctx, databasegen.FindEventsParams{
		AuthorID: authorID,
		Prev:     prev,
		Next:     next,
		Limit:    int32(size),
	})
	if err != nil {
		return err
	}

	items := make([]event.ShortTermView, len(events))
	for i, e := range events {
		// Turn topics (string[]) to EventTopic[]
		topics := make([]event.EventTopic, len(e.Event.Topics))
		for i, topic := range e.Event.Topics {
			topics[i] = event.EventTopic(topic)
		}

		items[i] = event.ShortTermView{
			BaseView: event.BaseView{
				ID:              e.Event.ID,
				EventType:       event.EventType(e.Event.EventType),
				Author:          e.Author,
				Name:            e.Event.Name,
				Description:     e.Event.Description,
				Media:           *e.Media,
				Topics:          topics,
				MaxParticipants: utils.NullInt32ToIntPtr(e.Event.MaxParticipants),
				Fee:             int(e.Event.Fee),
				StartAt:         utils.NullTimeToTimePtr(e.Event.StartAt),
				CreatedAt:       e.Event.CreatedAt,
				UpdatedAt:       e.Event.UpdatedAt,
			},
		}
	}
	return c.JSON(
		http.StatusOK,
		pnd.CursorPaginatedView[event.ShortTermView]{
			Items: items,
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
		return err
	}

	ctx := c.Request().Context()
	eventData, err := h.eventService.FindEvent(
		ctx,
		databasegen.FindEventParams{ID: uuid.NullUUID{UUID: id, Valid: true}},
	)
	if err != nil {
		return err
	}

	topics := make([]event.EventTopic, len(eventData.Event.Topics))
	for i, topic := range eventData.Event.Topics {
		topics[i] = event.EventTopic(topic)
	}
	view := event.ShortTermView{
		BaseView: event.BaseView{
			ID:              eventData.Event.ID,
			Author:          eventData.Author,
			EventType:       event.EventType(eventData.Event.EventType),
			Name:            eventData.Event.Name,
			Description:     eventData.Event.Description,
			Media:           *eventData.Media,
			Topics:          topics,
			MaxParticipants: utils.NullInt32ToIntPtr(eventData.Event.MaxParticipants),
			Fee:             int(eventData.Event.Fee),
			StartAt:         utils.NullTimeToTimePtr(eventData.Event.StartAt),
			CreatedAt:       eventData.Event.CreatedAt,
			UpdatedAt:       eventData.Event.UpdatedAt,
		},
	}

	return c.JSON(http.StatusOK, view)
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
		return err
	}
	authorID := foundUser.ID

	var reqBody event.CreateRequest
	if err := pnd.ParseBody(c, &reqBody); err != nil {
		return err
	}

	ctx := c.Request().Context()
	created, err := h.eventService.CreateEvent(ctx, authorID, reqBody)
	if err != nil {
		return err
	}

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
		return err
	}
	uid := foundUser.FirebaseUID

	var reqBody event.UpdateRequest
	if err := pnd.ParseBody(c, &reqBody); err != nil {
		return err
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
		return err
	}
	uid := foundUser.FirebaseUID

	id, err := pnd.ParseIDFromPath(c, "id")
	if err != nil {
		return err
	}

	log.Printf("uid: %s, id: %s", uid, id)
	// TODO: Implement delete event logic

	return c.JSON(http.StatusOK, nil)
}
